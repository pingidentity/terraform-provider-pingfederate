package protocolmetadatalifetimesettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &protocolMetadataLifetimeSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &protocolMetadataLifetimeSettingsDataSource{}
)

// ProtocolMetadataLifetimeSettingsDataSource is a helper function to simplify the provider implementation.
func NewProtocolMetadataLifetimeSettingsDataSource() datasource.DataSource {
	return &protocolMetadataLifetimeSettingsDataSource{}
}

// protocolMetadataLifetimeSettingsDataSource is the datasource implementation.
type protocolMetadataLifetimeSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type protocolMetadataLifetimeSettingsDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	CacheDuration types.Int64  `tfsdk:"cache_duration"`
	ReloadDelay   types.Int64  `tfsdk:"reload_delay"`
}

// GetSchema defines the schema for the datasource.
func (r *protocolMetadataLifetimeSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes a ProtocolMetadataLifetimeSettings.",
		Attributes: map[string]schema.Attribute{
			"cache_duration": schema.Int64Attribute{
				Description: "The validity of your metadata in minutes. The default value is 1440 (1 day).",
				Computed:    true,
				Optional:    false,
			},
			"reload_delay": schema.Int64Attribute{
				Description: "The frequency of automatic reloading of SAML metadata in minutes. The default value is 1440 (1 day).",
				Computed:    true,
				Optional:    false,
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *protocolMetadataLifetimeSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protocol_metadata_lifetime_settings"
}

func (r *protocolMetadataLifetimeSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readProtocolMetadataLifetimeSettingsResponseDataSource(ctx context.Context, r *client.MetadataLifetimeSettings, state *protocolMetadataLifetimeSettingsDataSourceModel) {
	state.Id = types.StringValue("protocolMetadataLifetimeSettingsId")
	state.CacheDuration = types.Int64Value(r.GetCacheDuration())
	state.ReloadDelay = types.Int64Value(r.GetReloadDelay())
}

func (r *protocolMetadataLifetimeSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state protocolMetadataLifetimeSettingsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadProtocolMetadataLifetimeSettings, httpResp, err := r.apiClient.ProtocolMetadataAPI.GetLifetimeSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Protocol Metadata Lifetime Settings", err, httpResp)
		return
	}

	// Read the response into the state
	readProtocolMetadataLifetimeSettingsResponseDataSource(ctx, apiReadProtocolMetadataLifetimeSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
