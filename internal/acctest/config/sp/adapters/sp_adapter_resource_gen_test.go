// Code generated by ping-terraform-plugin-framework-generator

package spadapters_test

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

const spAdapterAdapterId = "spAdapterAdapterId"

func TestAccSpAdapter_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spAdapter_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spAdapter_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					spAdapter_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSpAdapter_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spAdapter_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spAdapter_MinimalHCL(),
				Check:  spAdapter_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  spAdapter_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: spAdapter_CompleteHCL(),
				Check:  spAdapter_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: spAdapter_MinimalHCL(),
				Check:  spAdapter_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: spAdapter_CompleteHCL(),
				Check:  spAdapter_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               spAdapter_CompleteHCL(),
				ResourceName:                         "pingfederate_sp_adapter.example",
				ImportStateId:                        spAdapterAdapterId,
				ImportStateVerifyIdentifierAttribute: "adapter_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// configuration.fields has some sensitive values which can't be imported
				ImportStateVerifyIgnore: []string{"configuration.sensitive_fields"},
			},
		},
	})
}

// Minimal HCL with only required values set
func spAdapter_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_sp_adapter" "example" {
  adapter_id = "%s"
  configuration = {
    fields = [
      {
        "name" : "Password",
        "value" : "2FederateM0re"
      },
      {
        "name" : "Confirm Password",
        "value" : "2FederateM0re"
      },
    ]
  }
  name = "My sp adapter"
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.opentoken.SpAuthnAdapter"
  }
}
`, spAdapterAdapterId)
}

// Maximal HCL with all values set where possible
func spAdapter_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_sp_adapter" "example" {
  adapter_id = "%s"
  attribute_contract = {
    extended_attributes = [
      {
        name = "My extended attribute"
      }
    ]
  }
  configuration = {
    fields = [
      {
        "name" : "Transport Mode",
        "value" : "2"
      },
      {
        "name" : "Token Name",
        "value" : "opentoken"
      },
      {
        "name" : "Cipher Suite",
        "value" : "2"
      },
      {
        "name" : "Authentication Service",
        "value" : ""
      },
      {
        "name" : "Account Link Service",
        "value" : ""
      },
      {
        "name" : "Logout Service",
        "value" : ""
      },
      {
        "name" : "SameSite Cookie",
        "value" : "3"
      },
      {
        "name" : "Cookie Domain",
        "value" : ""
      },
      {
        "name" : "Cookie Path",
        "value" : "/"
      },
      {
        "name" : "Token Lifetime",
        "value" : "300"
      },
      {
        "name" : "Session Lifetime",
        "value" : "43200"
      },
      {
        "name" : "Not Before Tolerance",
        "value" : "0"
      },
      {
        "name" : "Force SunJCE Provider",
        "value" : "false"
      },
      {
        "name" : "Use Verbose Error Messages",
        "value" : "false"
      },
      {
        "name" : "Obfuscate Password",
        "value" : "true"
      },
      {
        "name" : "Session Cookie",
        "value" : "false"
      },
      {
        "name" : "Secure Cookie",
        "value" : "true"
      },
      {
        "name" : "HTTP Only Flag",
        "value" : "true"
      },
      {
        "name" : "Send Subject as Query Parameter",
        "value" : "false"
      },
      {
        "name" : "Subject Query Parameter                 ",
        "value" : ""
      },
      {
        "name" : "Send Extended Attributes",
        "value" : ""
      },
      {
        "name" : "Skip Trimming of Trailing Backslashes",
        "value" : "false"
      },
      {
        "name" : "URL Encode Cookie Values",
        "value" : "true"
      },
    ]
    sensitive_fields = [
      {
        "name" : "Password",
        "value" : "2FederateM0re"
      },
      {
        "name" : "Confirm Password",
        "value" : "2FederateM0re"
      },
    ]
    tables = []
  }
  name = "My sp adapter"
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.opentoken.SpAuthnAdapter"
  }
  target_application_info = {
    application_icon_url = "https://www.example.com/icon.png"
    application_name     = "My application name"
  }
}
`, spAdapterAdapterId)
}

// Validate any computed values when applying minimal HCL
func spAdapter_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_adapter.example", "attribute_contract.core_attributes.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_sp_adapter.example", "attribute_contract.core_attributes.0.name", "subject"),
		resource.TestCheckResourceAttr("pingfederate_sp_adapter.example", "attribute_contract.extended_attributes.#", "0"),
		resource.TestCheckNoResourceAttr("pingfederate_sp_adapter.example", "target_application_info.application_icon_url"),
		resource.TestCheckNoResourceAttr("pingfederate_sp_adapter.example", "target_application_info.application_name"),
	)
}

// Validate any computed values when applying complete HCL
func spAdapter_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_adapter.example", "attribute_contract.core_attributes.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_sp_adapter.example", "attribute_contract.core_attributes.0.name", "subject"),
	)
}

// Delete the resource
func spAdapter_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.SpAdaptersAPI.DeleteSpAdapter(acctest.TestBasicAuthContext(), spAdapterAdapterId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func spAdapter_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.SpAdaptersAPI.DeleteSpAdapter(acctest.TestBasicAuthContext(), spAdapterAdapterId).Execute()
	if err == nil {
		return fmt.Errorf("sp_adapter still exists after tests. Expected it to be destroyed")
	}
	return nil
}
