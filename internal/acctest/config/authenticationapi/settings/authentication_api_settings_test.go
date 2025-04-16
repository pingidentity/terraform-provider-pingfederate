// Copyright © 2025 Ping Identity Corporation

package authenticationapisettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccAuthenticationApiSettings_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: authenticationApiSettings_MinimalHCL(),
				Check:  authenticationApiSettings_CheckComputedValuesMinimal(),
			},
			{
				// Update to a complete model. No computed values to check.
				Config: authenticationApiSettings_CompleteHCL(),
			},
			{
				// Test importing the resource
				Config:                               authenticationApiSettings_CompleteHCL(),
				ResourceName:                         "pingfederate_authentication_api_settings.example",
				ImportStateVerifyIdentifierAttribute: "api_enabled",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				// Back to minimal model
				Config: authenticationApiSettings_MinimalHCL(),
				Check:  authenticationApiSettings_CheckComputedValuesMinimal(),
			},
		},
	})
}

// Minimal HCL with only required values set
func authenticationApiSettings_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_api_settings" "example" {
}
data "pingfederate_authentication_api_settings" "example" {
  depends_on = [
    pingfederate_authentication_api_settings.example
  ]
}
`)
}

// Maximal HCL with all values set where possible
func authenticationApiSettings_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_api_application" "example" {
  application_id = "settingsTestApp"
  name = "authApiApp"
  url = "https://example.com"
}

resource "pingfederate_authentication_api_settings" "example" {
  api_enabled = true
  default_application_ref = {
    id = pingfederate_authentication_api_application.example.id
  }
  enable_api_descriptions = true
  include_request_context = true
  restrict_access_to_redirectless_mode = true
  # Ensures this resource will be updated before deleting the authentication api application
  lifecycle {
    create_before_destroy = true
  }
}
data "pingfederate_authentication_api_settings" "example" {
  depends_on = [
    pingfederate_authentication_api_settings.example
  ]
}
`)
}

// Validate any computed values when applying minimal HCL
func authenticationApiSettings_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_authentication_api_settings.example", "api_enabled", "false"),
		resource.TestCheckResourceAttr("pingfederate_authentication_api_settings.example", "enable_api_descriptions", "false"),
		resource.TestCheckResourceAttr("pingfederate_authentication_api_settings.example", "include_request_context", "false"),
		resource.TestCheckResourceAttr("pingfederate_authentication_api_settings.example", "restrict_access_to_redirectless_mode", "false"),
	)
}
