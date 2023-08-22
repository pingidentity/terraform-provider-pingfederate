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

const idpAdaptersId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type idpAdaptersResourceModel struct {
	id string
	authnCtxClassRef string	
	id string	
	name string	
	pluginDescriptorRef 	
	parentRef 	
	configuration 	
	attributeMapping 	
	attributeContract 
}

func TestAccIdpAdapters(t *testing.T) {
	resourceName := "myIdpAdapters"
	initialResourceModel := idpAdaptersResourceModel{
		authnCtxClassRef: fill in test value,	
		id: fill in test value,	
		name: fill in test value,	
		pluginDescriptorRef: fill in test value,	
		parentRef: fill in test value,	
		configuration: fill in test value,	
		attributeMapping: fill in test value,	
		attributeContract: fill in test value,
	}
	updatedResourceModel := idpAdaptersResourceModel{
		authnCtxClassRef: fill in test value,	
		id: fill in test value,	
		name: fill in test value,	
		pluginDescriptorRef: fill in test value,	
		parentRef: fill in test value,	
		configuration: fill in test value,	
		attributeMapping: fill in test value,	
		attributeContract: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckIdpAdaptersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdpAdapters(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpAdaptersAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIdpAdapters(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIdpAdaptersAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccIdpAdapters(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_idp_adapters." + resourceName,
				ImportStateId:           idpAdaptersId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testAccIdpAdapters(resourceName string, resourceModel idpAdaptersResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_adapters" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedIdpAdaptersAttributes(config idpAdaptersResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "IdpAdapters"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetIdpAdapters(ctx, idpAdaptersId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckIdpAdaptersDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeleteIdpAdapters(ctx, idpAdaptersId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("IdpAdapters", idpAdaptersId)
	}
	return nil
}
