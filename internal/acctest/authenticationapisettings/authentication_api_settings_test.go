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
type authenticationApiSettingsResourceModel struct {
	apiEnabled                       bool
	enableApiDescriptions            bool
	restrictAccessToRedirectlessMode bool
	includeRequestContext            bool
	defaultApplicationRef            string
}

func TestAccAuthenticationApiSettings(t *testing.T) {
	resourceName := "myAuthenticationApiSettings"
	// Use values that differ from the resource defaults
	updatedResourceModel := authenticationApiSettingsResourceModel{
		apiEnabled:                       true,
		enableApiDescriptions:            true,
		restrictAccessToRedirectlessMode: true,
		includeRequestContext:            true,
		defaultApplicationRef:            "myauthenticationapiapplication",
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationApiSettings(resourceName, nil),
			},
			{
				Config: testAccAuthenticationApiSettings(resourceName, &updatedResourceModel),
				Check:  testAccCheckExpectedAuthenticationApiSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationApiSettings(resourceName, &updatedResourceModel),
				ResourceName:      "pingfederate_authentication_api_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthenticationApiSettings(resourceName, nil),
			},
		},
	})
}

func testAccAuthenticationApiSettings(resourceName string, resourceModel *authenticationApiSettingsResourceModel) string {
	if resourceModel == nil {
		// Use resource defaults
		return fmt.Sprintf(`
resource "pingfederate_authentication_api_settings" "%[1]s" {
}

data "pingfederate_authentication_api_settings" "%[1]s" {
  depends_on = [
    pingfederate_authentication_api_settings.%[1]s
  ]
}`, resourceName)
	}

	return fmt.Sprintf(`
resource "pingfederate_authentication_api_settings" "%[1]s" {
  api_enabled                          = %[2]t
  enable_api_descriptions              = %[3]t
  restrict_access_to_redirectless_mode = %[4]t
  include_request_context              = %[5]t
  default_application_ref = {
    id = "%[6]s"
  }
}

data "pingfederate_authentication_api_settings" "%[1]s" {
  depends_on = [
    pingfederate_authentication_api_settings.%[1]s
  ]
}`, resourceName,
		resourceModel.apiEnabled,
		resourceModel.enableApiDescriptions,
		resourceModel.restrictAccessToRedirectlessMode,
		resourceModel.includeRequestContext,
		resourceModel.defaultApplicationRef,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationApiSettingsAttributes(config authenticationApiSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationApiSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationApiAPI.GetAuthenticationApiSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "api_enabled",
			config.apiEnabled, *response.ApiEnabled)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_api_descriptions",
			config.enableApiDescriptions, *response.EnableApiDescriptions)
		if err != nil {
			return err
		}
		return nil
	}
}
