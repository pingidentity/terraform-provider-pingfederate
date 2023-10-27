package oauth

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopeentry"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopegroupentry"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &oauthAuthServerSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &oauthAuthServerSettingsDataSource{}
)

// Create a Administrative Account data source
func NewOauthAuthServerSettingsDataSource() datasource.DataSource {
	return &oauthAuthServerSettingsDataSource{}
}

// oauthAuthServerSettingsDataSource is the datasource implementation.
type oauthAuthServerSettingsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAuthServerSettingsDataSourceModel struct {
	Id                                          types.String `tfsdk:"id"`
	DefaultScopeDescription                     types.String `tfsdk:"default_scope_description"`
	Scopes                                      types.Set    `tfsdk:"scopes"`
	ScopeGroups                                 types.Set    `tfsdk:"scope_groups"`
	ExclusiveScopes                             types.Set    `tfsdk:"exclusive_scopes"`
	ExclusiveScopeGroups                        types.Set    `tfsdk:"exclusive_scope_groups"`
	AuthorizationCodeTimeout                    types.Int64  `tfsdk:"authorization_code_timeout"`
	AuthorizationCodeEntropy                    types.Int64  `tfsdk:"authorization_code_entropy"`
	DisallowPlainPKCE                           types.Bool   `tfsdk:"disallow_plain_pkce"`
	IncludeIssuerInAuthorizationResponse        types.Bool   `tfsdk:"include_issuer_in_authorization_response"`
	TrackUserSessionsForLogout                  types.Bool   `tfsdk:"track_user_sessions_for_logout"`
	TokenEndpointBaseUrl                        types.String `tfsdk:"token_endpoint_base_url"`
	PersistentGrantLifetime                     types.Int64  `tfsdk:"persistent_grant_lifetime"`
	PersistentGrantLifetimeUnit                 types.String `tfsdk:"persistent_grant_lifetime_unit"`
	PersistentGrantIdleTimeout                  types.Int64  `tfsdk:"persistent_grant_idle_timeout"`
	PersistentGrantIdleTimeoutTimeUnit          types.String `tfsdk:"persistent_grant_idle_timeout_time_unit"`
	RefreshTokenLength                          types.Int64  `tfsdk:"refresh_token_length"`
	RollRefreshTokenValues                      types.Bool   `tfsdk:"roll_refresh_token_values"`
	RefreshTokenRollingGracePeriod              types.Int64  `tfsdk:"refresh_token_rolling_grace_period"`
	RefreshRollingInterval                      types.Int64  `tfsdk:"refresh_rolling_interval"`
	PersistentGrantReuseGrantTypes              types.Set    `tfsdk:"persistent_grant_reuse_grant_types"`
	PersistentGrantContract                     types.Object `tfsdk:"persistent_grant_contract"`
	BypassAuthorizationForApprovedGrants        types.Bool   `tfsdk:"bypass_authorization_for_approved_grants"`
	AllowUnidentifiedClientROCreds              types.Bool   `tfsdk:"allow_unidentified_client_ro_creds"`
	AllowUnidentifiedClientExtensionGrants      types.Bool   `tfsdk:"allow_unidentified_client_extension_grants"`
	AdminWebServicePcvRef                       types.Object `tfsdk:"admin_web_service_pcv_ref"`
	AtmIdForOAuthGrantManagement                types.String `tfsdk:"atm_id_for_oauth_grant_management"`
	ScopeForOAuthGrantManagement                types.String `tfsdk:"scope_for_oauth_grant_management"`
	AllowedOrigins                              types.List   `tfsdk:"allowed_origins"`
	UserAuthorizationUrl                        types.String `tfsdk:"user_authorization_url"`
	BypassActivationCodeConfirmation            types.Bool   `tfsdk:"bypass_activation_code_confirmation"`
	RegisteredAuthorizationPath                 types.String `tfsdk:"registered_authorization_path"`
	PendingAuthorizationTimeout                 types.Int64  `tfsdk:"pending_authorization_timeout"`
	DevicePollingInterval                       types.Int64  `tfsdk:"device_polling_interval"`
	ActivationCodeCheckMode                     types.String `tfsdk:"activation_code_check_mode"`
	UserAuthorizationConsentPageSetting         types.String `tfsdk:"user_authorization_consent_page_setting"`
	UserAuthorizationConsentAdapter             types.String `tfsdk:"user_authorization_consent_adapter"`
	ApprovedScopesAttribute                     types.String `tfsdk:"approved_scopes_attribute"`
	ApprovedAuthorizationDetailAttribute        types.String `tfsdk:"approved_authorization_detail_attribute"`
	ParReferenceTimeout                         types.Int64  `tfsdk:"par_reference_timeout"`
	ParReferenceLength                          types.Int64  `tfsdk:"par_reference_length"`
	ParStatus                                   types.String `tfsdk:"par_status"`
	ClientSecretRetentionPeriod                 types.Int64  `tfsdk:"client_secret_retention_period"`
	JwtSecuredAuthorizationResponseModeLifetime types.Int64  `tfsdk:"jwt_secured_authorization_response_mode_lifetime"`
}

// GetSchema defines the schema for the datasource.
func (r *oauthAuthServerSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages OAuth Auth Server Settings",
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
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
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
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
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
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
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
							Validators: []validator.String{
								configvalidators.NoWhitespace(),
							},
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
				Description: "Determines whether PKCE's 'plain' code challenge method will be disallowed. The default value is false.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_issuer_in_authorization_response": schema.BoolAttribute{
				Description: "Determines whether the authorization server's issuer value is added to the authorization response or not. The default value is false.",
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
			"persistent_grant_lifetime": schema.Int64Attribute{
				Description: "The persistent grant lifetime. The default value is indefinite. -1 indicates an indefinite amount of time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_lifetime_unit": schema.StringAttribute{
				Description: "The persistent grant lifetime unit.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MINUTES", "DAYS", "HOURS"}...),
				},
			},
			"persistent_grant_idle_timeout": schema.Int64Attribute{
				Description: "The persistent grant idle timeout. The default value is 30 (days). -1 indicates an indefinite amount of time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"persistent_grant_idle_timeout_time_unit": schema.StringAttribute{
				Description: "The persistent grant idle timeout time unit. The default value is DAYS",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MINUTES", "DAYS", "HOURS"}...),
				},
			},
			"refresh_token_length": schema.Int64Attribute{
				Description: "The refresh token length in number of characters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"roll_refresh_token_values": schema.BoolAttribute{
				Description: "The roll refresh token values default policy. The default value is true.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"refresh_token_rolling_grace_period": schema.Int64Attribute{
				Description: "The grace period that a rolled refresh token remains valid in seconds. The default value is 0.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"refresh_rolling_interval": schema.Int64Attribute{
				Description: "The minimum interval to roll refresh tokens, in hours.",
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
				Description: "Bypass authorization for previously approved persistent grants. The default value is false.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_unidentified_client_ro_creds": schema.BoolAttribute{
				Description: "Allow unidentified clients to request resource owner password credentials grants. The default value is false.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_unidentified_client_extension_grants": schema.BoolAttribute{
				Description: "Allow unidentified clients to request extension grants. The default value is false.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"admin_web_service_pcv_ref": schema.SingleNestedAttribute{
				Description: "The password credential validator reference that is used for authenticating access to the OAuth Administrative Web Service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Attributes:  resourcelink.ToDataSourceSchema(),
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
			"allowed_origins": schema.ListAttribute{
				Description: "The list of allowed origins.",
				ElementType: types.StringType,
				Required:    false,
				Optional:    false,
				Computed:    true,
				Validators: []validator.List{
					configvalidators.ValidUrls(),
				},
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
				Validators: []validator.String{
					configvalidators.StartsWith("/"),
				},
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
				Description: "Determines whether the user is prompted to enter or confirm the activation code after authenticating or before. The default is AFTER_AUTHENTICATION.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"AFTER_AUTHENTICATION", "BEFORE_AUTHENTICATION"}...),
				},
			},
			"bypass_activation_code_confirmation": schema.BoolAttribute{
				Description: "Indicates if the Activation Code Confirmation page should be bypassed if 'verification_url_complete' is used by the end user to authorize a device.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_authorization_consent_page_setting": schema.StringAttribute{
				Description: "User Authorization Consent Page setting to use PingFederate's internal consent page or an external system",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"INTERNAL", "ADAPTER"}...),
				},
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
				Description: "The timeout, in seconds, of the pushed authorization request reference. The default value is 60.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"par_reference_length": schema.Int64Attribute{
				Description: "The entropy of pushed authorization request references, in bytes. The default value is 24.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"par_status": schema.StringAttribute{
				Description: "The status of pushed authorization request support. The default value is ENABLED.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"DISABLED", "ENABLED", "REQUIRED"}...),
				},
			},
			"client_secret_retention_period": schema.Int64Attribute{
				Description: "The length of time in minutes that client secrets will be retained as secondary secrets after secret change. The default value is 0, which will disable secondary client secret retention.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jwt_secured_authorization_response_mode_lifetime": schema.Int64Attribute{
				Description: "The lifetime, in seconds, of the JWT Secured authorization response. The default value is 600.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	id.ToSchema(&schemaDef)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *oauthAuthServerSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings"
}

// Configure adds the provider configured client to the data source.
func (r *oauthAuthServerSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a OauthAuthServerSettingsResponse object into the model struct
func readOauthAuthServerSettingsResponseDataSource(ctx context.Context, r *client.AuthorizationServerSettings, state *oauthAuthServerSettingsDataSourceModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = id.GenerateUUIDToState(existingId)
	state.DefaultScopeDescription = types.StringValue(r.DefaultScopeDescription)
	state.Scopes, respDiags = scopeentry.ToState(ctx, r.Scopes)
	diags.Append(respDiags...)
	state.ScopeGroups, respDiags = scopegroupentry.ToState(ctx, r.ScopeGroups)
	diags.Append(respDiags...)
	state.ExclusiveScopes, respDiags = scopeentry.ToState(ctx, r.ExclusiveScopes)
	diags.Append(respDiags...)
	state.ExclusiveScopeGroups, respDiags = scopegroupentry.ToState(ctx, r.ExclusiveScopeGroups)
	diags.Append(respDiags...)
	persistentGrantContract, respDiags := types.ObjectValueFrom(ctx, persistentGrantObjContractTypes, r.PersistentGrantContract)
	diags.Append(respDiags...)

	state.PersistentGrantContract = persistentGrantContract
	state.AuthorizationCodeTimeout = types.Int64Value(r.AuthorizationCodeTimeout)
	state.AuthorizationCodeEntropy = types.Int64Value(r.AuthorizationCodeEntropy)
	state.DisallowPlainPKCE = types.BoolPointerValue(r.DisallowPlainPKCE)
	state.IncludeIssuerInAuthorizationResponse = types.BoolPointerValue(r.IncludeIssuerInAuthorizationResponse)
	state.TrackUserSessionsForLogout = types.BoolPointerValue(r.TrackUserSessionsForLogout)
	state.TokenEndpointBaseUrl = types.StringPointerValue(r.TokenEndpointBaseUrl)
	state.PersistentGrantLifetime = types.Int64PointerValue(r.PersistentGrantLifetime)
	state.PersistentGrantLifetimeUnit = types.StringPointerValue(r.PersistentGrantLifetimeUnit)
	state.PersistentGrantIdleTimeout = types.Int64PointerValue(r.PersistentGrantIdleTimeout)
	state.PersistentGrantIdleTimeoutTimeUnit = types.StringPointerValue(r.PersistentGrantIdleTimeoutTimeUnit)
	state.RefreshTokenLength = types.Int64Value(r.RefreshTokenLength)
	state.RollRefreshTokenValues = types.BoolPointerValue(r.RollRefreshTokenValues)
	state.RefreshTokenRollingGracePeriod = types.Int64PointerValue(r.RefreshTokenRollingGracePeriod)
	state.RefreshRollingInterval = types.Int64Value(r.RefreshRollingInterval)
	state.PersistentGrantReuseGrantTypes = internaltypes.GetStringSet(r.PersistentGrantReuseGrantTypes)
	state.BypassAuthorizationForApprovedGrants = types.BoolPointerValue(r.BypassAuthorizationForApprovedGrants)
	state.AllowUnidentifiedClientROCreds = types.BoolPointerValue(r.AllowUnidentifiedClientROCreds)
	state.AllowUnidentifiedClientExtensionGrants = types.BoolPointerValue(r.AllowUnidentifiedClientExtensionGrants)
	state.AdminWebServicePcvRef, respDiags = resourcelink.ToDataSourceState(ctx, r.AdminWebServicePcvRef)
	diags.Append(respDiags...)
	state.AtmIdForOAuthGrantManagement = types.StringPointerValue(r.AtmIdForOAuthGrantManagement)
	state.ScopeForOAuthGrantManagement = types.StringPointerValue(r.ScopeForOAuthGrantManagement)
	state.AllowedOrigins = internaltypes.GetStringList(r.AllowedOrigins)
	state.UserAuthorizationUrl = types.StringPointerValue(r.UserAuthorizationUrl)
	state.RegisteredAuthorizationPath = types.StringValue(r.RegisteredAuthorizationPath)
	state.PendingAuthorizationTimeout = types.Int64Value(r.PendingAuthorizationTimeout)
	state.DevicePollingInterval = types.Int64Value(r.DevicePollingInterval)
	state.ActivationCodeCheckMode = types.StringPointerValue(r.ActivationCodeCheckMode)
	state.BypassActivationCodeConfirmation = types.BoolValue(r.BypassActivationCodeConfirmation)
	state.UserAuthorizationConsentPageSetting = types.StringPointerValue(r.UserAuthorizationConsentPageSetting)
	state.UserAuthorizationConsentAdapter = types.StringPointerValue(r.UserAuthorizationConsentAdapter)
	state.ApprovedScopesAttribute = types.StringPointerValue(r.ApprovedScopesAttribute)
	state.ApprovedAuthorizationDetailAttribute = types.StringPointerValue(r.ApprovedAuthorizationDetailAttribute)
	state.ParReferenceTimeout = types.Int64PointerValue(r.ParReferenceTimeout)
	state.ParReferenceLength = types.Int64PointerValue(r.ParReferenceLength)
	state.ParStatus = types.StringPointerValue(r.ParStatus)
	state.ClientSecretRetentionPeriod = types.Int64PointerValue(r.ClientSecretRetentionPeriod)
	state.JwtSecuredAuthorizationResponseModeLifetime = types.Int64PointerValue(r.JwtSecuredAuthorizationResponseModeLifetime)
	return diags
}

// Read resource information
func (r *oauthAuthServerSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state oauthAuthServerSettingsDataSourceModel

	var diags diag.Diagnostics

	apiReadOauthAuthServerSettings, httpResp, err := r.apiClient.OauthAuthServerSettingsAPI.GetAuthorizationServerSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Auth Server Settings", err, httpResp)
		return
	}

	// Read the response into the state
	var id = "id"
	diags = readOauthAuthServerSettingsResponseDataSource(ctx, apiReadOauthAuthServerSettings, &state, &id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
