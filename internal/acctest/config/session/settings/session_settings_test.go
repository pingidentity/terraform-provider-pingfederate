package sessionsettings_test

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
type sessionSettingsResourceModel struct {
	trackAdapterSessionsForLogout bool
	revokeUserSessionOnLogout     bool
	sessionRevocationLifetime     int64
}

func TestAccSessionSettings(t *testing.T) {
	resourceName := "mySessionSettings"
	updatedResourceModel := sessionSettingsResourceModel{
		trackAdapterSessionsForLogout: true,
		revokeUserSessionOnLogout:     false,
		sessionRevocationLifetime:     60,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSessionSettings(resourceName, nil),
				Check:  testAccCheckExpectedSessionSettingsAttributes(nil),
			},
			{
				// Test updating some fields
				Config: testAccSessionSettings(resourceName, &updatedResourceModel),
				Check:  testAccCheckExpectedSessionSettingsAttributes(&updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSessionSettings(resourceName, &updatedResourceModel),
				ResourceName:      "pingfederate_session_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccSessionSettings(resourceName, nil),
				Check:  testAccCheckExpectedSessionSettingsAttributes(nil),
			},
		},
	})
}

func testAccSessionSettings(resourceName string, resourceModel *sessionSettingsResourceModel) string {
	optionalHcl := ""
	if resourceModel != nil {
		optionalHcl = fmt.Sprintf(`
		  track_adapter_sessions_for_logout = %t
		  revoke_user_session_on_logout     = %t
		  session_revocation_lifetime       = %d
		`, resourceModel.trackAdapterSessionsForLogout,
			resourceModel.revokeUserSessionOnLogout,
			resourceModel.sessionRevocationLifetime,
		)
	}

	return fmt.Sprintf(`
resource "pingfederate_session_settings" "%s" {
  %s
}

data "pingfederate_server_settings_general_settings" "myServerSettings" {
  depends_on = [pingfederate_session_settings.%[1]s]
}`, resourceName,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSessionSettingsAttributes(config *sessionSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SessionSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		stateAttributes := s.RootModule().Resources["pingfederate_session_settings.mySessionSettings"].Primary.Attributes
		response, _, err := testClient.SessionAPI.GetSessionSettings(ctx).Execute()

		if err != nil {
			return err
		}

		if config == nil {
			return nil
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "track_adapter_sessions_for_logout",
			config.trackAdapterSessionsForLogout, *response.TrackAdapterSessionsForLogout)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "track_adapter_sessions_for_logout", config.trackAdapterSessionsForLogout)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "revoke_user_session_on_logout",
			config.revokeUserSessionOnLogout, *response.RevokeUserSessionOnLogout)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "revoke_user_session_on_logout", config.revokeUserSessionOnLogout)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "session_revocation_lifetime",
			config.sessionRevocationLifetime, *response.SessionRevocationLifetime)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "session_revocation_lifetime", config.sessionRevocationLifetime)
		if err != nil {
			return err
		}
		return nil
	}
}
