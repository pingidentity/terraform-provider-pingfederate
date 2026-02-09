// Copyright © 2026 Ping Identity Corporation

package certificatesca_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccCertificatesCAExport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: certificatesCAExport_MinimalHCL(),
				Check:  certificatesCAExport_CheckComputedValues(),
			},
		},
	})
}

// Only the ca_id attribute can be set on this resource
func certificatesCAExport_MinimalHCL() string {
	return `
data "pingfederate_certificates_ca_export" "example" {
  ca_id = "gdxuvcw6p95rex3go7eb3ctsb"
}
`
}

// Validate any computed values when applying HCL
func certificatesCAExport_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.pingfederate_certificates_ca_export.example", "exported_certificate"),
		resource.TestCheckResourceAttr("data.pingfederate_certificates_ca_export.example", "id", "gdxuvcw6p95rex3go7eb3ctsb"),
	)
}
