// Copyright © 2025 Ping Identity Corporation

package authenticationselector_test

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

const authenticationSelectorSelectorId = "authenticationSelectorSelectorId"

func TestAccAuthenticationSelector_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: authenticationSelector_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: authenticationSelector_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					authenticationSelector_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAuthenticationSelector_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: authenticationSelector_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: authenticationSelector_MinimalHCL(),
				Check:  authenticationSelector_CheckComputedValues(),
			},
			{
				// Delete the minimal model
				Config:  authenticationSelector_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: authenticationSelector_CompleteHCL(),
				Check:  authenticationSelector_CheckComputedValues(),
			},
			{
				// Back to minimal model
				Config: authenticationSelector_MinimalHCL(),
				Check:  authenticationSelector_CheckComputedValues(),
			},
			{
				// Back to complete model
				Config: authenticationSelector_CompleteHCL(),
				Check:  authenticationSelector_CheckComputedValues(),
			},
			{
				// Test importing the resource
				Config:                               authenticationSelector_CompleteHCL(),
				ResourceName:                         "pingfederate_authentication_selector.example",
				ImportStateId:                        authenticationSelectorSelectorId,
				ImportStateVerifyIdentifierAttribute: "selector_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore: []string{
					"configuration.tables",
					"configuration.fields",
				},
			},
		},
	})
}

// Minimal HCL with only required values set
func authenticationSelector_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_selector" "example" {
  selector_id = "%s"
  attribute_contract = {
    extended_attributes = [
      {
        name = "extendedattr"
      }
    ]
  }
  configuration = {
  }
  name = "myInitialSelector"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.selectors.saml.SamlAuthnContextAdapterSelector"
  }
}
`, authenticationSelectorSelectorId)
}

// Maximal HCL with all values set where possible
func authenticationSelector_CompleteHCL() string {
	versionedConfigurationFields := ""
	if acctest.VersionAtLeast("11.3.0") {
		versionedConfigurationFields = `
		{
			name  = "Override AuthN Context for Flow"
			value = "true"
		},
		`
	}
	return fmt.Sprintf(`
resource "pingfederate_authentication_selector" "example" {
  selector_id = "%s"
  attribute_contract = {
    extended_attributes = [
	  {
	    name = "another"
	  },
      {
        name = "extendedattr"
      }
    ]
  }
  configuration = {
    tables = []
    fields = [
		%s
      {
        name  = "Add or Update AuthN Context Attribute"
        value = "true"
      },
      {
        name  = "Enable 'No Match' Result Value"
        value = "true"
      },
      {
        name  = "Enable 'Not in Request' Result Value"
        value = "true"
      }
    ]
  }
  name = "myUpdatedSelector"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.selectors.saml.SamlAuthnContextAdapterSelector"
  }
}
`, authenticationSelectorSelectorId, versionedConfigurationFields)
}

// Validate any computed values when applying HCL
func authenticationSelector_CheckComputedValues() resource.TestCheckFunc {
	expectedFieldCount := "3"
	if acctest.VersionAtLeast("11.3.0") {
		expectedFieldCount = "4"
	}
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_authentication_selector.example", "configuration.tables_all.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_authentication_selector.example", "configuration.fields_all.#", expectedFieldCount),
	)
}

// Delete the resource
func authenticationSelector_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.AuthenticationSelectorsAPI.DeleteAuthenticationSelector(acctest.TestBasicAuthContext(), authenticationSelectorSelectorId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func authenticationSelector_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.AuthenticationSelectorsAPI.DeleteAuthenticationSelector(acctest.TestBasicAuthContext(), authenticationSelectorSelectorId).Execute()
	if err == nil {
		return fmt.Errorf("authentication_selector still exists after tests. Expected it to be destroyed")
	}
	return nil
}
