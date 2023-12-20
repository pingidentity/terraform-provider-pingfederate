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
type authenticationPoliciesSettingsResourceModel struct {
	enableIdpAuthnSelection bool
	enableSpAuthnSelection  bool
}

func TestAccAuthenticationPoliciesSettings(t *testing.T) {
	resourceName := "myAuthenticationPoliciesSettings"
	initialResourceModel := authenticationPoliciesSettingsResourceModel{
		enableIdpAuthnSelection: false,
		enableSpAuthnSelection:  true,
	}
	updatedResourceModel := authenticationPoliciesSettingsResourceModel{
		enableIdpAuthnSelection: true,
		enableSpAuthnSelection:  false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPoliciesSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationPoliciesSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAuthenticationPoliciesSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAuthenticationPoliciesSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationPoliciesSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_authentication_policies_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAuthenticationPoliciesSettings(resourceName string, resourceModel authenticationPoliciesSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_settings" "%[1]s" {
  enable_idp_authn_selection = %[2]t
  enable_sp_authn_selection  = %[3]t
}

data "pingfederate_authentication_policies_settings" "%[1]s" {
  depends_on = [pingfederate_authentication_policies_settings.%[1]s]
}`, resourceName,
		resourceModel.enableIdpAuthnSelection,
		resourceModel.enableSpAuthnSelection,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPoliciesSettingsAttributes(config authenticationPoliciesSettingsResourceModel) resource.TestCheckFunc {
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
			config.enableIdpAuthnSelection, *response.EnableIdpAuthnSelection)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_sp_authn_selection",
			config.enableSpAuthnSelection, *response.EnableSpAuthnSelection)
		if err != nil {
			return err
		}

		return nil
	}
}
