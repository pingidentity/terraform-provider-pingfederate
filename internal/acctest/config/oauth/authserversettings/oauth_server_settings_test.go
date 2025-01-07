package oauthauthserversettings_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// Attributes to test with. Add optional properties to test here if desired.
type oauthAuthServerSettingsResourceModel struct {
	activationCodeCheckMode                     string
	allowedOrigins                              []string
	allowUnidentifiedClientExtensionGrants      bool
	allowUnidentifiedClientRoCreds              bool
	atmIdForOauthGrantManagement                string
	authorizationCodeTimeout                    int64
	authorizationCodeEntropy                    int64
	bypassActivationCodeConfirmation            bool
	bypassAuthorizationForApprovedGrants        bool
	clientSecretRetentionPeriod                 int64
	defaultScopeDescription                     string
	devicePollingInterval                       int64
	disallowPlainPkce                           bool
	includeIssuerInAuthorizationResponse        bool
	jwtSecuredAuthorizationResponseModeLifetime int64
	parStatus                                   string
	pendingAuthorizationTimeout                 int64
	persistentGrantIdleTimeout                  int64
	persistentGrantLifetime                     int64
	persistentGrantLifetimeUnit                 string
	persistentGrantReuseGrantTypes              []string
	refreshRollingInterval                      int64
	refreshTokenLength                          int64
	registeredAuthorizationPath                 string
	rollRefreshTokenValues                      bool
	refreshTokenRollingGracePeriod              int64
	scopeForOauthGrantManagement                string
	tokenEndpointBaseUrl                        string
	trackUserSessionsForLogout                  bool
	userAuthorizationConsentPageSetting         string
	userAuthorizationUrl                        string
}

func TestAccOauthAuthServerSettings(t *testing.T) {
	resourceName := "myOauthAuthServerSettings"
	initialResourceModel := oauthAuthServerSettingsResourceModel{
		authorizationCodeTimeout: 50,
		authorizationCodeEntropy: 20,
		refreshTokenLength:       40,
		refreshRollingInterval:   1,
	}

	updatedResourceModel := oauthAuthServerSettingsResourceModel{
		activationCodeCheckMode:                     "BEFORE_AUTHENTICATION",
		allowedOrigins:                              []string{"https://example.com:*"},
		allowUnidentifiedClientExtensionGrants:      true,
		allowUnidentifiedClientRoCreds:              true,
		atmIdForOauthGrantManagement:                "jwt",
		authorizationCodeTimeout:                    60,
		authorizationCodeEntropy:                    30,
		bypassActivationCodeConfirmation:            true,
		bypassAuthorizationForApprovedGrants:        true,
		clientSecretRetentionPeriod:                 60,
		disallowPlainPkce:                           true,
		defaultScopeDescription:                     "example updated scope description",
		devicePollingInterval:                       3,
		includeIssuerInAuthorizationResponse:        true,
		jwtSecuredAuthorizationResponseModeLifetime: 60,
		refreshTokenLength:                          50,
		refreshRollingInterval:                      2,
		registeredAuthorizationPath:                 "/example",
		parStatus:                                   "ENABLED",
		pendingAuthorizationTimeout:                 650,
		persistentGrantIdleTimeout:                  60,
		persistentGrantLifetime:                     -1,
		persistentGrantLifetimeUnit:                 "DAYS",
		persistentGrantReuseGrantTypes:              []string{"AUTHORIZATION_CODE"},
		rollRefreshTokenValues:                      true,
		refreshTokenRollingGracePeriod:              0,
		scopeForOauthGrantManagement:                "examplescope",
		tokenEndpointBaseUrl:                        "https://example.com",
		trackUserSessionsForLogout:                  true,
		userAuthorizationConsentPageSetting:         "INTERNAL",
		userAuthorizationUrl:                        "https://example.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOauthAuthServerSettings(resourceName, initialResourceModel, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOauthAuthServerSettingsAttributes(initialResourceModel),
					checkPf121ComputedAttrs(resourceName),
				),
			},
			{
				// Test updating some fields
				Config: testAccOauthAuthServerSettings(resourceName, updatedResourceModel, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOauthAuthServerSettingsAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "activation_code_check_mode", updatedResourceModel.activationCodeCheckMode),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "allow_unidentified_client_extension_grants", fmt.Sprintf("%t", updatedResourceModel.allowUnidentifiedClientExtensionGrants)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "allow_unidentified_client_ro_creds", fmt.Sprintf("%t", updatedResourceModel.allowUnidentifiedClientRoCreds)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "allowed_origins.0", updatedResourceModel.allowedOrigins[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "atm_id_for_oauth_grant_management", updatedResourceModel.atmIdForOauthGrantManagement),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "bypass_activation_code_confirmation", fmt.Sprintf("%t", updatedResourceModel.bypassActivationCodeConfirmation)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "client_secret_retention_period", fmt.Sprintf("%d", updatedResourceModel.clientSecretRetentionPeriod)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "default_scope_description", updatedResourceModel.defaultScopeDescription),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "device_polling_interval", fmt.Sprintf("%d", updatedResourceModel.devicePollingInterval)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "disallow_plain_pkce", fmt.Sprintf("%t", updatedResourceModel.disallowPlainPkce)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "include_issuer_in_authorization_response", fmt.Sprintf("%t", updatedResourceModel.includeIssuerInAuthorizationResponse)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "jwt_secured_authorization_response_mode_lifetime", fmt.Sprintf("%d", updatedResourceModel.jwtSecuredAuthorizationResponseModeLifetime)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "par_status", updatedResourceModel.parStatus),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "persistent_grant_reuse_grant_types.0", updatedResourceModel.persistentGrantReuseGrantTypes[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "pending_authorization_timeout", fmt.Sprintf("%d", updatedResourceModel.pendingAuthorizationTimeout)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "persistent_grant_lifetime", fmt.Sprintf("%d", updatedResourceModel.persistentGrantLifetime)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "persistent_grant_lifetime_unit", updatedResourceModel.persistentGrantLifetimeUnit),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "persistent_grant_idle_timeout", fmt.Sprintf("%d", updatedResourceModel.persistentGrantIdleTimeout)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "refresh_token_rolling_grace_period", "60"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "registered_authorization_path", updatedResourceModel.registeredAuthorizationPath),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "roll_refresh_token_values", fmt.Sprintf("%t", updatedResourceModel.rollRefreshTokenValues)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "scope_for_oauth_grant_management", updatedResourceModel.scopeForOauthGrantManagement),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "token_endpoint_base_url", updatedResourceModel.tokenEndpointBaseUrl),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "track_user_sessions_for_logout", fmt.Sprintf("%t", updatedResourceModel.trackUserSessionsForLogout)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "user_authorization_consent_page_setting", updatedResourceModel.userAuthorizationConsentPageSetting),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "user_authorization_url", updatedResourceModel.userAuthorizationUrl),
				),
			},
			{
				// Test importing the resource
				Config:                               testAccOauthAuthServerSettings(resourceName, updatedResourceModel, true),
				ResourceName:                         "pingfederate_oauth_server_settings." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "refresh_token_length",
			},
			{
				// Back to minimal model
				Config: testAccOauthAuthServerSettings(resourceName, initialResourceModel, false),
				Check:  testAccCheckExpectedOauthAuthServerSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func checkPf121ComputedAttrs(resourceName string) resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1210) {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "require_offline_access_scope_to_issue_refresh_tokens", "false"),
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "offline_access_require_consent_prompt", "false"),
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "refresh_rolling_interval_time_unit", "HOURS"),
			resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "enable_cookieless_user_authorization_authentication_api", "false"),
		)
	}
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "require_offline_access_scope_to_issue_refresh_tokens"),
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "offline_access_require_consent_prompt"),
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "refresh_rolling_interval_time_unit"),
		resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_oauth_server_settings.%s", resourceName), "enable_cookieless_user_authorization_authentication_api"),
	)
}

func testAccOauthAuthServerSettings(resourceName string, resourceModel oauthAuthServerSettingsResourceModel, includeAllAttributes bool) string {
	addUpdatedResourceModelFields := []string{}

	allowedOrigins := fmt.Sprintf("allowed_origins = [%s]", acctest.StringSliceToString(resourceModel.allowedOrigins))
	addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, allowedOrigins)

	persistentGrantReuseGrantTypes := fmt.Sprintf("persistent_grant_reuse_grant_types = [%s]", acctest.StringSliceToString(resourceModel.persistentGrantReuseGrantTypes))
	addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, persistentGrantReuseGrantTypes)

	if resourceModel.bypassActivationCodeConfirmation == true {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("bypass_activation_code_confirmation = %t", resourceModel.bypassActivationCodeConfirmation))
	}
	if resourceModel.defaultScopeDescription != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("default_scope_description = \"%s\"", resourceModel.defaultScopeDescription))
	}
	if resourceModel.devicePollingInterval == 3 {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("device_polling_interval = %d", resourceModel.devicePollingInterval))
	}
	if resourceModel.pendingAuthorizationTimeout == 650 {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("pending_authorization_timeout = %d", resourceModel.pendingAuthorizationTimeout))
	}
	if resourceModel.registeredAuthorizationPath != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("registered_authorization_path = \"%s\"", resourceModel.registeredAuthorizationPath))
	}

	if resourceModel.activationCodeCheckMode != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("activation_code_check_mode = \"%s\"", resourceModel.activationCodeCheckMode))
	}

	if resourceModel.allowUnidentifiedClientExtensionGrants {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("allow_unidentified_client_extension_grants = %t", resourceModel.allowUnidentifiedClientExtensionGrants))
	}

	if resourceModel.allowUnidentifiedClientRoCreds {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("allow_unidentified_client_ro_creds = %t", resourceModel.allowUnidentifiedClientRoCreds))
	}

	if resourceModel.atmIdForOauthGrantManagement != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("atm_id_for_oauth_grant_management = \"%s\"", resourceModel.atmIdForOauthGrantManagement))
	}

	if resourceModel.persistentGrantLifetime != 0 {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("persistent_grant_lifetime = %d", resourceModel.persistentGrantLifetime))
	}

	if resourceModel.persistentGrantLifetimeUnit != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("persistent_grant_lifetime_unit = \"%s\"", resourceModel.persistentGrantLifetimeUnit))
	}

	if resourceModel.persistentGrantIdleTimeout != 0 {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("persistent_grant_idle_timeout = %d", resourceModel.persistentGrantIdleTimeout))
	}

	if resourceModel.rollRefreshTokenValues {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("roll_refresh_token_values = %t", resourceModel.rollRefreshTokenValues))
	}

	if resourceModel.refreshTokenRollingGracePeriod != 0 {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("refresh_token_rolling_grace_period = %d", resourceModel.refreshTokenRollingGracePeriod))
	}

	if resourceModel.scopeForOauthGrantManagement != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("scope_for_oauth_grant_management = \"%s\"", resourceModel.scopeForOauthGrantManagement))
	}

	if resourceModel.tokenEndpointBaseUrl != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("token_endpoint_base_url = \"%s\"", resourceModel.tokenEndpointBaseUrl))
	}

	if resourceModel.trackUserSessionsForLogout {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("track_user_sessions_for_logout = %t", resourceModel.trackUserSessionsForLogout))
	}

	if resourceModel.userAuthorizationConsentPageSetting != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("user_authorization_consent_page_setting = \"%s\"", resourceModel.userAuthorizationConsentPageSetting))
	}

	if resourceModel.userAuthorizationUrl != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("user_authorization_url = \"%s\"", resourceModel.userAuthorizationUrl))
	}

	if resourceModel.parStatus != "" {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("par_status = \"%s\"", resourceModel.parStatus))
	}

	if resourceModel.disallowPlainPkce {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("disallow_plain_pkce = %t", resourceModel.disallowPlainPkce))
	}

	if resourceModel.includeIssuerInAuthorizationResponse {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("include_issuer_in_authorization_response = %t", resourceModel.includeIssuerInAuthorizationResponse))
	}

	if resourceModel.jwtSecuredAuthorizationResponseModeLifetime != 0 {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("jwt_secured_authorization_response_mode_lifetime = %d", resourceModel.jwtSecuredAuthorizationResponseModeLifetime))
	}

	if resourceModel.clientSecretRetentionPeriod != 0 {
		addUpdatedResourceModelFields = append(addUpdatedResourceModelFields, fmt.Sprintf("client_secret_retention_period = %d", resourceModel.clientSecretRetentionPeriod))
	}

	updatedResourceModelFields := strings.Join(addUpdatedResourceModelFields[:], "\n")

	optionalHcl := ""
	if includeAllAttributes {
		optionalHcl = `
  scopes = [
    {
      name        = "examplescope",
      description = "example scope",
      dynamic     = false
    }
  ]
  scope_groups = [
    {
      name        = "examplescopegroup",
      description = "example scope group"
      scopes      = ["examplescope"]
    }
  ]
  exclusive_scopes = [
    {
      name        = "exampleexclusivescope",
      description = "example scope",
      dynamic     = false
    }
  ]
  exclusive_scope_groups = [
    {
      name        = "exampleexclusivescopegroup",
      description = "example exclusive scope group"
      scopes      = ["exampleexclusivescope"]
    }
  ]
  persistent_grant_contract = {
    extended_attributes = [
      {
        name = "example_extended_attribute"
      }
    ]
  }
  par_reference_timeout                            = 60
  par_reference_length                             = 24
		`
		if acctest.VersionAtLeast(version.PingFederate1130) {
			optionalHcl += `
  dpop_proof_require_nonce = true
  dpop_proof_lifetime_seconds = 60
  dpop_proof_enforce_replay_prevention = false
			`
		}

		if acctest.VersionAtLeast(version.PingFederate1200) {
			optionalHcl += `
  bypass_authorization_for_approved_consents = true
  consent_lifetime_days = 5
			`
		}

		if acctest.VersionAtLeast(version.PingFederate1210) {
			optionalHcl += `
  require_offline_access_scope_to_issue_refresh_tokens = true
  offline_access_require_consent_prompt = true
  refresh_rolling_interval_time_unit = "MINUTES"
  enable_cookieless_user_authorization_authentication_api = true
			`
		}
	}

	return fmt.Sprintf(`
resource "pingfederate_oauth_server_settings" "%[1]s" {
  authorization_code_entropy = %[2]d
  authorization_code_timeout = %[3]d
  refresh_rolling_interval   = %[4]d
  refresh_token_length       = %[5]d
  %[6]s
	%[7]s
}
data "pingfederate_oauth_server_settings" "%[1]s" {
  depends_on = [
    pingfederate_oauth_server_settings.%[1]s
  ]
}`, resourceName,
		resourceModel.authorizationCodeEntropy,
		resourceModel.authorizationCodeTimeout,
		resourceModel.refreshRollingInterval,
		resourceModel.refreshTokenLength,
		optionalHcl,
		updatedResourceModelFields,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthAuthServerSettingsAttributes(config oauthAuthServerSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAuthServerSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAuthServerSettingsAPI.GetAuthorizationServerSettings(ctx).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchStringSlice(resourceType, nil, "allowed_origins",
			config.allowedOrigins, response.AllowedOrigins)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "authorization_code_entropy",
			config.authorizationCodeEntropy, response.AuthorizationCodeEntropy)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "authorization_code_timeout",
			config.authorizationCodeTimeout, response.AuthorizationCodeTimeout)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "registered_authorization_path",
			config.registeredAuthorizationPath, *response.RegisteredAuthorizationPath)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "default_scope_description",
			config.defaultScopeDescription, *response.DefaultScopeDescription)
		if err != nil {
			return err
		}

		if config.devicePollingInterval != 0 {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "device_polling_interval",
				config.devicePollingInterval, *response.DevicePollingInterval)
			if err != nil {
				return err
			}
		}

		if config.pendingAuthorizationTimeout != 0 {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "pending_authorization_timeout",
				config.pendingAuthorizationTimeout, *response.PendingAuthorizationTimeout)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "refresh_rolling_interval",
			config.refreshRollingInterval, response.RefreshRollingInterval)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "refresh_token_length",
			config.refreshTokenLength, response.RefreshTokenLength)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "bypass_activation_code_confirmation",
			config.bypassActivationCodeConfirmation, *response.BypassActivationCodeConfirmation)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "roll_refresh_token_values",
			config.rollRefreshTokenValues, *response.RollRefreshTokenValues)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "disallow_plain_pkce",
			config.disallowPlainPkce, *response.DisallowPlainPKCE)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "include_issuer_in_authorization_response",
			config.includeIssuerInAuthorizationResponse, *response.IncludeIssuerInAuthorizationResponse)
		if err != nil {
			return err
		}

		if config.jwtSecuredAuthorizationResponseModeLifetime != 0 {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "jwt_secured_authorization_response_mode_lifetime",
				config.jwtSecuredAuthorizationResponseModeLifetime, *response.JwtSecuredAuthorizationResponseModeLifetime)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "atm_id_for_oauth_grant_management",
			config.atmIdForOauthGrantManagement, *response.AtmIdForOAuthGrantManagement)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "scope_for_oauth_grant_management",
			config.scopeForOauthGrantManagement, *response.ScopeForOAuthGrantManagement)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "token_endpoint_base_url",
			config.tokenEndpointBaseUrl, *response.TokenEndpointBaseUrl)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "track_user_sessions_for_logout",
			config.trackUserSessionsForLogout, *response.TrackUserSessionsForLogout)
		if err != nil {
			return err
		}

		if config.userAuthorizationConsentPageSetting == "" {
			err = acctest.TestAttributesMatchString(resourceType, nil, "user_authorization_consent_page_setting",
				"INTERNAL", *response.UserAuthorizationConsentPageSetting)
			if err != nil {
				return err
			}
		} else {
			err = acctest.TestAttributesMatchString(resourceType, nil, "user_authorization_consent_page_setting",
				config.userAuthorizationConsentPageSetting, *response.UserAuthorizationConsentPageSetting)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "user_authorization_url",
			config.userAuthorizationUrl, *response.UserAuthorizationUrl)
		if err != nil {
			return err
		}

		if config.parStatus == "" {
			err = acctest.TestAttributesMatchString(resourceType, nil, "par_status",
				"ENABLED", *response.ParStatus)
			if err != nil {
				return err
			}
		} else {
			err = acctest.TestAttributesMatchString(resourceType, nil, "par_status",
				config.parStatus, *response.ParStatus)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "allow_unidentified_client_extension_grants",
			config.allowUnidentifiedClientExtensionGrants, *response.AllowUnidentifiedClientExtensionGrants)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "allow_unidentified_client_ro_creds",
			config.allowUnidentifiedClientRoCreds, *response.AllowUnidentifiedClientROCreds)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "client_secret_retention_period",
			config.clientSecretRetentionPeriod, *response.ClientSecretRetentionPeriod)
		if err != nil {
			return err
		}

		if config.persistentGrantIdleTimeout != 0 {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "persistent_grant_idle_timeout",
				config.persistentGrantIdleTimeout, *response.PersistentGrantIdleTimeout)
			if err != nil {
				return err
			}
		}

		if config.persistentGrantLifetime != 0 {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "persistent_grant_lifetime",
				config.persistentGrantLifetime, *response.PersistentGrantLifetime)
			if err != nil {
				return err
			}
		}

		if config.persistentGrantLifetimeUnit == "" {
			err = acctest.TestAttributesMatchString(resourceType, nil, "persistent_grant_lifetime_unit",
				"DAYS", *response.PersistentGrantLifetimeUnit)
			if err != nil {
				return err
			}
		} else {
			err = acctest.TestAttributesMatchString(resourceType, nil, "persistent_grant_lifetime_unit",
				config.persistentGrantLifetimeUnit, *response.PersistentGrantLifetimeUnit)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchStringSlice(resourceType, nil, "persistent_grant_reuse_grant_types",
			config.persistentGrantReuseGrantTypes, response.PersistentGrantReuseGrantTypes)
		if err != nil {
			return err
		}

		if config.refreshTokenRollingGracePeriod != 0 {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "refresh_token_rolling_grace_period",
				config.refreshTokenRollingGracePeriod, *response.RefreshTokenRollingGracePeriod)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
