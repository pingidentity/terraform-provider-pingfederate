package idpspconnection_test

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

const idpSpConnId = "deletedriftconn"

func TestAccIdpSpConnection_RemovalDrift(t *testing.T) {
	t.SkipNow()
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spIdpConnection_SimpleCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpSpConnection_SamlMinimalHCL(idpSpConnId),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					spIdpConnection_Delete(t, idpSpConnId)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Delete the resource
func spIdpConnection_Delete(t *testing.T, id string) {
	testClient := acctest.TestClient()
	_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(acctest.TestBasicAuthContext(), id).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func spIdpConnection_SimpleCheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(acctest.TestBasicAuthContext(), idpSpConnId).Execute()
	if err == nil {
		return fmt.Errorf("sp_idp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
