package serversettingsgeneralsettings

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
	_ datasource.DataSource              = &serverSettingsGeneralSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsGeneralSettingsDataSource{}
)

// ServerSettingsGeneralSettingsDataSource is a helper function to simplify the provider implementation.
func NewServerSettingsGeneralSettingsDataSource() datasource.DataSource {
	return &serverSettingsGeneralSettingsDataSource{}
}

// serverSettingsGeneralSettingsDataSource is the datasource implementation.
type serverSettingsGeneralSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsGeneralSettingsDataSourceModel struct {
	Id                                      types.String `tfsdk:"id"`
	DisableAutomaticConnectionValidation    types.Bool   `tfsdk:"disable_automatic_connection_validation"`
	IdpConnectionTransactionLoggingOverride types.String `tfsdk:"idp_connection_transaction_logging_override"`
	SpConnectionTransactionLoggingOverride  types.String `tfsdk:"sp_connection_transaction_logging_override"`
	DatastoreValidationIntervalSecs         types.Int64  `tfsdk:"datastore_validation_interval_secs"`
	RequestHeaderForCorrelationId           types.String `tfsdk:"request_header_for_correlation_id"`
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsGeneralSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *serverSettingsGeneralSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_general_settings"
}

func (r *serverSettingsGeneralSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsGeneralSettingsDataSource(ctx context.Context, r *client.GeneralSettings, state *serverSettingsGeneralSettingsDataSourceModel) {
	state.Id = types.StringValue("server_settings_general_settings_id")
	state.DisableAutomaticConnectionValidation = types.BoolPointerValue(r.DisableAutomaticConnectionValidation)
	state.IdpConnectionTransactionLoggingOverride = types.StringPointerValue(r.IdpConnectionTransactionLoggingOverride)
	state.SpConnectionTransactionLoggingOverride = types.StringPointerValue(r.SpConnectionTransactionLoggingOverride)
	state.DatastoreValidationIntervalSecs = types.Int64PointerValue(r.DatastoreValidationIntervalSecs)
	state.RequestHeaderForCorrelationId = types.StringPointerValue(r.RequestHeaderForCorrelationId)
}

func (r *serverSettingsGeneralSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsGeneralSettingsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadServerSettingsGeneralSettings, httpResp, err := r.apiClient.ServerSettingsAPI.GetGeneralSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings General Settings", err, httpResp)
		return
	}

	// Read the response into the state

	readServerSettingsGeneralSettingsDataSource(ctx, apiReadServerSettingsGeneralSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
