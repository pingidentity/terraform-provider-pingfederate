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

const spConnWsTrustId = "wstrustspconn"

func TestAccIdpSpConnection_WsTrustMinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpSpConnection_WsTrustCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpSpConnection_WsTrustMinimalHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsTrustMinimal(),
			},
			{
				// Delete the minimal model
				Config:  idpSpConnection_WsTrustMinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: idpSpConnection_WsTrustCompleteHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsTrustComplete(),
			},
			{
				// Back to minimal model
				Config: idpSpConnection_WsTrustMinimalHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsTrustMinimal(),
			},
			{
				// Back to complete model
				Config: idpSpConnection_WsTrustCompleteHCL(),
				Check:  idpSpConnection_CheckComputedValuesWsTrustComplete(),
			},
			{
				// Test importing the resource
				Config:            idpSpConnection_WsTrustCompleteHCL(),
				ResourceName:      "pingfederate_idp_sp_connection.example",
				ImportStateId:     spConnWsTrustId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func idpSpConnection_WsTrustDependencyHCL() string {
	return `
resource "pingfederate_idp_sts_request_parameters_contract" "example" {
  contract_id = "mycontract"
  name        = "My Contract"
  parameters = [
    "firstparam",
    "secondparam",
    "thirdparam"
  ]
}
	`
}

// Minimal HCL with only required values set
func idpSpConnection_WsTrustMinimalHCL() string {
	return fmt.Sprintf(`
%s

resource "pingfederate_idp_sp_connection" "example" {
  connection_id = "%s"
  credentials = {
    signing_settings = {
      algorithm = "SHA256withRSA"
      signing_key_pair_ref = {
        id = "rsaprevious"
      }
    }
  }
  entity_id = "partner:entity:id"
  name      = "minimalWsTrust"
  ws_trust = {
    attribute_contract = {
      core_attributes = [
        {
          name = "TOKEN_SUBJECT"
        },
      ]
    }
    partner_service_ids = ["serviceid"]
    token_processor_mappings = [
      {
        attribute_contract_fulfillment = {
          TOKEN_SUBJECT = {
            source = {
              type = "CONTEXT"
            }
            value = "ClientIp"
          }
        }
        idp_token_processor_ref = {
          id = "UsernameTokenProcessor"
        }
      },
    ]
  }
}
`, idpSpConnection_WsTrustDependencyHCL(), spConnWsTrustId)
}

// Maximal HCL with all values set where possible
func idpSpConnection_WsTrustCompleteHCL() string {
	return fmt.Sprintf(`
%s

resource "pingfederate_idp_sp_connection" "example" {
  active                 = false
  application_icon_url   = "https://example.com/logo.png"
  application_name       = "App"
  base_url               = "https://example.com"
  connection_id          = "%s"
  connection_target_type = "STANDARD"
  contact_info = {
    company    = "Example Corp"
    email      = "name@example.com"
    first_name = "Name"
    last_name  = "Last"
    phone      = jsonencode(5555555555)
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
          file_data = "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjugAwIBAgIQPEkZGqCnSpsZf0jWCWxJ5jALBgkqhkiG9w0BAQswRzELMAkGA1UEBgwC\nVVMxHDAaBgNVBAoME0V4YW1wbGUgQ29ycG9yYXRpb24xGjAYBgNVBAMMEUV4YW1wbGUgQXV0aG9y\naXR5MB4XDTI0MDEwMTAwMDAwMFoXDTQzMTIyNzAwMDAwMFowRzELMAkGA1UEBgwCVVMxHDAaBgNV\nBAoME0V4YW1wbGUgQ29ycG9yYXRpb24xGjAYBgNVBAMMEUV4YW1wbGUgQXV0aG9yaXR5MIIBIjAN\nBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzdxT13IA0xJ8rB1hkqxa/JTTrmLnNjTRnVJdwagG\nKThQxpWqh0DchNMqXTaFNGgxia3hPB53ew7nEMuIf+Qfq4meexKL3yRg86Ng56BrGbsuK6z4ptRp\ndnmsxgkpEfGytdmUFkPXAGE6j4Td/UrAWByz7C9yl7qzFYeorWq5nABcIiOlLxBYXX3fOu3a44SN\nexNgl5dDJAtn8mosQ19wJcjm08fKRqHeWYvBV99kQlhWa7WiTxdrbUZOrUMHYRuKO/JD732dcpns\nar9HfjQi+PH3gCgw4NJNuBKzLv6t8DzZnNxaiKgZ+5cxdhhRAe98MF0QeTbymjVLyoFBpMrRDQID\nAQABoz0wOzAdBgNVHQ4EFgQUGNJsUqA63OVS8ouwVUkzaEP5vawwDAYDVR0TBAUwAwEB/zAMBgNV\nHQ8EBQMDBwYAMAsGCSqGSIb3DQEBCwOCAQEAGFvsWv35ipg0NNnq0x+e7Gtugn9OBhxkeTWoQ1IU\nR7CL9zMRdlErIx5waptJhlPZFZANVpuvYa+yRz7oz2txH8yf/0N+F0bTeNU/qZHenvp9RXzimxTF\nDoCkx7ESpW9b7IKSSZA6Zut6w7XzJeXRrNKfCSSrUGPfkCq4hOtAm9QzUVE7eJ5a7T3+O50gZdox\njdojPhh9h5E1b+bmexrfQKlVl/gL+KPacBJDbSxbiKECt5QGRdDGFFfoInhK1RiW7a/hQBhMWRsM\niOFtu0YpfxfwIyIaK5QfHZBCCC7JaJKg19njrnkjfmiBGoev7XiYWYt/WvYAiZR4nJn/cFrW1A==\n-----END CERTIFICATE-----\n"
          id        = "f1f4rt7f288buvljrwl3sqsle"
        }
      },
    ]
    key_transport_algorithm = "RSA_OAEP"
    signing_settings = {
      algorithm = "SHA256withRSA"
      alternative_signing_key_pair_refs = [
        {
          id = "ec256previous"
        },
      ]
      include_cert_in_signature    = true
      include_raw_key_in_signature = true
      signing_key_pair_ref = {
        id = "rsaprevious"
      }
    }
  }
  default_virtual_entity_id = "ex1"
  entity_id                 = "partner:entity:id"
  extended_properties = {
    authNexp = {
      values = ["val1"]
    }
    useAuthnApi = {
      values = ["val2"]
    }
  }
  logging_mode       = "ENHANCED"
  name               = "minimalWsTrust"
  virtual_entity_ids = ["ex1", "ex2"]
  ws_trust = {
    abort_if_not_fulfilled_from_request = true
    attribute_contract = {
      core_attributes = [
        {
          name = "TOKEN_SUBJECT"
        },
      ]
      extended_attributes = [
      ]
    }
    default_token_type      = "SAML20"
    encrypt_saml2_assertion = true
    generate_key            = true
    message_customizations = [
    ]
    minutes_after            = 30
    minutes_before           = 5
    oauth_assertion_profiles = true
    partner_service_ids      = ["anotherId", "serviceid"]
    request_contract_ref = {
      id = pingfederate_idp_sts_request_parameters_contract.example.id
    }
    token_processor_mappings = [
      {
        attribute_contract_fulfillment = {
          TOKEN_SUBJECT = {
            source = {
              id   = null
              type = "CONTEXT"
            }
            value = "ClientIp"
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
        idp_token_processor_ref = {
          id = "UsernameTokenProcessor"
        }
        issuance_criteria = {
      conditional_criteria = [
            {
              attribute_name = "TOKEN_SUBJECT"
              condition      = "MULTIVALUE_CONTAINS_DN"
              source = {
                type = "MAPPED_ATTRIBUTES"
              }
              value = "cn=Example,dc=example,dc=com"
            },
      ]
          expression_criteria = null
        }
        restricted_virtual_entity_ids = []
      },
    ]
  }
}
`, idpSpConnection_WsTrustDependencyHCL(), spConnWsTrustId)
}

// Validate any computed values when applying minimal HCL
func idpSpConnection_CheckComputedValuesWsTrustMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "active", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "connection_target_type", "STANDARD"),
		resource.TestCheckResourceAttrSet("pingfederate_idp_sp_connection.example", "creation_date"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.signing_settings.include_cert_in_signature", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.signing_settings.include_raw_key_in_signature", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "id", spConnWsTrustId),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "logging_mode", "STANDARD"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "virtual_entity_ids.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.attribute_contract.core_attributes.0.namespace", ""),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.encrypt_saml2_assertion", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.generate_key", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.message_customizations.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.minutes_after", "30"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.minutes_before", "5"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.oauth_assertion_profiles", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.token_processor_mappings.0.attribute_sources.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.token_processor_mappings.0.issuance_criteria.conditional_criteria.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.token_processor_mappings.0.restricted_virtual_entity_ids.#", "0"),
	)
}

// Validate any computed values when applying complete HCL
func idpSpConnection_CheckComputedValuesWsTrustComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("pingfederate_idp_sp_connection.example", "creation_date"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.expires", "2043-12-27T00:00:00Z"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.id", "f1f4rt7f288buvljrwl3sqsle"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.serial_number", "80133226587660155953237711066254559718"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.sha1_fingerprint", "D2B4B2033511D50BABE289E0AF2C17B8DA15FCC3"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.sha256_fingerprint", "6E61905D5223D6667B68D8E600E779B1F13DA404041A8BA7CDF07DB01179A897"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.status", "VALID"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.subject_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.valid_from", "2024-01-01T00:00:00Z"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.0.cert_view.version", "3"),
		resource.TestCheckResourceAttrSet("pingfederate_idp_sp_connection.example", "credentials.certs.0.x509_file.formatted_file_data"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "id", spConnWsTrustId),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "ws_trust.attribute_contract.core_attributes.0.namespace", ""),
	)
}

// Test that any objects created by the test are destroyed
func idpSpConnection_WsTrustCheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(acctest.TestBasicAuthContext(), spConnWsTrustId).Execute()
	if err == nil {
		return fmt.Errorf("pingfederate_idp_sp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
