package oauth

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
	_ datasource.DataSource              = &oauthAuthServerSettingsScopesExclusiveScopesDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthAuthServerSettingsScopesExclusiveScopesDataSource{}
)

// Create a Administrative Account data source
func NewOauthAuthServerSettingsScopesExclusiveScopesDataSource() datasource.DataSource {
	return &oauthAuthServerSettingsScopesExclusiveScopesDataSource{}
}

// oauthAuthServerSettingsScopesExclusiveScopesDataSource is the datasource implementation.
type oauthAuthServerSettingsScopesExclusiveScopesDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAuthServerSettingsScopesExclusiveScopesDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Dynamic     types.Bool   `tfsdk:"dynamic"`
}

// GetSchema defines the schema for the datasource.
func (r *oauthAuthServerSettingsScopesExclusiveScopesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a OauthAuthServerSettingsScopesCommonScopes.",
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
	id.ToDataSourceSchema(&schemaDef, false, "The ID of this resource.")
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *oauthAuthServerSettingsScopesExclusiveScopesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings_scopes_exclusive_scope"
}

// Configure adds the provider configured client to the data source.
func (r *oauthAuthServerSettingsScopesExclusiveScopesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a OauthAuthServerSettingsScopesExclusiveScopesResponse object into the model struct
func readOauthAuthServerSettingsScopesExclusiveScopesResponseDataSource(ctx context.Context, r *client.ScopeEntry, state *oauthAuthServerSettingsScopesExclusiveScopesDataSourceModel) {
	state.Id = types.StringValue(r.Name)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringValue(r.Description)
	state.Dynamic = types.BoolPointerValue(r.Dynamic)
}

// Read resource information
func (r *oauthAuthServerSettingsScopesExclusiveScopesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthAuthServerSettingsScopesExclusiveScopesDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOauthAuthServerSettingsScopesExclusiveScopes, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetExclusiveScope(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an OAuth Auth Server Settings Scopes Exclusive Scope", err, httpResp)
		return
	}

	// Read the response into the state
	readOauthAuthServerSettingsScopesExclusiveScopesResponseDataSource(ctx, apiReadOauthAuthServerSettingsScopesExclusiveScopes, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
