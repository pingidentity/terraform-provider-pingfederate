// Code generated by ping-terraform-plugin-framework-generator

package resource_sp_idp_connection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccSpIdpConnection_SamlMinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spIdpConnection_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spIdpConnection_SamlMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesSamlMinimal(),
			},
			{
				// Delete the minimal model
				Config:  spIdpConnection_SamlMinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: spIdpConnection_SamlCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesSamlComplete(),
			},
			{
				// Back to minimal model
				Config: spIdpConnection_SamlMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesSamlMinimal(),
			},
			{
				// Back to complete model
				Config: spIdpConnection_SamlCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesSamlComplete(),
			},
			{
				// Test importing the resource
				Config:            spIdpConnection_SamlCompleteHCL(),
				ResourceName:      "pingfederate_sp_idp_connection.example",
				ImportStateId:     spIdpConnectionConnectionId,
				ImportState:       true,
				ImportStateVerify: true,
				// file_data gets formatted by PF so it won't match, and passwords won't be returned by the API
				// encrypted_passwords change on each get
				ImportStateVerifyIgnore: []string{
					"credentials.certs.0.x509_file.file_data",
					"credentials.inbound_back_channel_auth.http_basic_credentials.password",
					"credentials.inbound_back_channel_auth.http_basic_credentials.encrypted_password",
					"credentials.outbound_back_channel_auth.http_basic_credentials.password",
					"credentials.outbound_back_channel_auth.http_basic_credentials.encrypted_password",
				},
			},
		},
	})
}

// Minimal HCL with only required values set
func spIdpConnection_SamlMinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_sp_idp_connection" "example" {
  connection_id                             = "%s"
  credentials = {
    outbound_back_channel_auth = {
      http_basic_credentials = {
        password = "2FederateM0re"
        username           = "user"
      }
    }
  }
  entity_id                 = "partnersec:entity:id"
  error_page_msg_id         = "errorDetail.spSsoFailure"
  idp_browser_sso = {
    artifact = {
      resolver_locations = [
        {
          index = 1
          url   = "https://example.com/endpoint"
        },
      ]
    }
    authentication_policy_contract_mappings = [
      {
        attribute_contract_fulfillment = {
          subject = {
            source = {
              type = "CONTEXT"
            }
            value = "ClientIp"
          }
        }
        authentication_policy_contract_ref = {
          id = "QGxlec5CX693lBQL"
        }
      },
    ]
    enabled_profiles                         = ["IDP_INITIATED_SSO"]
    idp_identity_mapping                     = "ACCOUNT_MAPPING"
    incoming_bindings                        = ["ARTIFACT"]
    protocol                                 = "SAML20"
  }
  name                              = "minimalSaml2"
}
`, spIdpConnectionConnectionId)
}

// Maximal HCL with all values set where possible
func spIdpConnection_SamlCompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "apc1" {
  contract_id = "sp_idp1"
  name = "Example sp_idp1"
  extended_attributes = [
    {
      name = "directory_id"
    },
    {
      name = "given_name"
    },
    {
      name = "family_name"
    },
    {
      name = "email"
    }
  ]
}

resource "pingfederate_authentication_policy_contract" "apc2" {
  contract_id = "sp_idp2"
  name = "Example sp_idp2"
  extended_attributes = [
    {
      name = "directory_id"
    },
    {
      name = "given_name"
    },
    {
      name = "family_name"
    },
    {
      name = "email"
    }
  ]
}

resource "pingfederate_metadata_url" "metadataUrl" {
  url_id             = "myUrlId"
  name               = "My Metadata Url"
  url                = "https://example.com/metadata"
}

resource "pingfederate_sp_idp_connection" "example" {
  active                                    = true
  attribute_query = {
    name_mappings = [
      {
        local_name  = "attr"
        remote_name = "attribute"
      },
    ]
    policy = {
      encrypt_name_id             = true
      mask_attribute_values       = false
      require_encrypted_assertion = true
      require_signed_assertion    = false
      require_signed_response     = true
      sign_attribute_query        = false
    }
    url = "https://example.com"
  }
  base_url      = "https://bxretail.org"
  connection_id = "%s"
  contact_info = {
    company    = "Ping Identity"
    email      = "test@test.com"
    first_name = "test"
    last_name  = "test"
    phone      = "555-5555"
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
          crypto_provider = null
          file_data       = "-----BEGIN CERTIFICATE-----\nMIIDOjCCAiICCQCjbB7XBVkxCzANBgkqhkiG9w0BAQsFADBfMRIwEAYDVQQDDAlsb2NhbGhvc3Qx\nDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoM\nA0NEUjELMAkGA1UEBhMCVVMwHhcNMjMwNzE0MDI1NDUzWhcNMjQwNzEzMDI1NDUzWjBfMRIwEAYD\nVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsM\nBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\nggEKAoIBAQC5yFrh9VR2wk9IjzMz+Ei80K453g1j1/Gv3EQ/SC9h7HZBI6aV9FaEYhGnaquRT5q8\n7p8lzCphKNXVyeL6T/pDJOW70zXItkl8Ryoc0tIaknRQmj8+YA0Hr9GDdmYev2yrxSoVS7s5Bl8p\noasn3DljgnWT07vsQz+hw3NY4SPp7IFGP2PpGUBBIIvrOaDWpPGsXeznBxSFtis6Qo+JiEoaVql9\nb9/XyKZj65wOsVyZhFWeM1nCQITSP9OqOc9FSoDFYQ1AVogm4A2AzUrkMnT1SrN2dCuTmNbeVw7g\nOMqMrVf0CiTv9hI0cATbO5we1sPAlJxscSkJjsaI+sQfjiAnAgMBAAEwDQYJKoZIhvcNAQELBQAD\nggEBACgwoH1qklPF1nI9+WbIJ4K12Dl9+U3ZMZa2lP4hAk1rMBHk9SHboOU1CHDQKT1Z6uxi0NI4\nJZHmP1qP8KPNEWTI8Q76ue4Q3aiA53EQguzGb3SEtyp36JGBq05Jor9erEebFftVl83NFvio72Fn\n0N2xvu8zCnlylf2hpz9x1i01Xnz5UNtZ2ppsf2zzT+4U6w3frH+pkp0RDPuoe9mnBF001AguP31h\nSBZyZzWcwQltuNELnSRCcgJl4kC2h3mAgaVtYalrFxLRa3tA2XF2BHRHmKgocedVhTq+81xrqj+W\nQuDmUe06DnrS3Ohmyj3jhsCCluznAolmrBhT/SaDuGg=\n-----END CERTIFICATE-----\n"
          id              = "4qrossmq1vxa4p836kyqzp48h"
        }
      },
    ]
    decryption_key_pair_ref = {
      id = "rsaprevious"
    }
    inbound_back_channel_auth = {
      certs             = null
      digital_signature = false
      http_basic_credentials = {
        password           = "2FederateM0re"
        username           = "usertwo"
      }
      require_ssl             = true
      verification_issuer_dn  = null
      verification_subject_dn = null
    }
    key_transport_algorithm = "RSA_OAEP"
    outbound_back_channel_auth = {
      digital_signature = false
      http_basic_credentials = {
        encrypted_password = "OBF:JWE:eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMi4xLjMuMCJ9..s29bbqx2Qnf9LLiYYlgjmA.RhccGAZaRPaDPXADkX3VCA.0TU1p0PG4UoPFmjUO8aYyw"
        username           = "user"
      }
      ssl_auth_key_pair_ref = {
        id = "sslclientcert"
      }
      validate_partner_cert = true
    }
    secondary_decryption_key_pair_ref = {
      id = "419x9yg43rlawqwq9v6az997k"
    }
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
  entity_id                 = "partner:entity:id"
  error_page_msg_id         = "errorDetail.spSsoFailure"
  extended_properties = {
    authNexp = {
      values = ["val1"]
    }
    useAuthnApi = {
      values = ["val2"]
    }
  }
  idp_browser_sso = {
    adapter_mappings = [
      {
        attribute_sources = []
        attribute_contract_fulfillment = {
          subject = {
            source = {
              type = "NO_MAPPING"
            }
              value = ""
          }
        }
        issuance_criteria = {
          conditional_criteria = []
        }
        restrict_virtual_entity_ids   = false
        restricted_virtual_entity_ids = []
        sp_adapter_ref = {
          id = "spadapter"
        }
      }
    ]
    always_sign_artifact_response = false
    artifact = {
      lifetime = 60
      resolver_locations = [
        {
          index = 1
          url = "/artifact"
        }
      ]
      source_id = null
    }
    assertions_signed = false
    attribute_contract = {
      extended_attributes = [
        {
          masked = false
          name   = "another"
        },
        {
          masked = true
          name   = "anotherone"
        },
      ]
    }
    authentication_policy_contract_mappings = [
      {
        attribute_contract_fulfillment = {
          directory_id = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          email = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          family_name = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          given_name = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          subject = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
        }
        attribute_sources = [
        ]
        authentication_policy_contract_ref = {
          id = pingfederate_authentication_policy_contract.apc1.id
        }
        issuance_criteria = {
          conditional_criteria = [
            {
              attribute_name = "SAML_SUBJECT"
              condition      = "EQUALS"
              error_result   = "error"
              source = {
                id   = null
                type = "ASSERTION"
              }
              value = "value"
            },
          ]
          expression_criteria = null
        }
        restrict_virtual_server_ids   = true
        restricted_virtual_server_ids = ["virtual_server_id_1", "virtual_server_id_2"]
      },
      {
        attribute_contract_fulfillment = {
          directory_id = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          email = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          family_name = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          given_name = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
          subject = {
            source = {
              id   = null
              type = "NO_MAPPING"
            }
              value = ""
          }
        }
        attribute_sources = [
        ]
        authentication_policy_contract_ref = {
          id = pingfederate_authentication_policy_contract.apc2.id
        }
        issuance_criteria = {
          conditional_criteria = [
            {
              attribute_name = "SAML_SUBJECT"
              condition      = "EQUALS"
              error_result   = "error"
              source = {
                id   = null
                type = "ASSERTION"
              }
              value = "value"
            },
          ]
          expression_criteria = null
        }
        restrict_virtual_server_ids   = true
        restricted_virtual_server_ids = ["virtual_server_id_3"]
      },
    ]
    authn_context_mappings = [
      {
        local = "asdf"
        remote = "sdfg"
      }
    ]
    decryption_policy = {
      assertion_encrypted           = false
      attributes_encrypted          = false
      slo_encrypt_subject_name_id   = false
      slo_subject_name_id_encrypted = false
      subject_name_id_encrypted     = false
    }
    default_target_url                       = "https://example.com"
    enabled_profiles                         = ["IDP_INITIATED_SLO", "IDP_INITIATED_SSO"]
    idp_identity_mapping                     = "ACCOUNT_MAPPING"
    incoming_bindings                        = ["POST", "SOAP", "ARTIFACT"]
    jit_provisioning                         = {
      error_handling = "ABORT_SSO"
      event_trigger = "NEW_USER_ONLY"
      user_attributes = {
      //TODO attribute_contract default here SAML_SUBJECT, another, anotherone
        do_attribute_query = false
      }
      user_repository = {
        ldap = {
          data_store_ref = {
            id = "pingdirectory"
          }
          unique_user_id_filter = "uid=john,ou=org"
          base_dn = "dc=example,dc=com"
          jit_repository_attribute_mapping = {
            USER_KEY = {
              source = {
                id   = null
                type = "NO_MAPPING"
              }
              value = ""
            }
            USER_NAME = {
              source = {
                id   = null
                type = "TEXT"
              }
              value = "asdf"
            }
          }
        }
      }
    }
    message_customizations                   = [
      {
        context_name = "authn-request"
        message_expression = "asdf"
      }
    ]
    oauth_authentication_policy_contract_ref = null
    oidc_provider_settings                   = null
    protocol                                 = "SAML20"
    sign_authn_requests                      = false
    slo_service_endpoints = [
      {
        binding      = "ARTIFACT"
        response_url = "/response"
        url          = "/artifact"
      },
    ]
    sso_oauth_mapping = {
      attribute_contract_fulfillment = {
        USER_KEY = {
          source = {
            id   = null
            type = "NO_MAPPING"
          }
          value = ""
        }
        USER_NAME = {
          source = {
            id   = null
            type = "NO_MAPPING"
          }
          value = ""
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
            filter      = "$${SAML_SUBJECT}"
            id          = null
            schema      = "INFORMATION_SCHEMA"
            table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
          }
          ldap_attribute_source = null
        },
      ]
      issuance_criteria = {
        conditional_criteria = [
        ]
        expression_criteria = null
      }
    }
    sso_service_endpoints = null
    url_whitelist_entries = null
  }
  idp_oauth_grant_attribute_mapping = {
    access_token_manager_mappings = [
      {
        attribute_contract_fulfillment = {
          "Username" = {
            source = {
              type = "NO_MAPPING"
            }
              value = ""
          }
          "OrgName" = {
            source = {
              type = "NO_MAPPING"
            }
              value = ""
          }
        }
          attribute_sources = []
      issuance_criteria = {
        conditional_criteria = [
        ]
        expression_criteria = null
      }
        access_token_manager_ref = {
          id = "jwt"
        }
      }
    ]
    idp_oauth_attribute_contract = {
      extended_attributes = [
        {
      masked = false
          name = "asdf"
        },
        {
      masked = false
          name = "asdfd"
        }
      ]
    }
  }
  inbound_provisioning              = null
  license_connection_group          = null
  logging_mode                      = "STANDARD"
  metadata_reload_settings          = {
    metadata_url_ref = {
      id = pingfederate_metadata_url.metadataUrl.id
    }
    enable_auto_metadata_update = false
  }
  name                              = "My Partner IdP Connection"
  oidc_client_credentials           = null
  virtual_entity_ids                = ["virtual_server_id_1", "virtual_server_id_2", "virtual_server_id_3"]
  ws_trust                          = null
}
`, spIdpConnectionConnectionId)
}

// Validate any computed values when applying minimal HCL
func spIdpConnection_CheckComputedValuesSamlMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "active", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.active_verification_cert", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.encryption_cert", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.primary_verification_cert", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.secondary_verification_cert", "false"),
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
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.encryption_cert", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.primary_verification_cert", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.secondary_verification_cert", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "logging_mode", "STANDARD"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "ws_trust.attribute_contract.core_attributes.0.name", "TOKEN_SUBJECT"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "ws_trust.attribute_contract.core_attributes.0.masked", "false"),
	)
}

// Validate any computed values when applying complete HCL
func spIdpConnection_CheckComputedValuesSamlComplete() resource.TestCheckFunc {
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
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "inbound_back_channel_auth.http_basic_credentials.encrypted_password"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.attribute_contract.core_attributes.0.name", "SAML_SUBJECT"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "idp_browser_sso.jit_provisioning.user_attributes.attribute_contract.#", "3"),
	)
}
