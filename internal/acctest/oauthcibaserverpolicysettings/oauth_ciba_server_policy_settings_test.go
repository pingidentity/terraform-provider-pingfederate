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

const oauthCibaServerPolicySettingsId = "exampleCibaReqPolicy"

// Attributes to test with. Add optional properties to test here if desired.
type oauthCibaServerPolicySettingsResourceModel struct {
	defaultRequestPolicyRefId string
}

func TestAccOauthCibaServerPolicySettings(t *testing.T) {
	resourceName := "myOauthCibaServerPolicySettings"
	initialResourceModel := oauthCibaServerPolicySettingsResourceModel{
		defaultRequestPolicyRefId: oauthCibaServerPolicySettingsId,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOauthCibaServerPolicySettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthCibaServerPolicySettingsAttributes(initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthCibaServerPolicySettings(resourceName, initialResourceModel),
				ResourceName:      "pingfederate_oauth_ciba_server_policy_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOauthCibaServerPolicySettings(resourceName string, resourceModel oauthCibaServerPolicySettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_ciba_server_policy_settings" "%[1]s" {
  default_request_policy_ref = {
    id = "%[2]s"
  }
}`, resourceName,
		resourceModel.defaultRequestPolicyRefId,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthCibaServerPolicySettingsAttributes(config oauthCibaServerPolicySettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthCibaServerPolicySettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthCibaServerPolicyAPI.GetCibaServerPolicySettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "id", config.defaultRequestPolicyRefId, response.DefaultRequestPolicyRef.Id)
		if err != nil {
			return err
		}

		return nil
	}
}
