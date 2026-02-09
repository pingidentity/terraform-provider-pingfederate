// Copyright © 2026 Ping Identity Corporation

package notificationpublisherssettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccNotificationPublisherSettings(t *testing.T) {
	resourceName := acctest.ResourceIdGen()

	var steps []resource.TestStep
	if acctest.VersionAtLeast(version.PingFederate1200) {
		steps = testAccNotificationPublisherSettingsPf120(resourceName)
	} else {
		steps = testAccNotificationPublisherSettingsPrePf120(resourceName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: steps,
	})
}

// Prior to PF 12.0 it isn't possible to delete the final notification publisher from the server config,
// because it is always in use.
func testAccNotificationPublisherSettingsPrePf120(resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			// Set to the existing default
			Config: testAccNotificationPublisherSettingsExistingDefault(resourceName, "acctestNotificationPublisher"),
		},
		{
			// Test importing the resource
			Config:                               testAccNotificationPublisherSettingsExistingDefault(resourceName, "acctestNotificationPublisher"),
			ResourceName:                         "pingfederate_notification_publisher_settings." + resourceName,
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "default_notification_publisher_ref.id",
		},
	}
}

func testAccNotificationPublisherSettingsPf120(resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			// No policies configured and no default
			Config: testAccNotificationPublisherSettingsEmpty(resourceName),
		},
		{
			// Set a default policy
			Config: testAccNotificationPublisherSettingsBuildDefault(resourceName),
		},
		{
			// Test importing the resource
			Config:                               testAccNotificationPublisherSettingsBuildDefault(resourceName),
			ResourceName:                         "pingfederate_notification_publisher_settings." + resourceName,
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "default_notification_publisher_ref.id",
		},
		{
			// Reset back to no policies
			Config: testAccNotificationPublisherSettingsEmpty(resourceName),
		},
	}
}

func testAccNotificationPublisherSettingsEmpty(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_notification_publisher_settings" "%[1]s" {
}`, resourceName)
}

func testAccNotificationPublisherSettingsBuildDefault(resourceName string) string {
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

func testAccNotificationPublisherSettingsExistingDefault(resourceName, defaultPublisherId string) string {
	return fmt.Sprintf(`
resource "pingfederate_notification_publisher_settings" "%[1]s" {
  default_notification_publisher_ref = {
    id = "%[2]s"
  }
}`, resourceName, defaultPublisherId)
}
