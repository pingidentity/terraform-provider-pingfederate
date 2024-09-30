// Code generated by ping-terraform-plugin-framework-generator

package keypairsoauthopenidconnect_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccKeypairsOauthOpenidConnect_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: keypairsOauthOpenidConnect_MinimalHCL(),
				Check:  keypairsOauthOpenidConnect_CheckComputedValuesMinimal(),
			},
			{
				// Update to a complete model
				Config: keypairsOauthOpenidConnect_CompleteHCL(),
				Check:  keypairsOauthOpenidConnect_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               keypairsOauthOpenidConnect_CompleteHCL(),
				ResourceName:                         "pingfederate_keypairs_oauth_openid_connect.example",
				ImportStateVerifyIdentifierAttribute: "static_jwks_enabled",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				// Back to minimal model
				Config: keypairsOauthOpenidConnect_MinimalHCL(),
				Check:  keypairsOauthOpenidConnect_CheckComputedValuesMinimal(),
			},
		},
	})
}

func keyspairsOauthOpenidConnect_VersionRestrictedHCL() string {
	if acctest.VersionAtLeast(version.PingFederate1201) {
		return `
  p256_active_key_id = "ec256active"
  p256_decryption_active_key_id = "ec256decryptactive"
  p256_decryption_previous_key_id = "ec256decryptprevious"
  p256_previous_key_id = "ec256previous"
  p384_active_key_id = "ec384active"
  p384_decryption_active_key_id = "ec384decryptactive"
  p384_decryption_previous_key_id = "ec384decryptprevious"
  p384_previous_key_id = "ec384previous"
  p521_active_key_id = "ec521active"
  p521_decryption_active_key_id = "ec521decryptactive"
  p521_decryption_previous_key_id = "ec521decryptprevious"
  p521_previous_key_id = "ec521previous"
  rsa_active_key_id = "rsaactive"
  rsa_decryption_active_key_id = "rsadecryptactive"
  rsa_decryption_previous_key_id = "rsadecryptprevious"
  rsa_previous_key_id = "rsaprevious"
  rsa_algorithm_active_key_ids = [
    {
      key_id       = "rsalistactive"
      rsa_alg_type = "RS256"
    }
  ]
  rsa_algorithm_previous_key_ids = [
    {
      key_id       = "rsalistpreviousone"
      rsa_alg_type = "RS384"
    },
    {
      key_id       = "rsalistprevioustwo"
      rsa_alg_type = "RS512"
    }
  ]
`
	}
	return ""
}

// Minimal HCL with only required values set
func keypairsOauthOpenidConnect_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_oauth_openid_connect" "example" {
  static_jwks_enabled = false
}
`)
}

// Maximal HCL with all values set where possible
func keypairsOauthOpenidConnect_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_oauth_openid_connect" "example" {
  p256_active_cert_ref = {
    id = "ec256active"
  }
  p256_decryption_active_cert_ref = {
    id = "ec256active"
  }
  p256_decryption_previous_cert_ref = {
    id = "ec256previous"
  }
  p256_decryption_publish_x5c_parameter = true
  p256_previous_cert_ref = {
    id = "ec256previous"
  }
  p384_active_cert_ref = {
    id = "ec384active"
  }
  p384_decryption_active_cert_ref = {
    id = "ec384active"
  }
  p384_decryption_previous_cert_ref = {
    id = "ec384previous"
  }
  p384_decryption_publish_x5c_parameter = true
  p384_previous_cert_ref = {
    id = "ec384previous"
  }
  p521_active_cert_ref = {
    id = "ec521active"
  }
  p521_decryption_active_cert_ref = {
    id = "ec521active"
  }
  p521_decryption_previous_cert_ref = {
    id = "ec521previous"
  }
  p521_previous_cert_ref = {
    id = "ec521previous"
  }
  p521_publish_x5c_parameter = true
  rsa_active_cert_ref = {
    id = "419x9yg43rlawqwq9v6az997k"
  }
  rsa_decryption_active_cert_ref = {
    id = "419x9yg43rlawqwq9v6az997k"
  }
  rsa_decryption_previous_cert_ref = {
    id = "rsaprevious"
  }
  rsa_previous_cert_ref = {
    id = "rsaprevious"
  }
  rsa_publish_x5c_parameter = true
  static_jwks_enabled       = true
  %s
}
`, keyspairsOauthOpenidConnect_VersionRestrictedHCL())
}

// Validate any computed values when applying minimal HCL
func keypairsOauthOpenidConnect_CheckComputedValuesMinimal() resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1201) {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_algorithm_active_key_ids.#", "0"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_algorithm_previous_key_ids.#", "0"),
		)
	}
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_algorithm_active_key_ids"),
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_algorithm_previous_key_ids"),
	)
}

// Validate any computed values when applying complete HCL
func keypairsOauthOpenidConnect_CheckComputedValuesComplete() resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1201) {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "p256_publish_x5c_parameter", "false"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "p384_publish_x5c_parameter", "false"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "p521_decryption_publish_x5c_parameter", "false"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_decryption_publish_x5c_parameter", "false"),
		)
	}
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "p256_publish_x5c_parameter", "false"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "p384_publish_x5c_parameter", "false"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "p521_decryption_publish_x5c_parameter", "false"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_decryption_publish_x5c_parameter", "false"),
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_algorithm_active_key_ids"),
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect.example", "rsa_algorithm_previous_key_ids"),
	)
}
