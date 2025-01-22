package idpspconnection_test

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

const spConnWsFedId = "wsfedspconn"

func TestAccIdpSpConnection_WsFedMinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpSpConnection_WsFedCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpSpConnection_WsFedMinimalHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsFedMinimal(),
			},
			{
				// Delete the minimal model
				Config:  idpSpConnection_WsFedMinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: idpSpConnection_WsFedCompleteHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsFedComplete(),
			},
			{
				// Back to minimal model
				Config: idpSpConnection_WsFedMinimalHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsFedMinimal(),
			},
			{
				// Back to complete model
				Config: idpSpConnection_WsFedCompleteHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsFedComplete(),
			},
			{
				// Test importing the resource
				Config:            idpSpConnection_WsFedCompleteHCL(),
				ResourceName:      "pingfederate_idp_sp_connection.example",
				ImportStateId:     spConnWsFedId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Minimal HCL with only required values set
func idpSpConnection_WsFedMinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "example" {
  base_url      = "https://example.com"
  connection_id = "%s"
  entity_id     = "myEntity"
  name          = "mySpConn"
  credentials = {
    signing_settings = {
      signing_key_pair_ref = {
        id = "419x9yg43rlawqwq9v6az997k"
      }
      algorithm = "SHA256withRSA"
    }
  }
  sp_browser_sso = {
    protocol = "WSFED"
    sso_service_endpoints = [
      {
        url = "/sp/prpwrong.wsf"
      }
    ]
    sp_ws_fed_identity_mapping = "EMAIL_ADDRESS"
    assertion_lifetime = {
      minutes_before = 5
      minutes_after  = 5
    }
    attribute_contract = {
      core_attributes = [
        {
          name = "SAML_SUBJECT"
        }
      ]
    }
    encryption_policy = {
    }
    adapter_mappings = [
      {
        attribute_contract_fulfillment = {
          "SAML_SUBJECT" = {
            source = {
              type = "ADAPTER"
            }
            value = "subject"
          }
        }
        idp_adapter_ref = {
          id = "OTIdPJava"
        }
      }
    ]
    ws_fed_token_type = "SAML20"
  }
}
`, spConnWsFedId)
}

// Maximal HCL with all values set where possible
func idpSpConnection_WsFedCompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "example" {
  active                 = false
  application_icon_url   = "https://example.com/icon.png"
  application_name       = "MyApp"
  base_url               = "https://example.com"
  connection_id          = "%s"
  connection_target_type = "STANDARD"
  contact_info = {
    company    = "Example Corp"
    email      = "bugsbunny@example.com"
    first_name = "Bugs"
    last_name  = "Bunny"
    phone      = jsonencode(5555555)
  }
  credentials = {
    block_encryption_algorithm = "AES_128"
    certs = [
      {
        active_verification_cert    = true
        encryption_cert             = true
        primary_verification_cert   = true
        secondary_verification_cert = false
        x509_file = {
          file_data = "-----BEGIN CERTIFICATE-----\nMIIDOjCCAiICCQCjbB7XBVkxCzANBgkqhkiG9w0BAQsFADBfMRIwEAYDVQQDDAlsb2NhbGhvc3Qx\nDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoM\nA0NEUjELMAkGA1UEBhMCVVMwHhcNMjMwNzE0MDI1NDUzWhcNMjQwNzEzMDI1NDUzWjBfMRIwEAYD\nVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsM\nBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\nggEKAoIBAQC5yFrh9VR2wk9IjzMz+Ei80K453g1j1/Gv3EQ/SC9h7HZBI6aV9FaEYhGnaquRT5q8\n7p8lzCphKNXVyeL6T/pDJOW70zXItkl8Ryoc0tIaknRQmj8+YA0Hr9GDdmYev2yrxSoVS7s5Bl8p\noasn3DljgnWT07vsQz+hw3NY4SPp7IFGP2PpGUBBIIvrOaDWpPGsXeznBxSFtis6Qo+JiEoaVql9\nb9/XyKZj65wOsVyZhFWeM1nCQITSP9OqOc9FSoDFYQ1AVogm4A2AzUrkMnT1SrN2dCuTmNbeVw7g\nOMqMrVf0CiTv9hI0cATbO5we1sPAlJxscSkJjsaI+sQfjiAnAgMBAAEwDQYJKoZIhvcNAQELBQAD\nggEBACgwoH1qklPF1nI9+WbIJ4K12Dl9+U3ZMZa2lP4hAk1rMBHk9SHboOU1CHDQKT1Z6uxi0NI4\nJZHmP1qP8KPNEWTI8Q76ue4Q3aiA53EQguzGb3SEtyp36JGBq05Jor9erEebFftVl83NFvio72Fn\n0N2xvu8zCnlylf2hpz9x1i01Xnz5UNtZ2ppsf2zzT+4U6w3frH+pkp0RDPuoe9mnBF001AguP31h\nSBZyZzWcwQltuNELnSRCcgJl4kC2h3mAgaVtYalrFxLRa3tA2XF2BHRHmKgocedVhTq+81xrqj+W\nQuDmUe06DnrS3Ohmyj3jhsCCluznAolmrBhT/SaDuGg=\n-----END CERTIFICATE-----\n"
          id        = "4qrossmq1vxa4p836kyqzp48h"
        }
      },
    ]
    key_transport_algorithm = "RSA_OAEP"
    signing_settings = {
      algorithm = "SHA256withRSA"
      alternative_signing_key_pair_refs = [
        {
          id = "rsaprevious"
        },
      ]
      include_cert_in_signature    = false
      include_raw_key_in_signature = true
      signing_key_pair_ref = {
        id = "419x9yg43rlawqwq9v6az997k"
      }
    }
  }
  default_virtual_entity_id = "example2"
  entity_id                 = "myEntity"
  extended_properties = {
    authNexp = {
      values = ["val1"]
    }
    useAuthnApi = {
      values = ["val2"]
    }
  }
  logging_mode = "STANDARD"
  name         = "mySpConn"
  sp_browser_sso = {
    adapter_mappings = [
      {
        abort_sso_transaction_as_fail_safe = false
        attribute_contract_fulfillment = {
          SAML_SUBJECT = {
            source = {
              id   = null
              type = "ADAPTER"
            }
            value = "subject"
          }
        }
        attribute_sources = [
          {
            custom_attribute_source = null
            jdbc_attribute_source = {
              attribute_contract_fulfillment = null
              column_names                   = ["GRANTEE"]
              data_store_ref = {
                id = "ProvisionerDS"
              }
              description = "JDBC"
              filter      = "example"
              id          = "jdbcattrsource"
              schema      = "INFORMATION_SCHEMA"
              table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
            }
            ldap_attribute_source = null
          },
        ]
        idp_adapter_ref = {
          id = "OTIdPJava"
        }
        issuance_criteria = {
      conditional_criteria = [
            {
              attribute_name = "SAML_SUBJECT"
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
        restrict_virtual_entity_ids   = false
        restricted_virtual_entity_ids = []
      },
    ]
    always_sign_artifact_response = false
    assertion_lifetime = {
      minutes_after  = 5
      minutes_before = 5
    }
    attribute_contract = {
      core_attributes = [
        {
          name = "SAML_SUBJECT"
        },
      ]
      extended_attributes = [
      ]
    }
    authentication_policy_contract_assertion_mappings = [
    ]
    encryption_policy = {
      encrypt_assertion             = false
      encrypt_slo_subject_name_id   = false
      encrypted_attributes          = ["SAML_SUBJECT"]
      slo_subject_name_id_encrypted = false
    }
    message_customizations = [
    ]
    protocol                      = "WSFED"
    require_signed_authn_requests = false
    sign_assertions               = false
    slo_service_endpoints = [
    ]
    sp_ws_fed_identity_mapping = "EMAIL_ADDRESS"
    sso_service_endpoints = [
      {
        is_default = false
        url        = "https://httpbin.org/anything"
      },
    ]
    url_whitelist_entries = [
      {
        allow_query_and_fragment = true
        require_https            = true
        valid_domain             = "example.com"
        valid_path               = "/path"
      },
    ]
    ws_fed_token_type = "SAML20"
    ws_trust_version  = "WSTRUST13"
  }
  virtual_entity_ids = ["example1", "example2"]
}
`, spConnWsFedId)
}

// Validate any computed values when applying minimal HCL
func idpSpConnection_CheckComputedValuesWsFedMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "active", "false"),
		resource.TestCheckResourceAttrSet("pingfederate_idp_sp_connection.example", "creation_date"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.signing_settings.include_cert_in_signature", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.signing_settings.include_raw_key_in_signature", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "id", spConnWsFedId),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "logging_mode", "STANDARD"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.adapter_mappings.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.adapter_mappings.0.abort_sso_transaction_as_fail_safe", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.adapter_mappings.0.attribute_sources.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.adapter_mappings.0.issuance_criteria.conditional_criteria.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.adapter_mappings.0.restricted_virtual_entity_ids.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.always_sign_artifact_response", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.attribute_contract.extended_attributes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.authentication_policy_contract_assertion_mappings.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.encryption_policy.encrypt_assertion", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.encryption_policy.encrypt_slo_subject_name_id", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.encryption_policy.encrypted_attributes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.encryption_policy.slo_subject_name_id_encrypted", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.message_customizations.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.require_signed_authn_requests", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.sign_assertions", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.slo_service_endpoints.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.sso_service_endpoints.0.is_default", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "sp_browser_sso.ws_trust_version", "WSTRUST12"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "virtual_entity_ids.#", "0"),
	)
}

// Validate any computed values when applying complete HCL
func idpSpConnection_CheckComputedValuesWsFedComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("pingfederate_idp_sp_connection.example", "creation_date"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.expires", "2024-07-13T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.id", "4qrossmq1vxa4p836kyqzp48h"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.issuer_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.serial_number", "11775821034523537675"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.sha1_fingerprint", "3CFE421ED628F7CEFE08B02DEB3EB4FB5DE9B92D"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.sha256_fingerprint", "633FF42A14E808AEEE5810D78F2C68358AD27787CDDADA302A7E201BA7F2A046"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.status", "EXPIRED"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.subject_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.valid_from", "2023-07-14T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.version", "1"),
		resource.TestCheckResourceAttrSet("pingfederate_idp_sp_connection.example", "credentials.certs.0.x509_file.formatted_file_data"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "id", spConnWsFedId),
	)
}

// Test that any objects created by the test are destroyed
func idpSpConnection_WsFedCheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(acctest.TestBasicAuthContext(), spConnWsFedId).Execute()
	if err == nil {
		return fmt.Errorf("pingfederate_idp_sp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
