// Copyright © 2025 Ping Identity Corporation

package oauthauthserversettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopeentry"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/scopegroupentry"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	attributeAttrTypes = map[string]attr.Type{
		"name": types.StringType,
	}
	attributeSetElementType = types.ObjectType{AttrTypes: attributeAttrTypes}
)

var (
	nameAttributeType = map[string]attr.Type{
		"name": types.StringType,
	}
	persistentGrantObjContractTypes = map[string]attr.Type{
		"core_attributes":     types.SetType{ElemType: types.ObjectType{AttrTypes: nameAttributeType}},
		"extended_attributes": types.SetType{ElemType: types.ObjectType{AttrTypes: nameAttributeType}},
	}
)

type oauthServerSettingsModel struct {
	DefaultScopeDescription                            types.String `tfsdk:"default_scope_description"`
	Scopes                                             types.Set    `tfsdk:"scopes"`
	ScopeGroups                                        types.Set    `tfsdk:"scope_groups"`
	ExclusiveScopes                                    types.Set    `tfsdk:"exclusive_scopes"`
	ExclusiveScopeGroups                               types.Set    `tfsdk:"exclusive_scope_groups"`
	AuthorizationCodeTimeout                           types.Int64  `tfsdk:"authorization_code_timeout"`
	AuthorizationCodeEntropy                           types.Int64  `tfsdk:"authorization_code_entropy"`
	DisallowPlainPKCE                                  types.Bool   `tfsdk:"disallow_plain_pkce"`
	IncludeIssuerInAuthorizationResponse               types.Bool   `tfsdk:"include_issuer_in_authorization_response"`
	TrackUserSessionsForLogout                         types.Bool   `tfsdk:"track_user_sessions_for_logout"`
	TokenEndpointBaseUrl                               types.String `tfsdk:"token_endpoint_base_url"`
	RequireOfflineAccessScopeToIssueRefreshTokens      types.Bool   `tfsdk:"require_offline_access_scope_to_issue_refresh_tokens"`
	OfflineAccessRequireConsentPrompt                  types.Bool   `tfsdk:"offline_access_require_consent_prompt"`
	PersistentGrantLifetime                            types.Int64  `tfsdk:"persistent_grant_lifetime"`
	PersistentGrantLifetimeUnit                        types.String `tfsdk:"persistent_grant_lifetime_unit"`
	PersistentGrantIdleTimeout                         types.Int64  `tfsdk:"persistent_grant_idle_timeout"`
	PersistentGrantIdleTimeoutTimeUnit                 types.String `tfsdk:"persistent_grant_idle_timeout_time_unit"`
	RefreshTokenLength                                 types.Int64  `tfsdk:"refresh_token_length"`
	RollRefreshTokenValues                             types.Bool   `tfsdk:"roll_refresh_token_values"`
	RefreshTokenRollingGracePeriod                     types.Int64  `tfsdk:"refresh_token_rolling_grace_period"`
	RefreshRollingInterval                             types.Int64  `tfsdk:"refresh_rolling_interval"`
	RefreshRollingIntervalTimeUnit                     types.String `tfsdk:"refresh_rolling_interval_time_unit"`
	PersistentGrantReuseGrantTypes                     types.Set    `tfsdk:"persistent_grant_reuse_grant_types"`
	PersistentGrantContract                            types.Object `tfsdk:"persistent_grant_contract"`
	BypassAuthorizationForApprovedGrants               types.Bool   `tfsdk:"bypass_authorization_for_approved_grants"`
	AllowUnidentifiedClientROCreds                     types.Bool   `tfsdk:"allow_unidentified_client_ro_creds"`
	AllowUnidentifiedClientExtensionGrants             types.Bool   `tfsdk:"allow_unidentified_client_extension_grants"`
	AdminWebServicePcvRef                              types.Object `tfsdk:"admin_web_service_pcv_ref"`
	AtmIdForOAuthGrantManagement                       types.String `tfsdk:"atm_id_for_oauth_grant_management"`
	ScopeForOAuthGrantManagement                       types.String `tfsdk:"scope_for_oauth_grant_management"`
	AllowedOrigins                                     types.Set    `tfsdk:"allowed_origins"`
	UserAuthorizationUrl                               types.String `tfsdk:"user_authorization_url"`
	BypassActivationCodeConfirmation                   types.Bool   `tfsdk:"bypass_activation_code_confirmation"`
	EnableCookielessUserAuthorizationAuthenticationApi types.Bool   `tfsdk:"enable_cookieless_user_authorization_authentication_api"`
	RegisteredAuthorizationPath                        types.String `tfsdk:"registered_authorization_path"`
	PendingAuthorizationTimeout                        types.Int64  `tfsdk:"pending_authorization_timeout"`
	DevicePollingInterval                              types.Int64  `tfsdk:"device_polling_interval"`
	ActivationCodeCheckMode                            types.String `tfsdk:"activation_code_check_mode"`
	UserAuthorizationConsentPageSetting                types.String `tfsdk:"user_authorization_consent_page_setting"`
	UserAuthorizationConsentAdapter                    types.String `tfsdk:"user_authorization_consent_adapter"`
	ApprovedScopesAttribute                            types.String `tfsdk:"approved_scopes_attribute"`
	ApprovedAuthorizationDetailAttribute               types.String `tfsdk:"approved_authorization_detail_attribute"`
	ParReferenceTimeout                                types.Int64  `tfsdk:"par_reference_timeout"`
	ParReferenceLength                                 types.Int64  `tfsdk:"par_reference_length"`
	ParStatus                                          types.String `tfsdk:"par_status"`
	ClientSecretRetentionPeriod                        types.Int64  `tfsdk:"client_secret_retention_period"`
	JwtSecuredAuthorizationResponseModeLifetime        types.Int64  `tfsdk:"jwt_secured_authorization_response_mode_lifetime"`
	DpopProofRequireNonce                              types.Bool   `tfsdk:"dpop_proof_require_nonce"`
	DpopProofLifetimeSeconds                           types.Int64  `tfsdk:"dpop_proof_lifetime_seconds"`
	DpopProofEnforceReplayPrevention                   types.Bool   `tfsdk:"dpop_proof_enforce_replay_prevention"`
	BypassAuthorizationForApprovedConsents             types.Bool   `tfsdk:"bypass_authorization_for_approved_consents"`
	ConsentLifetimeDays                                types.Int64  `tfsdk:"consent_lifetime_days"`
	ReturnIdTokenOnOpenIdWithDeviceAuthzGrant          types.Bool   `tfsdk:"return_id_token_on_open_id_with_device_authz_grant"`
}

func readOauthServerSettingsResponse(ctx context.Context, r *client.AuthorizationServerSettings, state *oauthServerSettingsModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.DefaultScopeDescription = types.StringPointerValue(r.DefaultScopeDescription)
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
	state.RequireOfflineAccessScopeToIssueRefreshTokens = types.BoolPointerValue(r.RequireOfflineAccessScopeToIssueRefreshTokens)
	state.OfflineAccessRequireConsentPrompt = types.BoolPointerValue(r.OfflineAccessRequireConsentPrompt)
	state.PersistentGrantLifetime = types.Int64PointerValue(r.PersistentGrantLifetime)
	state.PersistentGrantLifetimeUnit = types.StringPointerValue(r.PersistentGrantLifetimeUnit)
	state.PersistentGrantIdleTimeout = types.Int64PointerValue(r.PersistentGrantIdleTimeout)
	state.PersistentGrantIdleTimeoutTimeUnit = types.StringPointerValue(r.PersistentGrantIdleTimeoutTimeUnit)
	state.RefreshTokenLength = types.Int64Value(r.RefreshTokenLength)
	state.RollRefreshTokenValues = types.BoolPointerValue(r.RollRefreshTokenValues)
	state.RefreshTokenRollingGracePeriod = types.Int64PointerValue(r.RefreshTokenRollingGracePeriod)
	state.RefreshRollingInterval = types.Int64Value(r.RefreshRollingInterval)
	state.RefreshRollingIntervalTimeUnit = types.StringPointerValue(r.RefreshRollingIntervalTimeUnit)
	state.PersistentGrantReuseGrantTypes = internaltypes.GetStringSet(r.PersistentGrantReuseGrantTypes)
	state.BypassAuthorizationForApprovedGrants = types.BoolPointerValue(r.BypassAuthorizationForApprovedGrants)
	state.AllowUnidentifiedClientROCreds = types.BoolPointerValue(r.AllowUnidentifiedClientROCreds)
	state.AllowUnidentifiedClientExtensionGrants = types.BoolPointerValue(r.AllowUnidentifiedClientExtensionGrants)
	state.AdminWebServicePcvRef, respDiags = resourcelink.ToState(ctx, r.AdminWebServicePcvRef)
	diags.Append(respDiags...)
	state.AtmIdForOAuthGrantManagement = types.StringPointerValue(r.AtmIdForOAuthGrantManagement)
	state.ScopeForOAuthGrantManagement = types.StringPointerValue(r.ScopeForOAuthGrantManagement)
	state.AllowedOrigins = internaltypes.GetStringSet(r.AllowedOrigins)
	state.UserAuthorizationUrl = types.StringPointerValue(r.UserAuthorizationUrl)
	state.RegisteredAuthorizationPath = types.StringPointerValue(r.RegisteredAuthorizationPath)
	state.PendingAuthorizationTimeout = types.Int64PointerValue(r.PendingAuthorizationTimeout)
	state.DevicePollingInterval = types.Int64PointerValue(r.DevicePollingInterval)
	state.ActivationCodeCheckMode = types.StringPointerValue(r.ActivationCodeCheckMode)
	state.BypassActivationCodeConfirmation = types.BoolPointerValue(r.BypassActivationCodeConfirmation)
	state.EnableCookielessUserAuthorizationAuthenticationApi = types.BoolPointerValue(r.EnableCookielessUserAuthorizationAuthenticationApi)
	state.UserAuthorizationConsentPageSetting = types.StringPointerValue(r.UserAuthorizationConsentPageSetting)
	state.UserAuthorizationConsentAdapter = types.StringPointerValue(r.UserAuthorizationConsentAdapter)
	state.ApprovedScopesAttribute = types.StringPointerValue(r.ApprovedScopesAttribute)
	state.ApprovedAuthorizationDetailAttribute = types.StringPointerValue(r.ApprovedAuthorizationDetailAttribute)
	state.ParReferenceTimeout = types.Int64PointerValue(r.ParReferenceTimeout)
	state.ParReferenceLength = types.Int64PointerValue(r.ParReferenceLength)
	state.ParStatus = types.StringPointerValue(r.ParStatus)
	state.ClientSecretRetentionPeriod = types.Int64PointerValue(r.ClientSecretRetentionPeriod)
	state.JwtSecuredAuthorizationResponseModeLifetime = types.Int64PointerValue(r.JwtSecuredAuthorizationResponseModeLifetime)
	state.DpopProofRequireNonce = types.BoolPointerValue(r.DpopProofRequireNonce)
	state.DpopProofLifetimeSeconds = types.Int64PointerValue(r.DpopProofLifetimeSeconds)
	state.DpopProofEnforceReplayPrevention = types.BoolPointerValue(r.DpopProofEnforceReplayPrevention)
	state.BypassAuthorizationForApprovedConsents = types.BoolPointerValue(r.BypassAuthorizationForApprovedConsents)
	state.ConsentLifetimeDays = types.Int64PointerValue(r.ConsentLifetimeDays)
	state.ReturnIdTokenOnOpenIdWithDeviceAuthzGrant = types.BoolPointerValue(r.ReturnIdTokenOnOpenIdWithDeviceAuthzGrant)
	return diags
}
