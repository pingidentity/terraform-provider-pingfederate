package oauthclient

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

var (
	jwksSettingsAttrType = map[string]attr.Type{
		"jwks_url": types.StringType,
		"jwks":     types.StringType,
	}

	oidcPolicyAttrType = map[string]attr.Type{
		"id_token_signing_algorithm":                  types.StringType,
		"id_token_encryption_algorithm":               types.StringType,
		"id_token_content_encryption_algorithm":       types.StringType,
		"policy_group":                                types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"grant_access_session_revocation_api":         types.BoolType,
		"grant_access_session_session_management_api": types.BoolType,
		"ping_access_logout_capable":                  types.BoolType,
		"logout_uris":                                 types.SetType{ElemType: types.StringType},
		"pairwise_identifier_user_type":               types.BoolType,
		"sector_identifier_uri":                       types.StringType,
		"logout_mode":                                 types.StringType,
		"back_channel_logout_uri":                     types.StringType,
		"post_logout_redirect_uris":                   types.SetType{ElemType: types.StringType},
	}

	oidcPolicyDefaultAttrValue = map[string]attr.Value{
		"id_token_signing_algorithm":                  types.StringNull(),
		"id_token_encryption_algorithm":               types.StringNull(),
		"id_token_content_encryption_algorithm":       types.StringNull(),
		"policy_group":                                types.ObjectNull(resourcelink.AttrType()),
		"grant_access_session_revocation_api":         types.BoolValue(false),
		"grant_access_session_session_management_api": types.BoolValue(false),
		"ping_access_logout_capable":                  types.BoolValue(false),
		"logout_uris":                                 types.SetNull(types.StringType),
		"pairwise_identifier_user_type":               types.BoolValue(false),
		"sector_identifier_uri":                       types.StringNull(),
		"logout_mode":                                 types.StringNull(),
		"back_channel_logout_uri":                     types.StringNull(),
		"post_logout_redirect_uris":                   types.SetNull(types.StringType),
	}

	secondarySecretsAttrType = map[string]attr.Type{
		"secret":      types.StringType,
		"expiry_time": types.StringType,
	}

	clientAuthAttrType = map[string]attr.Type{
		"type":                                  types.StringType,
		"secret":                                types.StringType,
		"secondary_secrets":                     types.SetType{ElemType: types.ObjectType{AttrTypes: secondarySecretsAttrType}},
		"client_cert_issuer_dn":                 types.StringType,
		"client_cert_subject_dn":                types.StringType,
		"enforce_replay_prevention":             types.BoolType,
		"token_endpoint_auth_signing_algorithm": types.StringType,
	}

	clientAuthDefaultAttrValue = map[string]attr.Value{
		"type":                                  types.StringValue("NONE"),
		"secret":                                types.StringNull(),
		"secondary_secrets":                     secondarySecretsEmptySet,
		"client_cert_issuer_dn":                 types.StringNull(),
		"client_cert_subject_dn":                types.StringNull(),
		"enforce_replay_prevention":             types.BoolNull(),
		"token_endpoint_auth_signing_algorithm": types.StringNull(),
	}

	extendedParametersAttrType = map[string]attr.Type{
		"values": types.SetType{ElemType: types.StringType},
	}
	extendedParametersObjAttrType = types.ObjectType{AttrTypes: extendedParametersAttrType}
)

type oauthClientModel struct {
	Id                                                            types.String `tfsdk:"id"`
	ClientId                                                      types.String `tfsdk:"client_id"`
	Enabled                                                       types.Bool   `tfsdk:"enabled"`
	RedirectUris                                                  types.Set    `tfsdk:"redirect_uris"`
	GrantTypes                                                    types.Set    `tfsdk:"grant_types"`
	Name                                                          types.String `tfsdk:"name"`
	Description                                                   types.String `tfsdk:"description"`
	ModificationDate                                              types.String `tfsdk:"modification_date"`
	CreationDate                                                  types.String `tfsdk:"creation_date"`
	LogoUrl                                                       types.String `tfsdk:"logo_url"`
	DefaultAccessTokenManagerRef                                  types.Object `tfsdk:"default_access_token_manager_ref"`
	RestrictToDefaultAccessTokenManager                           types.Bool   `tfsdk:"restrict_to_default_access_token_manager"`
	ValidateUsingAllEligibleAtms                                  types.Bool   `tfsdk:"validate_using_all_eligible_atms"`
	PersistentGrantExpirationType                                 types.String `tfsdk:"persistent_grant_expiration_type"`
	PersistentGrantExpirationTime                                 types.Int64  `tfsdk:"persistent_grant_expiration_time"`
	PersistentGrantExpirationTimeUnit                             types.String `tfsdk:"persistent_grant_expiration_time_unit"`
	PersistentGrantIdleTimeoutType                                types.String `tfsdk:"persistent_grant_idle_timeout_type"`
	PersistentGrantIdleTimeout                                    types.Int64  `tfsdk:"persistent_grant_idle_timeout"`
	PersistentGrantIdleTimeoutTimeUnit                            types.String `tfsdk:"persistent_grant_idle_timeout_time_unit"`
	PersistentGrantReuseType                                      types.String `tfsdk:"persistent_grant_reuse_type"`
	PersistentGrantReuseGrantTypes                                types.Set    `tfsdk:"persistent_grant_reuse_grant_types"`
	AllowAuthenticationApiInit                                    types.Bool   `tfsdk:"allow_authentication_api_init"`
	EnableCookielessAuthenticationApi                             types.Bool   `tfsdk:"enable_cookieless_authentication_api"`
	BypassApprovalPage                                            types.Bool   `tfsdk:"bypass_approval_page"`
	RestrictScopes                                                types.Bool   `tfsdk:"restrict_scopes"`
	RestrictedScopes                                              types.Set    `tfsdk:"restricted_scopes"`
	ExclusiveScopes                                               types.Set    `tfsdk:"exclusive_scopes"`
	AuthorizationDetailTypes                                      types.Set    `tfsdk:"authorization_detail_types"`
	RestrictedResponseTypes                                       types.Set    `tfsdk:"restricted_response_types"`
	RequirePushedAuthorizationRequests                            types.Bool   `tfsdk:"require_pushed_authorization_requests"`
	RequireJwtSecuredAuthorizationResponseMode                    types.Bool   `tfsdk:"require_jwt_secured_authorization_response_mode"`
	RequireSignedRequests                                         types.Bool   `tfsdk:"require_signed_requests"`
	RequestObjectSigningAlgorithm                                 types.String `tfsdk:"request_object_signing_algorithm"`
	OidcPolicy                                                    types.Object `tfsdk:"oidc_policy"`
	ClientAuth                                                    types.Object `tfsdk:"client_auth"`
	JwksSettings                                                  types.Object `tfsdk:"jwks_settings"`
	ExtendedParameters                                            types.Map    `tfsdk:"extended_parameters"`
	DeviceFlowSettingType                                         types.String `tfsdk:"device_flow_setting_type"`
	UserAuthorizationUrlOverride                                  types.String `tfsdk:"user_authorization_url_override"`
	PendingAuthorizationTimeoutOverride                           types.Int64  `tfsdk:"pending_authorization_timeout_override"`
	DevicePollingIntervalOverride                                 types.Int64  `tfsdk:"device_polling_interval_override"`
	BypassActivationCodeConfirmationOverride                      types.Bool   `tfsdk:"bypass_activation_code_confirmation_override"`
	RequireProofKeyForCodeExchange                                types.Bool   `tfsdk:"require_proof_key_for_code_exchange"`
	CibaDeliveryMode                                              types.String `tfsdk:"ciba_delivery_mode"`
	CibaNotificationEndpoint                                      types.String `tfsdk:"ciba_notification_endpoint"`
	CibaPollingInterval                                           types.Int64  `tfsdk:"ciba_polling_interval"`
	CibaRequireSignedRequests                                     types.Bool   `tfsdk:"ciba_require_signed_requests"`
	CibaRequestObjectSigningAlgorithm                             types.String `tfsdk:"ciba_request_object_signing_algorithm"`
	CibaUserCodeSupported                                         types.Bool   `tfsdk:"ciba_user_code_supported"`
	RequestPolicyRef                                              types.Object `tfsdk:"request_policy_ref"`
	TokenExchangeProcessorPolicyRef                               types.Object `tfsdk:"token_exchange_processor_policy_ref"`
	RefreshRolling                                                types.String `tfsdk:"refresh_rolling"`
	RefreshTokenRollingIntervalType                               types.String `tfsdk:"refresh_token_rolling_interval_type"`
	RefreshTokenRollingInterval                                   types.Int64  `tfsdk:"refresh_token_rolling_interval"`
	RefreshTokenRollingIntervalTimeUnit                           types.String `tfsdk:"refresh_token_rolling_interval_time_unit"`
	RefreshTokenRollingGracePeriodType                            types.String `tfsdk:"refresh_token_rolling_grace_period_type"`
	RefreshTokenRollingGracePeriod                                types.Int64  `tfsdk:"refresh_token_rolling_grace_period"`
	ClientSecretRetentionPeriodType                               types.String `tfsdk:"client_secret_retention_period_type"`
	ClientSecretRetentionPeriod                                   types.Int64  `tfsdk:"client_secret_retention_period"`
	ClientSecretChangedTime                                       types.String `tfsdk:"client_secret_changed_time"`
	TokenIntrospectionSigningAlgorithm                            types.String `tfsdk:"token_introspection_signing_algorithm"`
	TokenIntrospectionEncryptionAlgorithm                         types.String `tfsdk:"token_introspection_encryption_algorithm"`
	TokenIntrospectionContentEncryptionAlgorithm                  types.String `tfsdk:"token_introspection_content_encryption_algorithm"`
	JwtSecuredAuthorizationResponseModeSigningAlgorithm           types.String `tfsdk:"jwt_secured_authorization_response_mode_signing_algorithm"`
	JwtSecuredAuthorizationResponseModeEncryptionAlgorithm        types.String `tfsdk:"jwt_secured_authorization_response_mode_encryption_algorithm"`
	JwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm types.String `tfsdk:"jwt_secured_authorization_response_mode_content_encryption_algorithm"`
	RequireDpop                                                   types.Bool   `tfsdk:"require_dpop"`
	RequireOfflineAccessScopeToIssueRefreshTokens                 types.String `tfsdk:"require_offline_access_scope_to_issue_refresh_tokens"`
	OfflineAccessRequireConsentPrompt                             types.String `tfsdk:"offline_access_require_consent_prompt"`
}

func readOauthClientResponseCommon(ctx context.Context, r *client.Client, state, plan *oauthClientModel, productVersion version.SupportedVersion, isImportRead bool) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue(r.ClientId)
	state.ClientId = types.StringValue(r.ClientId)
	state.Enabled = types.BoolPointerValue(r.Enabled)
	state.RedirectUris = internaltypes.GetStringSet(r.RedirectUris)
	state.GrantTypes = internaltypes.GetStringSet(r.GrantTypes)
	state.Name = types.StringValue(r.Name)
	description := r.Description
	if description == nil {
		state.Description = types.StringNull()
	} else {
		state.Description = types.StringPointerValue(description)
	}
	state.ModificationDate = types.StringValue(r.ModificationDate.Format(time.RFC3339Nano))
	state.CreationDate = types.StringValue(r.CreationDate.Format(time.RFC3339Nano))
	state.LogoUrl = types.StringPointerValue(r.LogoUrl)
	state.DefaultAccessTokenManagerRef, respDiags = resourcelink.ToState(ctx, r.DefaultAccessTokenManagerRef)
	diags.Append(respDiags...)
	state.RestrictToDefaultAccessTokenManager = types.BoolPointerValue(r.RestrictToDefaultAccessTokenManager)
	state.ValidateUsingAllEligibleAtms = types.BoolPointerValue(r.ValidateUsingAllEligibleAtms)
	state.RefreshRolling = types.StringPointerValue(r.RefreshRolling)

	// If this is an import read and the refresh_token_rolling_interval_type is set to "SERVER_DEFAULT", then set the other rolling refresh token fields to null,
	// since they can only be configured when the type is set to "OVERRIDE_SERVER_DEFAULT"
	state.RefreshTokenRollingIntervalType = types.StringPointerValue(r.RefreshTokenRollingIntervalType)
	if isImportRead && state.RefreshTokenRollingIntervalType.ValueString() == "SERVER_DEFAULT" {
		state.RefreshTokenRollingInterval = types.Int64Null()
		state.RefreshTokenRollingIntervalTimeUnit = types.StringNull()
	} else {
		state.RefreshTokenRollingInterval = types.Int64PointerValue(r.RefreshTokenRollingInterval)
		// This attribute is returned as empty string when set to its default, and it only exists on PF 12.1+
		compare, err := version.Compare(productVersion, version.PingFederate1210)
		if err != nil {
			resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to compare PingFederate versions: "+err.Error())
		}
		pfVersionAtLeast121 := compare >= 0
		if r.GetRefreshTokenRollingIntervalTimeUnit() == "" && pfVersionAtLeast121 {
			state.RefreshTokenRollingIntervalTimeUnit = types.StringValue("HOURS")
		} else {
			state.RefreshTokenRollingIntervalTimeUnit = types.StringPointerValue(r.RefreshTokenRollingIntervalTimeUnit)
		}
	}

	// If this is an import read and the persistent_grant_expiration_type is set to "SERVER_DEFAULT", then set the other persistent grant fields to null,
	// since they can only be configured when the type is set to "OVERRIDE_SERVER_DEFAULT"
	state.PersistentGrantExpirationType = types.StringPointerValue(r.PersistentGrantExpirationType)
	if isImportRead && state.PersistentGrantExpirationType.ValueString() == "SERVER_DEFAULT" {
		state.PersistentGrantExpirationTime = types.Int64Null()
		state.PersistentGrantExpirationTimeUnit = types.StringNull()
	} else {
		state.PersistentGrantExpirationTime = types.Int64PointerValue(r.PersistentGrantExpirationTime)
		if r.GetPersistentGrantExpirationTimeUnit() == "" {
			state.PersistentGrantExpirationTimeUnit = types.StringValue("DAYS")
		} else {
			state.PersistentGrantExpirationTimeUnit = types.StringPointerValue(r.PersistentGrantExpirationTimeUnit)
		}
	}

	state.PersistentGrantIdleTimeoutType = types.StringPointerValue(r.PersistentGrantIdleTimeoutType)
	state.PersistentGrantIdleTimeout = types.Int64PointerValue(r.PersistentGrantIdleTimeout)
	state.PersistentGrantIdleTimeoutTimeUnit = types.StringPointerValue(r.PersistentGrantIdleTimeoutTimeUnit)
	state.PersistentGrantReuseType = types.StringPointerValue(r.PersistentGrantReuseType)
	state.PersistentGrantReuseGrantTypes = internaltypes.GetStringSet(r.PersistentGrantReuseGrantTypes)
	state.AllowAuthenticationApiInit = types.BoolPointerValue(r.AllowAuthenticationApiInit)
	state.EnableCookielessAuthenticationApi = types.BoolPointerValue(r.EnableCookielessAuthenticationApi)
	state.BypassApprovalPage = types.BoolPointerValue(r.BypassApprovalPage)
	state.RestrictScopes = types.BoolPointerValue(r.RestrictScopes)
	restrictedScopesToSet, respDiags := types.SetValueFrom(ctx, types.StringType, r.RestrictedScopes)
	diags.Append(respDiags...)
	state.RestrictedScopes = restrictedScopesToSet
	state.ExclusiveScopes = internaltypes.GetStringSet(r.ExclusiveScopes)
	state.AuthorizationDetailTypes = internaltypes.GetStringSet(r.AuthorizationDetailTypes)
	state.RestrictedResponseTypes = internaltypes.GetStringSet(r.RestrictedResponseTypes)
	state.RequirePushedAuthorizationRequests = types.BoolPointerValue(r.RequirePushedAuthorizationRequests)
	state.RequireJwtSecuredAuthorizationResponseMode = types.BoolPointerValue(r.RequireJwtSecuredAuthorizationResponseMode)
	state.RequireSignedRequests = types.BoolPointerValue(r.RequireSignedRequests)
	state.RequestObjectSigningAlgorithm = types.StringPointerValue(r.RequestObjectSigningAlgorithm)
	state.DeviceFlowSettingType = types.StringPointerValue(r.DeviceFlowSettingType)
	state.UserAuthorizationUrlOverride = types.StringPointerValue(r.UserAuthorizationUrlOverride)
	state.PendingAuthorizationTimeoutOverride = types.Int64PointerValue(r.PendingAuthorizationTimeoutOverride)
	state.DevicePollingIntervalOverride = types.Int64PointerValue(r.DevicePollingIntervalOverride)
	state.BypassActivationCodeConfirmationOverride = types.BoolPointerValue(r.BypassActivationCodeConfirmationOverride)
	state.RequireProofKeyForCodeExchange = types.BoolPointerValue(r.RequireProofKeyForCodeExchange)
	state.CibaDeliveryMode = types.StringPointerValue(r.CibaDeliveryMode)
	state.CibaNotificationEndpoint = types.StringPointerValue(r.CibaNotificationEndpoint)
	state.CibaPollingInterval = types.Int64PointerValue(r.CibaPollingInterval)
	state.CibaRequireSignedRequests = types.BoolPointerValue(r.CibaRequireSignedRequests)
	state.CibaRequestObjectSigningAlgorithm = types.StringPointerValue(r.CibaRequestObjectSigningAlgorithm)
	state.CibaUserCodeSupported = types.BoolPointerValue(r.CibaUserCodeSupported)
	state.RefreshTokenRollingGracePeriodType = types.StringPointerValue(r.RefreshTokenRollingGracePeriodType)
	state.RefreshTokenRollingGracePeriod = types.Int64PointerValue(r.RefreshTokenRollingGracePeriod)
	state.ClientSecretRetentionPeriodType = types.StringPointerValue(r.ClientSecretRetentionPeriodType)
	state.ClientSecretRetentionPeriod = types.Int64PointerValue(r.ClientSecretRetentionPeriod)
	state.ClientSecretChangedTime = types.StringValue(r.GetClientSecretChangedTime().Format(time.RFC3339Nano))
	state.TokenIntrospectionSigningAlgorithm = types.StringPointerValue(r.TokenIntrospectionSigningAlgorithm)
	state.TokenIntrospectionEncryptionAlgorithm = types.StringPointerValue(r.TokenIntrospectionEncryptionAlgorithm)
	state.TokenIntrospectionContentEncryptionAlgorithm = types.StringPointerValue(r.TokenIntrospectionContentEncryptionAlgorithm)
	state.JwtSecuredAuthorizationResponseModeSigningAlgorithm = types.StringPointerValue(r.JwtSecuredAuthorizationResponseModeSigningAlgorithm)
	state.JwtSecuredAuthorizationResponseModeEncryptionAlgorithm = types.StringPointerValue(r.JwtSecuredAuthorizationResponseModeEncryptionAlgorithm)
	state.JwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm = types.StringPointerValue(r.JwtSecuredAuthorizationResponseModeContentEncryptionAlgorithm)
	state.RequireDpop = types.BoolPointerValue(r.RequireDpop)
	state.RequireOfflineAccessScopeToIssueRefreshTokens = types.StringPointerValue(r.RequireOfflineAccessScopeToIssueRefreshTokens)
	state.OfflineAccessRequireConsentPrompt = types.StringPointerValue(r.OfflineAccessRequireConsentPrompt)

	// state.OidcPolicy
	oidcPolicyToState, respDiags := types.ObjectValueFrom(ctx, oidcPolicyAttrType, r.OidcPolicy)
	diags.Append(respDiags...)
	state.OidcPolicy = oidcPolicyToState

	// state.JwksSettings
	jwksSettingsToState, respDiags := types.ObjectValueFrom(ctx, jwksSettingsAttrType, r.JwksSettings)
	diags.Append(respDiags...)
	state.JwksSettings = jwksSettingsToState

	// state.ExtendedParameters
	var extendedParametersToState basetypes.MapValue
	if r.ExtendedParameters != nil {
		extendedParametersToState, respDiags = types.MapValueFrom(ctx, extendedParametersObjAttrType, r.ExtendedParameters)
		diags.Append(respDiags...)
	} else {
		if plan != nil && internaltypes.IsDefined(plan.ExtendedParameters) {
			extendedParametersToState = plan.ExtendedParameters
		} else {
			extendedParametersToState = types.MapNull(extendedParametersObjAttrType)
		}
	}

	state.ExtendedParameters = extendedParametersToState

	// state.RequestPolicyRef
	requestPolicyRefToState, respDiags := resourcelink.ToState(ctx, r.RequestPolicyRef)
	diags.Append(respDiags...)
	state.RequestPolicyRef = requestPolicyRefToState

	// state.TokenExchangeProcessorPolicyRef
	tokenExchangeProcessorPolicyRefToState, respDiags := resourcelink.ToState(ctx, r.TokenExchangeProcessorPolicyRef)
	diags.Append(respDiags...)
	state.TokenExchangeProcessorPolicyRef = tokenExchangeProcessorPolicyRefToState

	return diags
}
