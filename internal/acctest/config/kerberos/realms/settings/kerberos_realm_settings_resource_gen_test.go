// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package kerberosrealmssettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccKerberosRealmSettings_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: kerberosRealmSettings_MinimalHCL(),
				Check:  kerberosRealmSettings_CheckComputedValuesMinimal(),
			},
			{
				// Update to a complete model
				Config: kerberosRealmSettings_CompleteHCL(),
			},
			{
				// Test importing the resource
				Config:                               kerberosRealmSettings_CompleteHCL(),
				ResourceName:                         "pingfederate_kerberos_realm_settings.example",
				ImportStateVerifyIdentifierAttribute: "debug_log_output",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				// Back to minimal model
				Config: kerberosRealmSettings_MinimalHCL(),
				Check:  kerberosRealmSettings_CheckComputedValuesMinimal(),
			},
		},
	})
}

func kerberosRealmSettings_DependencyHCL() string {
	return `
resource "pingfederate_kerberos_realm" "kerberosRealmExample" {
  realm_id            = "kerberosRealm"
  kerberos_realm_name = "kerberosRealm"
  kerberos_username   = "kerberosUsername"
  kerberos_password   = "kerberosPassword"
}
`
}

// Minimal HCL with only required values set
func kerberosRealmSettings_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_kerberos_realm_settings" "example" {
  depends_on  = [pingfederate_kerberos_realm.kerberosRealmExample]
  kdc_retries = 3
  kdc_timeout = 4
}
%s
`, kerberosRealmSettings_DependencyHCL())
}

// Maximal HCL with all values set where possible
func kerberosRealmSettings_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_kerberos_realm_settings" "example" {
  depends_on                    = [pingfederate_kerberos_realm.kerberosRealmExample]
  debug_log_output              = true
  force_tcp                     = true
  kdc_retries                   = 4
  kdc_timeout                   = 3
  key_set_retention_period_mins = 360
}
%s
`, kerberosRealmSettings_DependencyHCL())
}

// Validate any computed values when applying minimal HCL
func kerberosRealmSettings_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_kerberos_realm_settings.example", "debug_log_output", "false"),
		resource.TestCheckResourceAttr("pingfederate_kerberos_realm_settings.example", "force_tcp", "false"),
		resource.TestCheckResourceAttr("pingfederate_kerberos_realm_settings.example", "key_set_retention_period_mins", "610"),
	)
}
