// Copyright © 2025 Ping Identity Corporation

package notificationpublisherssettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccNotificationPublisherSettings(t *testing.T) {
	resourceName := acctest.ResourceIdGen()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationPublisherSettingsEmpty(resourceName),
			},
			{
				// Test updating some fields
				Config: testAccNotificationPublisherSettingsWithDefault(resourceName),
			},
			{
				// Test importing the resource
				Config:                               testAccNotificationPublisherSettingsWithDefault(resourceName),
				ResourceName:                         "pingfederate_notification_publisher_settings." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "default_notification_publisher_ref.id",
			},
			// Test putting back the original values
			{
				Config: testAccNotificationPublisherSettingsEmpty(resourceName),
			},
		},
	})
}

func testAccNotificationPublisherSettingsEmpty(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_notification_publisher_settings" "%[1]s" {
}`, resourceName)
}

func testAccNotificationPublisherSettingsWithDefault(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_notification_publisher" "%[1]sPub" {
  configuration = {
    fields = [
      {
        name  = "Connection Timeout"
        value = "30"
      },
      {
        name  = "Email Server"
        value = "example.com"
      },
      {
        name  = "Enable SMTP Debugging Messages"
        value = "false"
      },
      {
        name  = "Encryption Method"
        value = "NONE"
      },
      {
        name  = "From Address"
        value = "example@pingidentity.com"
      },
      {
        name  = "SMTP Port"
        value = "25"
      },
      {
        name  = "SMTPS Port"
        value = "465"
      },
      {
        name  = "UTF-8 Message Header Support"
        value = "false"
      },
      {
        name  = "Verify Hostname"
        value = "true"
      },
    ]
  }
  name = "%[1]sPub"
  plugin_descriptor_ref = {
    id = "com.pingidentity.email.SmtpNotificationPlugin"
  }
  publisher_id = "%[1]sPub"
}

resource "pingfederate_notification_publisher_settings" "%[1]s" {
  default_notification_publisher_ref = {
    id = pingfederate_notification_publisher.%[1]sPub.id
  }
}`, resourceName)
}
