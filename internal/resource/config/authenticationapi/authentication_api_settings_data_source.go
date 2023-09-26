package authenticationapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &authenticationApiSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationApiSettingsDataSource{}
)

// Create a Authentication Api Settings data source
func NewAuthenticationApiSettingsDataSource() datasource.DataSource {
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

type authenticationApiSettingsDataSourceModel struct {
	Id                               types.String `tfsdk:"id"`
	ApiEnabled                       types.Bool   `tfsdk:"api_enabled"`
	EnableApiDescriptions            types.Bool   `tfsdk:"enable_api_descriptions"`
	RestrictAccessToRedirectlessMode types.Bool   `tfsdk:"restrict_access_to_redirectless_mode"`
	IncludeRequestContext            types.Bool   `tfsdk:"include_request_context"`
	DefaultApplicationRef            types.Object `tfsdk:"default_application_ref"`
}

// GetSchema defines the schema for the datasource.
func (r *authenticationApiSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a AuthenticationApiSettings.",
		Attributes: map[string]schema.Attribute{
			"api_enabled": schema.BoolAttribute{
				Description: "Enable Authentication API",
				Optional:    true,
				Computed:    true,
			},
			"enable_api_descriptions": schema.BoolAttribute{
				Description: "Enable API descriptions",
				Optional:    true,
				Computed:    true,
			},
			"default_application_ref": schema.SingleNestedAttribute{
				Description: "Enable API descriptions",
				Computed:    true,
				Optional:    true,
			},
			"restrict_access_to_redirectless_mode": schema.BoolAttribute{
				Description: "Enable restrict access to redirectless mode",
				Optional:    true,
				Computed:    true,
			},
			"include_request_context": schema.BoolAttribute{
				Description: "Includes request context in API responses",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a DseeCompatAuthenticationApiSettingsResponse object into the model struct
func readAuthenticationApiSettingsResponseDataSource(ctx context.Context, r *client.AuthnApiSettings, state *authenticationApiSettingsDataSourceModel, expectedValues *authenticationApiSettingsDataSourceModel, diags *diag.Diagnostics) {
	state.Id = types.StringValue("id")
	state.ApiEnabled = types.BoolValue(*r.ApiEnabled)
	state.EnableApiDescriptions = types.BoolValue(*r.EnableApiDescriptions)
	state.RestrictAccessToRedirectlessMode = types.BoolValue(*r.RestrictAccessToRedirectlessMode)
	state.IncludeRequestContext = types.BoolValue(*r.IncludeRequestContext)
	resourceLinkObjectValue := internaltypes.ToStateResourceLink(ctx, r.GetDefaultApplicationRef())
	state.DefaultApplicationRef = resourceLinkObjectValue
}

// Read resource information
func (r *authenticationApiSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationApiSettingsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAuthenticationApiSettings, httpResp, err := r.apiClient.AuthenticationApiApi.GetAuthenticationApiSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Authentication Api Settings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAuthenticationApiSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAuthenticationApiSettingsResponseDataSource(ctx, apiReadAuthenticationApiSettings, &state, &state, &diags)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
