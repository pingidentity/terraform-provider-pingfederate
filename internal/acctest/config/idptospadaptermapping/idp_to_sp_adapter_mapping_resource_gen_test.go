// Code generated by ping-terraform-plugin-framework-generator

package idptospadaptermapping_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const idpAdapterId = "OTIdPJava"
const spAdapterId = "spadapter"
const idpToSpAdapterMappingMappingId = "OTIdPJava|spadapter"

func TestAccIdpToSpAdapterMapping_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpToSpAdapterMapping_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpToSpAdapterMapping_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					idpToSpAdapterMapping_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIdpToSpAdapterMapping_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpToSpAdapterMapping_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpToSpAdapterMapping_MinimalHCL(),
				Check:  idpToSpAdapterMapping_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  idpToSpAdapterMapping_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: idpToSpAdapterMapping_CompleteHCL(),
				Check:  idpToSpAdapterMapping_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: idpToSpAdapterMapping_MinimalHCL(),
				Check:  idpToSpAdapterMapping_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: idpToSpAdapterMapping_CompleteHCL(),
				Check:  idpToSpAdapterMapping_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               idpToSpAdapterMapping_CompleteHCL(),
				ResourceName:                         "pingfederate_idp_to_sp_adapter_mapping.example",
				ImportStateId:                        idpToSpAdapterMappingMappingId,
				ImportStateVerifyIdentifierAttribute: "mapping_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

// Minimal HCL with only required values set
func idpToSpAdapterMapping_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_to_sp_adapter_mapping" "example" {
  attribute_contract_fulfillment = {
    "subject" = {
      source = {
        type = "ADAPTER"
      }
      value = "subject"
    }
  }
  source_id = "%s"
  target_id = "%s"
}
`, idpAdapterId, spAdapterId)
}

// Maximal HCL with all values set where possible
func idpToSpAdapterMapping_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_to_sp_adapter_mapping" "example" {
  application_icon_url = "https://example.com/icon.png"
  application_name     = "My Application"
  attribute_contract_fulfillment = {
    "subject" = {
      source = {
        type = "ADAPTER"
      }
      value = "subject"
    }
  }
  // attribute_sources
  %s
  default_target_resource = "https://example.com"
  // issuance_criteria
  %s
  source_id = "%s"
  target_id = "%s"
}
`, attributesources.Hcl(nil, attributesources.LdapClientStruct("(cn=Example)", "SUBTREE", *client.NewResourceLink("pingdirectory"))),
		issuancecriteria.Hcl(issuancecriteria.ConditionalCriteria()),
		idpAdapterId, spAdapterId)
}

// Validate any computed values when applying minimal HCL
func idpToSpAdapterMapping_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_idp_to_sp_adapter_mapping.example", "mapping_id", idpToSpAdapterMappingMappingId),
		resource.TestCheckResourceAttr("pingfederate_idp_to_sp_adapter_mapping.example", "attribute_sources.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_to_sp_adapter_mapping.example", "issuance_criteria.conditional_criteria.#", "0"),
	)
}

// Validate any computed values when applying complete HCL
func idpToSpAdapterMapping_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_idp_to_sp_adapter_mapping.example", "mapping_id", idpToSpAdapterMappingMappingId),
	)
}

// Delete the resource
func idpToSpAdapterMapping_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.IdpToSpAdapterMappingAPI.DeleteIdpToSpAdapterMappingsById(acctest.TestBasicAuthContext(), idpToSpAdapterMappingMappingId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func idpToSpAdapterMapping_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.IdpToSpAdapterMappingAPI.DeleteIdpToSpAdapterMappingsById(acctest.TestBasicAuthContext(), idpToSpAdapterMappingMappingId).Execute()
	if err == nil {
		return fmt.Errorf("idp_to_sp_adapter_mapping still exists after tests. Expected it to be destroyed")
	}
	return nil
}
