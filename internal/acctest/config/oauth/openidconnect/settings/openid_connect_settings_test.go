package oauthopenidconnectsettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

type openIdConnectSettingsResourceModel struct {
	sessionSettings *client.OIDCSessionSettings
}

func TestAccOpenIdConnectSettings(t *testing.T) {
	resourceName := "myOpenIdConnectSettings"
	// send empty model to start
	initialResourceModel := openIdConnectSettingsResourceModel{}
	updatedResourceModel := openIdConnectSettingsResourceModel{
		sessionSettings: &client.OIDCSessionSettings{
			TrackUserSessionsForLogout: pointers.Bool(true),
			RevokeUserSessionOnLogout:  pointers.Bool(true),
			SessionRevocationLifetime:  pointers.Int64(180),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOpenIdConnectSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOpenIdConnectSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOpenIdConnectSettings(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOpenIdConnectSettingsAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_settings.%s", resourceName), "session_settings.track_user_sessions_for_logout", fmt.Sprintf("%t", *updatedResourceModel.sessionSettings.TrackUserSessionsForLogout)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_settings.%s", resourceName), "session_settings.revoke_user_session_on_logout", fmt.Sprintf("%t", *updatedResourceModel.sessionSettings.RevokeUserSessionOnLogout)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_openid_connect_settings.%s", resourceName), "session_settings.session_revocation_lifetime", fmt.Sprintf("%d", *updatedResourceModel.sessionSettings.SessionRevocationLifetime)),
				),
			},
			{
				// Test importing the resource
				Config:            testAccOpenIdConnectSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_openid_connect_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOpenIdConnectSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOpenIdConnectSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func generateDependentHcl(resourceModel openIdConnectSettingsResourceModel) string {
	if resourceModel.sessionSettings == nil {
		return ""
	} else {
		return fmt.Sprintf(`
		session_settings = {
			track_user_sessions_for_logout = %t
			revoke_user_session_on_logout = %t
			session_revocation_lifetime = %d
		}`, *resourceModel.sessionSettings.TrackUserSessionsForLogout, *resourceModel.sessionSettings.RevokeUserSessionOnLogout, *resourceModel.sessionSettings.SessionRevocationLifetime)
	}
}

func testAccOpenIdConnectSettings(resourceName string, resourceModel openIdConnectSettingsResourceModel) string {
	// The dependent OIDC policy is not created in this test because prior to PF 12.1 it isn't possible to delete
	// the final OIDC policy from the server config, because it is always in use.
	return fmt.Sprintf(`
resource "pingfederate_openid_connect_settings" "%[1]s" {
  default_policy_ref = {
    id = "oidcSettingsTestPolicy"
  }
	%[2]s
}`, resourceName,
		generateDependentHcl(resourceModel),
	)
}

func testAccCheckExpectedOpenIdConnectSettingsAttributes(config openIdConnectSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OpenIdConnectSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthOpenIdConnectAPI.GetOIDCSettings(ctx).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "id", "oidcSettingsTestPolicy", response.DefaultPolicyRef.Id)
		if err != nil {
			return err
		}

		if config.sessionSettings != nil {
			err = acctest.TestAttributesMatchBool(resourceType, nil, "track_user_sessions_for_logout", *config.sessionSettings.TrackUserSessionsForLogout, *response.SessionSettings.TrackUserSessionsForLogout)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "revoke_user_session_on_logout", *config.sessionSettings.RevokeUserSessionOnLogout, *response.SessionSettings.RevokeUserSessionOnLogout)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchInt(resourceType, nil, "session_revocation_lifetime", *config.sessionSettings.SessionRevocationLifetime, *response.SessionSettings.SessionRevocationLifetime)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
