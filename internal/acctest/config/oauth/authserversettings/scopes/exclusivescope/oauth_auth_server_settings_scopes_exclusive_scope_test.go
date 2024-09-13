package oauthauthserversettingsscopesexclusivescope_test

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

const oauthAuthServerSettingsScopesExclusiveScopesId = "*exampleExclusiveScope"

// Attributes to test with. Add optional properties to test here if desired.
type oauthAuthServerSettingsScopesExclusiveScopesResourceModel struct {
	id          string
	name        string
	description string
	dynamic     bool
}

func TestAccOauthAuthServerSettingsScopesExclusiveScopes(t *testing.T) {
	t.SkipNow()
	resourceName := "myOauthAuthServerSettingsScopesExclusiveScopes"
	initialResourceModel := oauthAuthServerSettingsScopesExclusiveScopesResourceModel{
		id:          oauthAuthServerSettingsScopesExclusiveScopesId,
		name:        oauthAuthServerSettingsScopesExclusiveScopesId,
		description: "example",
	}
	updatedResourceModel := oauthAuthServerSettingsScopesExclusiveScopesResourceModel{
		id:          oauthAuthServerSettingsScopesExclusiveScopesId,
		name:        oauthAuthServerSettingsScopesExclusiveScopesId,
		description: "updated description",
		dynamic:     true,
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedOauthAuthServerSettingsScopesExclusiveScopesAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_oauth_auth_server_settings_scopes_exclusive_scope.%s", resourceName), "dynamic", fmt.Sprintf("%t", updatedResourceModel.dynamic)),
				),
			},
			{
				// Test importing the resource
				Config:            testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_auth_server_settings_scopes_exclusive_scope." + resourceName,
				ImportStateId:     oauthAuthServerSettingsScopesExclusiveScopesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsScopesExclusiveScopesAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OauthAuthServerSettingsAPI.RemoveExclusiveScope(ctx, updatedResourceModel.id).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAuthServerSettingsScopesExclusiveScopesAttributes(initialResourceModel),
			},
		},
	})
}

func testAccOauthAuthServerSettingsScopesExclusiveScopes(resourceName string, resourceModel oauthAuthServerSettingsScopesExclusiveScopesResourceModel) string {
	dynamicHcl := ""
	// Leave off dynamic if false to test not including it, since it is optional
	if resourceModel.dynamic {
		dynamicHcl = "dynamic = true"
	}
	return fmt.Sprintf(`
resource "pingfederate_oauth_auth_server_settings_scopes_exclusive_scope" "%[1]s" {
  description = "%[2]s"
  name        = "%[3]s"
  %[4]s
}
data "pingfederate_oauth_auth_server_settings_scopes_exclusive_scope" "%[1]s" {
  name = pingfederate_oauth_auth_server_settings_scopes_exclusive_scope.%[1]s.name
}`, resourceName,
		resourceModel.description,
		resourceModel.name,
		dynamicHcl,
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

		if config.dynamic {
			err = acctest.TestAttributesMatchBool(resourceType, &config.id, "dynamic", config.dynamic, *response.Dynamic)
			if err != nil {
				return err
			}
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
