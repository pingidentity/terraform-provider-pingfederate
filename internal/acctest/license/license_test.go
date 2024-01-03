package acctest_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccLicense(t *testing.T) {
	resourceName := "myLicense"

	licenseVar := "PF_TF_ACC_TEST_LICENSE_11"
	if acctest.VersionAtLeast(version.PingFederate1200) {
		licenseVar = "PF_TF_ACC_TEST_LICENSE_12"
	}
	licenseFileData := os.Getenv(licenseVar)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if licenseFileData == "" {
				t.Fatal(licenseVar + " must be set for acceptance tests on this PingFederate version")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLicense(resourceName, licenseFileData),
			},
			{
				// Test importing the resource
				Config:            testAccLicense(resourceName, licenseFileData),
				ResourceName:      "pingfederate_license." + resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccLicense(resourceName string, licenseFileData string) string {
	return fmt.Sprintf(`
resource "pingfederate_license" "%[1]s" {
  file_data = "%[2]s"
}

data "pingfederate_license" "%[1]s" {
  depends_on = [
    pingfederate_license.%[1]s
  ]
}

resource "pingfederate_license" "licenseExample" {
  file_data = "%[2]s"
}`, resourceName,
		licenseFileData,
	)
}
