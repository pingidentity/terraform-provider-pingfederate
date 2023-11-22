package sessionauthenticationsessionpoliciesglobal

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
	_ datasource.DataSource              = &sessionAuthenticationSessionPoliciesGlobalDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionAuthenticationSessionPoliciesGlobalDataSource{}
)

// SessionAuthenticationSessionPoliciesGlobalResource is a helper function to simplify the provider implementation.
func NewSessionAuthenticationSessionPoliciesGlobalDataSource() datasource.DataSource {
	return &sessionAuthenticationSessionPoliciesGlobalDataSource{}
}

// sessionAuthenticationSessionPoliciesGlobalResource is the resource implementation.
type sessionAuthenticationSessionPoliciesGlobalDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type sessionAuthenticationSessionPoliciesGlobalDataSourceModel struct {
	Id                         types.String `tfsdk:"id"`
	EnableSessions             types.Bool   `tfsdk:"enable_sessions"`
	PersistentSessions         types.Bool   `tfsdk:"persistent_sessions"`
	HashUniqueUserKeyAttribute types.Bool   `tfsdk:"hash_unique_user_key_attribute"`
	IdleTimeoutMins            types.Int64  `tfsdk:"idle_timeout_mins"`
	IdleTimeoutDisplayUnit     types.String `tfsdk:"idle_timeout_display_unit"`
	MaxTimeoutMins             types.Int64  `tfsdk:"max_timeout_mins"`
	MaxTimeoutDisplayUnit      types.String `tfsdk:"max_timeout_display_unit"`
}

// GetSchema defines the schema for the resource.
func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the global settings for authentication session policies.",
		Attributes: map[string]schema.Attribute{
			"enable_sessions": schema.BoolAttribute{
				Description: "Determines whether authentication sessions are enabled globally.",
				Computed:    true,
				Optional:    false,
			},
			"persistent_sessions": schema.BoolAttribute{
				Description: "Determines whether authentication sessions are persistent by default. Persistent sessions are linked to a persistent cookie and stored in a data store. This field is ignored if enableSessions is false.",
				Computed:    true,
				Optional:    false,
			},
			"hash_unique_user_key_attribute": schema.BoolAttribute{
				Description: "Determines whether to hash the value of the unique user key attribute.",
				Computed:    true,
				Optional:    false,
			},
			"idle_timeout_mins": schema.Int64Attribute{
				Description: "The idle timeout period, in minutes. If set to -1, the idle timeout will be set to the maximum timeout. The default is 60.",
				Computed:    true,
				Optional:    false,
			},
			"idle_timeout_display_unit": schema.StringAttribute{
				Description: "The display unit for the idle timeout period in the PingFederate administrative console. When the display unit is HOURS or DAYS, the timeout value in minutes must correspond to a whole number value for the specified unit. [ MINUTES, HOURS, DAYS ]",
				Computed:    true,
				Optional:    false,
			},
			"max_timeout_mins": schema.Int64Attribute{
				Description: "The maximum timeout period, in minutes. If set to -1, sessions do not expire. The default is 480.",
				Computed:    true,
				Optional:    false,
			},
			"max_timeout_display_unit": schema.StringAttribute{
				Description: "The display unit for the maximum timeout period in the PingFederate administrative console. When the display unit is HOURS or DAYS, the timeout value in minutes must correspond to a whole number value for the specified unit.",
				Computed:    true,
				Optional:    false,
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_authentication_session_policies_global"
}

func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readSessionAuthenticationSessionPoliciesGlobalDataSource(ctx context.Context, r *client.GlobalAuthenticationSessionPolicy, state *sessionAuthenticationSessionPoliciesGlobalDataSourceModel) {
	state.Id = types.StringValue("session_authentication_session_policies_global_id")
	state.EnableSessions = types.BoolValue(r.EnableSessions)
	state.PersistentSessions = types.BoolPointerValue(r.PersistentSessions)
	state.HashUniqueUserKeyAttribute = types.BoolPointerValue(r.HashUniqueUserKeyAttribute)
	state.IdleTimeoutMins = types.Int64PointerValue(r.IdleTimeoutMins)
	state.IdleTimeoutDisplayUnit = types.StringPointerValue(r.IdleTimeoutDisplayUnit)
	state.MaxTimeoutMins = types.Int64PointerValue(r.MaxTimeoutMins)
	state.MaxTimeoutDisplayUnit = types.StringPointerValue(r.MaxTimeoutDisplayUnit)
}

func (r *sessionAuthenticationSessionPoliciesGlobalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionAuthenticationSessionPoliciesGlobalDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionAuthenticationSessionPoliciesGlobal, httpResp, err := r.apiClient.SessionAPI.GetGlobalPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the global session authentication session policies settings.", err, httpResp)
		return
	}

	// Read the response into the state
	readSessionAuthenticationSessionPoliciesGlobalDataSource(ctx, apiReadSessionAuthenticationSessionPoliciesGlobal, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
