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
type serverSettingsLogSettingsResourceModel struct {
	logCategoriesEnabled bool
}

func TestAccServerSettingsLogSettings(t *testing.T) {
	resourceName := "myServerSettingsLogSettings"
	initialResourceModel := serverSettingsLogSettingsResourceModel{
		logCategoriesEnabled: false,
	}
	updatedResourceModel := serverSettingsLogSettingsResourceModel{
		logCategoriesEnabled: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServerSettingsLogSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedServerSettingsLogSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccServerSettingsLogSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedServerSettingsLogSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccServerSettingsLogSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_server_settings_log_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccServerSettingsLogSettings(resourceName string, resourceModel serverSettingsLogSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_server_settings_log_settings" "%[1]s" {
  log_categories = [
    {
      id      = "policytree"
      enabled = false
    },
    {
      id      = "core"
      enabled = true
    },
    {
      id      = "trustedcas"
      enabled = true
    },
    {
      id      = "xmlsig"
      enabled = %[2]t
    },
    {
      id      = "requestheaders"
      enabled = false
    },
    {
      id      = "requestparams"
      enabled = true
    },
    {
      id      = "restdatastore"
      enabled = true
    },
  ]
}

data "pingfederate_server_settings_log_settings" "%[1]s" {
  depends_on = [pingfederate_server_settings_log_settings.%[1]s]
}`, resourceName,
		resourceModel.logCategoriesEnabled,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsLogSettingsAttributes(config serverSettingsLogSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettingsLogSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ServerSettingsAPI.GetLogSettings(ctx).Execute()

		if err != nil {
			return err
		}
		var logCategoryEnabledVal *bool
		logCategories := response.GetLogCategories()
		for i := 0; i < len(logCategories); i++ {
			logCategoryId := logCategories[i].Id
			if logCategoryId == "xmlsig" {
				logCategoryEnabledVal = logCategories[i].Enabled
			}
		}
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enabled", config.logCategoriesEnabled, *logCategoryEnabledVal)
		if err != nil {
			return err
		}

		return nil
	}
}
