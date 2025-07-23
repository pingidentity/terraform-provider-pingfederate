// Copyright © 2025 Ping Identity Corporation

package oauthauthserversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	resourcelinkdatasource "github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthServerSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthServerSettingsDataSource{}
)

// Create a Administrative Account data source
func OauthServerSettingsDataSource() datasource.DataSource {
	return &oauthServerSettingsDataSource{}
}

// oauthServerSettingsDataSource is the datasource implementation.
type oauthServerSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *oauthServerSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes the OAuth authorization server settings.",
		Attributes: map[string]schema.Attribute{
			"default_scope_description": schema.StringAttribute{
				Description: "The default scope description.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scopes": schema.SetNestedAttribute{
				Description: "The list of common scopes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope.",
							Required:    false,
							Optional:    false,
							Computed:    true,
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
				},
			},
			"scope_groups": schema.SetNestedAttribute{
				Description: "The list of common scope groups.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"scopes": schema.SetAttribute{
							Description: "The set of scopes for this scope group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"exclusive_scopes": schema.SetNestedAttribute{
				Description: "The list of exclusive scopes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope that appears when the user is prompted for authorization.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"dynamic": schema.BoolAttribute{
							Description: "True if the scope is dynamic.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
			"exclusive_scope_groups": schema.SetNestedAttribute{
				Description: "The list of exclusive scope groups.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the scope group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the scope group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"scopes": schema.SetAttribute{
							Description: "The set of scopes for this scope group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"authorization_code_timeout": schema.Int64Attribute{
				Description: "The authorization code timeout, in seconds.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"authorization_code_entropy": schema.Int64Attribute{
				Description: "The authorization code entropy, in bytes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disallow_plain_pkce": schema.BoolAttribute{
				Description: "Determines whether PKCE's 'plain' code challenge method will be disallowed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_issuer_in_authorization_response": schema.BoolAttribute{
				Description: "Determines whether the authorization server's issuer value is added to the authorization response or not.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"track_user_sessions_for_logout": schema.BoolAttribute{
				Description: "Determines whether user sessions are tracked for logout. If this property is not provided on a PUT, the setting is left unchanged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"token_endpoint_base_url": schema.StringAttribute{
				Description: "The token endpoint base URL used to validate the 'aud' claim during Private Key JWT Client Authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_offline_access_scope_to_issue_refresh_tokens": schema.BoolAttribute{
				Description: "Determines whether offline_access scope is required to issue refresh tokens or not.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"offline_access_require_consent_prompt": schema.BoolAttribute{
				Description: "Determines whether offline_access requires the prompt parameter value be 'consent' or not.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_lifetime": schema.Int64Attribute{
				Description: "The persistent grant lifetime. `-1` indicates an indefinite amount of time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_lifetime_unit": schema.StringAttribute{
				Description: "The persistent grant lifetime unit.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout. `-1` indicates an indefinite amount of time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"refresh_token_length": schema.Int64Attribute{
				Description: "The refresh token length in number of characters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"roll_refresh_token_values": schema.BoolAttribute{
				Description: "The roll refresh token values default policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"refresh_token_rolling_grace_period": schema.Int64Attribute{
				Description: "The grace period that a rolled refresh token remains valid in seconds.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"refresh_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"refresh_rolling_interval_time_unit": schema.StringAttribute{
				Description: "The refresh token rolling interval time unit.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_reuse_grant_types": schema.SetAttribute{
				Description: "The grant types that the OAuth AS can reuse rather than creating a new grant for each request. Only 'IMPLICIT' or 'AUTHORIZATION_CODE' or 'RESOURCE_OWNER_CREDENTIALS' are valid grant types.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"persistent_grant_contract": schema.SingleNestedAttribute{
				Description: "The persistent grant contract defines attributes that are associated with OAuth persistent grants.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "This is a read-only list of persistent grant attributes and includes USER_KEY and USER_NAME. Changes to this field will be ignored.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes for the persistent grant contract.",
						Required:    false,
						Optional:    false,
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    false,
									Optional:    false,
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"bypass_authorization_for_approved_grants": schema.BoolAttribute{
				Description: "Bypass authorization for previously approved persistent grants. ",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_unidentified_client_ro_creds": schema.BoolAttribute{
				Description: "Allow unidentified clients to request resource owner password credentials grants.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_unidentified_client_extension_grants": schema.BoolAttribute{
				Description: "Allow unidentified clients to request extension grants.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"admin_web_service_pcv_ref": schema.SingleNestedAttribute{
				Description: "The password credential validator reference that is used for authenticating access to the OAuth Administrative Web Service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelinkdatasource.ToDataSourceSchema(),
			},
			"atm_id_for_oauth_grant_management": schema.StringAttribute{
				Description: "The ID of the Access Token Manager used for OAuth enabled grant management.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scope_for_oauth_grant_management": schema.StringAttribute{
				Description: "The OAuth scope to validate when accessing grant management service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_origins": schema.SetAttribute{
				Description: "The list of allowed origins.",
				ElementType: types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_authorization_url": schema.StringAttribute{
				Description: "The URL used to generate 'verification_url' and 'verification_url_complete' values in a Device Authorization request",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"registered_authorization_path": schema.StringAttribute{
				Description: "The Registered Authorization Path is concatenated to PingFederate base URL to generate 'verification_url' and 'verification_url_complete' values in a Device Authorization request. PingFederate listens to this path if specified",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pending_authorization_timeout": schema.Int64Attribute{
				Description: "The 'device_code' and 'user_code' timeout, in seconds.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"device_polling_interval": schema.Int64Attribute{
				Description: "The amount of time client should wait between polling requests, in seconds.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"activation_code_check_mode": schema.StringAttribute{
				Description: "Determines whether the user is prompted to enter or confirm the activation code after authenticating or before.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bypass_activation_code_confirmation": schema.BoolAttribute{
				Description: "Indicates if the Activation Code Confirmation page should be bypassed if 'verification_url_complete' is used by the end user to authorize a device.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_cookieless_user_authorization_authentication_api": schema.BoolAttribute{
				Description: "Indicates if cookies should be used for state tracking when the user authorization endpoint is operating in authentication API redirectless mode",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_authorization_consent_page_setting": schema.StringAttribute{
				Description: "User Authorization Consent Page setting to use PingFederate's internal consent page or an external system",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_authorization_consent_adapter": schema.StringAttribute{
				Description: "Adapter ID of the external consent adapter to be used for the consent page user interface.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"approved_scopes_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved scopes returned by the external consent page.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"approved_authorization_detail_attribute": schema.StringAttribute{
				Description: "Attribute from the external consent adapter's contract, intended for storing approved authorization details returned by the external consent page.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"par_reference_timeout": schema.Int64Attribute{
				Description: "The timeout, in seconds, of the pushed authorization request reference.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"par_reference_length": schema.Int64Attribute{
				Description: "The entropy of pushed authorization request references, in bytes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"par_status": schema.StringAttribute{
				Description: "The status of pushed authorization request support.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is 0, which will disable secondary client secret retention.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jwt_secured_authorization_response_mode_lifetime": schema.Int64Attribute{
				Description: "The lifetime, in seconds, of the JWT Secured authorization response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dpop_proof_require_nonce": schema.BoolAttribute{
				Description: "Determines whether nonce is required in the Demonstrating Proof-of-Possession (DPoP) proof JWT. Supported in PF version `11.3` or later.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dpop_proof_lifetime_seconds": schema.Int64Attribute{
				Description: "The lifetime, in seconds, of the Demonstrating Proof-of-Possession (DPoP) proof JWT. Supported in PF version `11.3` or later.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dpop_proof_enforce_replay_prevention": schema.BoolAttribute{
				Description: "Determines whether Demonstrating Proof-of-Possession (DPoP) proof JWT replay prevention is enforced. Supported in PF version `11.3` or later.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bypass_authorization_for_approved_consents": schema.BoolAttribute{
				Description: "Bypass authorization for previously approved consents. Supported in PF version 12.0 or later.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"consent_lifetime_days": schema.Int64Attribute{
				Description: "The consent lifetime in days. The default value is indefinite. -1 indicates an indefinite amount of time. Supported in PF version 12.0 or later.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"return_id_token_on_open_id_with_device_authz_grant": schema.BoolAttribute{
				Description: "Indicates if an ID token should be returned during the device authorization grant flow when the 'openid' scope is approved. The default is `false`. Supported in PF version `12.2` or later.",
				Computed:    true,
			},
		},
	}
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *oauthServerSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_server_settings"
}

// Configure adds the provider configured client to the data source.
func (r *oauthServerSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read resource information
func (r *oauthServerSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthServerSettingsModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadOauthServerSettings, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetAuthorizationServerSettings(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Auth Server Settings", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readOauthServerSettingsResponse(ctx, apiReadOauthServerSettings, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
