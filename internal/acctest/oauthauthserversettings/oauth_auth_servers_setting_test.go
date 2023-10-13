package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
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
		defaultScopeDescription:          "example scope description",
		authorizationCodeTimeout:         50,
		authorizationCodeEntropy:         20,
		refreshTokenLength:               40,
		refreshRollingInterval:           1,
		registeredAuthorizationPath:      "/example",
		pendingAuthorizationTimeout:      550,
		devicePollingInterval:            4,
		bypassActivationCodeConfirmation: false,
	}
	updatedResourceModel := oauthAuthServerSettingsResourceModel{
		defaultScopeDescription:          "example updated scope description",
		authorizationCodeTimeout:         60,
		authorizationCodeEntropy:         30,
		refreshTokenLength:               50,
		refreshRollingInterval:           2,
		registeredAuthorizationPath:      "/updatedexample",
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
				Config: testAccOauthAuthServerSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthAuthServerSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthAuthServerSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_auth_server_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOauthAuthServerSettings(resourceName string, resourceModel oauthAuthServerSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_auth_server_settings" "%[1]s" {
  authorization_code_entropy          = %[2]d
  authorization_code_timeout          = %[3]d
  registered_authorization_path       = "%[4]s"
  default_scope_description           = "%[5]s"
  device_polling_interval             = %[6]d
  pending_authorization_timeout       = %[7]d
  refresh_rolling_interval            = %[8]d
  refresh_token_length                = %[9]d
  bypass_activation_code_confirmation = %[10]t
}`, resourceName,
		resourceModel.authorizationCodeEntropy,
		resourceModel.authorizationCodeTimeout,
		resourceModel.registeredAuthorizationPath,
		resourceModel.defaultScopeDescription,
		resourceModel.devicePollingInterval,
		resourceModel.pendingAuthorizationTimeout,
		resourceModel.refreshRollingInterval,
		resourceModel.refreshTokenLength,
		resourceModel.bypassActivationCodeConfirmation,
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
			config.registeredAuthorizationPath, response.RegisteredAuthorizationPath)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, nil, "default_scope_description",
			config.defaultScopeDescription, response.DefaultScopeDescription)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "device_polling_interval",
			config.devicePollingInterval, response.DevicePollingInterval)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, nil, "pending_authorization_timeout",
			config.pendingAuthorizationTimeout, response.PendingAuthorizationTimeout)
		if err != nil {
			return err
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
			config.bypassActivationCodeConfirmation, response.BypassActivationCodeConfirmation)
		if err != nil {
			return err
		}

		return nil
	}
}
