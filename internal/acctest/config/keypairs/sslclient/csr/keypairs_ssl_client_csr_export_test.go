// Copyright © 2025 Ping Identity Corporation

package keypairssslclientcsr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccKeypairsSslClientCsrExport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: keypairsSslClientCsrExport_MinimalHCL(),
				Check:  keypairsSslClientCsrExport_CheckComputedValues(),
			},
			{
				// Expect no additional rotation
				Config: keypairsSslClientCsrExport_NoExportHCL(),
				Check:  keypairsSslClientCsrExport_CheckComputedValues(),
			},
			{
				// Expect rotation
				Config: keypairsSslClientCsrExport_SecondExportHCL(),
				Check:  keypairsSslClientCsrExport_CheckComputedValues(),
			},
			{
				// Expect no additional rotation
				Config: keypairsSslClientCsrExport_SecondNoExportHCL(),
				Check:  keypairsSslClientCsrExport_CheckComputedValues(),
			},
			{
				// Back to the original with no trigger values
				Config: keypairsSslClientCsrExport_MinimalHCL(),
				Check:  keypairsSslClientCsrExport_CheckComputedValues(),
			},
		},
	})
}

// TODO update once the ssl client key can be created in this test with HCL
// Minimal HCL with only required values set
func keypairsSslClientCsrExport_MinimalHCL() string {
	return `
resource "pingfederate_keypairs_ssl_client_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
}
`
}

func keypairsSslClientCsrExport_NoExportHCL() string {
	return `
resource "pingfederate_keypairs_ssl_client_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger" = "false"
  }
}
`
}

func keypairsSslClientCsrExport_SecondExportHCL() string {
	return `
resource "pingfederate_keypairs_ssl_client_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger"    = "updated"
    "newtrigger" = "new"
  }
}
`
}

func keypairsSslClientCsrExport_SecondNoExportHCL() string {
	return `
resource "pingfederate_keypairs_ssl_client_csr_export" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  export_trigger_values = {
    "trigger" = "updated"
  }
}
`
}

// Validate any computed values when applying HCL
func keypairsSslClientCsrExport_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_client_csr_export.example", "id", "419x9yg43rlawqwq9v6az997k"),
		resource.TestCheckResourceAttrSet("pingfederate_keypairs_ssl_client_csr_export.example", "exported_csr"),
	)
}
