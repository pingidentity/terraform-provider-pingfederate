package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccServerSettingsLogSettings(t *testing.T) {
	resourceName := "myServerSettingsLogSettings"
	//TODO when the plugin framework fixes issues with Set plans, we can test this resource with a
	// minimal model. For now, just testing setting all values. See the schema
	// in server_settings_log_settings_resource.go for details.
	logCategoriesEnabledInitial := pointers.Bool(true)
	logCategoriesEnabled := pointers.Bool(true)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServerSettingsLogSettings(resourceName, logCategoriesEnabledInitial),
				Check:  testAccCheckExpectedServerSettingsLogSettingsAttributes(logCategoriesEnabledInitial),
			},
			{
				// Test updating some fields
				Config: testAccServerSettingsLogSettings(resourceName, logCategoriesEnabled),
				Check:  testAccCheckExpectedServerSettingsLogSettingsAttributes(logCategoriesEnabled),
			},
			{
				// Test importing the resource
				Config:            testAccServerSettingsLogSettings(resourceName, logCategoriesEnabled),
				ResourceName:      "pingfederate_server_settings_log_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccServerSettingsLogSettings(resourceName, logCategoriesEnabledInitial),
				Check:  testAccCheckExpectedServerSettingsLogSettingsAttributes(logCategoriesEnabledInitial),
			},
		},
	})
}

func testAccServerSettingsLogSettings(resourceName string, logCategoriesEnabled *bool) string {
	logCategoriesHcl := ""
	if logCategoriesEnabled != nil {
		logCategoriesHcl = fmt.Sprintf(`
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
      enabled = %t
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
`, *logCategoriesEnabled)
		if acctest.VersionAtLeast(version.PingFederate1200) {
			logCategoriesHcl += `
		{
			id = "protocolrequestresponse"
			enabled = false
		},
			`
		}
		logCategoriesHcl += `
	]
	`
	}

	return fmt.Sprintf(`
resource "pingfederate_server_settings_log_settings" "%s" {
	%s
}

data "pingfederate_server_settings_log_settings" "%[1]s" {
  depends_on = [pingfederate_server_settings_log_settings.%[1]s]
}`, resourceName,
		logCategoriesHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsLogSettingsAttributes(logCategoriesEnabled *bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettingsLogSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ServerSettingsAPI.GetLogSettings(ctx).Execute()

		if err != nil {
			return err
		}

		if logCategoriesEnabled == nil {
			return nil
		}

		var logCategoryEnabledVal *bool
		logCategories := response.GetLogCategories()
		for i := 0; i < len(logCategories); i++ {
			logCategoryId := logCategories[i].Id
			if logCategoryId == "xmlsig" {
				logCategoryEnabledVal = logCategories[i].Enabled
			}
		}
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enabled", *logCategoriesEnabled, *logCategoryEnabledVal)
		if err != nil {
			return err
		}

		return nil
	}
}
