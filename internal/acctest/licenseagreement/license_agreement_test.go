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

const licenseAgreementId = "id"
const licenseAgreementUrlVal = "https://localhost:9999/pf-admin-api/license-agreement"
const acceptedVal = true

// Attributes to test with. Add optional properties to test here if desired.
type licenseAgreementResourceModel struct {
	id                  string
	licenseAgreementUrl string
	accepted            bool
}

func TestAccLicenseAgreement(t *testing.T) {
	resourceName := "myLicenseAgreement"
	initialResourceModel := licenseAgreementResourceModel{
		id:                  licenseAgreementId,
		licenseAgreementUrl: licenseAgreementUrlVal,
		accepted:            acceptedVal,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseAgreement(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLicenseAgreementAttributes(initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLicenseAgreement(resourceName, initialResourceModel),
				ResourceName:      "pingfederate_license_agreement." + resourceName,
				ImportStateId:     licenseAgreementId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccLicenseAgreement(resourceName string, resourceModel licenseAgreementResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_license_agreement" "%[1]s" {
  license_agreement_url = "%[2]s"
  accepted              = %[3]t
}

data "pingfederate_license_agreement" "%[1]s" {
  depends_on = [
    pingfederate_license_agreement.%[1]s
  ]
}`, resourceName,
		resourceModel.licenseAgreementUrl,
		resourceModel.accepted,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedLicenseAgreementAttributes(config licenseAgreementResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "LicenseAgreement"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LicenseAPI.GetLicenseAgreement(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, nil, "license_agreement_url",
			config.licenseAgreementUrl, response.GetLicenseAgreementUrl())
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "accepted",
			config.accepted, response.GetAccepted())
		if err != nil {
			return err
		}
		return nil
	}
}
