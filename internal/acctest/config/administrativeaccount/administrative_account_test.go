// Copyright © 2025 Ping Identity Corporation

package administrativeaccount_test

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

const administrativeAccountUsername = "administrativeAccountUsername"

func TestAccAdministrativeAccount_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: administrativeAccount_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: administrativeAccount_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					administrativeAccount_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAdministrativeAccount_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: administrativeAccount_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: administrativeAccount_MinimalHCL(),
				Check:  administrativeAccount_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  administrativeAccount_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: administrativeAccount_CompleteHCL(),
				Check:  administrativeAccount_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: administrativeAccount_MinimalHCL(),
				Check:  administrativeAccount_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: administrativeAccount_CompleteHCL(),
				Check:  administrativeAccount_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               administrativeAccount_CompleteHCL(),
				ResourceName:                         "pingfederate_administrative_account.example",
				ImportStateId:                        administrativeAccountUsername,
				ImportStateVerifyIdentifierAttribute: "username",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// Password can't be imported, and encrypted_password will change each time it is reaad
				ImportStateVerifyIgnore: []string{"password", "encrypted_password"},
			},
		},
	})
}

// Minimal HCL with only required values set
func administrativeAccount_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_administrative_account" "example" {
  username = "%s"
  roles = ["USER_ADMINISTRATOR"]
  password = "2FederateM0re!"
}
`, administrativeAccountUsername)
}

// Maximal HCL with all values set where possible
func administrativeAccount_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_administrative_account" "example" {
  username = "%s"
  active = false
  auditor = true
  department = "mydepartment"
  description = "mydescription"
  email_address = "aggie@draynorvillage.example.com"
  password = "2FederateM0re!"
  phone_number = "555-555-5555"
  roles = []
}
`, administrativeAccountUsername)
}

// Validate any computed values when applying minimal HCL
func administrativeAccount_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_administrative_account.example", "active", "false"),
		resource.TestCheckResourceAttr("pingfederate_administrative_account.example", "auditor", "false"),
		resource.TestCheckNoResourceAttr("pingfederate_administrative_account.example", "department"),
		resource.TestCheckNoResourceAttr("pingfederate_administrative_account.example", "description"),
		resource.TestCheckNoResourceAttr("pingfederate_administrative_account.example", "email_address"),
		resource.TestCheckResourceAttrSet("pingfederate_administrative_account.example", "encrypted_password"),
		resource.TestCheckNoResourceAttr("pingfederate_administrative_account.example", "phone_number"),
	)
}

// Validate any computed values when applying complete HCL
func administrativeAccount_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("pingfederate_administrative_account.example", "encrypted_password"),
	)
}

// Delete the resource
func administrativeAccount_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.AdministrativeAccountsAPI.DeleteAccount(acctest.TestBasicAuthContext(), administrativeAccountUsername).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func administrativeAccount_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.AdministrativeAccountsAPI.DeleteAccount(acctest.TestBasicAuthContext(), administrativeAccountUsername).Execute()
	if err == nil {
		return fmt.Errorf("administrative_account still exists after tests. Expected it to be destroyed")
	}
	return nil
}
