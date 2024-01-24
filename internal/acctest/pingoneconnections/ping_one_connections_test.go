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

const pingOneConnectionId = "2"
const pingOneConnectionName = "myPingOneConnection"

// Attributes to test with. Add optional properties to test here if desired.
type pingOneConnectionResourceModel struct {
	id string	
	name string	
	description string	
	active bool	
}

func TestAccPingOneConnection(t *testing.T) {
	resourceName := "myPingOneConnection"
	initialResourceModel := pingOneConnectionResourceModel{
		id: pingOneConnectionId,
		name: pingOneConnectionName,
		description:,
		active:,
	}
	updatedResourceModel := pingOneConnectionResourceModel{
		id: pingOneConnectionId,
		name: pingOneConnectionName,
		description:,
		active:,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckPingOneConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOneConnection(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPingOneConnectionAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPingOneConnection(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPingOneConnectionAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccPingOneConnection(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_ping_one_connections." + resourceName,
				ImportStateId:           pingOneConnectionId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
			{
				Config: testAccPingOneConnection(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPingOneConnectionAttributes(initialResourceModel),
			},
		},
	})
}

func testAccPingOneConnection(resourceName string, resourceModel pingOneConnectionResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_ping_one_connections" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedPingOneConnectionAttributes(config pingOneConnectionResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "PingOneConnection"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetPingOneConnection(ctx, pingOneConnectionId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPingOneConnectionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeletePingOneConnection(ctx, pingOneConnectionId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("PingOneConnection", pingOneConnectionId)
	}
	return nil
}
