package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const licenseId = "id"
const fileData = "SUQ9MDA0ODk3MzEKV1NUcnVzdFNUUz10cnVlCk9BdXRoPXRydWUKU2Fhc1Byb3Zpc2lvbmluZz10cnVlClByb2R1Y3Q9UGluZ0ZlZGVyYXRlClZlcnNpb249MTEuMgpFbmZvcmNlbWVudFR5cGU9MwpUaWVyPUZyZWUKSXNzdWVEYXRlPTIwMjMtMDYtMDcKRXhwaXJhdGlvbkRhdGU9MjAyNC0wNi0wNwpEZXBsb3ltZW50TWV0aG9kPURvY2tlcgpPcmdhbml6YXRpb249UGluZyBJZGVudGl0eSBDb3Jwb3JhdGlvbgpTaWduQ29kZT1GRjBGClNpZ25hdHVyZT0zMDJDMDIxNDM4QkIwQzk5RjYwQUY1RkE4MzBBRUQ4NjEzOENGRENCNTAxNDYzNzUwMjE0NENCODc3MEI3N0ZDNzgwMUQ0M0QwNjQwMTVDNjIwOTJBNDY1RjZEMA=="

// Attributes to test with. Add optional properties to test here if desired.
type licenseResourceModel struct {
	id       string
	fileData string
}

func TestAccLicense(t *testing.T) {
	resourceName := "myLicense"
	initialResourceModel := licenseResourceModel{
		id:       licenseId,
		fileData: fileData,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccLicense(resourceName, initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLicense(resourceName, initialResourceModel),
				ResourceName:      "pingfederate_license." + resourceName,
				ImportStateId:     licenseId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccLicense(resourceName string, resourceModel licenseResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_license" "%[1]s" {
  file_data = "%[2]s"
}`, resourceName,
		fileData,
	)
}
