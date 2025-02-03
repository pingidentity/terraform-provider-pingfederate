// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package protocolmetadatasigningsettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// No computed values in this resource
func TestAccProtocolMetadataSigningSettings_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: protocolMetadataSigningSettings_MinimalHCL(),
			},
			{
				// Update to a complete model
				Config: protocolMetadataSigningSettings_CompleteHCL(),
			},
			{
				// Test importing the resource
				Config:                               protocolMetadataSigningSettings_CompleteHCL(),
				ResourceName:                         "pingfederate_protocol_metadata_signing_settings.example",
				ImportStateVerifyIdentifierAttribute: "signature_algorithm",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				// Back to minimal model
				Config: protocolMetadataSigningSettings_MinimalHCL(),
			},
		},
	})
}

// Minimal HCL with only required values set
func protocolMetadataSigningSettings_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_protocol_metadata_signing_settings" "example" {
}
`)
}

// Maximal HCL with all values set where possible
func protocolMetadataSigningSettings_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_protocol_metadata_signing_settings" "example" {
  signature_algorithm = "SHA256withRSA"
  signing_key_ref = {
    id = "419x9yg43rlawqwq9v6az997k"
  }
}
`)
}
