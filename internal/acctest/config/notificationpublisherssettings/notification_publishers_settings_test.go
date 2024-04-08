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
type notificationPublishersSettingsResourceModel struct {
	defaultNotificationPublisherRefId string
}

func TestAccNotificationPublishersSettings(t *testing.T) {
	resourceName := "myNotificationPublishersSettings"
	initialResourceModel := notificationPublishersSettingsResourceModel{
		defaultNotificationPublisherRefId: "exampleSmtpPublisher",
	}
	updatedResourceModel := notificationPublishersSettingsResourceModel{
		defaultNotificationPublisherRefId: "exampleSmtpPublisher2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationPublishersSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedNotificationPublishersSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccNotificationPublishersSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedNotificationPublishersSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccNotificationPublishersSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_notification_publishers_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test putting back the original values
			{
				Config: testAccNotificationPublishersSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedNotificationPublishersSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccNotificationPublishersSettings(resourceName string, resourceModel notificationPublishersSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_notification_publishers_settings" "%[1]s" {
  default_notification_publisher_ref = {
    id = "%[2]s"
  }
}`, resourceName,
		resourceModel.defaultNotificationPublisherRefId,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedNotificationPublishersSettingsAttributes(config notificationPublishersSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "NotificationPublishersSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.NotificationPublishersAPI.GetNotificationPublishersSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "id", config.defaultNotificationPublisherRefId, response.DefaultNotificationPublisherRef.Id)
		if err != nil {
			return err
		}

		return nil
	}
}
