// Copyright © 2025 Ping Identity Corporation

package keypairssigningcertificate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccKeypairsSigningCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: keypairsSigningCertificate_MinimalHCL(),
				Check:  keypairsSigningCertificate_CheckComputedValues(),
			},
		},
	})
}

// Only the key_id attribute can be set on this resource
func keypairsSigningCertificate_MinimalHCL() string {
	return `
data "pingfederate_keypairs_signing_certificate" "example" {
  key_id = "419x9yg43rlawqwq9v6az997k"
}
`
}

// Validate any computed values when applying HCL
func keypairsSigningCertificate_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.pingfederate_keypairs_signing_certificate.example", "exported_certificate"),
		resource.TestCheckResourceAttr("data.pingfederate_keypairs_signing_certificate.example", "id", "419x9yg43rlawqwq9v6az997k"),
	)
}
