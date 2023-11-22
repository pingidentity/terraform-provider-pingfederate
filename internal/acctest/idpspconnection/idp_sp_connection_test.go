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

const spConnectionId = "spConnId"

type spConnectionResourceModel struct {
	name     string
	entityId string
}

func TestAccIdpSpConnection(t *testing.T) {
	const resourceName = "mySpConnection"
	initialResourceModel := spConnectionResourceModel{
		name:     "spConnName",
		entityId: "myEntity",
	}

	updatedResourceModel := spConnectionResourceModel{
		name:     "spConnNameUpdated",
		entityId: "myEntityUpdated",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckIdpAdapterDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal model
				Config: testAccIdpAdapter(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpAdapterAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIdpAdapter(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIdpAdapterAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccIdpAdapter(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_idp_adapter." + resourceName,
				ImportStateId:     spConnectionId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to the initial minimal model
				Config: testAccIdpAdapter(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpAdapterAttributes(initialResourceModel),
			},
		},
	})
}

func testAccIdpAdapter(resourceName string, resourceModel spConnectionResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = "%[1]s"
  entity_id = "%s"
  name       = "%s"
}`, resourceName,
		resourceModel.entityId,
		resourceModel.name,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedIdpAdapterAttributes(config spConnectionResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//resourceType := "IdP SP Connection"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.IdpSpConnectionsAPI.GetSpConnection(ctx, spConnectionId).Execute()

		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckIdpAdapterDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(ctx, spConnectionId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("IdP SP Connection", spConnectionId)
	}
	return nil
}
