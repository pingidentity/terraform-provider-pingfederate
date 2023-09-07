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

const serverSettingsId = "2"

// Attributes to test with. Add optional properties to test here if desired.
// The serverSettingsResourceModel struct represents a model for server settings resources.
// It defines the fields that can be used to configure various aspects of the server settings.
type serverSettingsResourceModel struct {
	id                string
	contactInfo       string
	notifications     bool
	rolesAndProtocols []string
	federationInfo    map[string]string
	emailServer       EmailServer
	captchaSettings   CaptchaSettings
}

func TestAccServerSettings(t *testing.T) {
	resourceName := "myServerSettings"
	initialResourceModel := serverSettingsResourceModel{
		contactInfo: fill in test value,	
		notifications: fill in test value,	
		rolesAndProtocols: fill in test value,	
		federationInfo: fill in test value,	
		emailServer: fill in test value,	
		captchaSettings: fill in test value,
	}
	updatedResourceModel := serverSettingsResourceModel{
		contactInfo: fill in test value,	
		notifications: fill in test value,	
		rolesAndProtocols: fill in test value,	
		federationInfo: fill in test value,	
		emailServer: fill in test value,	
		captchaSettings: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckServerSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedServerSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccServerSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedServerSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccServerSettings(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_server_settings." + resourceName,
				ImportStateId:           serverSettingsId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testAccServerSettings(resourceName string, resourceModel serverSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_server_settings" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsAttributes(config serverSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetServerSettings(ctx, serverSettingsId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckServerSettingsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeleteServerSettings(ctx, serverSettingsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("ServerSettings", serverSettingsId)
	}
	return nil
}
