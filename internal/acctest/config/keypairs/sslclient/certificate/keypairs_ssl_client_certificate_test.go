package keypairssslclientcertificate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccKeypairsSslClientCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: keypairsSslClientCertificate_MinimalHCL(),
				Check:  keypairsSslClientCertificate_CheckComputedValues(),
			},
		},
	})
}

// TODO when the ssl_client_key resource is supported, create the dependency in this test
// Only the key_id attribute can be set on this resource
func keypairsSslClientCertificate_MinimalHCL() string {
	return `
data "pingfederate_keypairs_ssl_client_certificate" "example" {
  key_id = "sslclientcert"
}
`
}

// Validate any computed values when applying HCL
func keypairsSslClientCertificate_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.pingfederate_keypairs_ssl_client_certificate.example", "id", "sslclientcert"),
		resource.TestCheckResourceAttrSet("data.pingfederate_keypairs_ssl_client_certificate.example", "exported_certificate"),
	)
}
