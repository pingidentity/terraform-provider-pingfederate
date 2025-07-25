// Copyright © 2025 Ping Identity Corporation

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
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
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
				Description: "The persistent, unique ID for the data store. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.",
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
	// Compare to version 12.1 of PF
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingFederate1210)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast121 := compare >= 0

	if !pfVersionAtLeast121 {
		if internaltypes.IsDefined(plan.LdapDataStore) {
			useStartTls := plan.LdapDataStore.Attributes()["use_start_tls"]
			if internaltypes.IsDefined(useStartTls) {
				resp.Diagnostics.AddAttributeError(
					path.Root("ldap_data_store").AtMapKey("use_start_tls"),
					providererror.InvalidProductVersionAttribute,
					"Attribute 'use_start_tls' not supported for LDAP data stores by PingFederate version "+string(r.providerConfig.ProductVersion)+". PF 12.1 or later required")
			}
		}
		if internaltypes.IsDefined(plan.PingOneLdapGatewayDataStore) {
			useStartTls := plan.PingOneLdapGatewayDataStore.Attributes()["use_start_tls"]
			if internaltypes.IsDefined(useStartTls) {
				resp.Diagnostics.AddAttributeError(
					path.Root("ping_one_ldap_gateway_data_store").AtMapKey("use_start_tls"),
					providererror.InvalidProductVersionAttribute,
					"Attribute 'use_start_tls' not supported for LDAP gateway data stores by PingFederate version "+string(r.providerConfig.ProductVersion)+". PF 12.1 or later required")
			}
		}
	}

	// Check for parent_ref, which had support removed in version 12.0
	compare, err = version.Compare(r.providerConfig.ProductVersion, version.PingFederate1200)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		return
	}
	pfVersionAtLeast120 := compare >= 0
	if pfVersionAtLeast120 && internaltypes.IsDefined(plan.CustomDataStore) && internaltypes.IsDefined(plan.CustomDataStore.Attributes()["parent_ref"]) {
		resp.Diagnostics.AddAttributeError(
			path.Root("custom_data_store").AtMapKey("parent_ref"),
			providererror.InvalidProductVersionAttribute,
			"Attribute 'parent_ref' not supported for custom data stores by PingFederate version "+string(r.providerConfig.ProductVersion)+". PF 11.3 or earlier required")
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
				if !tag.IsUnknown() {
					tagAttrs := tag.(types.Object).Attributes()
					if tagAttrs["default_source"].(types.Bool).ValueBool() && !tagAttrs["connection_url"].IsUnknown() {
						jdbcDataStore["connection_url"] = types.StringValue(tagAttrs["connection_url"].(types.String).ValueString())
						break
					}
				}
			}
		}

		// If connection_url_tags is not set, build it based on connection_url
		if jdbcDataStore["connection_url_tags"].IsUnknown() && !jdbcDataStore["connection_url"].IsUnknown() {
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
		if jdbcDataStore["name"].IsUnknown() && !jdbcDataStore["connection_url"].IsUnknown() && !jdbcDataStore["user_name"].IsUnknown() {
			var nameStr strings.Builder
			nameStr.WriteString(jdbcDataStore["connection_url"].(types.String).ValueString())
			nameStr.WriteString(" (")
			if !jdbcDataStore["user_name"].IsNull() {
				nameStr.WriteString(jdbcDataStore["user_name"].(types.String).ValueString())
			} else {
				nameStr.WriteString("null")
			}
			nameStr.WriteString(")")
			jdbcDataStore["name"] = types.StringValue(nameStr.String())
		}

		// If password value has changed, mark encrypted_password value as unknown
		if state != nil {
			stateJdbcDataStore := state.JdbcDataStore.Attributes()
			if !jdbcDataStore["password"].Equal(stateJdbcDataStore["password"]) {
				jdbcDataStore["encrypted_password"] = types.StringUnknown()
			}
		}

		plan.JdbcDataStore, respDiags = types.ObjectValue(jdbcDataStoreAttrType, jdbcDataStore)
		resp.Diagnostics.Append(respDiags...)
	}

	if internaltypes.IsDefined(plan.LdapDataStore) {
		ldapDataStore := plan.LdapDataStore.Attributes()
		// If hostnames is not set, build it based on hostnames_tags
		if ldapDataStore["hostnames"].IsUnknown() && !ldapDataStore["hostnames_tags"].IsUnknown() {
			// Find the hostnames_tags with default_source set to true
			var defaultHostnames []attr.Value
			for _, tag := range ldapDataStore["hostnames_tags"].(types.Set).Elements() {
				if tag.IsUnknown() {
					continue
				}
				tagAttrs := tag.(types.Object).Attributes()
				if tagAttrs["default_source"].(types.Bool).ValueBool() && !tagAttrs["hostnames"].IsUnknown() {
					defaultHostnames = tagAttrs["hostnames"].(types.List).Elements()
					ldapDataStore["hostnames"], respDiags = types.ListValue(types.StringType, defaultHostnames)
					resp.Diagnostics.Append(respDiags...)
					break
				}
			}
		}

		// If hostnames_tags is not set, build it based on hostnames
		if ldapDataStore["hostnames_tags"].IsUnknown() && !ldapDataStore["hostnames"].IsUnknown() {
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
		if ldapDataStore["name"].IsUnknown() && !ldapDataStore["hostnames"].IsUnknown() && !ldapDataStore["user_dn"].IsUnknown() {
			var nameStr strings.Builder
			anyUnknown := false
			for _, hostname := range ldapDataStore["hostnames"].(types.List).Elements() {
				if hostname.IsUnknown() {
					anyUnknown = true
					break
				}
				nameStr.WriteString(hostname.(types.String).ValueString())
				nameStr.WriteString(" ")
			}
			if !anyUnknown {
				nameStr.WriteString("(")
				if !ldapDataStore["user_dn"].IsNull() {
					nameStr.WriteString(ldapDataStore["user_dn"].(types.String).ValueString())
				} else {
					nameStr.WriteString("null")
				}
				nameStr.WriteString(")")
				ldapDataStore["name"] = types.StringValue(nameStr.String())
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

		// If password value has changed, mark encrypted_password value as unknown
		if state != nil {
			stateLdapDataStore := state.LdapDataStore.Attributes()
			if !ldapDataStore["password"].Equal(stateLdapDataStore["password"]) {
				ldapDataStore["encrypted_password"] = types.StringUnknown()
			}
		}

		plan.LdapDataStore, respDiags = types.ObjectValue(ldapDataStoreAttrType, ldapDataStore)
		resp.Diagnostics.Append(respDiags...)

		// Ensure one and only one of the authentication attributes is set
		if internaltypes.IsDefined(ldapDataStore["client_tls_certificate_ref"]) {
			// If client_tls_certificate_ref is defined, ensure use_start_tls or use_ssl is set to true
			if !ldapDataStore["use_start_tls"].IsUnknown() && !ldapDataStore["use_start_tls"].(types.Bool).ValueBool() &&
				!ldapDataStore["use_ssl"].IsUnknown() && !ldapDataStore["use_ssl"].(types.Bool).ValueBool() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ldap_data_store").AtMapKey("client_tls_certificate_ref"),
					providererror.InvalidAttributeConfiguration,
					"Attribute 'client_tls_certificate_ref' requires either 'use_start_tls' or 'use_ssl' to be set to true")
			}
			// user_dn and bind_anonymously should not be set
			if internaltypes.IsDefined(ldapDataStore["user_dn"]) {
				resp.Diagnostics.AddAttributeError(
					path.Root("ldap_data_store").AtMapKey("user_dn"),
					providererror.InvalidAttributeConfiguration,
					"Attribute 'client_tls_certificate_ref' requires 'user_dn' to be null")
			}
			if ldapDataStore["bind_anonymously"].(types.Bool).ValueBool() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ldap_data_store").AtMapKey("bind_anonymously"),
					providererror.InvalidAttributeConfiguration,
					"Attribute 'client_tls_certificate_ref' requires 'bind_anonymously' to be false")
			}
		}
		if internaltypes.IsDefined(ldapDataStore["user_dn"]) {
			if ldapDataStore["bind_anonymously"].(types.Bool).ValueBool() {
				resp.Diagnostics.AddAttributeError(
					path.Root("ldap_data_store").AtMapKey("bind_anonymously"),
					providererror.InvalidAttributeConfiguration,
					"Attribute 'user_dn' requires 'bind_anonymously' to be false")
			}
		}
		// If password or encrypted_password is set, then user_dn must be set
		if ((internaltypes.IsDefined(ldapDataStore["password"]) || internaltypes.IsDefined(ldapDataStore["encrypted_password"])) && ldapDataStore["user_dn"].IsNull()) ||
			(ldapDataStore["password"].IsNull() && ldapDataStore["encrypted_password"].IsNull() && internaltypes.IsDefined(ldapDataStore["user_dn"])) {
			resp.Diagnostics.AddAttributeError(
				path.Root("ldap_data_store"),
				providererror.InvalidAttributeConfiguration,
				"'password' (or 'encrypted_password') and 'user_dn' must be set together")
		}
	}

	if internaltypes.IsDefined(plan.PingOneLdapGatewayDataStore) {
		pingOneLdapGatewayDataStore := plan.PingOneLdapGatewayDataStore.Attributes()
		if pingOneLdapGatewayDataStore["name"].IsUnknown() &&
			!pingOneLdapGatewayDataStore["ping_one_connection_ref"].IsUnknown() &&
			!pingOneLdapGatewayDataStore["ping_one_connection_ref"].(types.Object).Attributes()["id"].IsUnknown() &&
			!pingOneLdapGatewayDataStore["ping_one_environment_id"].IsUnknown() &&
			!pingOneLdapGatewayDataStore["ping_one_ldap_gateway_id"].IsUnknown() {
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

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
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
	dataStoreGetReq, httpResp, err := r.apiClient.DataStoresAPI.GetDataStore(config.AuthContext(ctx, r.providerConfig), state.DataStoreId.ValueString()).Execute()
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
	resource.ImportStatePassthroughID(ctx, path.Root("data_store_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
