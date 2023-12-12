package oauthauthserversettingsscopescommonscope

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
	_ datasource.DataSource              = &oauthAuthServerSettingsScopesCommonScopeDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthAuthServerSettingsScopesCommonScopeDataSource{}
)

// Create a Administrative Account data source
func OauthAuthServerSettingsScopesCommonScopeDataSource() datasource.DataSource {
	return &oauthAuthServerSettingsScopesCommonScopeDataSource{}
}

// oauthAuthServerSettingsScopesCommonScopeDataSource is the datasource implementation.
type oauthAuthServerSettingsScopesCommonScopeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *oauthAuthServerSettingsScopesCommonScopeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a common scope in the authorization server settings.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the scope.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the scope that appears when the user is prompted for authorization.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dynamic": schema.BoolAttribute{
				Description: "True if the scope is dynamic. (Defaults to false)",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	id.ToDataSourceSchema(&schemaDef)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *oauthAuthServerSettingsScopesCommonScopeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings_scopes_common_scope"
}

// Configure adds the provider configured client to the data source.
func (r *oauthAuthServerSettingsScopesCommonScopeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read resource information
func (r *oauthAuthServerSettingsScopesCommonScopeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthAuthServerSettingsScopesCommonScopeModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOauthAuthServerSettingsScopesCommonScope, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetCommonScope(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Auth Server Settings Scopes Common Scope", err, httpResp)
		return
	}

	// Read the response into the state
	readOauthAuthServerSettingsScopesCommonScopeResponse(ctx, apiReadOauthAuthServerSettingsScopesCommonScope, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
