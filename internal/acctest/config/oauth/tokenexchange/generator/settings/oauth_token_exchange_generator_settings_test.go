package oauthtokenexchangegeneratorsettings_test

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
type oauthTokenExchangeGeneratorSettingsResourceModel struct {
	defaultGeneratorGroupRefId string
}

func TestAccOauthTokenExchangeGeneratorSettings(t *testing.T) {
	resourceName := "myOauthTokenExchangeGeneratorSettings"
	initialResourceModel := oauthTokenExchangeGeneratorSettingsResourceModel{
		defaultGeneratorGroupRefId: "exampleGeneratorGroup",
	}

	updatedResourceModel := oauthTokenExchangeGeneratorSettingsResourceModel{
		defaultGeneratorGroupRefId: "exampleGeneratorGroup2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOauthTokenExchangeGeneratorSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeGeneratorSettingsAttributes(initialResourceModel),
			},
			{
				Config: testAccOauthTokenExchangeGeneratorSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeGeneratorSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthTokenExchangeGeneratorSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_token_exchange_generator_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOauthTokenExchangeGeneratorSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeGeneratorSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccOauthTokenExchangeGeneratorSettings(resourceName string, resourceModel oauthTokenExchangeGeneratorSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_token_exchange_generator_settings" "%[1]s" {
  default_generator_group_ref = {
    id = "%[2]s"
  }
}`, resourceName,
		resourceModel.defaultGeneratorGroupRefId,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthTokenExchangeGeneratorSettingsAttributes(config oauthTokenExchangeGeneratorSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthTokenExchangeGeneratorSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthTokenExchangeGeneratorAPI.GetOauthTokenExchangeSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "id", config.defaultGeneratorGroupRefId, response.DefaultGeneratorGroupRef.Id)
		if err != nil {
			return err
		}

		return nil
	}
}
