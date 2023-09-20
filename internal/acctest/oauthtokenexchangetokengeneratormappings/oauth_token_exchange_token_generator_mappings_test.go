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

const oauthTokenExchangeTokenGeneratorMappingsId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type oauthTokenExchangeTokenGeneratorMappingsResourceModel struct {
	id string
	attributeSources []correct_type	
	attributeContractFulfillment determine this value manually	
	issuanceCriteria 	
	id string	
	sourceId string	
	targetId string	
	licenseConnectionGroupAssignment string
}

func TestAccOauthTokenExchangeTokenGeneratorMappings(t *testing.T) {
	resourceName := "myOauthTokenExchangeTokenGeneratorMappings"
	initialResourceModel := oauthTokenExchangeTokenGeneratorMappingsResourceModel{
		attributeSources: fill in test value,	
		attributeContractFulfillment: fill in test value,	
		issuanceCriteria: fill in test value,	
		id: fill in test value,	
		sourceId: fill in test value,	
		targetId: fill in test value,	
		licenseConnectionGroupAssignment: fill in test value,
	}
	updatedResourceModel := oauthTokenExchangeTokenGeneratorMappingsResourceModel{
		attributeSources: fill in test value,	
		attributeContractFulfillment: fill in test value,	
		issuanceCriteria: fill in test value,	
		id: fill in test value,	
		sourceId: fill in test value,	
		targetId: fill in test value,	
		licenseConnectionGroupAssignment: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckOauthTokenExchangeTokenGeneratorMappingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthTokenExchangeTokenGeneratorMappings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeTokenGeneratorMappingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthTokenExchangeTokenGeneratorMappings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeTokenGeneratorMappingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccOauthTokenExchangeTokenGeneratorMappings(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_oauth_token_exchange_token_generator_mappings." + resourceName,
				ImportStateId:           oauthTokenExchangeTokenGeneratorMappingsId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testAccOauthTokenExchangeTokenGeneratorMappings(resourceName string, resourceModel oauthTokenExchangeTokenGeneratorMappingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_token_exchange_token_generator_mappings" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthTokenExchangeTokenGeneratorMappingsAttributes(config oauthTokenExchangeTokenGeneratorMappingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthTokenExchangeTokenGeneratorMappings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetOauthTokenExchangeTokenGeneratorMappings(ctx, oauthTokenExchangeTokenGeneratorMappingsId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOauthTokenExchangeTokenGeneratorMappingsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeleteOauthTokenExchangeTokenGeneratorMappings(ctx, oauthTokenExchangeTokenGeneratorMappingsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthTokenExchangeTokenGeneratorMappings", oauthTokenExchangeTokenGeneratorMappingsId)
	}
	return nil
}
