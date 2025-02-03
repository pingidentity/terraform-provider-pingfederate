// Copyright © 2025 Ping Identity Corporation

package licenseagreement_test

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

func TestAccLicenseAgreement(t *testing.T) {
	resourceName := "myLicenseAgreement"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseAgreement(resourceName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedLicenseAgreementAttributes(true),
					resource.TestCheckResourceAttr("pingfederate_license_agreement."+resourceName, "accepted", "true"),
				),
			},
			{
				// Test importing the resource
				Config:                               testAccLicenseAgreement(resourceName, true),
				ResourceName:                         "pingfederate_license_agreement." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "accepted",
			},
		},
	})
}

func testAccLicenseAgreement(resourceName string, accepted bool) string {
	return fmt.Sprintf(`
resource "pingfederate_license_agreement" "%[1]s" {
  accepted = %[2]t
}

data "pingfederate_license_agreement" "%[1]s" {
  depends_on = [
    pingfederate_license_agreement.%[1]s
  ]
}`, resourceName,
		accepted,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedLicenseAgreementAttributes(accepted bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "LicenseAgreement"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LicenseAPI.GetLicenseAgreement(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "accepted",
			accepted, response.GetAccepted())
		if err != nil {
			return err
		}
		return nil
	}
}
