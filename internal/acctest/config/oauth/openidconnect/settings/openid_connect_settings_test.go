package oauthopenidconnectsettings_test

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

func TestAccOpenIdConnectSettings(t *testing.T) {
	resourceName := "myOpenIdConnectSettings"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOpenIdConnectSettings(resourceName),
				Check:  testAccCheckExpectedOpenIdConnectSettingsAttributes(),
			},
			{
				// Test importing the resource
				Config:                               testAccOpenIdConnectSettings(resourceName),
				ResourceName:                         "pingfederate_openid_connect_settings." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "default_policy_ref.id",
			},
		},
	})
}

func testAccOpenIdConnectSettings(resourceName string) string {
	// The dependent OIDC policy is not created in this test because prior to PF 12.1 it isn't possible to delete
	// the final OIDC policy from the server config, because it is always in use.
	return fmt.Sprintf(`
resource "pingfederate_openid_connect_settings" "%s" {
  default_policy_ref = {
    id = "oidcSettingsTestPolicy"
  }
}`, resourceName,
	)
}

func testAccCheckExpectedOpenIdConnectSettingsAttributes() resource.TestCheckFunc {
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

		return nil
	}
}
