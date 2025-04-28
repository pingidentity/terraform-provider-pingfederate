// Copyright © 2025 Ping Identity Corporation

package licenseagreement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccLicenseAgreement(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseAgreement(),
				Check:  resource.TestCheckResourceAttrSet("pingfederate_license_agreement.example", "license_agreement_url"),
			},
			{
				// Test importing the resource
				Config:                               testAccLicenseAgreement(),
				ResourceName:                         "pingfederate_license_agreement.example",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "accepted",
			},
		},
	})
}

func testAccLicenseAgreement() string {
	return fmt.Sprintf(`
resource "pingfederate_license_agreement" "example" {
  accepted = true
}

data "pingfederate_license_agreement" "example" {
  depends_on = [
    pingfederate_license_agreement.example
  ]
}`)
}
