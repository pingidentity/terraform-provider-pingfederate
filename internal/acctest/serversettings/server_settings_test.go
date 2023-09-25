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

const sessionSettingsId = "2"

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
	baseUrl string
}

type emailServerResourceModel struct {
	sourceAddr  string
	emailServer string
}

type serverSettingsResourceModel struct {
	contactInfo    contactInfoResourceModel
	federationInfo federationInfoResourceModel
	emailServer    emailServerResourceModel
}

func TestAccServerSettings(t *testing.T) {
	resourceName := "myServerSettings"
	initialResourceModel := serverSettingsResourceModel{
		contactInfo: contactInfoResourceModel{
			company:   "initial company",
			email:     "initialAdmin@example.com",
			firstName: "Jane",
			lastName:  "Admin",
			phone:     "555-555-1111",
		},
		federationInfo: federationInfoResourceModel{
			baseUrl: "https://localhost:9999",
		},
		emailServer: emailServerResourceModel{
			sourceAddr:  "initialEmailServerAdmin@example.com",
			emailServer: "initialEmailserver.example.com",
		},
	}

	updatedResourceModel := serverSettingsResourceModel{
		contactInfo: contactInfoResourceModel{
			company:   "updated company",
			email:     "updatedAdminemail@example.com",
			firstName: "Jane2",
			lastName:  "Admin2",
			phone:     "555-555-2222",
		},
		federationInfo: federationInfoResourceModel{
			baseUrl: "https://localhost2:9999",
		},
		emailServer: emailServerResourceModel{
			sourceAddr:  "updatedEmailServerAdmin@example.com",
			emailServer: "updatedEmailserver.example.com",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
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
				ImportStateId:           sessionSettingsId,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"email_server"},
			},
		},
	})
}

func testAccServerSettings(resourceName string, resourceModel serverSettingsResourceModel) string {
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
    base_url = "%[7]s"
  }

  email_server = {
    source_addr  = "%[8]s"
    email_server = "%[9]s"
    username     = "EmailServerAdmin"
    password     = "EmailServerAdminPassword"
  }
}`, resourceName,
		resourceModel.contactInfo.company,
		resourceModel.contactInfo.email,
		resourceModel.contactInfo.firstName,
		resourceModel.contactInfo.lastName,
		resourceModel.contactInfo.phone,
		resourceModel.federationInfo.baseUrl,
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
		response, _, err := testClient.ServerSettingsApi.GetServerSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
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

		err = acctest.TestAttributesMatchString(resourceType, nil, "base_url",
			config.federationInfo.baseUrl, *response.FederationInfo.BaseUrl)
		if err != nil {
			return err
		}

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

		return nil
	}
}
