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

// Validate any computed values when applying HCL
func keypairsSslServerCsrExport_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_keypairs_ssl_server_csr_export.example", "id", "419x9yg43rlawqwq9v6az997k"),
		resource.TestCheckResourceAttrSet("pingfederate_keypairs_ssl_server_csr_export.example", "exported_csr"),
	)
}
