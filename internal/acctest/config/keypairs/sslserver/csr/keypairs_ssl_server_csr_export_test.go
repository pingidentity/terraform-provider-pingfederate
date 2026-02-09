// Copyright © 2026 Ping Identity Corporation

package keypairssslservercsr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccKeypairsSslServerCsrExport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: keypairsSslServerCsrExport_MinimalHCL(),
				Check:  keypairsSslServerCsrExport_CheckComputedValues(),
			},
			{
				// Expect no additional rotation
				Config: keypairsSslServerCsrExport_NoExportHCL(),
				Check:  keypairsSslServerCsrExport_CheckComputedValues(),
			},
			{
				// Expect rotation
				Config: keypairsSslServerCsrExport_SecondExportHCL(),
				Check:  keypairsSslServerCsrExport_CheckComputedValues(),
			},
			{
				// Expect no additional rotation
				Config: keypairsSslServerCsrExport_SecondNoExportHCL(),
				Check:  keypairsSslServerCsrExport_CheckComputedValues(),
			},
			{
				// Back to the original with no trigger values
				Config: keypairsSslServerCsrExport_MinimalHCL(),
				Check:  keypairsSslServerCsrExport_CheckComputedValues(),
			},
		},
	})
}

// TODO update once the ssl client key can be created in this test with HCL
// Minimal HCL with only required values set
func keypairsSslServerCsrExport_MinimalHCL() string {
	return `
resource "pingfederate_keypairs_ssl_server_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
}
`
}

func keypairsSslServerCsrExport_NoExportHCL() string {
	return `
resource "pingfederate_keypairs_ssl_server_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger" = "false"
  }
}
`
}

func keypairsSslServerCsrExport_SecondExportHCL() string {
	return `
resource "pingfederate_keypairs_ssl_server_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger"    = "updated"
    "newtrigger" = "new"
  }
}
`
}

func keypairsSslServerCsrExport_SecondNoExportHCL() string {
	return `
resource "pingfederate_keypairs_ssl_server_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger" = "updated"
  }
}
`
}

// Validate any computed values when applying HCL
func keypairsSslServerCsrExport_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_server_csr_export.example", "id", "419x9yg43rlawqwq9v6az997k"),
		resource.TestCheckResourceAttrSet("pingfederate_keypairs_ssl_server_csr_export.example", "exported_csr"),
	)
}
