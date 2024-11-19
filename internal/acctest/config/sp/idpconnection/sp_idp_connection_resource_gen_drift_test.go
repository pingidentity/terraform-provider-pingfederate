// Code generated by ping-terraform-plugin-framework-generator

package resource_sp_idp_connection_test

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

const spIdpConnectionConnectionId = "sp_idp_connection_connection_id"

func TestAccSpIdpConnection_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spIdpConnection_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spIdpConnection_WsTrustStsMinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					spIdpConnection_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Delete the resource
func spIdpConnection_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.SpIdpConnectionsAPI.DeleteConnection(acctest.TestBasicAuthContext(), spIdpConnectionConnectionId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func spIdpConnection_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.SpIdpConnectionsAPI.DeleteConnection(acctest.TestBasicAuthContext(), spIdpConnectionConnectionId).Execute()
	if err == nil {
		return fmt.Errorf("sp_idp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}