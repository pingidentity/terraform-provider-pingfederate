package authenticationpolicycontract_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

var stateId string

// Attributes to test with. Add optional properties to test here if desired.
type authenticationPolicyContractResourceModel struct {
	id                 string
	name               string
	extendedAttributes []string
}

func TestAccAuthenticationPolicyContract(t *testing.T) {
	resourceName := "myAuthenticationPolicyContract"
	initialResourceModel := authenticationPolicyContractResourceModel{
		name:               "example",
		extendedAttributes: []string{},
	}
	updatedResourceModel := authenticationPolicyContractResourceModel{
		name:               "example",
		extendedAttributes: []string{"extended_attribute", "extended_attribute2", "extendedwith\\\"escaped\\\"quotes"},
	}

	minimalResourceModelWithId := authenticationPolicyContractResourceModel{
		name:               "example",
		extendedAttributes: []string{},
		id:                 "myAuthenticationPolicyContract",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckAuthenticationPolicyContractDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal model
				Config: testAccAuthenticationPolicyContract(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationPolicyContractAttributes(initialResourceModel),
			},
			{
				// More complete model
				Config: testAccAuthenticationPolicyContract(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAuthenticationPolicyContractAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policy_contract.%s", resourceName), "extended_attributes.0.name", updatedResourceModel.extendedAttributes[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policy_contract.%s", resourceName), "extended_attributes.1.name", updatedResourceModel.extendedAttributes[1]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policy_contract.%s", resourceName), "extended_attributes.2.name", "extendedwith\"escaped\"quotes"),
				),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationPolicyContract(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_authentication_policy_contract." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccAuthenticationPolicyContract(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationPolicyContractAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.AuthenticationPolicyContractsAPI.DeleteAuthenticationPolicyContract(ctx, stateId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Minimal model
				Config: testAccAuthenticationPolicyContract(resourceName, minimalResourceModelWithId),
				Check:  testAccCheckExpectedAuthenticationPolicyContractAttributes(minimalResourceModelWithId),
			},
		},
	})
}

func testAccAuthenticationPolicyContract(resourceName string, resourceModel authenticationPolicyContractResourceModel) string {
	extendedAttrsHcl := ""
	if len(resourceModel.extendedAttributes) > 0 {
		extendedAttrsHcl = fmt.Sprintf("extended_attributes = %[1]s",
			acctest.ObjectSliceOfKvStringsToTerraformString("name", resourceModel.extendedAttributes))
	}
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "%[1]s" {
  %[2]s
  name = "%[3]s"
  %[4]s
}
data "pingfederate_authentication_policy_contract" "authenticationPolicyContractExample" {
  contract_id = pingfederate_authentication_policy_contract.%[1]s.contract_id
}`, resourceName,
		acctest.AddIdHcl("contract_id", resourceModel.id),
		resourceModel.name,
		extendedAttrsHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPolicyContractAttributes(config authenticationPolicyContractResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPolicyContract"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		stateAttributes := s.RootModule().Resources["pingfederate_authentication_policy_contract.myAuthenticationPolicyContract"].Primary.Attributes
		response, _, err := testClient.AuthenticationPolicyContractsAPI.GetAuthenticationPolicyContract(ctx, stateAttributes["id"]).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name",
			config.name, *response.Name)
		if err != nil {
			return err
		}

		extendedAttrs := response.GetExtendedAttributes()
		extendedAttrNames := make([]string, len(extendedAttrs))
		for i := 0; i < len(extendedAttrNames); i++ {
			// Need to "re-escape" quotes here since they had to be escaped to form the terraform HCL
			extendedAttrNames[i] = strings.ReplaceAll(extendedAttrs[i].Name, "\"", "\\\"")
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "extended_attributes",
			config.extendedAttributes, extendedAttrNames)
		if err != nil {
			return err
		}

		stateId = stateAttributes["id"]
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAuthenticationPolicyContractDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AuthenticationPolicyContractsAPI.GetAuthenticationPolicyContract(ctx, stateId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AuthenticationPolicyContract", stateId)
	}
	return nil
}
