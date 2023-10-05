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

const oauthAuthServerSettingsScopesExclusiveScopesId = "exampleExclusiveScope"
const dynamicVal = false

// Attributes to test with. Add optional properties to test here if desired.
type oauthAuthServerSettingsScopesExclusiveScopesResourceModel struct {
	id          string
	name        string
	description string
	dynamic     bool
}

func TestAccOauthAuthServerSettingsScopesExclusiveScopes(t *testing.T) {
	resourceName := "myOauthAuthServerSettingsScopesExclusiveScopes"
	initialResourceModel := oauthAuthServerSettingsScopesExclusiveScopesResourceModel{
		id:          oauthAuthServerSettingsScopesExclusiveScopesId,
		name:        oauthAuthServerSettingsScopesExclusiveScopesId,
		description: "example",
		dynamic:     dynamicVal,
	}
	updatedResourceModel := oauthAuthServerSettingsScopesExclusiveScopesResourceModel{
		id:          oauthAuthServerSettingsScopesExclusiveScopesId,
		name:        oauthAuthServerSettingsScopesExclusiveScopesId,
		description: "updated description",
		dynamic:     dynamicVal,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOauthAuthServerSettingsScopesExclusiveScopesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsScopesExclusiveScopesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsScopesExclusiveScopesAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_auth_server_settings_scopes_exclusive_scope." + resourceName,
				ImportStateId:     oauthAuthServerSettingsScopesExclusiveScopesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName string, resourceModel oauthAuthServerSettingsScopesExclusiveScopesResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_auth_server_settings_scopes_exclusive_scope" "%[1]s" {
  dynamic     = %[2]t
  description = "%[3]s"
  name        = "%[4]s"
}`, resourceName,
		resourceModel.dynamic,
		resourceModel.description,
		resourceModel.name,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthAuthServerSettingsScopesExclusiveScopesAttributes(config oauthAuthServerSettingsScopesExclusiveScopesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAuthServerSettingsScopesExclusiveScopes"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAuthServerSettingsAPI.GetExclusiveScope(ctx, oauthAuthServerSettingsScopesExclusiveScopesId).Execute()

		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, &config.id, "description", config.description, response.Description)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.name, response.Name)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOauthAuthServerSettingsScopesExclusiveScopesDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthAuthServerSettingsAPI.RemoveExclusiveScope(ctx, oauthAuthServerSettingsScopesExclusiveScopesId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAuthServerSettingsScopesExclusiveScopes", oauthAuthServerSettingsScopesExclusiveScopesId)
	}
	return nil
}
