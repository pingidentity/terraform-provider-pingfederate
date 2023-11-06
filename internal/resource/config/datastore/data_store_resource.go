package datastore

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type dataStoreResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	CustomId                    types.String `tfsdk:"custom_id"`
	MaskAttributeValues         types.Bool   `tfsdk:"mask_attribute_values"`
	CustomDataStore             types.Object `tfsdk:"custom_data_store"`
	JdbcDataStore               types.Object `tfsdk:"jdbc_data_store"`
	LdapDataStore               types.Object `tfsdk:"ldap_data_store"`
	PingOneLdapGatewayDataStore types.Object `tfsdk:"ping_one_ldap_gateway_data_store"`
}

// GetSchema defines the schema for the resource.
func (r *dataStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Data Store",
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
	id.ToSchemaCustomId(&schema, false,
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

func createDataStore(dataStore configurationapi.DataStoreAggregation, dsr *dataStoreResource, con context.Context, resp *resource.CreateResponse) (*client.DataStoreAggregation, *http.Response, error) {
	apiCreateDataStore := dsr.apiClient.DataStoresAPI.CreateDataStore(config.ProviderBasicAuthContext(con, dsr.providerConfig))
	apiCreateDataStore = apiCreateDataStore.Body(dataStore)
	return dsr.apiClient.DataStoresAPI.CreateDataStoreExecute(apiCreateDataStore)
}

func (r *dataStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dataStoreResourceModel

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
	var state dataStoreResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	dataStoreGetReq, httpResp, err := r.apiClient.DataStoresAPI.GetDataStore(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Data Store", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  Data Store", err, httpResp)
		}
	}

	if dataStoreGetReq.CustomDataStore != nil {
		diags = readCustomDataStoreResponse(ctx, dataStoreGetReq, &state, &state.CustomDataStore)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.JdbcDataStore != nil {
		diags = readJdbcDataStoreResponse(ctx, dataStoreGetReq, &state, &state)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.LdapDataStore != nil {
		diags = readLdapDataStoreResponse(ctx, dataStoreGetReq, &state, &state.LdapDataStore)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.PingOneLdapGatewayDataStore != nil {
		diags = readPingOneLdapGatewayDataStoreResponse(ctx, dataStoreGetReq, &state, &state.PingOneLdapGatewayDataStore)
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

	var plan dataStoreResourceModel

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
	var state dataStoreResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.DataStoresAPI.DeleteDataStore(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Data Store", err, httpResp)
		return
	}
}

func (r *dataStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
