// Code generated by ping-terraform-plugin-framework-generator

package resource_sp_idp_connection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const idpConnOAuthAssertionGrantId = "oauthgrantconn"

func TestAccSpIdpConnection_OAuthAssertionGrantMinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spIdpConnection_OAuthAssertionGrantCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spIdpConnection_OAuthAssertionGrantMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesOAuthAssertionGrantMinimal(),
			},
			{
				// Delete the minimal model
				Config:  spIdpConnection_OAuthAssertionGrantMinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: spIdpConnection_OAuthAssertionGrantCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesOAuthAssertionGrantComplete(),
			},
			{
				// Back to minimal model
				Config: spIdpConnection_OAuthAssertionGrantMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesOAuthAssertionGrantMinimal(),
			},
			{
				// Back to complete model
				Config: spIdpConnection_OAuthAssertionGrantCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesOAuthAssertionGrantComplete(),
			},
			{
				// Test importing the resource
				Config:            spIdpConnection_OAuthAssertionGrantCompleteHCL(),
				ResourceName:      "pingfederate_sp_idp_connection.example",
				ImportStateId:     idpConnOAuthAssertionGrantId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Minimal HCL with only required values set
func spIdpConnection_OAuthAssertionGrantMinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_sp_idp_connection" "example" {
  connection_id = "%s"
  credentials = {
  }
  entity_id = "docker"
  idp_oauth_grant_attribute_mapping = {
    access_token_manager_mappings = [
      {
        access_token_manager_ref = {
          id = "jwt"
        }
        attribute_contract_fulfillment = {
          OrgName = {
            source = {
              id   = null
              type = "TEXT"
            }
            value = "Ping Identity Corporation"
          }
          Username = {
            source = {
              id   = null
              type = "ASSERTION"
            }
            value = "TOKEN_SUBJECT"
          }
        }
      },
    ]
    idp_oauth_attribute_contract = {
    }
  }
  name = "docker"
}
`, idpConnOAuthAssertionGrantId)
}

// Maximal HCL with all values set where possible
func spIdpConnection_OAuthAssertionGrantCompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_sp_idp_connection" "example" {
  active        = true
  base_url      = "https://localhost:9031"
  connection_id = "%s"
  contact_info = {
    company    = "Ping Identity"
    email      = "test@test.com"
    first_name = "test"
    last_name  = "test"
    phone      = "555-5555"
  }
  credentials = {
    certs = [
      {
        active_verification_cert    = true
        encryption_cert             = false
        primary_verification_cert   = true
        secondary_verification_cert = false
        x509_file = {
          crypto_provider = null
          file_data       = "-----BEGIN CERTIFICATE-----\nMIIDOjCCAiICCQCjbB7XBVkxCzANBgkqhkiG9w0BAQsFADBfMRIwEAYDVQQDDAlsb2NhbGhvc3Qx\nDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoM\nA0NEUjELMAkGA1UEBhMCVVMwHhcNMjMwNzE0MDI1NDUzWhcNMjQwNzEzMDI1NDUzWjBfMRIwEAYD\nVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsM\nBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\nggEKAoIBAQC5yFrh9VR2wk9IjzMz+Ei80K453g1j1/Gv3EQ/SC9h7HZBI6aV9FaEYhGnaquRT5q8\n7p8lzCphKNXVyeL6T/pDJOW70zXItkl8Ryoc0tIaknRQmj8+YA0Hr9GDdmYev2yrxSoVS7s5Bl8p\noasn3DljgnWT07vsQz+hw3NY4SPp7IFGP2PpGUBBIIvrOaDWpPGsXeznBxSFtis6Qo+JiEoaVql9\nb9/XyKZj65wOsVyZhFWeM1nCQITSP9OqOc9FSoDFYQ1AVogm4A2AzUrkMnT1SrN2dCuTmNbeVw7g\nOMqMrVf0CiTv9hI0cATbO5we1sPAlJxscSkJjsaI+sQfjiAnAgMBAAEwDQYJKoZIhvcNAQELBQAD\nggEBACgwoH1qklPF1nI9+WbIJ4K12Dl9+U3ZMZa2lP4hAk1rMBHk9SHboOU1CHDQKT1Z6uxi0NI4\nJZHmP1qP8KPNEWTI8Q76ue4Q3aiA53EQguzGb3SEtyp36JGBq05Jor9erEebFftVl83NFvio72Fn\n0N2xvu8zCnlylf2hpz9x1i01Xnz5UNtZ2ppsf2zzT+4U6w3frH+pkp0RDPuoe9mnBF001AguP31h\nSBZyZzWcwQltuNELnSRCcgJl4kC2h3mAgaVtYalrFxLRa3tA2XF2BHRHmKgocedVhTq+81xrqj+W\nQuDmUe06DnrS3Ohmyj3jhsCCluznAolmrBhT/SaDuGg=\n-----END CERTIFICATE-----\n"
          id              = "4qrossmq1vxa4p836kyqzp48h"
        }
      },
    ]
    signing_settings = {
      signing_key_pair_ref = {
        id = "419x9yg43rlawqwq9v6az997k"
      }
      include_raw_key_in_signature = false
      include_cert_in_signature    = false
      algorithm                    = "SHA256withRSA"
    }
  }
  default_virtual_entity_id = "virtual_server_id_1"
  entity_id                 = "docker"
  extended_properties = {
    authNexp = {
      values = ["val1"]
    }
    useAuthnApi = {
      values = ["val2"]
    }
  }
  idp_oauth_grant_attribute_mapping = {
    access_token_manager_mappings = [
      {
        access_token_manager_ref = {
          id = "jwt"
        }
        attribute_contract_fulfillment = {
          OrgName = {
            source = {
              id   = null
              type = "TEXT"
            }
            value = "Ping Identity Corporation"
          }
          Username = {
            source = {
              id   = null
              type = "ASSERTION"
            }
            value = "TOKEN_SUBJECT"
          }
        }
        attribute_sources = [
          {
            ldap_attribute_source = {
              attribute_contract_fulfillment = null
              base_dn                        = "ou=Applications,ou=Ping,ou=Groups,dc=dm,dc=example,dc=com"
              binary_attribute_settings      = null
              data_store_ref = {
                id = "pingdirectory"
              }
              description            = "PingDirectory"
              id                     = "LDAP"
              member_of_nested_group = false
              search_attributes      = ["Subject DN"]
              search_filter          = "(&(memberUid=uid)(cn=Postman))"
              search_scope           = "SUBTREE"
              type                   = "LDAP"
            }
          },
        ]
        issuance_criteria = {
      conditional_criteria = [
            {
              attribute_name = "Username"
              condition      = "MULTIVALUE_CONTAINS_DN"
              error_result   = "myerrorresult"
              source = {
                type = "MAPPED_ATTRIBUTES"
              }
              value = "cn=Example,dc=example,dc=com"
            },
      ]
          expression_criteria = null
        }
      },
    ]
    idp_oauth_attribute_contract = {
      extended_attributes = [
        {
          masked = false
          name   = "asdf"
        },
        {
          masked = false
          name   = "asdfd"
        }
      ]
    }
  }
  logging_mode       = "STANDARD"
  name               = "docker"
  virtual_entity_ids = ["virtual_server_id_1", "virtual_server_id_2", "virtual_server_id_3"]
}
`, idpConnOAuthAssertionGrantId)
}

// Validate any computed values when applying minimal HCL
func spIdpConnection_CheckComputedValuesOAuthAssertionGrantMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "active", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "id", idpConnOAuthAssertionGrantId),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.access_token_manager_mappings.0.attribute_sources.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.access_token_manager_mappings.0.issuance_criteria.conditional_criteria.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.0.name", "TOKEN_SUBJECT"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.extended_attributes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "logging_mode", "STANDARD"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "virtual_entity_ids.#", "0"),
	)
}

// Validate any computed values when applying complete HCL
func spIdpConnection_CheckComputedValuesOAuthAssertionGrantComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.expires", "2024-07-13T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.id", "4qrossmq1vxa4p836kyqzp48h"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.issuer_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.serial_number", "11775821034523537675"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.sha1_fingerprint", "3CFE421ED628F7CEFE08B02DEB3EB4FB5DE9B92D"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.sha256_fingerprint", "633FF42A14E808AEEE5810D78F2C68358AD27787CDDADA302A7E201BA7F2A046"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.status", "EXPIRED"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.subject_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.valid_from", "2023-07-14T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.version", "1"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "credentials.certs.0.x509_file.formatted_file_data"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "id", idpConnOAuthAssertionGrantId),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.0.name", "TOKEN_SUBJECT"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.core_attributes.0.masked", "false"),
	)
}

// Test that any objects created by the test are destroyed
func spIdpConnection_OAuthAssertionGrantCheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.SpIdpConnectionsAPI.DeleteConnection(acctest.TestBasicAuthContext(), idpConnOAuthAssertionGrantId).Execute()
	if err == nil {
		return fmt.Errorf("sp_idp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
