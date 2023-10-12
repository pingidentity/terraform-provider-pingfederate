package authenticationapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
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
				Attributes:  resourcelink.AddResourceLinkDataSourceSchema(),
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
	id.DataSourceSchema(&schemaDef, false, "The ID of this resource.")
	resp.Schema = schemaDef
}

// Read a AuthenticationApiSettingsResponse object into the model struct
func readAuthenticationApiSettingsResponseDataSource(ctx context.Context, r *client.AuthnApiSettings, state *authenticationApiSettingsDataSourceModel, expectedValues *authenticationApiSettingsDataSourceModel) diag.Diagnostics {
	//TODO different placeholder?
	state.Id = types.StringValue("id")
	state.ApiEnabled = types.BoolValue(*r.ApiEnabled)
	state.EnableApiDescriptions = types.BoolValue(*r.EnableApiDescriptions)
	state.RestrictAccessToRedirectlessMode = types.BoolValue(*r.RestrictAccessToRedirectlessMode)
	state.IncludeRequestContext = types.BoolValue(*r.IncludeRequestContext)
	var valueFromDiags diag.Diagnostics
	resourceLinkObjectValue, valueFromDiags := resourcelink.ToState(ctx, r.DefaultApplicationRef)
	state.DefaultApplicationRef = resourceLinkObjectValue
	return valueFromDiags
}

// Read resource information
func (r *authenticationApiSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationApiSettingsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAuthenticationApiSettings, httpResp, err := r.apiClient.AuthenticationApiAPI.GetAuthenticationApiSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Authentication Api Settings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, responseErr := apiReadAuthenticationApiSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	} else {
		diags.AddError("There was an issue retrieving the response of the Authentication API Settings: %s", responseErr.Error())
	}

	// Read the response into the state
	readAuthenticationApiSettingsResponseDataSource(ctx, apiReadAuthenticationApiSettings, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
