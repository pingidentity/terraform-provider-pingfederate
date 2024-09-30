package sessionapplicationsessionpolicy_test

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
type sessionApplicationPolicyResourceModel struct {
	idleTimeoutMins int64
	maxTimeoutMins  int64
}

func TestAccSessionApplicationPolicy(t *testing.T) {
	resourceName := "mySessionApplicationPolicy"
	updatedResourceModel := sessionApplicationPolicyResourceModel{
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
				Config: testAccSessionApplicationPolicy(resourceName, nil),
				Check:  testAccCheckExpectedSessionApplicationPolicyAttributes(nil),
			},
			{
				// Test updating some fields
				Config: testAccSessionApplicationPolicy(resourceName, &updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSessionApplicationPolicyAttributes(&updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_application_policy.%s", resourceName), "idle_timeout_mins", fmt.Sprintf("%d", updatedResourceModel.idleTimeoutMins)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_application_policy.%s", resourceName), "max_timeout_mins", fmt.Sprintf("%d", updatedResourceModel.maxTimeoutMins)),
				),
			},
			{
				// Test importing the resource
				Config:                               testAccSessionApplicationPolicy(resourceName, &updatedResourceModel),
				ResourceName:                         "pingfederate_session_application_policy." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "idle_timeout_mins",
			},
			{
				// Back to minimal model
				Config: testAccSessionApplicationPolicy(resourceName, nil),
				Check:  testAccCheckExpectedSessionApplicationPolicyAttributes(nil),
			},
		},
	})
}

func testAccSessionApplicationPolicy(resourceName string, resourceModel *sessionApplicationPolicyResourceModel) string {
	optionalHcl := ""
	if resourceModel != nil {
		optionalHcl = fmt.Sprintf(`
  idle_timeout_mins = %d
  max_timeout_mins  = %d
`,
			resourceModel.idleTimeoutMins,
			resourceModel.maxTimeoutMins,
		)
	}
	return fmt.Sprintf(`
resource "pingfederate_session_application_policy" "%s" {
  %s
}
data "pingfederate_session_application_policy" "%[1]s" {
  depends_on = [pingfederate_session_application_policy.%[1]s]
}`,
		resourceName,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSessionApplicationPolicyAttributes(config *sessionApplicationPolicyResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SessionApplicationPolicy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SessionAPI.GetApplicationPolicy(ctx).Execute()

		if err != nil {
			return err
		}

		if config == nil {
			return nil
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
