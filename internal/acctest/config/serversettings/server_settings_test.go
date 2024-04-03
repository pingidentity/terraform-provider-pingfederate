package serversettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

type serverSettingsResourceModel struct {
	contactInfo     *client.ContactInfo
	federationInfo  *client.FederationInfo
	emailServer     *client.EmailServerSettings
	captchaSettings *client.CaptchaSettings
	notifications   *client.NotificationSettings
}

func TestAccServerSettings(t *testing.T) {
	resourceName := "myServerSettings"
	initialResourceModel := serverSettingsResourceModel{
		federationInfo: &client.FederationInfo{
			BaseUrl:       pointers.String("https://localhost:9999"),
			Saml2EntityId: pointers.String("initial.pingidentity.com"),
		},
	}

	updatedResourceModel := serverSettingsResourceModel{
		federationInfo: &client.FederationInfo{
			BaseUrl:        pointers.String("https://localhost2:9999"),
			Saml2EntityId:  pointers.String("updated.pingidentity.com"),
			Saml1xIssuerId: pointers.String("example.com"),
			WsfedRealm:     pointers.String("myrealm"),
		},
		contactInfo: &client.ContactInfo{
			Company:   pointers.String("updated company"),
			Email:     pointers.String("updatedAdminemail@example.com"),
			FirstName: pointers.String("Jane2"),
			LastName:  pointers.String("Admin2"),
			Phone:     pointers.String("555-555-2222"),
		},
		emailServer: &client.EmailServerSettings{
			SourceAddr:               "updatedEmailServerAdmin@example.com",
			EmailServer:              "updatedemailserver.example.com",
			UseSSL:                   pointers.Bool(true),
			VerifyHostname:           pointers.Bool(true),
			EnableUtf8MessageHeaders: pointers.Bool(true),
			UseDebugging:             pointers.Bool(false),
			Username:                 pointers.String("EmailServerAdmin"),
			Password:                 pointers.String("EmailServerAdminPassword"),
		},
		notifications: &client.NotificationSettings{
			LicenseEvents: &client.LicenseEventNotificationSettings{
				EmailAddress: "license-events-email@example.com",
				NotificationPublisherRef: &client.ResourceLink{
					Id: "exampleSmtpPublisher",
				},
			},
			CertificateExpirations: &client.CertificateExpirationNotificationSettings{
				EmailAddress:         "example@example.com",
				InitialWarningPeriod: pointers.Int64(45),
				FinalWarningPeriod:   7,
				NotificationPublisherRef: &client.ResourceLink{
					Id: "exampleSmtpPublisher",
				},
			},
			NotifyAdminUserPasswordChanges: pointers.Bool(true),
			AccountChangesNotificationPublisherRef: &client.ResourceLink{
				Id: "exampleSmtpPublisher",
			},
			MetadataNotificationSettings: &client.MetadataEventNotificationSettings{
				EmailAddress: "metadata-notification@example.com",
				NotificationPublisherRef: &client.ResourceLink{
					Id: "exampleSmtpPublisher",
				},
			},
		},
		captchaSettings: &client.CaptchaSettings{
			SiteKey:   pointers.String("exampleCaptchaProviderV2"),
			SecretKey: pointers.String("2FederateM0re"),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServerSettingsMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedServerSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccServerSettingsComplete(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedServerSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccServerSettingsComplete(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_server_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Email server details are not returned by PF
				ImportStateVerifyIgnore: []string{"email_server", "captcha_settings"},
			},
			{
				// Back to minimal model
				Config: testAccServerSettingsMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedServerSettingsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccServerSettingsMinimal(resourceName string, resourceModel serverSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_server_settings" "%s" {
  federation_info = {
    base_url         = "%s"
    saml_2_entity_id = "%s"
  }
}`, resourceName,
		*resourceModel.federationInfo.BaseUrl,
		*resourceModel.federationInfo.Saml2EntityId,
	)
}

func testAccServerSettingsComplete(resourceName string, resourceModel serverSettingsResourceModel) string {
	versionedHcl := ""
	if acctest.VersionAtLeast(version.PingFederate1130) {
		versionedHcl = `
      notification_mode = "NOTIFICATION_PUBLISHER"
		`
	}
	if acctest.VersionAtLeast(version.PingFederate1200) {
		versionedHcl = `
	  expired_certificate_administrative_console_warning_days = 10
	  expiring_certificate_administrative_console_warning_days = 11
	  thread_pool_exhaustion_notification_settings = {
		thread_dump_enabled = false
		notification_publisher_ref = {
			id = "exampleSmtpPublisher"
		}
		notification_mode = "LOGGING_ONLY"
	  }
		`
	}
	return fmt.Sprintf(`
resource "pingfederate_server_settings" "%[1]s" {
  contact_info = {
    company    = "%[2]s"
    email      = "%[3]s"
    first_name = "%[4]s"
    last_name  = "%[5]s"
    phone      = "%[6]s"
  }

  federation_info = {
    base_url          = "%[7]s"
    saml_2_entity_id  = "%[8]s"
    saml_1x_issuer_id = "%[9]s"
    wsfed_realm       = "%[10]s"
  }

  email_server = {
    source_addr                 = "%[11]s"
    email_server                = "%[12]s"
    use_ssl                     = %[13]t
    verify_hostname             = %[14]t
    enable_utf8_message_headers = %[15]t
    use_debugging               = %[16]t
    username                    = "%[17]s"
    password                    = "%[18]s"
  }

  notifications = {
    license_events = {
      email_address = "%[19]s"
      notification_publisher_ref = {
        id = "%[20]s"
      }
    }
    certificate_expirations = {
      initial_warning_period = %[21]d
      final_warning_period   = %[22]d
      notification_publisher_ref = {
        id = "%[23]s"
      }
      email_address = "%[24]s"
	  %[25]s
    }
    notify_admin_user_password_changes = %[26]t
    account_changes_notification_publisher_ref = {
      id = "%[27]s"
    }
    metadata_notification_settings = {
      email_address = "%[28]s"
      notification_publisher_ref = {
        id = "%[29]s"
      }
    }
  }

  captcha_settings = {
    site_key   = "%[30]s"
    secret_key = "%[31]s"
  }
}
data "pingfederate_server_settings" "%[1]s" {
  depends_on = [pingfederate_server_settings.%[1]s]
}`, resourceName,
		*resourceModel.contactInfo.Company,
		*resourceModel.contactInfo.Email,
		*resourceModel.contactInfo.FirstName,
		*resourceModel.contactInfo.LastName,
		*resourceModel.contactInfo.Phone,
		*resourceModel.federationInfo.BaseUrl,
		*resourceModel.federationInfo.Saml2EntityId,
		*resourceModel.federationInfo.Saml1xIssuerId,
		*resourceModel.federationInfo.WsfedRealm,
		resourceModel.emailServer.SourceAddr,
		resourceModel.emailServer.EmailServer,
		*resourceModel.emailServer.UseSSL,
		*resourceModel.emailServer.VerifyHostname,
		*resourceModel.emailServer.EnableUtf8MessageHeaders,
		*resourceModel.emailServer.UseDebugging,
		*resourceModel.emailServer.Username,
		*resourceModel.emailServer.Password,
		resourceModel.notifications.LicenseEvents.EmailAddress,
		resourceModel.notifications.LicenseEvents.NotificationPublisherRef.Id,
		*resourceModel.notifications.CertificateExpirations.InitialWarningPeriod,
		resourceModel.notifications.CertificateExpirations.FinalWarningPeriod,
		resourceModel.notifications.CertificateExpirations.NotificationPublisherRef.Id,
		resourceModel.notifications.CertificateExpirations.EmailAddress,
		versionedHcl,
		*resourceModel.notifications.NotifyAdminUserPasswordChanges,
		resourceModel.notifications.AccountChangesNotificationPublisherRef.Id,
		resourceModel.notifications.MetadataNotificationSettings.EmailAddress,
		resourceModel.notifications.MetadataNotificationSettings.NotificationPublisherRef.Id,
		*resourceModel.captchaSettings.SiteKey,
		*resourceModel.captchaSettings.SecretKey,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsAttributes(config serverSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		stateAttributes := s.RootModule().Resources["pingfederate_server_settings.myServerSettings"].Primary.Attributes
		response, _, err := testClient.ServerSettingsAPI.GetServerSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "base_url",
			*config.federationInfo.BaseUrl, *response.FederationInfo.BaseUrl)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "federation_info.base_url", *config.federationInfo.BaseUrl)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "saml_2_entity_id",
			*config.federationInfo.Saml2EntityId, response.FederationInfo.Saml2EntityId)
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributes, "federation_info.saml_2_entity_id", *config.federationInfo.Saml2EntityId)
		if err != nil {
			return err
		}

		if config.federationInfo.Saml1xIssuerId != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "saml_1x_issuer_id", *config.federationInfo.Saml1xIssuerId, *response.FederationInfo.Saml1xIssuerId)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "federation_info.saml_1x_issuer_id", *config.federationInfo.Saml1xIssuerId)
			if err != nil {
				return err
			}
		}

		if config.federationInfo.WsfedRealm != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "wsfed_realm", *config.federationInfo.WsfedRealm, *response.FederationInfo.WsfedRealm)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "federation_info.wsfed_realm", *config.federationInfo.WsfedRealm)
			if err != nil {
				return err
			}
		}

		if config.contactInfo != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "company",
				*config.contactInfo.Company, *response.ContactInfo.Company)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "contact_info.company", *config.contactInfo.Company)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "email",
				*config.contactInfo.Email, *response.ContactInfo.Email)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "contact_info.email", *config.contactInfo.Email)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "first_name",
				*config.contactInfo.FirstName, *response.ContactInfo.FirstName)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "contact_info.first_name", *config.contactInfo.FirstName)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "last_name",
				*config.contactInfo.LastName, *response.ContactInfo.LastName)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "contact_info.last_name", *config.contactInfo.LastName)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "phone",
				*config.contactInfo.Phone, *response.ContactInfo.Phone)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "contact_info.phone", *config.contactInfo.Phone)
			if err != nil {
				return err
			}
		}

		if config.emailServer != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "source_addr",
				config.emailServer.SourceAddr, response.EmailServer.SourceAddr)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "email_server.source_addr", config.emailServer.SourceAddr)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "email_server",
				config.emailServer.EmailServer, response.EmailServer.EmailServer)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "email_server.email_server", config.emailServer.EmailServer)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "use_ssl", *config.emailServer.UseSSL, *response.EmailServer.UseSSL)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "email_server.use_ssl", *config.emailServer.UseSSL)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "verify_hostname", *config.emailServer.VerifyHostname, *response.EmailServer.VerifyHostname)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "email_server.verify_hostname", *config.emailServer.VerifyHostname)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_utf8_message_headers", *config.emailServer.EnableUtf8MessageHeaders, *response.EmailServer.EnableUtf8MessageHeaders)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "email_server.enable_utf8_message_headers", *config.emailServer.EnableUtf8MessageHeaders)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "use_debugging", *config.emailServer.UseDebugging, *response.EmailServer.UseDebugging)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "email_server.use_debugging", *config.emailServer.UseDebugging)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "username", *config.emailServer.Username, *response.EmailServer.Username)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "email_server.username", *config.emailServer.Username)
			if err != nil {
				return err
			}

		}

		if config.notifications != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "email_address", config.notifications.LicenseEvents.EmailAddress, response.Notifications.LicenseEvents.EmailAddress)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.license_events.email_address", config.notifications.LicenseEvents.EmailAddress)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "id", config.notifications.LicenseEvents.NotificationPublisherRef.Id, response.Notifications.LicenseEvents.NotificationPublisherRef.Id)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.license_events.notification_publisher_ref.id", config.notifications.LicenseEvents.NotificationPublisherRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "email_address", config.notifications.CertificateExpirations.EmailAddress, response.Notifications.CertificateExpirations.EmailAddress)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.certificate_expirations.email_address", config.notifications.CertificateExpirations.EmailAddress)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchInt(resourceType, nil, "initial_warning_period", *config.notifications.CertificateExpirations.InitialWarningPeriod, *response.Notifications.CertificateExpirations.InitialWarningPeriod)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.certificate_expirations.initial_warning_period", *config.notifications.CertificateExpirations.InitialWarningPeriod)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchInt(resourceType, nil, "final_warning_period", config.notifications.CertificateExpirations.FinalWarningPeriod, response.Notifications.CertificateExpirations.FinalWarningPeriod)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.certificate_expirations.final_warning_period", config.notifications.CertificateExpirations.FinalWarningPeriod)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "notify_admin_user_password_changes", *config.notifications.NotifyAdminUserPasswordChanges, *response.Notifications.NotifyAdminUserPasswordChanges)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.notify_admin_user_password_changes", *config.notifications.NotifyAdminUserPasswordChanges)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "id", config.notifications.AccountChangesNotificationPublisherRef.Id, response.Notifications.AccountChangesNotificationPublisherRef.Id)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.account_changes_notification_publisher_ref.id", config.notifications.AccountChangesNotificationPublisherRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "email_address", config.notifications.MetadataNotificationSettings.EmailAddress, response.Notifications.MetadataNotificationSettings.EmailAddress)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.metadata_notification_settings.email_address", config.notifications.MetadataNotificationSettings.EmailAddress)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "id", config.notifications.MetadataNotificationSettings.NotificationPublisherRef.Id, response.Notifications.MetadataNotificationSettings.NotificationPublisherRef.Id)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "notifications.metadata_notification_settings.notification_publisher_ref.id", config.notifications.MetadataNotificationSettings.NotificationPublisherRef.Id)
			if err != nil {
				return err
			}
		}
		if config.captchaSettings != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "site_key",
				*config.captchaSettings.SiteKey, *response.CaptchaSettings.SiteKey)
			if err != nil {
				return err
			}

			err = acctest.VerifyStateAttributeValue(stateAttributes, "captcha_settings.site_key", *config.captchaSettings.SiteKey)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
