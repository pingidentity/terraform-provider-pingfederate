package authenticationpoliciessettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &authenticationPoliciesSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &authenticationPoliciesSettingsDataSource{}
)

// AuthenticationPoliciesSettingsDataSource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesSettingsDataSource() datasource.DataSource {
	return &authenticationPoliciesSettingsDataSource{}
}

// authenticationPoliciesSettingsDataSource is the datasource implementation.
type authenticationPoliciesSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *authenticationPoliciesSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Authentication Policies Settings",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the resource.",
				Computed:    true,
				Optional:    false,
			},
			"enable_idp_authn_selection": schema.BoolAttribute{
				Description: "Enable IdP authentication policies.",
				Optional:    true,
				Computed:    true,
			},
			"enable_sp_authn_selection": schema.BoolAttribute{
				Description: "Enable SP authentication policies.",
				Optional:    true,
				Computed:    true,
			},
		},
	}

	// // Set attributes in string list
	// if setOptionalToComputed {
	// 	config.SetAllAttributesToOptionalAndComputed(&schema, []string{"FIX_ME"})
	// }
	// config.AddCommonSchema(&schema, false)
	id.ToDataSourceSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the datasource type name.
func (r *authenticationPoliciesSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies_settings"
}

func (r *authenticationPoliciesSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *authenticationPoliciesSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state authenticationPoliciesSettingsModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAuthenticationPoliciesSettings, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.GetAuthenticationPolicySettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the authentication policies settings", err, httpResp)
			return
		}

		readAuthenticationPoliciesSettings(ctx, apiReadAuthenticationPoliciesSettings, &state, id)

		// Set refreshed state
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
	}
}
