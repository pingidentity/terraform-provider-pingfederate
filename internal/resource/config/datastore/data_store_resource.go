package datastore

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

var (
	_ resource.Resource                = &dataStoreResource{}
	_ resource.ResourceWithConfigure   = &dataStoreResource{}
	_ resource.ResourceWithImportState = &dataStoreResource{}

	customId = "data_store_id"
)

func DataStoreResource() resource.Resource {
	return &dataStoreResource{}
}

type dataStoreResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *dataStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Resource to create and manage data stores.",
		Attributes: map[string]schema.Attribute{
			"mask_attribute_values": schema.BoolAttribute{
				Description: "Whether attribute values should be masked in the log. Default value is `false`.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"data_store_id": schema.StringAttribute{
				Description: "The persistent, unique ID for the data store. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					configvalidators.PingFederateId(),
				},
			},
			"custom_data_store":                toSchemaCustomDataStore(),
			"jdbc_data_store":                  toSchemaJdbcDataStore(),
			"ldap_data_store":                  toSchemaLdapDataStore(),
			"ping_one_ldap_gateway_data_store": toSchemaPingOneLdapGatewayDataStore(),
		},
	}
	id.ToSchema(&schema)

	resp.Schema = schema
}

func (r *dataStoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (r *dataStoreResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *dataStoreResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *dataStoreModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	var respDiags diag.Diagnostics

	if plan == nil {
		return
	}

	// Validating attributes that depend on a specific version of PF
	// Compare to version 11.3 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1130)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	// Compare to version 12.1 of PF
	pfVersionAtLeast113 := compare >= 0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1210)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast121 := compare >= 0

	if !pfVersionAtLeast113 {
		// Prior to 11.3, the user_name field is required for jdbc data stores
		if internaltypes.IsDefined(plan.JdbcDataStore) {
			username := plan.JdbcDataStore.Attributes()["user_name"]
			if !internaltypes.IsDefined(username) {
				resp.Diagnostics.AddError("'user_name' is required for JDBC data stores prior to PingFederate 11.3", "")
			}
		}
		// The ldap data store client_tls_certificate_ref and retry_failed_operations attributes require PF 11.3
		if internaltypes.IsDefined(plan.LdapDataStore) {
			ldapDataStoreAttrs := plan.LdapDataStore.Attributes()
			clientTlsCertificateRef := ldapDataStoreAttrs["client_tls_certificate_ref"]
			if internaltypes.IsDefined(clientTlsCertificateRef) {
				resp.Diagnostics.AddError("Attribute 'client_tls_certificate_ref' not supported for LDAP data stores by PingFederate version "+string(r.providerConfig.ProductVersion), "PF 11.3 or later required")
			}
			retryFailedOperations := ldapDataStoreAttrs["retry_failed_operations"].(types.Bool)
			if internaltypes.IsDefined(retryFailedOperations) {
				resp.Diagnostics.AddError("Attribute 'retry_failed_operations' not supported for LDAP data stores by PingFederate version "+string(r.providerConfig.ProductVersion), "PF 11.3 or later required")
			}
		}
	}

	if !pfVersionAtLeast121 {
		if internaltypes.IsDefined(plan.LdapDataStore) {
			useStartTls := plan.LdapDataStore.Attributes()["use_start_tls"]
			if internaltypes.IsDefined(useStartTls) {
				resp.Diagnostics.AddError("Attribute 'use_start_tls' not supported for LDAP data stores by PingFederate version "+string(r.providerConfig.ProductVersion), "PF 12.1 or later required")
			}
		}
		if internaltypes.IsDefined(plan.PingOneLdapGatewayDataStore) {
			useStartTls := plan.PingOneLdapGatewayDataStore.Attributes()["use_start_tls"]
			if internaltypes.IsDefined(useStartTls) {
				resp.Diagnostics.AddError("Attribute 'use_start_tls' not supported for LDAP gateway data stores by PingFederate version "+string(r.providerConfig.ProductVersion), "PF 12.1 or later required")
			}
		}
	}

	// Check for parent_ref, which had support removed in version 12.0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingFederate versions", err.Error())
		return
	}
	pfVersionAtLeast120 := compare >= 0
	if pfVersionAtLeast120 && internaltypes.IsDefined(plan.CustomDataStore) && internaltypes.IsDefined(plan.CustomDataStore.Attributes()["parent_ref"]) {
		resp.Diagnostics.AddError("Attribute 'parent_ref' not supported for custom data stores by PingFederate version "+string(r.providerConfig.ProductVersion), "PF 11.3 or earlier required")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Build name attribute for JDBC data stores
	if internaltypes.IsDefined(plan.JdbcDataStore) {
		jdbcDataStore := plan.JdbcDataStore.Attributes()
		// If connection_url is not set, build it based on connection_url_tags
		if jdbcDataStore["connection_url"].IsUnknown() {
			// Find the connection_url_tags with default_source set to true
			for _, tag := range jdbcDataStore["connection_url_tags"].(types.Set).Elements() {
				tagAttrs := tag.(types.Object).Attributes()
				if tagAttrs["default_source"].(types.Bool).ValueBool() {
					jdbcDataStore["connection_url"] = types.StringValue(tagAttrs["connection_url"].(types.String).ValueString())
					break
				}
			}
		}

		// If connection_url_tags is not set, build it based on connection_url
		if jdbcDataStore["connection_url_tags"].IsUnknown() {
			urlTag, respDiags := types.ObjectValue(jdbcTagConfigAttrType.AttrTypes, map[string]attr.Value{
				"connection_url": types.StringValue(jdbcDataStore["connection_url"].(types.String).ValueString()),
				"tags":           types.StringNull(),
				"default_source": types.BoolValue(true),
			})
			resp.Diagnostics.Append(respDiags...)
			jdbcDataStore["connection_url_tags"], respDiags = types.SetValue(jdbcTagConfigAttrType, []attr.Value{urlTag})
			resp.Diagnostics.Append(respDiags...)
		}

		// If name is not set, build it based on connection_url and user_name
		if jdbcDataStore["name"].IsUnknown() {
			var nameStr strings.Builder
			nameStr.WriteString(jdbcDataStore["connection_url"].(types.String).ValueString())
			nameStr.WriteString(" (")
			if internaltypes.IsDefined(jdbcDataStore["user_name"]) {
				nameStr.WriteString(jdbcDataStore["user_name"].(types.String).ValueString())
			} else {
				nameStr.WriteString("null")
			}
			nameStr.WriteString(")")
			jdbcDataStore["name"] = types.StringValue(nameStr.String())
		}

		plan.JdbcDataStore, respDiags = types.ObjectValue(jdbcDataStoreAttrType, jdbcDataStore)
		resp.Diagnostics.Append(respDiags...)
	}

	if internaltypes.IsDefined(plan.LdapDataStore) {
		ldapDataStore := plan.LdapDataStore.Attributes()
		// If hostnames is not set, build it based on hostnames_tags
		if ldapDataStore["hostnames"].IsUnknown() {
			// Find the hostnames_tags with default_source set to true
			var defaultHostnames []attr.Value
			for _, tag := range ldapDataStore["hostnames_tags"].(types.Set).Elements() {
				tagAttrs := tag.(types.Object).Attributes()
				if tagAttrs["default_source"].(types.Bool).ValueBool() {
					defaultHostnames = tagAttrs["hostnames"].(types.List).Elements()
					break
				}
			}
			if len(defaultHostnames) == 0 {
				resp.Diagnostics.AddError("No default hostnames found in hostnames_tags", "")
			} else {
				ldapDataStore["hostnames"], respDiags = types.ListValue(types.StringType, defaultHostnames)
				resp.Diagnostics.Append(respDiags...)
			}
		}

		// If hostnames_tags is not set, build it based on hostnames
		if ldapDataStore["hostnames_tags"].IsUnknown() {
			tagHostnames, respDiags := types.ListValue(types.StringType, ldapDataStore["hostnames"].(types.List).Elements())
			resp.Diagnostics.Append(respDiags...)
			tagAttr, respDiags := types.ObjectValue(ldapTagConfigAttrType.AttrTypes, map[string]attr.Value{
				"hostnames":      tagHostnames,
				"tags":           types.StringNull(),
				"default_source": types.BoolValue(true),
			})
			resp.Diagnostics.Append(respDiags...)
			ldapDataStore["hostnames_tags"], respDiags = types.SetValue(ldapTagConfigAttrType, []attr.Value{tagAttr})
			resp.Diagnostics.Append(respDiags...)
		}

		// If name is not set, build it based on hostnames and user_dn
		if ldapDataStore["name"].IsUnknown() {
			var nameStr strings.Builder
			for _, hostname := range ldapDataStore["hostnames"].(types.List).Elements() {
				nameStr.WriteString(hostname.(types.String).ValueString())
				nameStr.WriteString(" ")
			}
			nameStr.WriteString("(")
			if internaltypes.IsDefined(ldapDataStore["user_dn"]) {
				nameStr.WriteString(ldapDataStore["user_dn"].(types.String).ValueString())
			} else {
				nameStr.WriteString("null")
			}
			nameStr.WriteString(")")
			ldapDataStore["name"] = types.StringValue(nameStr.String())
		}

		// If PF version is at least 11.3 then set a default for retry_failed_operations.
		// If not ensure it is set to null rather than unknown.
		if ldapDataStore["retry_failed_operations"].IsUnknown() {
			if pfVersionAtLeast113 {
				ldapDataStore["retry_failed_operations"] = types.BoolValue(false)
			} else {
				ldapDataStore["retry_failed_operations"] = types.BoolNull()
			}
		}

		// use_start_tls attribute added in PF version 12.1
		if ldapDataStore["use_start_tls"].IsUnknown() {
			if pfVersionAtLeast121 {
				ldapDataStore["use_start_tls"] = types.BoolValue(false)
			} else {
				ldapDataStore["use_start_tls"] = types.BoolNull()
			}
		}

		plan.LdapDataStore, respDiags = types.ObjectValue(ldapDataStoreAttrType, ldapDataStore)
		resp.Diagnostics.Append(respDiags...)

		// Ensure one and only one of the authentication attributes is set
		if internaltypes.IsDefined(ldapDataStore["client_tls_certificate_ref"]) {
			// If client_tls_certificate_ref is defined, ensure use_start_tls or use_ssl is set to true
			if !ldapDataStore["use_start_tls"].(types.Bool).ValueBool() && !ldapDataStore["use_ssl"].(types.Bool).ValueBool() {
				resp.Diagnostics.AddError("Attribute 'client_tls_certificate_ref' requires either 'use_start_tls' or 'use_ssl' to be set to true", "")
			}
			// user_dn and bind_anonymously should not be set
			if internaltypes.IsDefined(ldapDataStore["user_dn"]) {
				resp.Diagnostics.AddError("Attribute 'client_tls_certificate_ref' requires 'user_dn' to be null", "")
			}
			if ldapDataStore["bind_anonymously"].(types.Bool).ValueBool() {
				resp.Diagnostics.AddError("Attribute 'client_tls_certificate_ref' requires 'bind_anonymously' to be false", "")
			}
		}
		if internaltypes.IsDefined(ldapDataStore["user_dn"]) {
			if ldapDataStore["bind_anonymously"].(types.Bool).ValueBool() {
				resp.Diagnostics.AddError("Attribute 'user_dn' requires 'bind_anonymously' to be false", "")
			}
		}
		// If password is set, then user_dn must be set
		if (internaltypes.IsDefined(ldapDataStore["password"]) && !internaltypes.IsDefined(ldapDataStore["user_dn"])) ||
			(!internaltypes.IsDefined(ldapDataStore["password"]) && internaltypes.IsDefined(ldapDataStore["user_dn"])) {
			resp.Diagnostics.AddError("'password' and 'user_dn' must be set together", "")
		}
	}

	if internaltypes.IsDefined(plan.PingOneLdapGatewayDataStore) {
		pingOneLdapGatewayDataStore := plan.PingOneLdapGatewayDataStore.Attributes()
		if pingOneLdapGatewayDataStore["name"].IsUnknown() {
			pingOneConnectionRefId := pingOneLdapGatewayDataStore["ping_one_connection_ref"].(types.Object).Attributes()["id"].(types.String).ValueString()
			pingOneEnvironmentId := pingOneLdapGatewayDataStore["ping_one_environment_id"].(types.String).ValueString()
			pingOneLdapGatewayId := pingOneLdapGatewayDataStore["ping_one_ldap_gateway_id"].(types.String).ValueString()
			pingOneLdapGatewayDataStore["name"] = types.StringValue(pingOneConnectionRefId + ":" + pingOneEnvironmentId + ":" + pingOneLdapGatewayId)
		}

		// use_start_tls attribute added in PF version 12.1
		if pingOneLdapGatewayDataStore["use_start_tls"].IsUnknown() {
			if pfVersionAtLeast121 {
				pingOneLdapGatewayDataStore["use_start_tls"] = types.BoolValue(false)
			} else {
				pingOneLdapGatewayDataStore["use_start_tls"] = types.BoolNull()
			}
		}

		plan.PingOneLdapGatewayDataStore, respDiags = types.ObjectValue(pingOneLdapGatewayDataStoreAttrType, pingOneLdapGatewayDataStore)
		resp.Diagnostics.Append(respDiags...)
	}

	// Mark the _all attributes as unknown when the corresponding non-_all attribute is changed for plugin configuration in custom_data_store
	if internaltypes.IsDefined(plan.CustomDataStore) && state != nil && internaltypes.IsDefined(state.CustomDataStore) {
		customDataStore := plan.CustomDataStore.Attributes()
		customDataStoreState := state.CustomDataStore.Attributes()
		if internaltypes.IsDefined(customDataStore["configuration"]) && internaltypes.IsDefined(customDataStoreState["configuration"]) {
			customDataStore["configuration"], respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(
				customDataStore["configuration"].(types.Object), customDataStoreState["configuration"].(types.Object))
			resp.Diagnostics.Append(respDiags...)
		}

		plan.CustomDataStore, respDiags = types.ObjectValue(plan.CustomDataStore.AttributeTypes(ctx), customDataStore)
		resp.Diagnostics.Append(respDiags...)
	}

	resp.Plan.Set(ctx, plan)
}

func createDataStore(dataStore client.DataStoreAggregation, dsr *dataStoreResource, con context.Context, resp *resource.CreateResponse) (*client.DataStoreAggregation, *http.Response, error) {
	apiCreateDataStore := dsr.apiClient.DataStoresAPI.CreateDataStore(config.AuthContext(con, dsr.providerConfig))
	apiCreateDataStore = apiCreateDataStore.Body(dataStore)
	return dsr.apiClient.DataStoresAPI.CreateDataStoreExecute(apiCreateDataStore)
}

func (r *dataStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataStoreModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if internaltypes.IsDefined(plan.CustomDataStore) {
		createCustomDataStore(plan, ctx, req, resp, r)
	}

	if internaltypes.IsDefined(plan.JdbcDataStore) {
		createJdbcDataStore(plan, ctx, req, resp, r)
	}

	if internaltypes.IsDefined(plan.LdapDataStore) {
		createLdapDataStore(plan, ctx, req, resp, r)
	}

	if internaltypes.IsDefined(plan.PingOneLdapGatewayDataStore) {
		createPingOneLdapGatewayDataStore(plan, ctx, req, resp, r)
	}

}

func (r *dataStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var state dataStoreModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	dataStoreGetReq, httpResp, err := r.apiClient.DataStoresAPI.GetDataStore(config.AuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Data Store", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while getting the data store", err, httpResp, &customId)
		}
		return
	}

	if dataStoreGetReq.CustomDataStore != nil {
		diags = readCustomDataStoreResponse(ctx, dataStoreGetReq, &state, &state.CustomDataStore, true, isImportRead)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.JdbcDataStore != nil {
		diags = readJdbcDataStoreResponse(ctx, dataStoreGetReq, &state, &state, true)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.LdapDataStore != nil {
		diags = readLdapDataStoreResponse(ctx, dataStoreGetReq, &state, &state.LdapDataStore, true)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.PingOneLdapGatewayDataStore != nil {
		diags = readPingOneLdapGatewayDataStoreResponse(ctx, dataStoreGetReq, &state, &state.PingOneLdapGatewayDataStore, true)
		resp.Diagnostics.Append(diags...)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func updateDataStore(dataStore client.DataStoreAggregation, dsr *dataStoreResource, con context.Context, resp *resource.UpdateResponse, id string) (*client.DataStoreAggregation, *http.Response, error) {
	updateDataStore := dsr.apiClient.DataStoresAPI.UpdateDataStore(config.AuthContext(con, dsr.providerConfig), id)
	updateDataStore = updateDataStore.Body(dataStore)
	return dsr.apiClient.DataStoresAPI.UpdateDataStoreExecute(updateDataStore)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dataStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan dataStoreModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if internaltypes.IsDefined(plan.CustomDataStore) {
		updateCustomDataStore(plan, ctx, req, resp, r)
	}

	if internaltypes.IsDefined(plan.JdbcDataStore) {
		updateJdbcDataStore(plan, ctx, req, resp, r)
	}

	if internaltypes.IsDefined(plan.LdapDataStore) {
		updateLdapDataStore(plan, ctx, req, resp, r)
	}

	if internaltypes.IsDefined(plan.PingOneLdapGatewayDataStore) {
		updatePingOneLdapGatewayDataStore(plan, ctx, req, resp, r)
	}

}

func (r *dataStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dataStoreModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.DataStoresAPI.DeleteDataStore(config.AuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting a data store", err, httpResp, &customId)
	}
}

func (r *dataStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
