// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package keypairsoauthopenidconnectadditionalkeysets_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

const keypairsOauthOpenidConnectAdditionalKeySetSetId = "keypairs_oauth_openid_connect_ad"

func TestAccKeypairsOauthOpenidConnectAdditionalKeySet_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: keypairsOauthOpenidConnectAdditionalKeySet_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: keypairsOauthOpenidConnectAdditionalKeySet_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					keypairsOauthOpenidConnectAdditionalKeySet_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKeypairsOauthOpenidConnectAdditionalKeySet_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: keypairsOauthOpenidConnectAdditionalKeySet_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: keypairsOauthOpenidConnectAdditionalKeySet_MinimalHCL(),
				Check:  keypairsOauthOpenidConnectAdditionalKeySet_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  keypairsOauthOpenidConnectAdditionalKeySet_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: keypairsOauthOpenidConnectAdditionalKeySet_CompleteHCL(),
				Check:  keypairsOauthOpenidConnectAdditionalKeySet_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: keypairsOauthOpenidConnectAdditionalKeySet_MinimalHCL(),
				Check:  keypairsOauthOpenidConnectAdditionalKeySet_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: keypairsOauthOpenidConnectAdditionalKeySet_CompleteHCL(),
				Check:  keypairsOauthOpenidConnectAdditionalKeySet_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               keypairsOauthOpenidConnectAdditionalKeySet_CompleteHCL(),
				ResourceName:                         "pingfederate_keypairs_oauth_openid_connect_additional_key_set.example",
				ImportStateId:                        keypairsOauthOpenidConnectAdditionalKeySetSetId,
				ImportStateVerifyIdentifierAttribute: "set_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
		},
	})
}

func keypairsOauthOpenidConnectAdditionalKeySet_VersionRestrictedHCL() string {
	if acctest.VersionAtLeast(version.PingFederate1201) {
		return `
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
  p256_active_key_id = "ec256active"
  p256_previous_key_id = "ec256previous"
  p384_active_key_id = "ec384active"
  p384_previous_key_id = "ec384previous"
  p521_active_key_id = "ec521active"
  p521_previous_key_id = "ec521previous"
  rsa_active_key_id = "rsaactive"
  rsa_previous_key_id = "rsaprevious"
`
	}
	return ""
}

func keypairsOauthOpenidConnectAdditionalKeySet_DependencyHCL() string {
	return `
resource "pingfederate_oauth_issuer" "myissuer" {
  issuer_id   = "myissuer"
  description = "my desc"
  host        = "hostname"
  name        = "my issuer"
  path        = "/example"
}
`
}

// Minimal HCL with only required values set
func keypairsOauthOpenidConnectAdditionalKeySet_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_oauth_openid_connect_additional_key_set" "example" {
  depends_on = [pingfederate_oauth_issuer.myissuer]
  set_id     = "%s"
  issuers = [
    {
      id = "myissuer"
    }
  ]
  name = "minimalname"
  signing_keys = {
    rsa_active_cert_ref = {
      id = "419x9yg43rlawqwq9v6az997k"
    }
  }
}
%s
`, keypairsOauthOpenidConnectAdditionalKeySetSetId,
		keypairsOauthOpenidConnectAdditionalKeySet_DependencyHCL())
}

// Maximal HCL with all values set where possible
func keypairsOauthOpenidConnectAdditionalKeySet_CompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_oauth_openid_connect_additional_key_set" "example" {
  depends_on  = [pingfederate_oauth_issuer.myissuer]
  set_id      = "%s"
  description = "this is my key set"
  issuers = [
    {
      id = "myissuer"
    }
  ]
  name = "mykeyset"
  signing_keys = {
    p256_active_cert_ref = {
      id = "ec256active"
    }
    p256_previous_cert_ref = {
      id = "ec256previous"
    }
    p384_active_cert_ref = {
      id = "ec384active"
    }
    p384_previous_cert_ref = {
      id = "ec384previous"
    }
    p521_active_cert_ref = {
      id = "ec521active"
    }
    p521_previous_cert_ref = {
      id = "ec521previous"
    }
    p521_publish_x5c_parameter = true
    rsa_active_cert_ref = {
      id = "419x9yg43rlawqwq9v6az997k"
    }
    rsa_previous_cert_ref = {
      id = "rsaprevious"
    }
    rsa_publish_x5c_parameter = true
  %s
  }
}
%s
`, keypairsOauthOpenidConnectAdditionalKeySetSetId,
		keypairsOauthOpenidConnectAdditionalKeySet_VersionRestrictedHCL(),
		keypairsOauthOpenidConnectAdditionalKeySet_DependencyHCL())
}

// Validate any computed values when applying minimal HCL
func keypairsOauthOpenidConnectAdditionalKeySet_CheckComputedValuesMinimal() resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1201) {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "id", keypairsOauthOpenidConnectAdditionalKeySetSetId),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.rsa_algorithm_active_key_ids.#", "0"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.rsa_algorithm_previous_key_ids.#", "0"),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p256_publish_x5c_parameter"),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p384_publish_x5c_parameter"),
		)
	} else {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "id", keypairsOauthOpenidConnectAdditionalKeySetSetId),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.rsa_algorithm_active_key_ids"),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.rsa_algorithm_previous_key_ids"),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p256_publish_x5c_parameter"),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p384_publish_x5c_parameter"),
		)
	}
}

// Validate any computed values when applying complete HCL
func keypairsOauthOpenidConnectAdditionalKeySet_CheckComputedValuesComplete() resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1201) {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "id", keypairsOauthOpenidConnectAdditionalKeySetSetId),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p256_publish_x5c_parameter", "false"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p384_publish_x5c_parameter", "false"),
		)
	} else {
		return resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "id", keypairsOauthOpenidConnectAdditionalKeySetSetId),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.rsa_algorithm_active_key_ids"),
			resource.TestCheckNoResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.rsa_algorithm_previous_key_ids"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p256_publish_x5c_parameter", "false"),
			resource.TestCheckResourceAttr("pingfederate_keypairs_oauth_openid_connect_additional_key_set.example", "signing_keys.p384_publish_x5c_parameter", "false"),
		)
	}
}

// Delete the resource
func keypairsOauthOpenidConnectAdditionalKeySet_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.KeyPairsOauthOpenIdConnectAPI.DeleteKeySet(acctest.TestBasicAuthContext(), keypairsOauthOpenidConnectAdditionalKeySetSetId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func keypairsOauthOpenidConnectAdditionalKeySet_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.KeyPairsOauthOpenIdConnectAPI.DeleteKeySet(acctest.TestBasicAuthContext(), keypairsOauthOpenidConnectAdditionalKeySetSetId).Execute()
	if err == nil {
		return fmt.Errorf("keypairs_oauth_openid_connect_additional_key_set still exists after tests. Expected it to be destroyed")
	}
	return nil
}
