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

// Attributes to test with. Add optional properties to test here if desired.
type sessionApplicationSessionPolicyResourceModel struct {
	idleTimeoutMins int64
	maxTimeoutMins  int64
}

func TestAccSessionApplicationSessionPolicy(t *testing.T) {
	resourceName := "mySessionApplicationSessionPolicy"
	initialResourceModel := sessionApplicationSessionPolicyResourceModel{
		idleTimeoutMins: 60,
		maxTimeoutMins:  90,
	}
	updatedResourceModel := sessionApplicationSessionPolicyResourceModel{
		idleTimeoutMins: -1,
		maxTimeoutMins:  60,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSessionApplicationSessionPolicy(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSessionApplicationSessionPolicyAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSessionApplicationSessionPolicy(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSessionApplicationSessionPolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSessionApplicationSessionPolicy(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_session_application_session_policy." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSessionApplicationSessionPolicy(resourceName string, resourceModel sessionApplicationSessionPolicyResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_session_application_session_policy" "%[1]s" {
  idle_timeout_mins = %[2]d
  max_timeout_mins  = %[3]d
}`, resourceName,
		resourceModel.idleTimeoutMins,
		resourceModel.maxTimeoutMins,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSessionApplicationSessionPolicyAttributes(config sessionApplicationSessionPolicyResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SessionApplicationSessionPolicy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SessionAPI.GetApplicationPolicy(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, nil, "idle_timeout_mins",
			config.idleTimeoutMins, *response.IdleTimeoutMins)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, nil, "max_timeout_mins",
			config.maxTimeoutMins, *response.MaxTimeoutMins)
		if err != nil {
			return err
		}

		return nil
	}
}
