// Copyright © 2026 Ping Identity Corporation

package license_test

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

	var licenseVar string
	if acctest.VersionAtLeast(version.PingFederate1300) {
		licenseVar = "PF_TF_ACC_TEST_LICENSE_13"
	} else if acctest.VersionAtLeast(version.PingFederate1200) {
		licenseVar = "PF_TF_ACC_TEST_LICENSE_12"
	} else {
		licenseVar = "PF_TF_ACC_TEST_LICENSE_11"
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
}`, resourceName,
		licenseFileData,
	)
}
