// Copyright © 2025 Ping Identity Corporation

package idpspconnection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccIdpSpConnection_CertificateId(t *testing.T) {
	connId := "certidconn"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpSpConnection_SimpleCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Test applying two certs without defined ids in credentials
				Config: idpSpConnection_CredentialsCertNoId(connId),
			},
		},
	})
}

func idpSpConnection_CredentialsCertNoId(id string) string {
	return fmt.Sprintf(`
%s

resource "pingfederate_idp_sp_connection" "example" {
  active               = false
  application_icon_url = "https://example.com/icon.png"
  application_name     = "MyApp"
  attribute_query = {
    attributes = ["cn"]
    attribute_contract_fulfillment = {
      "cn" = {
        source = {
          type = "TEXT"
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
    issuance_criteria = {
      conditional_criteria = [
        {
          attribute_name = "cn"
          condition      = "MULTIVALUE_CONTAINS_DN"
          source = {
            type = "MAPPED_ATTRIBUTES"
          }
          value = "cn=Example,dc=example,dc=com"
        },
      ]
      expression_criteria = null
    }
    policy = {
      encrypt_assertion              = false
      require_encrypted_name_id      = false
      require_signed_attribute_query = false
      sign_assertion                 = false
      sign_response                  = false
    }
  }
  base_url               = "https://example.com"
  connection_id          = "%s"
  connection_target_type = "STANDARD"
  contact_info = {
    company    = "Example Corp"
    first_name = "Bugs"
    phone      = "5555555"
    email      = "bugsbunny@example.com"
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
          # No id supplied for the cert
        }
      },
    ]
    key_transport_algorithm = "RSA_OAEP"
    inbound_back_channel_auth = {
      certs = [
        {
          active_verification_cert    = true
          encryption_cert             = false
          primary_verification_cert   = true
          secondary_verification_cert = false
          x509_file = {
            file_data = "-----BEGIN CERTIFICATE-----\nMIIDOjCCAiICCQCjbB7XBVkxCzANBgkqhkiG9w0BAQsFADBfMRIwEAYDVQQDDAlsb2NhbGhvc3Qx\nDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoM\nA0NEUjELMAkGA1UEBhMCVVMwHhcNMjMwNzE0MDI1NDUzWhcNMjQwNzEzMDI1NDUzWjBfMRIwEAYD\nVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsM\nBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\nggEKAoIBAQC5yFrh9VR2wk9IjzMz+Ei80K453g1j1/Gv3EQ/SC9h7HZBI6aV9FaEYhGnaquRT5q8\n7p8lzCphKNXVyeL6T/pDJOW70zXItkl8Ryoc0tIaknRQmj8+YA0Hr9GDdmYev2yrxSoVS7s5Bl8p\noasn3DljgnWT07vsQz+hw3NY4SPp7IFGP2PpGUBBIIvrOaDWpPGsXeznBxSFtis6Qo+JiEoaVql9\nb9/XyKZj65wOsVyZhFWeM1nCQITSP9OqOc9FSoDFYQ1AVogm4A2AzUrkMnT1SrN2dCuTmNbeVw7g\nOMqMrVf0CiTv9hI0cATbO5we1sPAlJxscSkJjsaI+sQfjiAnAgMBAAEwDQYJKoZIhvcNAQELBQAD\nggEBACgwoH1qklPF1nI9+WbIJ4K12Dl9+U3ZMZa2lP4hAk1rMBHk9SHboOU1CHDQKT1Z6uxi0NI4\nJZHmP1qP8KPNEWTI8Q76ue4Q3aiA53EQguzGb3SEtyp36JGBq05Jor9erEebFftVl83NFvio72Fn\n0N2xvu8zCnlylf2hpz9x1i01Xnz5UNtZ2ppsf2zzT+4U6w3frH+pkp0RDPuoe9mnBF001AguP31h\nSBZyZzWcwQltuNELnSRCcgJl4kC2h3mAgaVtYalrFxLRa3tA2XF2BHRHmKgocedVhTq+81xrqj+W\nQuDmUe06DnrS3Ohmyj3jhsCCluznAolmrBhT/SaDuGg=\n-----END CERTIFICATE-----\n"
            # No id supplied for the cert
          }
        },
      ]
      digital_signature = false
      http_basic_credentials = {
        encrypted_password = "OBF:JWE:eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMi4xLjMuMCJ9..Q7cuD9L9LT5W8VKdl32iOQ.4AqqdVvegqo6vQxJBc1sBZckQgjxaCrSssbiCaQV3B_ijayEtLePXRCtLUE8P9-U8526lbVee7t93rrByvapYw.PgMT-r6-kKm8TJrmP7-MHg"
        password           = null # sensitive
        username           = "anotheruser"
      }
      require_ssl             = true
      verification_issuer_dn  = null
      verification_subject_dn = null
    }
    outbound_back_channel_auth = {
      digital_signature = false
      http_basic_credentials = {
        password = "2FederateM0re"
        username = "user"
      }
      validate_partner_cert = true
    }
    signing_settings = {
      algorithm                    = "SHA256withRSA"
      include_cert_in_signature    = false
      include_raw_key_in_signature = false
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
  metadata_reload_settings = {
    metadata_url_ref = {
      id = pingfederate_metadata_url.metadataUrl.id
    }
    enable_auto_metadata_update = false
  }
  name = "mySpConn"
  sp_browser_sso = {
    adapter_mappings = [
    ]
    always_sign_artifact_response = true
    artifact = {
      resolver_locations = [
        {
          index = 1
          url   = "https://example.com/endpoint"
        },
      ]
    }
    assertion_lifetime = {
      minutes_after  = 5
      minutes_before = 5
    }
    attribute_contract = {
      core_attributes = [
        {
          name        = "SAML_SUBJECT"
          name_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
        },
      ]
      extended_attributes = [
      ]
    }
    authentication_policy_contract_assertion_mappings = [
      {
        abort_sso_transaction_as_fail_safe = false
        attribute_contract_fulfillment = {
          SAML_SUBJECT = {
            source = {
              type = "AUTHENTICATION_POLICY_CONTRACT"
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
        authentication_policy_contract_ref = {
          id = "QGxlec5CX693lBQL"
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
        }
        restrict_virtual_entity_ids   = false
        restricted_virtual_entity_ids = []
      },
    ]
    default_target_url = "https://example.com"
    enabled_profiles   = ["IDP_INITIATED_SLO", "IDP_INITIATED_SSO"]
    encryption_policy = {
      encrypt_assertion             = true
      encrypt_slo_subject_name_id   = false
      encrypted_attributes          = []
      slo_subject_name_id_encrypted = false
    }
    incoming_bindings = ["ARTIFACT", "POST", "REDIRECT", "SOAP"]
    message_customizations = [
    ]
    protocol                      = "SAML20"
    require_signed_authn_requests = false
    sign_assertions               = true
    sign_response_as_required     = true
    slo_service_endpoints = [
      {
        binding      = "POST"
        response_url = "/response"
        url          = "/artifact"
      },
    ]
    sp_saml_identity_mapping = "STANDARD"
    sso_service_endpoints = [
      {
        binding    = "POST"
        index      = 0
        is_default = true
        url        = "https://httpbin.org/anything"
      },
    ]
  }
  virtual_entity_ids = [
    "example1",
    "example2"
  ]
}
`, idpSpConnection_SamlDependencyHCL(), id)
}
