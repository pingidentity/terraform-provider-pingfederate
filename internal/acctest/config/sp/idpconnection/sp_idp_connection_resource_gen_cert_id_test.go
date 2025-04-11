// Copyright © 2025 Ping Identity Corporation

package resource_sp_idp_connection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/accesstokenmanager"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccSpIdpConnection_CertificateId(t *testing.T) {
	connId := "certidconn"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spIdpConnection_SimpleCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Test applying two certs without defined ids in credentials
				Config: spIdpConnection_CredentialsCertNoId(connId),
			},
		},
	})
}

func spIdpConnection_CredentialsCertNoId(id string) string {
	return fmt.Sprintf(`
  %s

resource "pingfederate_sp_idp_connection" "example" {
  connection_id             = "%s"
  active                    = true
  name                      = "connection name"
  entity_id                 = "entity_id"
  logging_mode              = "STANDARD"
  virtual_entity_ids        = ["virtual_server_id"]
  base_url                  = "https://example.com"
  default_virtual_entity_id = "virtual_server_id"
  error_page_msg_id         = "errorDetail.spSsoFailure"

  attribute_query = {
    url = "https://example.com"
    name_mappings = [
      {
        local_name  = "local name"
        remote_name = "remote name"
      }
    ]
    policy = {
      sign_attribute_query        = true
      encrypt_name_id             = true
      require_signed_response     = true
      require_signed_assertion    = true
      require_encrypted_assertion = true
      mask_attribute_values       = true
    }
  }

  contact_info = {
    first_name = "test"
    last_name  = "test"
    phone      = "555-5555"
    email      = "test@test.com"
    company    = "Ping Identity"
  }

  idp_browser_sso = {
    protocol = "SAML20"
    enabled_profiles = [
      "IDP_INITIATED_SSO"
    ]
    incoming_bindings = [
      "POST"
    ]
    default_target_url            = "https://example.com"
    always_sign_artifact_response = false
    decryption_policy = {
      assertion_encrypted           = false
      subject_name_id_encrypted     = false
      attributes_encrypted          = false
      slo_encrypt_subject_name_id   = false
      slo_subject_name_id_encrypted = false
    }
    idp_identity_mapping = "ACCOUNT_MAPPING"
    attribute_contract = {
      extended_attributes = []
    }
    adapter_mappings = [
      {
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
              member_of_nested_group = false
              search_attributes      = ["Subject DN"]
              search_filter          = "(&(memberUid=uid)(cn=Postman))"
              search_scope           = "SUBTREE"
              type                   = "LDAP"
            }
          },
        ]
        attribute_contract_fulfillment = {
          subject = {
            source = {
              type = "NO_MAPPING"
            }
          }
        }
        restrict_virtual_entity_ids   = false
        restricted_virtual_entity_ids = []
        sp_adapter_ref = {
          id = "spadapter",
        }
      }
    ]
    authentication_policy_contract_mappings = [
      {
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
        attribute_contract_fulfillment = {
          "firstName" : {
            source = {
              type = "NO_MAPPING"
            }
          },
          "lastName" : {
            source = {
              type = "NO_MAPPING"
            }
          },
          "ImmutableID" : {
            source = {
              type = "NO_MAPPING"
            }
          },
          "mail" : {
            source = {
              type = "NO_MAPPING"
            }
          },
          "subject" : {
            source = {
              type = "NO_MAPPING"
            }
          },
          "SAML_AUTHN_CTX" : {
            source = {
              type = "NO_MAPPING"
            }
          }
        }
        issuance_criteria = {
          conditional_criteria = [
            {
              error_result = "error",
              source = {
                type = "ASSERTION"
              },
              attribute_name = "SAML_SUBJECT",
              condition      = "EQUALS",
              value          = "value"
            }
          ]
        }

        authentication_policy_contract_ref = {
          id = "default"
        }

        restrict_virtual_server_ids   = true
        restricted_virtual_server_ids = ["virtual_server_id"]
      }
    ]
    assertions_signed   = false
    sign_authn_requests = false

    sso_oauth_mapping = {
      attribute_sources = [
        {
          jdbc_attribute_source = {
            data_store_ref = {
              id = "ProvisionerDS",
            },
            description = "JDBC",
            schema      = "INFORMATION_SCHEMA",
            table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS",
            filter      = "example",
            column_names = [
              "GRANTEE"
            ]
          }
        }
      ]
      attribute_contract_fulfillment = {
        "USER_NAME" : {
          source = {
            type = "NO_MAPPING"
          }
        },
        "USER_KEY" : {
          source = {
            type = "NO_MAPPING"
          }
        }
      },
    }
  }

  credentials = {
    certs = [{
      x509_file = {
        # No id supplied for the cert
        file_data = "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjugAwIBAgIQPEkZGqCnSpsZf0jWCWxJ5jALBgkqhkiG9w0BAQswRzEL\nMAkGA1UEBgwCVVMxHDAaBgNVBAoME0V4YW1wbGUgQ29ycG9yYXRpb24xGjAYBgNV\nBAMMEUV4YW1wbGUgQXV0aG9yaXR5MB4XDTI0MDEwMTAwMDAwMFoXDTQzMTIyNzAw\nMDAwMFowRzELMAkGA1UEBgwCVVMxHDAaBgNVBAoME0V4YW1wbGUgQ29ycG9yYXRp\nb24xGjAYBgNVBAMMEUV4YW1wbGUgQXV0aG9yaXR5MIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEAzdxT13IA0xJ8rB1hkqxa/JTTrmLnNjTRnVJdwagGKThQ\nxpWqh0DchNMqXTaFNGgxia3hPB53ew7nEMuIf+Qfq4meexKL3yRg86Ng56BrGbsu\nK6z4ptRpdnmsxgkpEfGytdmUFkPXAGE6j4Td/UrAWByz7C9yl7qzFYeorWq5nABc\nIiOlLxBYXX3fOu3a44SNexNgl5dDJAtn8mosQ19wJcjm08fKRqHeWYvBV99kQlhW\na7WiTxdrbUZOrUMHYRuKO/JD732dcpnsar9HfjQi+PH3gCgw4NJNuBKzLv6t8DzZ\nnNxaiKgZ+5cxdhhRAe98MF0QeTbymjVLyoFBpMrRDQIDAQABoz0wOzAdBgNVHQ4E\nFgQUGNJsUqA63OVS8ouwVUkzaEP5vawwDAYDVR0TBAUwAwEB/zAMBgNVHQ8EBQMD\nBwYAMAsGCSqGSIb3DQEBCwOCAQEAGFvsWv35ipg0NNnq0x+e7Gtugn9OBhxkeTWo\nQ1IUR7CL9zMRdlErIx5waptJhlPZFZANVpuvYa+yRz7oz2txH8yf/0N+F0bTeNU/\nqZHenvp9RXzimxTFDoCkx7ESpW9b7IKSSZA6Zut6w7XzJeXRrNKfCSSrUGPfkCq4\nhOtAm9QzUVE7eJ5a7T3+O50gZdoxjdojPhh9h5E1b+bmexrfQKlVl/gL+KPacBJD\nbSxbiKECt5QGRdDGFFfoInhK1RiW7a/hQBhMWRsMiOFtu0YpfxfwIyIaK5QfHZBC\nCC7JaJKg19njrnkjfmiBGoev7XiYWYt/WvYAiZR4nJn/cFrW1A==\n-----END CERTIFICATE-----"
      }
      active_verification_cert    = true
      encryption_cert             = true
      primary_verification_cert   = true
      secondary_verification_cert = false
    }]

    inbound_back_channel_auth = {
      http_basic_credentials = {
        username = "admin"
        password = "2FederateM0re!"
      },
      digital_signature = true
      require_ssl       = true
      certs = [{
        x509_file = {
          # No id supplied for the cert
          file_data = "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjugAwIBAgIQPEkZGqCnSpsZf0jWCWxJ5jALBgkqhkiG9w0BAQswRzEL\nMAkGA1UEBgwCVVMxHDAaBgNVBAoME0V4YW1wbGUgQ29ycG9yYXRpb24xGjAYBgNV\nBAMMEUV4YW1wbGUgQXV0aG9yaXR5MB4XDTI0MDEwMTAwMDAwMFoXDTQzMTIyNzAw\nMDAwMFowRzELMAkGA1UEBgwCVVMxHDAaBgNVBAoME0V4YW1wbGUgQ29ycG9yYXRp\nb24xGjAYBgNVBAMMEUV4YW1wbGUgQXV0aG9yaXR5MIIBIjANBgkqhkiG9w0BAQEF\nAAOCAQ8AMIIBCgKCAQEAzdxT13IA0xJ8rB1hkqxa/JTTrmLnNjTRnVJdwagGKThQ\nxpWqh0DchNMqXTaFNGgxia3hPB53ew7nEMuIf+Qfq4meexKL3yRg86Ng56BrGbsu\nK6z4ptRpdnmsxgkpEfGytdmUFkPXAGE6j4Td/UrAWByz7C9yl7qzFYeorWq5nABc\nIiOlLxBYXX3fOu3a44SNexNgl5dDJAtn8mosQ19wJcjm08fKRqHeWYvBV99kQlhW\na7WiTxdrbUZOrUMHYRuKO/JD732dcpnsar9HfjQi+PH3gCgw4NJNuBKzLv6t8DzZ\nnNxaiKgZ+5cxdhhRAe98MF0QeTbymjVLyoFBpMrRDQIDAQABoz0wOzAdBgNVHQ4E\nFgQUGNJsUqA63OVS8ouwVUkzaEP5vawwDAYDVR0TBAUwAwEB/zAMBgNVHQ8EBQMD\nBwYAMAsGCSqGSIb3DQEBCwOCAQEAGFvsWv35ipg0NNnq0x+e7Gtugn9OBhxkeTWo\nQ1IUR7CL9zMRdlErIx5waptJhlPZFZANVpuvYa+yRz7oz2txH8yf/0N+F0bTeNU/\nqZHenvp9RXzimxTFDoCkx7ESpW9b7IKSSZA6Zut6w7XzJeXRrNKfCSSrUGPfkCq4\nhOtAm9QzUVE7eJ5a7T3+O50gZdoxjdojPhh9h5E1b+bmexrfQKlVl/gL+KPacBJD\nbSxbiKECt5QGRdDGFFfoInhK1RiW7a/hQBhMWRsMiOFtu0YpfxfwIyIaK5QfHZBC\nCC7JaJKg19njrnkjfmiBGoev7XiYWYt/WvYAiZR4nJn/cFrW1A==\n-----END CERTIFICATE-----"
        }
        active_verification_cert    = true
        encryption_cert             = false
        primary_verification_cert   = true
        secondary_verification_cert = false
      }]
    }

    decryption_key_pair_ref = {
      id = "419x9yg43rlawqwq9v6az997k"
    }

    signing_settings = {
      signing_key_pair_ref = {
        id = "419x9yg43rlawqwq9v6az997k"
      }
      algorithm                    = "SHA256withRSA"
      include_cert_in_signature    = false
      include_raw_key_in_signature = false
    }

    block_encryption_algorithm = "AES_128"
    key_transport_algorithm    = "RSA_OAEP"

    outbound_back_channel_auth = {
      http_basic_credentials = {
        username = "Administrator"
        password = "2FederateM0re!"
      }
      digital_signature     = false
      validate_partner_cert = true
    }
  }

  idp_oauth_grant_attribute_mapping = {
    idp_oauth_attribute_contract = {
      extended_attributes = []
    }
    access_token_manager_mappings = [
      {
        attribute_sources = [
          {
            custom_attribute_source = {
              data_store_ref = {
                id = "customDataStore"
              }
              description = "APIStubs"
              filter_fields = [
                {
                  name = "Authorization Header"
                },
                {
                  name = "Body"
                },
                {
                  name  = "Resource Path"
                  value = "/users/extid"
                },
              ]
              id = "APIStubs"
            }
          },
        ]
        attribute_contract_fulfillment = {
          "Username" = {
            source = {
              type = "NO_MAPPING"
            }
          }
          "OrgName" = {
            source = {
              type = "NO_MAPPING"
            }
          }
        }
        access_token_manager_ref = {
          id = pingfederate_oauth_access_token_manager.idpConnCertNoIdAtm.id
        }
      }
    ]
  }

  ws_trust = {
    attribute_contract = {
      core_attributes = [
        {
          name   = "TOKEN_SUBJECT"
          masked = false
        }
      ]
      extended_attributes = [
        {
          name   = "test"
          masked = false
        }
      ]
    }
    token_generator_mappings = [
      {
        attribute_contract_fulfillment = {
          "SAML_SUBJECT" = {
            source = {
              type = "NO_MAPPING"
            }
          }
        }
        attribute_sources = [
          {
            custom_attribute_source = {
              data_store_ref = {
                id = "customDataStore"
              }
              description = "APIStubs"
              filter_fields = [
                {
                  name = "Authorization Header"
                },
                {
                  name = "Body"
                },
                {
                  name  = "Resource Path"
                  value = "/users/extid"
                },
              ]
              id = "APIStubs"
            }
          },
        ]
        sp_token_generator_ref = {
          id = "tokengenerator"
        }
        default_mapping = true
      }
    ]
    generate_local_token = true
  }

  inbound_provisioning = {
    group_support = false

    user_repository = {
      identity_store = {
        identity_store_provisioner_ref = {
          id = "identityStoreProvisioner"
        }
      }
    }

    custom_schema = {
      namespace  = "urn:scim:schemas:extension:custom:1.0"
      attributes = []
    }

    users = {
      write_users = {
        attribute_fulfillment = {
          "username" = {
            source = {
              type = "TEXT"
            }
            value = "username"
          }
        }
      }

      read_users = {
        attribute_contract = {
          extended_attributes = [
            {
              name   = "userName"
              masked = false
            }
          ]
        }
        attributes = []
        attribute_fulfillment = {
          "userName" = {
            source = {
              type = "TEXT"
            }
            value = "username"
          }
        }
      }
    }
  }
  # Ensures this resource will be updated before deleting the oauth access token manager
  lifecycle {
    create_before_destroy = true
  }
}
  `, accesstokenmanager.AccessTokenManagerTestHCL("idpConnCertNoIdAtm"),
		wsTrustStsConnId)
}
