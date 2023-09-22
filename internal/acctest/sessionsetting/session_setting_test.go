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

const sessionSettingsId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type sessionSettingsResourceModel struct {
	trackAdapterSessionsForLogout bool
	revokeUserSessionOnLogout     bool
	sessionRevocationLifetime     int64
}

func TestAccSessionSettings(t *testing.T) {
	resourceName := "mySessionSettings"
	initialResourceModel := sessionSettingsResourceModel{
		trackAdapterSessionsForLogout: true,
		revokeUserSessionOnLogout:     true,
		sessionRevocationLifetime:     60,
	}
	updatedResourceModel := sessionSettingsResourceModel{
		trackAdapterSessionsForLogout: false,
		revokeUserSessionOnLogout:     true,
		sessionRevocationLifetime:     40,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSessionSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSessionSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSessionSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSessionSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSessionSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_session_setting." + resourceName,
				ImportStateId:     sessionSettingsId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccSessionSettings(resourceName string, resourceModel sessionSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_session_setting" "%[1]s" {
  track_adapter_sessions_for_logout = %[2]t
  revoke_user_session_on_logout     = %[3]t
  session_revocation_lifetime       = %[4]d
}`, resourceName,
		resourceModel.trackAdapterSessionsForLogout,
		resourceModel.revokeUserSessionOnLogout,
		resourceModel.sessionRevocationLifetime,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSessionSettingsAttributes(config sessionSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SessionSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SessionApi.GetSessionSettings(ctx).Execute()

		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "track_adapter_sessions_for_logout",
			config.trackAdapterSessionsForLogout, *response.TrackAdapterSessionsForLogout)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "revoke_user_session_on_logout",
			config.revokeUserSessionOnLogout, *response.RevokeUserSessionOnLogout)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "session_revocation_lifetime",
			config.sessionRevocationLifetime, *response.SessionRevocationLifetime)
		if err != nil {
			return err
		}
		return nil
	}
}
