// Code generated by ping-terraform-plugin-framework-generator

package sptokengenerators_test

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

const spTokenGeneratorGeneratorId = "spTokenGeneratorGeneratorId"

func TestAccSpTokenGenerator_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spTokenGenerator_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spTokenGenerator_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					spTokenGenerator_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSpTokenGenerator_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spTokenGenerator_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spTokenGenerator_MinimalHCL(),
				Check:  spTokenGenerator_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  spTokenGenerator_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: spTokenGenerator_CompleteHCL(),
				Check:  spTokenGenerator_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: spTokenGenerator_MinimalHCL(),
				Check:  spTokenGenerator_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: spTokenGenerator_CompleteHCL(),
				Check:  spTokenGenerator_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               spTokenGenerator_CompleteHCL(),
				ResourceName:                         "pingfederate_sp_token_generator.example",
				ImportStateId:                        spTokenGeneratorGeneratorId,
				ImportStateVerifyIdentifierAttribute: "generator_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

// Minimal HCL with only required values set
func spTokenGenerator_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_sp_token_generator" "example" {
  generator_id = "%s"
  attribute_contract = {
    core_attributes = [
      {
        name = "SAML_SUBJECT"
      }
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Minutes Before"
        value = "60"
      },
      {
        name  = "Minutes After"
        value = "60"
      },
      {
        name  = "Issuer"
        value = "issuer"
      },
      {
        name  = "Signing Certificate"
        value = "419x9yg43rlawqwq9v6az997k"
      },
      {
        name  = "Signing Algorithm"
        value = "SHA1"
      },
      {
        name  = "Include Certificate in KeyInfo"
        value = "false"
      },
      {
        name  = "Include Raw Key in KeyValue"
        value = "false"
      },
      {
        name  = "Audience"
        value = "audience"
      },
      {
        name  = "Confirmation Method"
        value = "urn:oasis:names:tc:SAML:2.0:cm:sender-vouches"
      }
    ]
  }
  name = "My token generator"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.generator.saml.Saml20TokenGenerator"
  }
}
`, spTokenGeneratorGeneratorId)
}

// Maximal HCL with all values set where possible
func spTokenGenerator_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_sp_token_generator" "example" {
  generator_id = "%s"
  attribute_contract = {
    core_attributes = [
      {
        name = "SAML_SUBJECT"
      }
    ]
    extended_attributes = [
      {
        name = "ExtendedAttr"
      }
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Minutes Before"
        value = "60"
      },
      {
        name  = "Minutes After"
        value = "60"
      },
      {
        name  = "Issuer"
        value = "issuer"
      },
      {
        name  = "Signing Certificate"
        value = "419x9yg43rlawqwq9v6az997k"
      },
      {
        name  = "Signing Algorithm"
        value = "SHA1"
      },
      {
        name  = "Include Certificate in KeyInfo"
        value = "false"
      },
      {
        name  = "Include Raw Key in KeyValue"
        value = "false"
      },
      {
        name  = "Audience"
        value = "audience"
      },
      {
        name  = "Confirmation Method"
        value = "urn:oasis:names:tc:SAML:2.0:cm:sender-vouches"
      },
      {
        name  = "Encryption Certificate"
        value = ""
      },
      {
        name  = "Message Customization Expression"
        value = ""
      }
    ]
    tables = []
  }
  name = "My updated token generator"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.generator.saml.Saml20TokenGenerator"
  }
}
`, spTokenGeneratorGeneratorId)
}

// Validate any computed values when applying minimal HCL
func spTokenGenerator_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_token_generator.example", "id", spTokenGeneratorGeneratorId),
		resource.TestCheckResourceAttr("pingfederate_sp_token_generator.example", "configuration.fields_all.#", "11"),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_sp_token_generator.example", "configuration.fields_all.*", map[string]string{
			"name":  "Encryption Certificate",
			"value": "",
		}),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_sp_token_generator.example", "configuration.fields_all.*", map[string]string{
			"name":  "Message Customization Expression",
			"value": "",
		}),
		resource.TestCheckResourceAttr("pingfederate_sp_token_generator.example", "configuration.tables_all.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_token_generator.example", "attribute_contract.extended_attributes.#", "0"),
	)
}

// Validate any computed values when applying complete HCL
func spTokenGenerator_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_token_generator.example", "id", spTokenGeneratorGeneratorId),
		resource.TestCheckResourceAttr("pingfederate_sp_token_generator.example", "configuration.fields_all.#", "11"),
		resource.TestCheckResourceAttr("pingfederate_sp_token_generator.example", "configuration.tables_all.#", "0"),
	)
}

// Delete the resource
func spTokenGenerator_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.SpTokenGeneratorsAPI.DeleteTokenGenerator(acctest.TestBasicAuthContext(), spTokenGeneratorGeneratorId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func spTokenGenerator_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.SpTokenGeneratorsAPI.DeleteTokenGenerator(acctest.TestBasicAuthContext(), spTokenGeneratorGeneratorId).Execute()
	if err == nil {
		return fmt.Errorf("sp_token_generator still exists after tests. Expected it to be destroyed")
	}
	return nil
}
