package acctest_test

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

const authenticationPolicyContractsId = "2"
const coreAttr = "subject"

// Attributes to test with. Add optional properties to test here if desired.
type authenticationPolicyContractsResourceModel struct {
	id                 string
	name               string
	coreAttributes     []string
	extendedAttributes []string
}

func TestAccAuthenticationPolicyContracts(t *testing.T) {
	resourceName := "myAuthenticationPolicyContracts"
	initialResourceModel := authenticationPolicyContractsResourceModel{
		id:                 authenticationPolicyContractsId,
		name:               "example",
		coreAttributes:     []string{coreAttr},
		extendedAttributes: []string{"extended_attribute", "extended_attribute2"},
	}
	updatedResourceModel := authenticationPolicyContractsResourceModel{
		id:                 authenticationPolicyContractsId,
		name:               "example",
		coreAttributes:     []string{coreAttr},
		extendedAttributes: []string{"extended_attribute"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckAuthenticationPolicyContractsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPolicyContracts(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationPolicyContractsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAuthenticationPolicyContracts(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAuthenticationPolicyContractsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationPolicyContracts(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_authentication_policy_contract." + resourceName,
				ImportStateId:     authenticationPolicyContractsId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAuthenticationPolicyContracts(resourceName string, resourceModel authenticationPolicyContractsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "%[1]s" {
  custom_id           = "%[2]s"
  core_attributes     = %[3]s
  extended_attributes = %[4]s
  name                = "%[5]s"
}`, resourceName,
		resourceModel.id,
		acctest.ObjectSliceOfKvStringsToTerraformString("name", resourceModel.coreAttributes),
		acctest.ObjectSliceOfKvStringsToTerraformString("name", resourceModel.extendedAttributes),
		resourceModel.name,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPolicyContractsAttributes(config authenticationPolicyContractsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPolicyContracts"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationPolicyContractsApi.GetAuthenticationPolicyContract(ctx, authenticationPolicyContractsId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name",
			config.name, *response.Name)
		if err != nil {
			return err
		}

		rEa := response.GetExtendedAttributes()
		newSet := make([]string, len(rEa))
		for i := 0; i < len(newSet); i++ {
			newSet[i] = rEa[i].Name
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "extended_properties",
			config.extendedAttributes, newSet)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAuthenticationPolicyContractsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, _, err := testClient.AuthenticationPolicyContractsApi.GetAuthenticationPolicyContract(ctx, authenticationPolicyContractsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AuthenticationPolicyContracts", authenticationPolicyContractsId)
	}
	return nil
}
