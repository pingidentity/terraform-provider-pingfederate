package idpspconnection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const spConnOutboundProvisionId = "outboundspconn"

var pingOneConnection, pingOneEnvironment string

func TestAccIdpSpConnection_OutboundProvisionMinimalMaximal(t *testing.T) {
	pingOneConnection = os.Getenv("PF_TF_P1_CONNECTION_ID")
	pingOneEnvironment = os.Getenv("PF_TF_P1_CONNECTION_ENV_ID")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if pingOneConnection == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ID must be set for the TestAccAuthenticationPoliciesFragment acceptance test")
			}
			if pingOneEnvironment == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ENV_ID must be set for the TestAccAuthenticationPoliciesFragment acceptance test")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: idpSpConnection_OutboundProvisionCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: idpSpConnection_OutboundProvisionMinimalHCL(),
				Check:  idpSpConnection_CheckComputedValuesOutboundProvisionMinimal(),
			},
			{
				// Delete the minimal model
				Config:  idpSpConnection_OutboundProvisionMinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: idpSpConnection_OutboundProvisionCompleteHCL(),
				Check:  idpSpConnection_CheckComputedValuesOutboundProvisionComplete(),
			},
			{
				// Back to minimal model
				Config: idpSpConnection_OutboundProvisionMinimalHCL(),
				Check:  idpSpConnection_CheckComputedValuesOutboundProvisionMinimal(),
			},
			{
				// Back to complete model
				Config: idpSpConnection_OutboundProvisionCompleteHCL(),
				Check:  idpSpConnection_CheckComputedValuesOutboundProvisionComplete(),
			},
			{
				// Test importing the resource
				Config:            idpSpConnection_OutboundProvisionCompleteHCL(),
				ResourceName:      "pingfederate_idp_sp_connection.example",
				ImportStateId:     spConnOutboundProvisionId,
				ImportState:       true,
				ImportStateVerify: true,
				// There is currently an issue where values of target_settings are not imported
				ImportStateVerifyIgnore: []string{"outbound_provision.target_settings"},
			},
		},
	})
}

// Minimal HCL with only required values set
func idpSpConnection_OutboundProvisionMinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "example" {
  connection_id                             = "%s"
  credentials = {
  }
  entity_id                 = "partnermin:entity:id"
  name                      = "minimalOutbound"
  outbound_provision = {
    channels = [
      {
        active = false
        attribute_mapping = [
          {
            field_name = "ZIPCode"
            saas_field_info = {
              attribute_names = ["postalCode"]
            }
          },
          {
            field_name = "accountID"
            saas_field_info = {
            }
          },
          {
            field_name = "authoritativeIdp"
            saas_field_info = {
            }
          },
          {
            field_name = "city"
            saas_field_info = {
              attribute_names = ["l"]
            }
          },
          {
            field_name = "country"
            saas_field_info = {
            }
          },
          {
            field_name = "email"
            saas_field_info = {
              attribute_names = ["mail"]
            }
          },
          {
            field_name = "externalID"
            saas_field_info = {
            }
          },
          {
            field_name = "firstName"
            saas_field_info = {
              attribute_names = ["givenName"]
            }
          },
          {
            field_name = "forceChangePassword"
            saas_field_info = {
            }
          },
          {
            field_name = "fullName"
            saas_field_info = {
            }
          },
          {
            field_name = "honorificPrefix"
            saas_field_info = {
            }
          },
          {
            field_name = "honorificSuffix"
            saas_field_info = {
            }
          },
          {
            field_name = "jobTitle"
            saas_field_info = {
              attribute_names = ["title"]
            }
          },
          {
            field_name = "lastName"
            saas_field_info = {
              attribute_names = ["sn"]
            }
          },
          {
            field_name = "locale"
            saas_field_info = {
            }
          },
          {
            field_name = "mfaDeviceEmail1"
            saas_field_info = {
              attribute_names = ["mail"]
            }
          },
          {
            field_name = "mfaDeviceEmail2"
            saas_field_info = {
            }
          },
          {
            field_name = "mfaDeviceEmail3"
            saas_field_info = {
            }
          },
          {
            field_name = "mfaDeviceSms1"
            saas_field_info = {
              attribute_names = ["mobile"]
            }
          },
          {
            field_name = "mfaDeviceSms2"
            saas_field_info = {
            }
          },
          {
            field_name = "mfaDeviceSms3"
            saas_field_info = {
            }
          },
          {
            field_name = "mfaDeviceVoice1"
            saas_field_info = {
              attribute_names = ["mobile"]
            }
          },
          {
            field_name = "mfaDeviceVoice2"
            saas_field_info = {
            }
          },
          {
            field_name = "mfaDeviceVoice3"
            saas_field_info = {
            }
          },
          {
            field_name = "mfaEnabled"
            saas_field_info = {
            }
          },
          {
            field_name = "middleName"
            saas_field_info = {
            }
          },
          {
            field_name = "mobilePhone"
            saas_field_info = {
              attribute_names = ["mobile"]
            }
          },
          {
            field_name = "nickname"
            saas_field_info = {
            }
          },
          {
            field_name = "password"
            saas_field_info = {
            }
          },
          {
            field_name = "populationID"
            saas_field_info = {
            }
          },
          {
            field_name = "preferredLanguage"
            saas_field_info = {
            }
          },
          {
            field_name = "primaryPhone"
            saas_field_info = {
              attribute_names = ["telephoneNumber"]
            }
          },
          {
            field_name = "profileImage"
            saas_field_info = {
            }
          },
          {
            field_name = "stateRegion"
            saas_field_info = {
              attribute_names = ["st"]
            }
          },
          {
            field_name = "streetAddress"
            saas_field_info = {
              attribute_names = ["streetAddress"]
            }
          },
          {
            field_name = "timezone"
            saas_field_info = {
            }
          },
          {
            field_name = "userType"
            saas_field_info = {
            }
          },
          {
            field_name = "username"
            saas_field_info = {
              attribute_names = ["uid"]
            }
          },
        ]
        channel_source = {
          account_management_settings = {
            account_status_algorithm      = "ACCOUNT_STATUS_ALGORITHM_FLAG"
            account_status_attribute_name = "nsaccountlock"
            flag_comparison_value         = jsonencode(true)
          }
          base_dn = "dc=example,dc=com"
          change_detection_settings = {
            changed_users_algorithm   = "TIMESTAMP_NO_NEGATION"
            group_object_class        = "groupOfUniqueNames"
            time_stamp_attribute_name = "modifyTimestamp"
            user_object_class         = "inetOrgPerson"
          }
          data_source = {
            id = "pingdirectory"
          }
          group_membership_detection = {
            group_member_attribute_name    = "uniqueMember"
          }
          guid_attribute_name = "entryUUID"
          guid_binary         = false
          user_source_location = {
            group_dn      = "o=group,dc=example,dc=com"
          }
        }
        name        = "min"
      },
    ]
    target_settings = [
      {
        name            = "PINGONE_ENVIRONMENT"
        value           = "%s|%s"
      },
    ]
    type = "PingOne"
  }
}
`, spConnOutboundProvisionId, pingOneConnection, pingOneEnvironment)
}

// Maximal HCL with all values set where possible
func idpSpConnection_OutboundProvisionCompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "example" {
  active                                    = false
  application_icon_url                      = "https://example.com/logo.png"
  application_name                          = "App"
  base_url                                  = "https://api.pingone.com/v1"
  connection_id                             = "%s"
  connection_target_type                    = "STANDARD"
  contact_info = {
    company    = "Example Corp"
    email      = "name@example.com"
    first_name = "Name"
    last_name  = "Last"
    phone      = jsonencode(5555555555)
  }
  credentials = {
    certs = [
      {
        active_verification_cert    = true
        encryption_cert             = false
        primary_verification_cert   = true
        secondary_verification_cert = false
        x509_file = {
          file_data       = "-----BEGIN CERTIFICATE-----\nMIIDUTCCAjugAwIBAgIQPEkZGqCnSpsZf0jWCWxJ5jALBgkqhkiG9w0BAQswRzELMAkGA1UEBgwC\nVVMxHDAaBgNVBAoME0V4YW1wbGUgQ29ycG9yYXRpb24xGjAYBgNVBAMMEUV4YW1wbGUgQXV0aG9y\naXR5MB4XDTI0MDEwMTAwMDAwMFoXDTQzMTIyNzAwMDAwMFowRzELMAkGA1UEBgwCVVMxHDAaBgNV\nBAoME0V4YW1wbGUgQ29ycG9yYXRpb24xGjAYBgNVBAMMEUV4YW1wbGUgQXV0aG9yaXR5MIIBIjAN\nBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzdxT13IA0xJ8rB1hkqxa/JTTrmLnNjTRnVJdwagG\nKThQxpWqh0DchNMqXTaFNGgxia3hPB53ew7nEMuIf+Qfq4meexKL3yRg86Ng56BrGbsuK6z4ptRp\ndnmsxgkpEfGytdmUFkPXAGE6j4Td/UrAWByz7C9yl7qzFYeorWq5nABcIiOlLxBYXX3fOu3a44SN\nexNgl5dDJAtn8mosQ19wJcjm08fKRqHeWYvBV99kQlhWa7WiTxdrbUZOrUMHYRuKO/JD732dcpns\nar9HfjQi+PH3gCgw4NJNuBKzLv6t8DzZnNxaiKgZ+5cxdhhRAe98MF0QeTbymjVLyoFBpMrRDQID\nAQABoz0wOzAdBgNVHQ4EFgQUGNJsUqA63OVS8ouwVUkzaEP5vawwDAYDVR0TBAUwAwEB/zAMBgNV\nHQ8EBQMDBwYAMAsGCSqGSIb3DQEBCwOCAQEAGFvsWv35ipg0NNnq0x+e7Gtugn9OBhxkeTWoQ1IU\nR7CL9zMRdlErIx5waptJhlPZFZANVpuvYa+yRz7oz2txH8yf/0N+F0bTeNU/qZHenvp9RXzimxTF\nDoCkx7ESpW9b7IKSSZA6Zut6w7XzJeXRrNKfCSSrUGPfkCq4hOtAm9QzUVE7eJ5a7T3+O50gZdox\njdojPhh9h5E1b+bmexrfQKlVl/gL+KPacBJDbSxbiKECt5QGRdDGFFfoInhK1RiW7a/hQBhMWRsM\niOFtu0YpfxfwIyIaK5QfHZBCCC7JaJKg19njrnkjfmiBGoev7XiYWYt/WvYAiZR4nJn/cFrW1A==\n-----END CERTIFICATE-----\n"
          id              = "f1f4rt7f288buvljrwl3sqsle"
        }
      },
    ]
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
  entity_id                 = "PingOne Connector"
  extended_properties = {
    authNexp = {
      values = ["val1"]
    }
    useAuthnApi = {
      values = ["val2"]
    }
  }
  logging_mode             = "ENHANCED"
  name                     = "PingOne Connector"
  outbound_provision = {
    channels = [
      {
        active = false
        attribute_mapping = [
          {
            field_name = "ZIPCode"
            saas_field_info = {
              attribute_names = ["postalCode"]
              character_case  = "NONE"
              create_only     = false
              default_value   = "12345"
              expression      = "post"
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "accountID"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "authoritativeIdp"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "city"
            saas_field_info = {
              attribute_names = ["l"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "country"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "email"
            saas_field_info = {
              attribute_names = ["mail"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "externalID"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "firstName"
            saas_field_info = {
              attribute_names = ["givenName"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "forceChangePassword"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "fullName"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "honorificPrefix"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "honorificSuffix"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "jobTitle"
            saas_field_info = {
              attribute_names = ["title"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "lastName"
            saas_field_info = {
              attribute_names = ["sn"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "locale"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = "en-US"
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceEmail1"
            saas_field_info = {
              attribute_names = ["mail"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceEmail2"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceEmail3"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceSms1"
            saas_field_info = {
              attribute_names = ["mobile"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceSms2"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceSms3"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceVoice1"
            saas_field_info = {
              attribute_names = ["mobile"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceVoice2"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaDeviceVoice3"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mfaEnabled"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "middleName"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "mobilePhone"
            saas_field_info = {
              attribute_names = ["mobile"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "nickname"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "password"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "populationID"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = "374fdb3c-4e94-4547-838a-0c200b9a7c70"
              expression      = "pop"
              masked          = true
              parser          = "EXTRACT_CN_FROM_DN"
              trim            = false
            }
          },
          {
            field_name = "preferredLanguage"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = "en-us"
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "primaryPhone"
            saas_field_info = {
              attribute_names = ["telephoneNumber"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "profileImage"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "stateRegion"
            saas_field_info = {
              attribute_names = ["st"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "streetAddress"
            saas_field_info = {
              attribute_names = ["streetAddress"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "timezone"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "userType"
            saas_field_info = {
              attribute_names = []
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
          {
            field_name = "username"
            saas_field_info = {
              attribute_names = ["uid"]
              character_case  = "NONE"
              create_only     = false
              default_value   = null
              expression      = null
              masked          = false
              parser          = "NONE"
              trim            = false
            }
          },
        ]
        channel_source = {
          account_management_settings = {
            account_status_algorithm      = "ACCOUNT_STATUS_ALGORITHM_FLAG"
            account_status_attribute_name = "nsaccountlock"
            default_status                = true
            flag_comparison_status        = false
            flag_comparison_value         = jsonencode(true)
          }
          base_dn = "dc=example,dc=com"
          change_detection_settings = {
            changed_users_algorithm   = "TIMESTAMP_NO_NEGATION"
            group_object_class        = "groupOfUniqueNames"
            time_stamp_attribute_name = "modifyTimestamp"
            user_object_class         = "inetOrgPerson"
          }
          data_source = {
            id = "pingdirectory"
          }
          group_membership_detection = {
            group_member_attribute_name    = "uniqueMember"
            member_of_group_attribute_name = "group"
          }
          group_source_location = {
            filter        = "group=second"
            group_dn      = "o=group,dc=example,dc=com"
            nested_search = true
          }
          guid_attribute_name = "entryUUID"
          guid_binary         = false
          user_source_location = {
            filter        = "group=test"
            group_dn      = "o=group,dc=example,dc=com"
            nested_search = false
          }
        }
        max_threads = 1
        name        = "chan1"
        timeout     = 60
      },
    ]
    target_settings = [
      {
        name            = "CREATE_USERS_PROV_OPT"
        value           = "true"
      },
      {
        name            = "DEFAULT_AUTH_METHOD"
        value           = "Email 1"
      },
      {
        name            = "MFA_USER_DEVICE_MANAGEMENT"
        value           = "Merge with devices in PingOne"
      },
      {
        name            = "PINGONE_ENVIRONMENT"
        value           = "%s|%s"
      },
      {
        name            = "PROVISION_DISABLED_USERS_PROV_OPT"
        value           = "true"
      },
      {
        name            = "Provisioning Options"
        value           = "opt"
      },
      {
        name            = "REMOVE_ACTION"
        value           = "Disable"
      },
      {
        name            = "REMOVE_USERS_PROV_OPT"
        value           = "true"
      },
      {
        name            = "UPDATE_USERS_PROV_OPT"
        value           = "true"
      },
    ]
    type = "PingOne"
  }
  virtual_entity_ids = ["ex1", "ex2"]
}
`, spConnOutboundProvisionId, pingOneConnection, pingOneEnvironment)
}

// Validate any computed values when applying minimal HCL
func idpSpConnection_CheckComputedValuesOutboundProvisionMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "active", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "connection_target_type", "STANDARD"),
		resource.TestCheckResourceAttrSet("pingfederate_idp_sp_connection.example", "creation_date"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "credentials.certs.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "id", spConnOutboundProvisionId),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "logging_mode", "STANDARD"),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.attribute_mapping.*",
			map[string]string{
				"field_name":                    "ZIPCode",
				"saas_field_info.create_only":   "false",
				"saas_field_info.default_value": "",
				"saas_field_info.expression":    "",
				"saas_field_info.masked":        "false",
				"saas_field_info.parser":        "NONE",
				"saas_field_info.trim":          "false",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.attribute_mapping.*",
			map[string]string{
				"field_name":                        "fullName",
				"saas_field_info.attribute_names.#": "0",
				"saas_field_info.create_only":       "false",
				"saas_field_info.default_value":     "",
				"saas_field_info.expression":        "",
				"saas_field_info.masked":            "false",
				"saas_field_info.parser":            "NONE",
				"saas_field_info.trim":              "false",
			},
		),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.channel_source.account_management_settings.default_status", "true"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.channel_source.account_management_settings.flag_comparison_status", "true"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.channel_source.group_source_location.nested_search", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.channel_source.user_source_location.nested_search", "false"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.max_threads", "1"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.channels.0.timeout", "60"),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.target_settings_all.#", "9"),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_idp_sp_connection.example", "outbound_provision.target_settings_all.*",
			map[string]string{
				"name":  "CREATE_USERS_PROV_OPT",
				"value": "true",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_idp_sp_connection.example", "outbound_provision.target_settings_all.*",
			map[string]string{
				"name":  "PROVISION_DISABLED_USERS_PROV_OPT",
				"value": "true",
			},
		),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "virtual_entity_ids.#", "0"),
	)
}

// Validate any computed values when applying complete HCL
func idpSpConnection_CheckComputedValuesOutboundProvisionComplete() resource.TestCheckFunc {
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
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "id", spConnOutboundProvisionId),
		resource.TestCheckResourceAttr("pingfederate_idp_sp_connection.example", "outbound_provision.target_settings_all.#", "9"),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_idp_sp_connection.example", "outbound_provision.target_settings_all.*",
			map[string]string{
				"name":  "CREATE_USERS_PROV_OPT",
				"value": "true",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_idp_sp_connection.example", "outbound_provision.target_settings_all.*",
			map[string]string{
				"name":  "PROVISION_DISABLED_USERS_PROV_OPT",
				"value": "true",
			},
		),
	)
}

// Test that any objects created by the test are destroyed
func idpSpConnection_OutboundProvisionCheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(acctest.TestBasicAuthContext(), spConnOutboundProvisionId).Execute()
	if err == nil {
		return fmt.Errorf("sp_idp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
