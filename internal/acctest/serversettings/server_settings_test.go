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
// The serverSettingsResourceModel struct represents a model for server settings resources.
// It defines the fields that can be used to configure various aspects of the server settings.
type contactInfoResourceModel struct {
	company   string
	email     string
	firstName string
	lastName  string
	phone     string
}

type federationInfoResourceModel struct {
	baseUrl       string
	saml2EntityId string
}

type emailServerResourceModel struct {
	sourceAddr  string
	emailServer string
}

type serverSettingsResourceModel struct {
	contactInfo    *contactInfoResourceModel
	federationInfo federationInfoResourceModel
	emailServer    *emailServerResourceModel
}

func TestAccServerSettings(t *testing.T) {
	resourceName := "myServerSettings"
	initialResourceModel := serverSettingsResourceModel{
		federationInfo: federationInfoResourceModel{
			baseUrl:       "https://localhost:9999",
			saml2EntityId: "initial.pingidentity.com",
		},
	}

	updatedResourceModel := serverSettingsResourceModel{
		contactInfo: &contactInfoResourceModel{
			company:   "updated company",
			email:     "updatedAdminemail@example.com",
			firstName: "Jane2",
			lastName:  "Admin2",
			phone:     "555-555-2222",
		},
		federationInfo: federationInfoResourceModel{
			baseUrl:       "https://localhost2:9999",
			saml2EntityId: "updated.pingidentity.com",
		},
		emailServer: &emailServerResourceModel{
			sourceAddr:  "updatedEmailServerAdmin@example.com",
			emailServer: "updatedemailserver.example.com",
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
				ImportStateVerifyIgnore: []string{"email_server"},
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
		resourceModel.federationInfo.baseUrl,
		resourceModel.federationInfo.saml2EntityId,
	)
}

func testAccServerSettingsComplete(resourceName string, resourceModel serverSettingsResourceModel) string {
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
    base_url         = "%[7]s"
    saml_2_entity_id = "%[8]s"
	saml_1x_issuer_id = "example.com"
    saml_1x_source_id = ""
    wsfed_realm       = "myrealm"
  }

  email_server = {
    source_addr  = "%[9]s"
    email_server = "%[10]s"
	use_ssl                     = true
    verify_hostname             = true
    enable_utf8_message_headers = true
    use_debugging               = false
    username     = "EmailServerAdmin"
    password     = "EmailServerAdminPassword"
  }

  notifications = {
    license_events = {
      email_address = "license-events-email@example.com"
      notification_publisher_ref = {
        id = "exampleSmtpPublisher"
      }
    }
    certificate_expirations = {
      email_address          = "cert-expire-notifications@example.com"
      initial_warning_period = 45
      final_warning_period   = 7
      notification_publisher_ref = {
        id = "exampleSmtpPublisher"
      }
    }
    notify_admin_user_password_changes = true
    account_changes_notification_publisher_ref = {
      id = "exampleSmtpPublisher"
    }
    metadata_notification_settings = {
      email_address = "metadata-notification@example.com"
      notification_publisher_ref = {
        id = "exampleSmtpPublisher"
      }
    }
  }
}`, resourceName,
		resourceModel.contactInfo.company,
		resourceModel.contactInfo.email,
		resourceModel.contactInfo.firstName,
		resourceModel.contactInfo.lastName,
		resourceModel.contactInfo.phone,
		resourceModel.federationInfo.baseUrl,
		resourceModel.federationInfo.saml2EntityId,
		resourceModel.emailServer.sourceAddr,
		resourceModel.emailServer.emailServer,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsAttributes(config serverSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ServerSettingsAPI.GetServerSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "base_url",
			config.federationInfo.baseUrl, *response.FederationInfo.BaseUrl)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "saml_2_entity_id",
			config.federationInfo.saml2EntityId, response.FederationInfo.Saml2EntityId)
		if err != nil {
			return err
		}

		if config.contactInfo != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "company",
				config.contactInfo.company, *response.ContactInfo.Company)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "email",
				config.contactInfo.email, *response.ContactInfo.Email)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "first_name",
				config.contactInfo.firstName, *response.ContactInfo.FirstName)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "last_name",
				config.contactInfo.lastName, *response.ContactInfo.LastName)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "phone",
				config.contactInfo.phone, *response.ContactInfo.Phone)
			if err != nil {
				return err
			}
		}

		if config.emailServer != nil {
			err = acctest.TestAttributesMatchString(resourceType, nil, "source_addr",
				config.emailServer.sourceAddr, response.EmailServer.SourceAddr)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "email_server",
				config.emailServer.emailServer, response.EmailServer.EmailServer)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
