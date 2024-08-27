package notificationpublisherssettings_test

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
type notificationPublisherSettingsResourceModel struct {
	defaultNotificationPublisherRefId string
}

func TestAccNotificationPublisherSettings(t *testing.T) {
	resourceName := "myNotificationPublisherSettings"
	initialResourceModel := notificationPublisherSettingsResourceModel{
		defaultNotificationPublisherRefId: "exampleSmtpPublisher",
	}
	updatedResourceModel := notificationPublisherSettingsResourceModel{
		defaultNotificationPublisherRefId: "exampleSmtpPublisher2",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationPublisherSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedNotificationPublisherSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccNotificationPublisherSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedNotificationPublisherSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccNotificationPublisherSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_notification_publisher_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test putting back the original values
			{
				Config: testAccNotificationPublisherSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedNotificationPublisherSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccNotificationPublisherSettings(resourceName string, resourceModel notificationPublisherSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_notification_publisher_settings" "%[1]s" {
  default_notification_publisher_ref = {
    id = "%[2]s"
  }
}`, resourceName,
		resourceModel.defaultNotificationPublisherRefId,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedNotificationPublisherSettingsAttributes(config notificationPublisherSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "NotificationPublisherSettings"
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
