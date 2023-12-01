package datastore

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &dataStoreResource{}
	_ resource.ResourceWithConfigure   = &dataStoreResource{}
	_ resource.ResourceWithImportState = &dataStoreResource{}
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
		Description: "Manages a data store resource",
		Attributes: map[string]schema.Attribute{
			"mask_attribute_values": schema.BoolAttribute{
				Description: "Whether attribute values should be masked in the log.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"custom_data_store":                toSchemaCustomDataStore(),
			"jdbc_data_store":                  toSchemaJdbcDataStore(),
			"ldap_data_store":                  toSchemaLdapDataStore(),
			"ping_one_ldap_gateway_data_store": toSchemaPingOneLdapGatewayDataStore(),
		},
	}
	id.ToSchema(&schema)
	id.ToSchemaCustomId(&schema,
		"data_store_id",
		false,
		"The persistent, unique ID for the data store. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.")

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
	var plan *dataStoreModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	var respDiags diag.Diagnostics

	if plan == nil {
		return
	}

	// Build name attribute for data stores that have it
	if internaltypes.IsDefined(plan.JdbcDataStore) {
		jdbcDataStore := plan.JdbcDataStore.Attributes()

		// Build connection_url and connection_url_tags attributes
		var hasTopConnectionUrl bool
		var topConnectionUrlVal string
		if internaltypes.IsDefined(jdbcDataStore["connection_url"]) {
			hasTopConnectionUrl = true
			topConnectionUrlVal = jdbcDataStore["connection_url"].(types.String).ValueString()
		} else {
			hasTopConnectionUrl = false
			topConnectionUrlVal = ""
		}

		var hasConnectionUrlTags bool
		var connectionUrlTagsConnectionUrlVal basetypes.StringValue
		var connectionUrlTagsTags basetypes.StringValue
		var connectionUrlTagsDefaultSource basetypes.BoolValue
		if internaltypes.IsDefined(jdbcDataStore["connection_url_tags"]) {
			connectionUrlTags := jdbcDataStore["connection_url_tags"].(types.Set)
			if len(connectionUrlTags.Elements()) > 0 {
				for _, elem := range connectionUrlTags.Elements() {
					objAttrs := elem.(types.Object).Attributes()
					connectionUrl := objAttrs["connection_url"].(types.String)
					defaultSource := objAttrs["default_source"].(types.Bool)
					var tags types.String
					if internaltypes.IsDefined(objAttrs["tags"]) {
						tags = objAttrs["tags"].(types.String)
					} else {
						tags = types.StringNull()
					}
					if internaltypes.IsNonEmptyString(connectionUrl) && defaultSource.ValueBool() {
						hasConnectionUrlTags = true
						connectionUrlTagsConnectionUrlVal = connectionUrl
						connectionUrlTagsTags = tags
						connectionUrlTagsDefaultSource = defaultSource
					}
				}
			}
		} else {
			hasConnectionUrlTags = false
			connectionUrlTagsConnectionUrlVal = types.StringNull()
			connectionUrlTagsTags = types.StringNull()
			connectionUrlTagsDefaultSource = types.BoolValue(true)
		}

		// If connection_url is not defined, use connection_url_tags connection_url value
		var connectionUrlVal basetypes.StringValue
		if !hasTopConnectionUrl {
			connectionUrlVal = connectionUrlTagsConnectionUrlVal
		} else {
			connectionUrlVal = types.StringValue(topConnectionUrlVal)
		}
		jdbcDataStore["connection_url"] = connectionUrlVal

		connectionUrlTagsAttrVal := map[string]attr.Value{
			"connection_url": connectionUrlVal,
			"tags":           connectionUrlTagsTags,
			"default_source": connectionUrlTagsDefaultSource,
		}
		connectionUrlTagsObj, respDiags := types.ObjectValue(jdbcTagConfigAttrType.AttrTypes, connectionUrlTagsAttrVal)
		resp.Diagnostics.Append(respDiags...)

		connectionUrlTagsSetAttrValue := []attr.Value{}
		connectionUrlTagsSetAttrValue = append(connectionUrlTagsSetAttrValue, connectionUrlTagsObj)
		connectionUrlTagsSetAttrValue = append(connectionUrlTagsSetAttrValue, jdbcDataStore["connection_url_tags"].(types.Set).Elements()...)
		connectionUrlTagsSet, respDiags := types.SetValue(jdbcTagConfigAttrType, connectionUrlTagsSetAttrValue)
		resp.Diagnostics.Append(respDiags...)
		jdbcDataStore["connection_url_tags"] = connectionUrlTagsSet

		//  Build name attribute if not defined
		var namePrefix basetypes.StringValue
		var prefix string
		if !internaltypes.IsDefined(jdbcDataStore["name"]) {
			if hasTopConnectionUrl {
				prefix = topConnectionUrlVal
			}
			if hasConnectionUrlTags {
				prefix = connectionUrlTagsConnectionUrlVal.ValueString()
			}
			userName := jdbcDataStore["user_name"].(types.String).ValueString()
			namePrefix = types.StringValue(prefix + " (" + userName + ")")
		} else {
			namePrefix = jdbcDataStore["name"].(types.String)
		}

		jdbcDataStore["name"] = namePrefix
		plan.JdbcDataStore, respDiags = types.ObjectValue(jdbcDataStoreAttrType, jdbcDataStore)
		resp.Diagnostics.Append(respDiags...)
	}

	if internaltypes.IsDefined(plan.LdapDataStore) {
		ldapDataStore := plan.LdapDataStore.Attributes()

		// Build hostnames and hostnames_tags attributes
		var hasHostnames bool
		var hostnamesVal []attr.Value
		if internaltypes.IsDefined(ldapDataStore["hostnames"]) {
			topHostNames := ldapDataStore["hostnames"].(types.Set)
			if len(topHostNames.Elements()) > 0 {
				hasHostnames = true
				hostnamesVal = topHostNames.Elements()
			}
		} else {
			hasHostnames = false
			hostnamesVal = nil
		}

		// Set hostname from the default source value in hostnames_tags, if defined
		var hasHostnamesTags bool
		var hostnamesTagsHostnamesVal []attr.Value
		var hostnamesTagsTags basetypes.StringValue
		var hostnamesTagsDefaultSource basetypes.BoolValue
		if internaltypes.IsDefined(ldapDataStore["hostnames_tags"]) {
			hostnamesTags := ldapDataStore["hostnames_tags"].(types.Set)
			if len(hostnamesTags.Elements()) > 0 {
				hostnamesTagsFirstElem := hostnamesTags.Elements()[0].(types.Object).Attributes()
				if len(hostnamesTagsFirstElem["hostnames"].(types.Set).Elements()) > 0 {
					var tags types.String
					if internaltypes.IsDefined(hostnamesTagsFirstElem["tags"]) {
						tags = hostnamesTagsFirstElem["tags"].(types.String)
					} else {
						tags = types.StringNull()
					}
					var defaultSource types.Bool
					if internaltypes.IsDefined(hostnamesTagsFirstElem["default_source"]) {
						defaultSource = hostnamesTagsFirstElem["default_source"].(types.Bool)
					} else {
						defaultSource = types.BoolValue(true)
					}
					hasHostnamesTags = true
					hostnamesTagsHostnamesVal = hostnamesTagsFirstElem["hostnames"].(types.Set).Elements()
					hostnamesTagsTags = tags
					hostnamesTagsDefaultSource = defaultSource
				}
			}
		} else {
			hasHostnamesTags = false
			hostnamesTagsHostnamesVal = nil
			hostnamesTagsTags = types.StringNull()
			hostnamesTagsDefaultSource = types.BoolValue(true)
		}

		// If hostnames is not defined, use hostnames_tags hostnames value
		var hostnamesSetVal []attr.Value
		if !hasHostnames {
			hostnamesSetVal = hostnamesTagsHostnamesVal
		} else {
			hostnamesSetVal = hostnamesVal
		}

		hostnamesBaseTypesSetValue, respDiags := types.SetValue(types.StringType, hostnamesSetVal)
		resp.Diagnostics.Append(respDiags...)
		ldapDataStore["hostnames"] = hostnamesBaseTypesSetValue

		hostnamesTagsAttrVal := map[string]attr.Value{
			"hostnames":      hostnamesBaseTypesSetValue,
			"tags":           hostnamesTagsTags,
			"default_source": hostnamesTagsDefaultSource,
		}
		hostnamesTagsObj, respDiags := types.ObjectValue(ldapTagConfigAttrType.AttrTypes, hostnamesTagsAttrVal)
		resp.Diagnostics.Append(respDiags...)

		hostnamesTagsSetAttrValue := []attr.Value{hostnamesTagsObj}
		hostnamesTagsSet, respDiags := types.SetValue(ldapTagConfigAttrType, hostnamesTagsSetAttrValue)
		resp.Diagnostics.Append(respDiags...)
		ldapDataStore["hostnames_tags"] = hostnamesTagsSet

		//  Build name attribute if not defined
		var namePrefix string
		if hasHostnames {
			namePrefix = hostnamesVal[0].(types.String).ValueString()
		}
		if hasHostnamesTags {
			namePrefix = hostnamesTagsHostnamesVal[0].(types.String).ValueString()
		}

		var nameValue basetypes.StringValue
		if ldapDataStore["name"].IsUnknown() {
			userDn := ldapDataStore["user_dn"].(types.String).ValueString()
			nameValue = types.StringValue(namePrefix + " (" + userDn + ")")
		} else {
			nameValue = ldapDataStore["name"].(types.String)
		}

		ldapDataStore["name"] = nameValue
		plan.LdapDataStore, respDiags = types.ObjectValue(ldapDataStoreAttrType, ldapDataStore)
		resp.Diagnostics.Append(respDiags...)
	}

	if internaltypes.IsDefined(plan.PingOneLdapGatewayDataStore) {
		pingOneLdapGatewayDataStore := plan.PingOneLdapGatewayDataStore.Attributes()
		if !internaltypes.IsDefined(pingOneLdapGatewayDataStore["name"]) {
			pingOneConnectionRefId := pingOneLdapGatewayDataStore["ping_one_connection_ref"].(types.Object).Attributes()["id"].(types.String).ValueString()
			pingOneEnvironmentId := pingOneLdapGatewayDataStore["ping_one_environment_id"].(types.String).ValueString()
			pingOneLdapGatewayId := pingOneLdapGatewayDataStore["ping_one_ldap_gateway_id"].(types.String).ValueString()
			pingOneLdapGatewayDataStore["name"] = types.StringValue(pingOneConnectionRefId + ":" + pingOneEnvironmentId + ":" + pingOneLdapGatewayId)
			plan.PingOneLdapGatewayDataStore, respDiags = types.ObjectValue(pingOneLdapGatewayDataStoreAttrType, pingOneLdapGatewayDataStore)
			resp.Diagnostics.Append(respDiags...)
		}
	}

	resp.Plan.Set(ctx, plan)
}

func createDataStore(dataStore configurationapi.DataStoreAggregation, dsr *dataStoreResource, con context.Context, resp *resource.CreateResponse) (*client.DataStoreAggregation, *http.Response, error) {
	apiCreateDataStore := dsr.apiClient.DataStoresAPI.CreateDataStore(config.ProviderBasicAuthContext(con, dsr.providerConfig))
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
	var state dataStoreModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	dataStoreGetReq, httpResp, err := r.apiClient.DataStoresAPI.GetDataStore(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the data store", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the data store", err, httpResp)
		}
	}

	if dataStoreGetReq.CustomDataStore != nil {
		diags = readCustomDataStoreResponse(ctx, dataStoreGetReq, &state, &state.CustomDataStore, true)
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

func updateDataStore(dataStore configurationapi.DataStoreAggregation, dsr *dataStoreResource, con context.Context, resp *resource.UpdateResponse, id string) (*client.DataStoreAggregation, *http.Response, error) {
	updateDataStore := dsr.apiClient.DataStoresAPI.UpdateDataStore(config.ProviderBasicAuthContext(con, dsr.providerConfig), id)
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
	httpResp, err := r.apiClient.DataStoresAPI.DeleteDataStore(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a data store", err, httpResp)
		return
	}
}

func (r *dataStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
