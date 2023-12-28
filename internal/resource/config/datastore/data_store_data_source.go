package datastore

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ datasource.DataSource              = &dataStoreDataSource{}
	_ datasource.DataSourceWithConfigure = &dataStoreDataSource{}
)

func DataStoreDataSource() datasource.DataSource {
	return &dataStoreDataSource{}
}

type dataStoreDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *dataStoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes a data store.",
		Attributes: map[string]schema.Attribute{
			"mask_attribute_values": schema.BoolAttribute{
				Description: "Whether attribute values should be masked in the log.",
				Computed:    true,
				Optional:    false,
			},
			"last_modified": schema.StringAttribute{
				Description: "The time at which the datastore instance was last changed. Supported in PF version 12.0 or later.",
				Optional:    false,
				Computed:    true,
			},
			"custom_data_store":                toDataSourceSchemaCustomDataStore(),
			"jdbc_data_store":                  toDataSourceSchemaJdbcDataStore(),
			"ldap_data_store":                  toDataSourceSchemaLdapDataStore(),
			"ping_one_ldap_gateway_data_store": toDataSourceSchemaPingOneLdapGatewayDataStore(),
		},
	}
	id.ToDataSourceSchema(&schema)
	id.ToDataSourceSchemaCustomId(&schema,
		"data_store_id",
		true,
		"Unique ID for the data store.")

	resp.Schema = schema
}

func (r *dataStoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (r *dataStoreDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *dataStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state dataStoreModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	dataStoreGetReq, httpResp, err := r.apiClient.DataStoresAPI.GetDataStore(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.DataStoreId.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the data store", err, httpResp)

	}

	if dataStoreGetReq.CustomDataStore != nil {
		diags = readCustomDataStoreResponse(ctx, dataStoreGetReq, &state, &state.CustomDataStore, false)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.JdbcDataStore != nil {
		diags = readJdbcDataStoreResponse(ctx, dataStoreGetReq, &state, &state, false)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.LdapDataStore != nil {
		diags = readLdapDataStoreResponse(ctx, dataStoreGetReq, &state, &state.LdapDataStore, false)
		resp.Diagnostics.Append(diags...)
	}

	if dataStoreGetReq.PingOneLdapGatewayDataStore != nil {
		diags = readPingOneLdapGatewayDataStoreResponse(ctx, dataStoreGetReq, &state, &state.PingOneLdapGatewayDataStore, false)
		resp.Diagnostics.Append(diags...)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
