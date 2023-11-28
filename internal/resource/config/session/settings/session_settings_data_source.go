package sessionsettings

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
	_ datasource.DataSource              = &sessionSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionSettingsDataSource{}
)

func SessionSettingsDataSource() datasource.DataSource {
	return &sessionSettingsDataSource{}
}

// sessionSettingsResource is the resource implementation.
type sessionSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type sessionSettingsDataSourceModel struct {
	Id                            types.String `tfsdk:"id"`
	TrackAdapterSessionsForLogout types.Bool   `tfsdk:"track_adapter_sessions_for_logout"`
	RevokeUserSessionOnLogout     types.Bool   `tfsdk:"revoke_user_session_on_logout"`
	SessionRevocationLifetime     types.Int64  `tfsdk:"session_revocation_lifetime"`
}

// GetSchema defines the schema for the resource.
func (r *sessionSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the general session management settings.",
		Attributes: map[string]schema.Attribute{
			"track_adapter_sessions_for_logout": schema.BoolAttribute{
				Description: "Determines whether adapter sessions are tracked for cleanup during single logout. The default is false.",
				Computed:    true,
				Optional:    false,
			},
			"revoke_user_session_on_logout": schema.BoolAttribute{
				Description: "Determines whether the user's session is revoked on logout.",
				Computed:    true,
				Optional:    false,
			},
			"session_revocation_lifetime": schema.Int64Attribute{
				Description: "How long a session revocation is tracked and stored, in minutes.",
				Computed:    true,
				Optional:    false,
			},
		},
	}
	id.ToDataSourceSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *sessionSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_settings"
}

func (r *sessionSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readSessionSettingsDataSource(ctx context.Context, r *client.SessionSettings, state *sessionSettingsDataSourceModel) {
	state.Id = types.StringValue("session_settings_id")
	state.TrackAdapterSessionsForLogout = types.BoolPointerValue(r.TrackAdapterSessionsForLogout)
	state.RevokeUserSessionOnLogout = types.BoolPointerValue(r.RevokeUserSessionOnLogout)
	state.SessionRevocationLifetime = types.Int64PointerValue(r.SessionRevocationLifetime)
}

func (r *sessionSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionSettingsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionSettings, httpResp, err := r.apiClient.SessionAPI.GetSessionSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the session settings", err, httpResp)
		return
	}

	// Read the response into the state
	readSessionSettingsDataSource(ctx, apiReadSessionSettings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
