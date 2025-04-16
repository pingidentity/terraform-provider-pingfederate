// Copyright © 2025 Ping Identity Corporation

package authenticationpolicycontract_test

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

const authenticationPolicyContractContractId = "authenticationPolicyContractCont"

func TestAccAuthenticationPolicyContract_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: authenticationPolicyContract_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: authenticationPolicyContract_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					authenticationPolicyContract_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAuthenticationPolicyContract_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: authenticationPolicyContract_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: authenticationPolicyContract_MinimalHCL(),
				Check:  authenticationPolicyContract_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  authenticationPolicyContract_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: authenticationPolicyContract_CompleteHCL(),
				Check:  authenticationPolicyContract_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: authenticationPolicyContract_MinimalHCL(),
				Check:  authenticationPolicyContract_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: authenticationPolicyContract_CompleteHCL(),
				Check:  authenticationPolicyContract_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               authenticationPolicyContract_CompleteHCL(),
				ResourceName:                         "pingfederate_authentication_policy_contract.example",
				ImportStateId:                        authenticationPolicyContractContractId,
				ImportStateVerifyIdentifierAttribute: "contract_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

// Minimal HCL with only required values set
func authenticationPolicyContract_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "example" {
  contract_id = "%s"
  name = "initialApc"
}
  data "pingfederate_authentication_policy_contract" "example" {
  contract_id = pingfederate_authentication_policy_contract.example.id
}
`, authenticationPolicyContractContractId)
}

// Maximal HCL with all values set where possible
func authenticationPolicyContract_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "example" {
  contract_id = "%s"
  extended_attributes = [
    {
      name = "extended_attribute"
    },
	{
					name = "extended_attribute2"
	},
	{
		name = "extendedwith\"escaped\"quotes"
	}
  ]
  name = "myApc"
}
  data "pingfederate_authentication_policy_contract" "example" {
  contract_id = pingfederate_authentication_policy_contract.example.id
}
`, authenticationPolicyContractContractId)
}

// Validate any computed values when applying minimal HCL
func authenticationPolicyContract_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_authentication_policy_contract.example", "core_attributes.0.name", "subject"),
		resource.TestCheckResourceAttr("pingfederate_authentication_policy_contract.example", "extended_attributes.#", "0"),
	)

}

// Validate any computed values when applying complete HCL
func authenticationPolicyContract_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_authentication_policy_contract.example", "core_attributes.0.name", "subject"),
	)
}

// Delete the resource
func authenticationPolicyContract_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.AuthenticationPolicyContractsAPI.DeleteAuthenticationPolicyContract(acctest.TestBasicAuthContext(), authenticationPolicyContractContractId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func authenticationPolicyContract_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.AuthenticationPolicyContractsAPI.DeleteAuthenticationPolicyContract(acctest.TestBasicAuthContext(), authenticationPolicyContractContractId).Execute()
	if err == nil {
		return fmt.Errorf("authentication_policy_contract still exists after tests. Expected it to be destroyed")
	}
	return nil
}
