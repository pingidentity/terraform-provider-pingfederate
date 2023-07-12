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

const serverSettingsLogSettingsId = "id"

// Attributes to test with. Add optional properties to test here if desired.
type serverSettingsLogSettingsResourceModel struct {
	id                   string
	logCategoriesEnabled bool
}

func TestAccServerSettingsLogSettings(t *testing.T) {
	resourceName := "myServerSettingsLogSettings"
	initialResourceModel := serverSettingsLogSettingsResourceModel{
		id:                   serverSettingsLogSettingsId,
		logCategoriesEnabled: false,
	}
	updatedResourceModel := serverSettingsLogSettingsResourceModel{
		id:                   serverSettingsLogSettingsId,
		logCategoriesEnabled: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
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
				ImportStateId:     serverSettingsLogSettingsId,
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
      id          = "core"
      name        = "Core"
      description = "Debug logging for core components."
      enabled     = true
    },
    {
      id          = "policytree"
      name        = "Policy Tree"
      description = "Policy tree debug logging."
      enabled     = false
    },
    {
      id          = "trustedcas"
      name        = "Trusted CAs"
      description = "Log PingFederate and JRE trusted CAs when they are loaded."
      enabled     = true
    },
    {
      id          = "xmlsig"
      name        = "XML Signatures"
      description = "Debug logging for XML signature operations."
      enabled     = %[2]t
    },
    {
      id          = "requestheaders"
      name        = "HTTP Request Headers"
      description = "Log HTTP request headers. Sensitive information, such as passwords, may be logged when this category is enabled."
      enabled     = false
    },
    {
      id          = "requestparams"
      name        = "HTTP Request Parameters"
      description = "Log HTTP GET request parameters. Sensitive information, such as passwords, may be logged when this category is enabled."
      enabled     = true
    },
    {
      id          = "restdatastore"
      name        = "REST Data Store Requests and Responses"
      description = "Log REST datastore requests and responses. Sensitive information, such as passwords, may be logged when this category is enabled."
      enabled     = true
    },
  ]
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
		response, _, err := testClient.ServerSettingsApi.GetLogSettings(ctx).Execute()

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
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enabled", config.logCategoriesEnabled, *logCategoryEnabledVal)
		if err != nil {
			return err
		}

		return nil
	}
}
