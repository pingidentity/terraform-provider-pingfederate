package acctest_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

var pingOneConnectionId = "myPingOneConnectionId"
var pingOneConnectionName = "myPingOneConnectionName"
var credentialData = os.Getenv("PF_TF_ACC_TEST_PING_ONE_CONNECTION_CREDENTIAL_DATA")
var pingOneEnvironmentId = os.Getenv("PF_TF_P1_CONNECTION_ENV_ID")

// Attributes to test with. Add optional properties to test here if desired.
type pingOneConnectionResourceModel struct {
	name        string
	description string
	active      bool
	credential  string
}

func TestAccPingOneConnection(t *testing.T) {
	resourceName := "myPingOneConnection"
	initialResourceModel := pingOneConnectionResourceModel{
		name:       pingOneConnectionName,
		credential: credentialData,
	}
	updatedResourceModel := pingOneConnectionResourceModel{
		name:        pingOneConnectionName,
		credential:  credentialData,
		description: "Updated description",
		active:      false,
	}

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
		CheckDestroy: testAccCheckPingOneConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOneConnection(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPingOneConnectionAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPingOneConnection(resourceName, updatedResourceModel),
				Check:  resource.ComposeTestCheckFunc(testAccCheckExpectedPingOneConnectionAttributes(updatedResourceModel)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(fmt.Sprintf("pingfederate_ping_one_connection.%s", resourceName), tfjsonpath.New("creation_date"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(fmt.Sprintf("pingfederate_ping_one_connection.%s", resourceName), tfjsonpath.New("credential_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(fmt.Sprintf("pingfederate_ping_one_connection.%s", resourceName), tfjsonpath.New("environment_id"), knownvalue.StringExact(pingOneEnvironmentId)),
					statecheck.ExpectKnownValue(fmt.Sprintf("pingfederate_ping_one_connection.%s", resourceName), tfjsonpath.New("region"), knownvalue.StringExact("North America")),
					statecheck.ExpectKnownValue(fmt.Sprintf("pingfederate_ping_one_connection.%s", resourceName), tfjsonpath.New("ping_one_authentication_api_endpoint"), knownvalue.StringExact("https://auth.pingone.com")),
					statecheck.ExpectKnownValue(fmt.Sprintf("pingfederate_ping_one_connection.%s", resourceName), tfjsonpath.New("ping_one_management_api_endpoint"), knownvalue.StringExact("https://api.pingone.com")),
				},
			},
			{
				// Test importing the resource
				Config:                  testAccPingOneConnection(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_ping_one_connection." + resourceName,
				ImportStateId:           pingOneConnectionId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential"},
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.PingOneConnectionsAPI.DeletePingOneConnection(ctx, pingOneConnectionId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccPingOneConnection(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPingOneConnectionAttributes(initialResourceModel),
			},
		},
	})
}

func optionalHcl(model pingOneConnectionResourceModel) string {
	if model.description != "" && model.active {
		return fmt.Sprintf(`
		description = %[1]s
		active = %[2]t
	`, model.description, model.active)
	}
	return ""
}

func testAccPingOneConnection(resourceName string, resourceModel pingOneConnectionResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_ping_one_connection" "%[1]s" {
  connection_id = "%[2]s"
  name          = "%[3]s"
  credential    = "%[4]s"
	%[5]s
}`, resourceName,
		pingOneConnectionId,
		resourceModel.name,
		resourceModel.credential,
		optionalHcl(resourceModel),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedPingOneConnectionAttributes(config pingOneConnectionResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "PingOneConnection"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PingOneConnectionsAPI.GetPingOneConnection(ctx, pingOneConnectionId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "name", config.name, response.Name)
		if err != nil {
			return err
		}

		if config.description != "" && response.Description != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "description", config.description, *response.Description)
			if err != nil {
				return err
			}
		}

		if config.active {
			err = acctest.TestAttributesMatchBool(resourceType, nil, "active", config.active, *response.Active)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPingOneConnectionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.PingOneConnectionsAPI.DeletePingOneConnection(ctx, pingOneConnectionId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("PingOneConnection", pingOneConnectionId)
	}
	return nil
}
