package keypairssigningcsr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccKeypairsSigningCsrDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: keypairsSigningCsrDataSource_MinimalHCL(),
				Check:  keypairsSigningCsrDataSource_CheckComputedValues(),
			},
		},
	})
}

// Only the keypair_id attribute can be set on this resource
func keypairsSigningCsrDataSource_MinimalHCL() string {
	return `
data "pingfederate_keypairs_signing_csr" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
}
`
}

// Validate any computed values when applying HCL
func keypairsSigningCsrDataSource_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.pingfederate_keypairs_signing_csr.example", "exported_csr"),
	)
}
