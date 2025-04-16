// Copyright © 2025 Ping Identity Corporation

package authenticationapiapplication_test

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

const authenticationApiApplicationApplicationId = "authenticationApiApplicationAppl"

func TestAccAuthenticationApiApplication_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: authenticationApiApplication_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: authenticationApiApplication_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					authenticationApiApplication_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAuthenticationApiApplication_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: authenticationApiApplication_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: authenticationApiApplication_MinimalHCL(),
				Check:  authenticationApiApplication_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  authenticationApiApplication_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model. No computed values to check in complete model
				Config: authenticationApiApplication_CompleteHCL(),
			},
			{
				// Back to minimal model
				Config: authenticationApiApplication_MinimalHCL(),
				Check:  authenticationApiApplication_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: authenticationApiApplication_CompleteHCL(),
			},
			{
				// Test importing the resource
				Config:                               authenticationApiApplication_CompleteHCL(),
				ResourceName:                         "pingfederate_authentication_api_application.example",
				ImportStateId:                        authenticationApiApplicationApplicationId,
				ImportStateVerifyIdentifierAttribute: "application_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

// Minimal HCL with only required values set
func authenticationApiApplication_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_api_application" "example" {
  application_id = "%s"
  name = "authApiApp"
  url = "https://example.com"
}
`, authenticationApiApplicationApplicationId)
}

// Maximal HCL with all values set where possible
func authenticationApiApplication_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_api_settings" "example" {
  restrict_access_to_redirectless_mode = true
}

	resource "pingfederate_oauth_client" "example" {
		client_id                     = "myOauthClientExample"
		name                          = "myOauthClientExample"
		grant_types                   = ["EXTENSION"]
		allow_authentication_api_init = true
	  }

resource "pingfederate_authentication_api_application" "example" {
  depends_on = [pingfederate_authentication_api_settings.example]
  application_id = "%s"
  additional_allowed_origins = [
	"https://example.com",
	"https://example2.com",
  ]
  client_for_redirectless_mode_ref = {
    id = pingfederate_oauth_client.example.id
  }
  description = "this is my app"
  name = "authApiAppUpdated"
  url = "https://changed.example.com"
}
`, authenticationApiApplicationApplicationId)
}

// Validate any computed values when applying minimal HCL
func authenticationApiApplication_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_authentication_api_application.example", "additional_allowed_origins.#", "0"),
		resource.TestCheckNoResourceAttr("pingfederate_authentication_api_application.example", "description"),
	)
}

// Delete the resource
func authenticationApiApplication_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.AuthenticationApiAPI.DeleteApplication(acctest.TestBasicAuthContext(), authenticationApiApplicationApplicationId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func authenticationApiApplication_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.AuthenticationApiAPI.DeleteApplication(acctest.TestBasicAuthContext(), authenticationApiApplicationApplicationId).Execute()
	if err == nil {
		return fmt.Errorf("authentication_api_application still exists after tests. Expected it to be destroyed")
	}
	return nil
}
