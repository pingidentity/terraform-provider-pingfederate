// Code generated by ping-terraform-plugin-framework-generator

package idptokenprocessors_test

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

const idpTokenProcessorProcessorId = "idpTokenProcessorProcessorId"

func TestAccIdpTokenProcessor_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpTokenProcessor_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpTokenProcessor_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					idpTokenProcessor_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIdpTokenProcessor_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpTokenProcessor_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpTokenProcessor_MinimalHCL(),
				Check:  idpTokenProcessor_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  idpTokenProcessor_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: idpTokenProcessor_CompleteHCL(),
				Check:  idpTokenProcessor_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: idpTokenProcessor_MinimalHCL(),
				Check:  idpTokenProcessor_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: idpTokenProcessor_CompleteHCL(),
				Check:  idpTokenProcessor_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               idpTokenProcessor_CompleteHCL(),
				ResourceName:                         "pingfederate_idp_token_processor.example",
				ImportStateId:                        idpTokenProcessorProcessorId,
				ImportStateVerifyIdentifierAttribute: "processor_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

// Minimal HCL with only required values set
func idpTokenProcessor_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_token_processor" "example" {
  processor_id = "%s"
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
        name  = "Audience",
        value = "myaudience"
      }
    ]
  }
  name = "My token processor"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.processor.saml.Saml20TokenProcessor"
  }
}
`, idpTokenProcessorProcessorId)
}

// Maximal HCL with all values set where possible
func idpTokenProcessor_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_token_processor" "example" {
  processor_id = "%s"
  attribute_contract = {
    core_attributes = [
      {
        name = "username"
      }
    ]
    extended_attributes = [
      {
        masked = true
        name   = "MySecretAttr"
      },
      {
        name = "MyClearAttr"
      }
    ]
    mask_ognl_values = true
  }
  configuration = {
    tables = [
      {
        name = "Credential Validators",
        rows = [
          {
            fields = [
              {
                name  = "Password Credential Validator Instance",
                value = "simple"
              }
            ]
          }
        ]
      }
    ]
    fields = [
      {
        name  = "Authentication Attempts"
        value = "3"
      }
    ]
  }
  name = "My updated token processor"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.tokenprocessors.username.UsernameTokenProcessor"
  }
}
`, idpTokenProcessorProcessorId)
}

// Validate any computed values when applying minimal HCL
func idpTokenProcessor_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "id", idpTokenProcessorProcessorId),
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "attribute_contract.extended_attributes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "attribute_contract.mask_ognl_values", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "configuration.fields_all.0.name", "Audience"),
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "configuration.fields_all.0.value", "myaudience"),
	)
}

// Validate any computed values when applying complete HCL
func idpTokenProcessor_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "id", idpTokenProcessorProcessorId),
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_token_processor.example", "attribute_contract.extended_attributes.1.masked", "false"),
	)
}

// Delete the resource
func idpTokenProcessor_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.IdpTokenProcessorsAPI.DeleteTokenProcessor(acctest.TestBasicAuthContext(), idpTokenProcessorProcessorId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func idpTokenProcessor_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.IdpTokenProcessorsAPI.DeleteTokenProcessor(acctest.TestBasicAuthContext(), idpTokenProcessorProcessorId).Execute()
	if err == nil {
		return fmt.Errorf("idp_token_processor still exists after tests. Expected it to be destroyed")
	}
	return nil
}
