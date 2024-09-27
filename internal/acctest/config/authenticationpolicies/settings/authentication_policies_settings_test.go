package authenticationpoliciessettings_test

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

func TestAccAuthenticationPoliciesSettings(t *testing.T) {
	resourceName := "myAuthenticationPoliciesSettings"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test empty object (defaults)
			{
				Config: testAccAuthenticationPoliciesSettings(resourceName, false),
				Check:  testAccCheckExpectedAuthenticationPoliciesSettingsAttributes(false),
			},
			// Test updating all fields
			{
				Config: testAccAuthenticationPoliciesSettings(resourceName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAuthenticationPoliciesSettingsAttributes(true),
					resource.TestCheckResourceAttr("pingfederate_authentication_policies_settings.myAuthenticationPoliciesSettings", "enable_idp_authn_selection", "true"),
					resource.TestCheckResourceAttr("pingfederate_authentication_policies_settings.myAuthenticationPoliciesSettings", "enable_sp_authn_selection", "true"),
				),
			},
			{
				// Test importing the resource
				Config:                               testAccAuthenticationPoliciesSettings(resourceName, true),
				ResourceName:                         "pingfederate_authentication_policies_settings." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "enable_idp_authn_selection",
			},
			{
				// Back to minimal model
				Config: testAccAuthenticationPoliciesSettings(resourceName, false),
				Check:  testAccCheckExpectedAuthenticationPoliciesSettingsAttributes(false),
			},
		},
	})
}

func testAccAuthenticationPoliciesSettings(resourceName string, includeAttributes bool) string {
	optionalHcl := ""
	if includeAttributes {
		optionalHcl = `
		enable_idp_authn_selection = true
		enable_sp_authn_selection  = true
		`
	}

	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_settings" "%[1]s" {
		%[2]s
}

data "pingfederate_authentication_policies_settings" "%[1]s" {
  depends_on = [pingfederate_authentication_policies_settings.%[1]s]
}`, resourceName,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPoliciesSettingsAttributes(includeAttributes bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPoliciesSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationPoliciesAPI.GetAuthenticationPolicySettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_idp_authn_selection",
			includeAttributes, *response.EnableIdpAuthnSelection)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_sp_authn_selection",
			includeAttributes, *response.EnableSpAuthnSelection)
		if err != nil {
			return err
		}

		return nil
	}
}
