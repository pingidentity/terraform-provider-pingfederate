package serversettingsgeneralsettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsGeneralDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsGeneralDataSource{}
)

// ServerSettingsGeneralDataSource is a helper function to simplify the provider implementation.
func ServerSettingsGeneralDataSource() datasource.DataSource {
	return &serverSettingsGeneralDataSource{}
}

// serverSettingsGeneralDataSource is the datasource implementation.
type serverSettingsGeneralDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsGeneralDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the general server settings.",
		Attributes: map[string]schema.Attribute{
			"datastore_validation_interval_secs": schema.Int64Attribute{
				Description: "How long (in seconds) the result of testing a datastore connection is cached.",
				Computed:    true,
				Optional:    false,
			},
			"disable_automatic_connection_validation": schema.BoolAttribute{
				Description: "Boolean that disables automatic connection validation when set to true.",
				Computed:    true,
				Optional:    false,
			},
			"idp_connection_transaction_logging_override": schema.StringAttribute{
				Description: "Describes the level of transaction logging for all identity provider connections. The default is DONT_OVERRIDE, in which case the logging level will be determined by each individual IdP connection [ DONT_OVERRIDE, NONE, FULL, STANDARD, ENHANCED ]",
				Computed:    true,
				Optional:    false,
			},
			"request_header_for_correlation_id": schema.StringAttribute{
				Description: "HTTP request header for retrieving correlation ID.",
				Computed:    true,
				Optional:    false,
			},
			"sp_connection_transaction_logging_override": schema.StringAttribute{
				Description: "Determines the level of transaction logging for all service provider connections. The default is DONT_OVERRIDE, in which case the logging level will be determined by each individual SP connection [ DONT_OVERRIDE, NONE, FULL, STANDARD, ENHANCED ]",
				Computed:    true,
				Optional:    false,
			},
		},
	}

	id.ToDataSourceSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *serverSettingsGeneralDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_general"
}

func (r *serverSettingsGeneralDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *serverSettingsGeneralDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsGeneralModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsGeneral, httpResp, err := r.apiClient.ServerSettingsAPI.GetGeneralSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings General Settings", err, httpResp)
		return
	}

	// Read the response into the state

	readServerSettingsGeneralResponse(ctx, apiReadServerSettingsGeneral, &state, pointers.String("server_settings_general_settings_id"))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
