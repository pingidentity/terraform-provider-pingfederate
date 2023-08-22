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

const oauthAuthServerSettingsScopesCommonScopesId = "exampleCommonScope"
const dynamicVal = false

// Attributes to test with. Add optional properties to test here if desired.
type oauthAuthServerSettingsScopesCommonScopesResourceModel struct {
	id          string
	name        string
	description string
	dynamic     bool
}

func TestAccOauthAuthServerSettingsScopesCommonScopes(t *testing.T) {
	resourceName := "myOauthAuthServerSettingsScopesCommonScopes"
	initialResourceModel := oauthAuthServerSettingsScopesCommonScopesResourceModel{
		id:          oauthAuthServerSettingsScopesCommonScopesId,
		name:        oauthAuthServerSettingsScopesCommonScopesId,
		description: "example",
		dynamic:     dynamicVal,
	}
	updatedResourceModel := oauthAuthServerSettingsScopesCommonScopesResourceModel{
		id:          oauthAuthServerSettingsScopesCommonScopesId,
		name:        oauthAuthServerSettingsScopesCommonScopesId,
		description: "updated description",
		dynamic:     dynamicVal,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckOauthAuthServerSettingsScopesCommonScopesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthAuthServerSettingsScopesCommonScopes(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsScopesCommonScopesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthAuthServerSettingsScopesCommonScopes(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsScopesCommonScopesAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthAuthServerSettingsScopesCommonScopes(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_auth_server_settings_scopes_common_scopes." + resourceName,
				ImportStateId:     oauthAuthServerSettingsScopesCommonScopesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOauthAuthServerSettingsScopesCommonScopes(resourceName string, resourceModel oauthAuthServerSettingsScopesCommonScopesResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_auth_server_settings_scopes_common_scopes" "%[1]s" {
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
func testAccCheckExpectedOauthAuthServerSettingsScopesCommonScopesAttributes(config oauthAuthServerSettingsScopesCommonScopesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAuthServerSettingsScopesCommonScopes"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAuthServerSettingsApi.GetCommonScope(ctx, oauthAuthServerSettingsScopesCommonScopesId).Execute()
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
func testAccCheckOauthAuthServerSettingsScopesCommonScopesDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthAuthServerSettingsApi.RemoveCommonScope(ctx, oauthAuthServerSettingsScopesCommonScopesId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAuthServerSettingsScopesCommonScopes", oauthAuthServerSettingsScopesCommonScopesId)
	}
	return nil
}
