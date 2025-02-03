// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package keypairssslserversettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// This resource has all attributes required, so no need to switch between different HCL models
func TestAccKeypairsSslServerSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: keypairsSslServerSettings_MinimalHCL(),
			},
			{
				// Test importing the resource
				Config:                               keypairsSslServerSettings_MinimalHCL(),
				ResourceName:                         "pingfederate_keypairs_ssl_server_settings.example",
				ImportStateVerifyIdentifierAttribute: "admin_console_cert_ref.id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

// Minimal HCL with only required values set
func keypairsSslServerSettings_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_ssl_server_settings" "example" {
  active_admin_console_certs = [
    {
      id = "sslservercert"
    }
  ]
  active_runtime_server_certs = [
    {
      id = "sslservercert"
    }
  ]
  admin_console_cert_ref = {
    id = "sslservercert"
  }
  runtime_server_cert_ref = {
    id = "sslservercert"
  }
}
`)
}
