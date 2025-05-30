// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

package pingoneconnection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const pingoneConnectionConnectionId = "pingoneConnectionConnectionId"

var credentialData = os.Getenv("PF_TF_ACC_TEST_PING_ONE_CONNECTION_CREDENTIAL_DATA")
var pingOneEnvironmentId = os.Getenv("PF_TF_P1_CONNECTION_ENV_ID")

func TestAccPingoneConnection_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if credentialData == "" {
				t.Fatal("PF_TF_ACC_TEST_PING_ONE_CONNECTION_CREDENTIAL_DATA must be set for acceptance tests")
			}
			if pingOneEnvironmentId == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ENV_ID must be set for acceptance tests")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: pingoneConnection_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: pingoneConnection_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					pingoneConnection_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPingoneConnection_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if credentialData == "" {
				t.Fatal("PF_TF_ACC_TEST_PING_ONE_CONNECTION_CREDENTIAL_DATA must be set for acceptance tests")
			}
			if pingOneEnvironmentId == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ENV_ID must be set for acceptance tests")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: pingoneConnection_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: pingoneConnection_MinimalHCL(),
				Check:  pingoneConnection_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  pingoneConnection_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: pingoneConnection_CompleteHCL(),
				Check:  pingoneConnection_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: pingoneConnection_MinimalHCL(),
				Check:  pingoneConnection_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: pingoneConnection_CompleteHCL(),
				Check:  pingoneConnection_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               pingoneConnection_CompleteHCL(),
				ResourceName:                         "pingfederate_pingone_connection.example",
				ImportStateId:                        pingoneConnectionConnectionId,
				ImportStateVerifyIdentifierAttribute: "connection_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// Credential will not be returned by the API
				ImportStateVerifyIgnore: []string{"credential"},
			},
		},
	})
}

// Minimal HCL with only required values set
func pingoneConnection_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_pingone_connection" "example" {
  connection_id = "%s"
  name          = "myconn"
  credential    = "%s"
}
`, pingoneConnectionConnectionId, credentialData)
}

// Maximal HCL with all values set where possible
func pingoneConnection_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_pingone_connection" "example" {
  connection_id = "%s"
  name          = "myconn"
  credential    = "%s"
  active        = false
  description   = "my conn desc"
}
`, pingoneConnectionConnectionId, credentialData)
}

// Validate any computed values when applying minimal HCL
func pingoneConnection_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "active", "true"),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "creation_date"),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "credential_id"),
		resource.TestCheckNoResourceAttr("pingfederate_pingone_connection.example", "description"),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "encrypted_credential"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "environment_id", pingOneEnvironmentId),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "id", pingoneConnectionConnectionId),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "organization_name"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "ping_one_authentication_api_endpoint", "https://auth.pingone.com"),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "ping_one_connection_id"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "ping_one_management_api_endpoint", "https://api.pingone.com"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "region", "North America"),
	)
}

// Validate any computed values when applying complete HCL
func pingoneConnection_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "creation_date"),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "credential_id"),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "encrypted_credential"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "environment_id", pingOneEnvironmentId),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "id", pingoneConnectionConnectionId),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "organization_name"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "ping_one_authentication_api_endpoint", "https://auth.pingone.com"),
		resource.TestCheckResourceAttrSet("pingfederate_pingone_connection.example", "ping_one_connection_id"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "ping_one_management_api_endpoint", "https://api.pingone.com"),
		resource.TestCheckResourceAttr("pingfederate_pingone_connection.example", "region", "North America"),
	)
}

// Delete the resource
func pingoneConnection_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.PingOneConnectionsAPI.DeletePingOneConnection(acctest.TestBasicAuthContext(), pingoneConnectionConnectionId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func pingoneConnection_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.PingOneConnectionsAPI.DeletePingOneConnection(acctest.TestBasicAuthContext(), pingoneConnectionConnectionId).Execute()
	if err == nil {
		return fmt.Errorf("pingone_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
