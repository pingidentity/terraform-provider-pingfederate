package keypairssigningcsr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccKeypairsSigningCsrExportResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: keypairsSigningCsrExportResource_MinimalHCL(),
				Check:  keypairsSigningCsrExportResource_CheckComputedValues(),
			},
			{
				// Expect no additional rotation
				Config: keypairsSigningCsrExportResource_NoExportHCL(),
				Check:  keypairsSigningCsrExportResource_CheckComputedValues(),
			},
			{
				// Expect rotation
				Config: keypairsSigningCsrExportResource_SecondExportHCL(),
				Check:  keypairsSigningCsrExportResource_CheckComputedValues(),
			},
			{
				// Expect no additional rotation
				Config: keypairsSigningCsrExportResource_SecondNoExportHCL(),
				Check:  keypairsSigningCsrExportResource_CheckComputedValues(),
			},
			{
				// Back to the original with no trigger values
				Config: keypairsSigningCsrExportResource_MinimalHCL(),
				Check:  keypairsSigningCsrExportResource_CheckComputedValues(),
			},
		},
	})
}

func keypairsSigningCsrExportResource_MinimalHCL() string {
	return `
resource "pingfederate_keypairs_signing_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
}
`
}

func keypairsSigningCsrExportResource_NoExportHCL() string {
	return `
resource "pingfederate_keypairs_signing_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger" = "false"
  }
}
`
}

func keypairsSigningCsrExportResource_SecondExportHCL() string {
	return `
resource "pingfederate_keypairs_signing_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger"    = "updated"
    "newtrigger" = "new"
  }
}
`
}

func keypairsSigningCsrExportResource_SecondNoExportHCL() string {
	return `
resource "pingfederate_keypairs_signing_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger" = "updated"
  }
}
`
}

// Validate any computed values when applying HCL
func keypairsSigningCsrExportResource_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("pingfederate_keypairs_signing_csr_export.example", "exported_csr"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_export.example", "id", "419x9yg43rlawqwq9v6az997k"),
	)
}
