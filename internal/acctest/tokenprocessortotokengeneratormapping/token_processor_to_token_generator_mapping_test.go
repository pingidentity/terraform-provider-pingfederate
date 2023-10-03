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

const tokenProcessorToTokenGeneratorMappingsId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type tokenProcessorToTokenGeneratorMappingsResourceModel struct {
	id string
	attributeSources []correct_type	
	attributeContractFulfillment determine this value manually	
	issuanceCriteria 	
	sourceId string	
	targetId string	
	id string	
	defaultTargetResource string	
	licenseConnectionGroupAssignment string
}

func TestAccTokenProcessorToTokenGeneratorMappings(t *testing.T) {
	resourceName := "myTokenProcessorToTokenGeneratorMappings"
	initialResourceModel := tokenProcessorToTokenGeneratorMappingsResourceModel{
		attributeSources: fill in test value,	
		attributeContractFulfillment: fill in test value,	
		issuanceCriteria: fill in test value,	
		sourceId: fill in test value,	
		targetId: fill in test value,	
		id: fill in test value,	
		defaultTargetResource: fill in test value,	
		licenseConnectionGroupAssignment: fill in test value,
	}
	updatedResourceModel := tokenProcessorToTokenGeneratorMappingsResourceModel{
		attributeSources: fill in test value,	
		attributeContractFulfillment: fill in test value,	
		issuanceCriteria: fill in test value,	
		sourceId: fill in test value,	
		targetId: fill in test value,	
		id: fill in test value,	
		defaultTargetResource: fill in test value,	
		licenseConnectionGroupAssignment: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckTokenProcessorToTokenGeneratorMappingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTokenProcessorToTokenGeneratorMappings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedTokenProcessorToTokenGeneratorMappingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccTokenProcessorToTokenGeneratorMappings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedTokenProcessorToTokenGeneratorMappingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccTokenProcessorToTokenGeneratorMappings(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_token_processor_to_token_generator_mappings." + resourceName,
				ImportStateId:           tokenProcessorToTokenGeneratorMappingsId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testAccTokenProcessorToTokenGeneratorMappings(resourceName string, resourceModel tokenProcessorToTokenGeneratorMappingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_token_processor_to_token_generator_mappings" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedTokenProcessorToTokenGeneratorMappingsAttributes(config tokenProcessorToTokenGeneratorMappingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "TokenProcessorToTokenGeneratorMappings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetTokenProcessorToTokenGeneratorMappings(ctx, tokenProcessorToTokenGeneratorMappingsId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckTokenProcessorToTokenGeneratorMappingsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeleteTokenProcessorToTokenGeneratorMappings(ctx, tokenProcessorToTokenGeneratorMappingsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("TokenProcessorToTokenGeneratorMappings", tokenProcessorToTokenGeneratorMappingsId)
	}
	return nil
}
