package authenticationapisettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	resourcelinkdatasource "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &authenticationApiSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationApiSettingsDataSource{}
)

// Create a Authentication Api Settings data source
func AuthenticationApiSettingsDataSource() datasource.DataSource {
	return &authenticationApiSettingsDataSource{}
}

// authenticationApiSettingsDataSource is the datasource implementation.
type authenticationApiSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *authenticationApiSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_api_settings"
}

// Configure adds the provider configured client to the data source.
func (r *authenticationApiSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// GetSchema defines the schema for the datasource.
func (r *authenticationApiSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes the authentication API application settings.",
		Attributes: map[string]schema.Attribute{
			"api_enabled": schema.BoolAttribute{
				Description: "Enable Authentication API",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_api_descriptions": schema.BoolAttribute{
				Description: "Enable API descriptions",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_application_ref": schema.SingleNestedAttribute{
				Description: "Enable API descriptions",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelinkdatasource.ToDataSourceSchema(),
			},
			"restrict_access_to_redirectless_mode": schema.BoolAttribute{
				Description: "Enable restrict access to redirectless mode",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_request_context": schema.BoolAttribute{
				Description: "Includes request context in API responses",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	id.ToDataSourceSchema(&schemaDef)
	resp.Schema = schemaDef
}

// Read resource information
func (r *authenticationApiSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationApiSettingsModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAuthenticationApiSettings, httpResp, err := r.apiClient.AuthenticationApiAPI.GetAuthenticationApiSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the authentication API settings", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readAuthenticationApiSettingsResponse(ctx, apiReadAuthenticationApiSettings, &state, pointers.String("authentication_api_settings_id"))
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
