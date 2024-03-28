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
	defaultScopeDescription          string
	authorizationCodeTimeout         int64
	authorizationCodeEntropy         int64
	refreshTokenLength               int64
	refreshRollingInterval           int64
	registeredAuthorizationPath      string
	pendingAuthorizationTimeout      int64
	devicePollingInterval            int64
	bypassActivationCodeConfirmation bool
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
		defaultScopeDescription:          "example updated scope description",
		authorizationCodeTimeout:         60,
		authorizationCodeEntropy:         30,
		refreshTokenLength:               50,
		refreshRollingInterval:           2,
		registeredAuthorizationPath:      "/example",
		pendingAuthorizationTimeout:      650,
		devicePollingInterval:            3,
		bypassActivationCodeConfirmation: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOauthAuthServerSettings(resourceName, initialResourceModel, false),
				Check:  testAccCheckExpectedOauthAuthServerSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthAuthServerSettings(resourceName, updatedResourceModel, true),
				Check:  testAccCheckExpectedOauthAuthServerSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthAuthServerSettings(resourceName, updatedResourceModel, true),
				ResourceName:      "pingfederate_oauth_auth_server_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccOauthAuthServerSettings(resourceName, initialResourceModel, false),
				Check:  testAccCheckExpectedOauthAuthServerSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccOauthAuthServerSettings(resourceName string, resourceModel oauthAuthServerSettingsResourceModel, includeAllAttributes bool) string {
	addUpdatedResourceModelFields := []string{}
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
  disallow_plain_pkce                      = false
  include_issuer_in_authorization_response = false
  persistent_grant_lifetime                = -1
  persistent_grant_lifetime_unit           = "DAYS"
  persistent_grant_idle_timeout            = 30
  persistent_grant_idle_timeout_time_unit  = "DAYS"
  roll_refresh_token_values                = true
  refresh_token_rolling_grace_period       = 0
  persistent_grant_reuse_grant_types       = []
  persistent_grant_contract = {
    extended_attributes = [
      {
        name = "example_extended_attribute"
      }
    ]
  }
  bypass_authorization_for_approved_grants         = false
  allow_unidentified_client_ro_creds               = false
  allow_unidentified_client_extension_grants       = false
  token_endpoint_base_url                          = ""
  user_authorization_url                           = ""
  activation_code_check_mode                       = "BEFORE_AUTHENTICATION"
  user_authorization_consent_page_setting          = "INTERNAL"
  atm_id_for_oauth_grant_management                = ""
  scope_for_oauth_grant_management                 = ""
  allowed_origins                                  = []
  track_user_sessions_for_logout                   = false
  par_reference_timeout                            = 60
  par_reference_length                             = 24
  par_status                                       = "ENABLED"
  client_secret_retention_period                   = 0
  jwt_secured_authorization_response_mode_lifetime = 600
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
	}

	return fmt.Sprintf(`
resource "pingfederate_oauth_auth_server_settings" "%[1]s" {
  authorization_code_entropy = %[2]d
  authorization_code_timeout = %[3]d
  refresh_rolling_interval   = %[4]d
  refresh_token_length       = %[5]d
  %[6]s
	%[7]s
}
data "pingfederate_oauth_auth_server_settings" "%[1]s" {
  depends_on = [
    pingfederate_oauth_auth_server_settings.%[1]s
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

		return nil
	}
}
