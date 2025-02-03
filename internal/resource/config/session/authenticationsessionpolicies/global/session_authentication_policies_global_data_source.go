// Copyright © 2025 Ping Identity Corporation

package sessionauthenticationsessionpoliciesglobal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sessionAuthenticationPoliciesGlobalDataSource{}
	_ datasource.DataSourceWithConfigure = &sessionAuthenticationPoliciesGlobalDataSource{}
)

// SessionAuthenticationPoliciesGlobalResource is a helper function to simplify the provider implementation.
func SessionAuthenticationPoliciesGlobalDataSource() datasource.DataSource {
	return &sessionAuthenticationPoliciesGlobalDataSource{}
}

// sessionAuthenticationPoliciesGlobalResource is the resource implementation.
type sessionAuthenticationPoliciesGlobalDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *sessionAuthenticationPoliciesGlobalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *sessionAuthenticationPoliciesGlobalDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_authentication_policies_global"
}

func (r *sessionAuthenticationPoliciesGlobalDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *sessionAuthenticationPoliciesGlobalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sessionAuthenticationPoliciesGlobalModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionAuthenticationPoliciesGlobal, httpResp, err := r.apiClient.SessionAPI.GetGlobalPolicy(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the global session authentication session policies settings.", err, httpResp)
		return
	}

	// Read the response into the state
	readSessionAuthenticationPoliciesGlobalResponse(ctx, apiReadSessionAuthenticationPoliciesGlobal, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
