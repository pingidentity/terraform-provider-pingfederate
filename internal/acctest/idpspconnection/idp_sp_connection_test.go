package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const (
	spConnectionId = "spConnId"
	resourceType   = "IdP SP Connection"
)

func TestAccIdpSpConnection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckSpConnectionDestroy,
		Steps: []resource.TestStep{
			{
				// Outbound provision connection, minimal
				Config: testAccSpConnectionOutboundProvision(spConnectionId),
				Check:  testAccCheckExpectedSpConnectionAttributesOutboundProvision(),
			},
			{
				// Browser SSO SAML connection minimal
				Config: testAccSpConnectionBrowserSso(spConnectionId, false),
				Check:  testAccCheckExpectedSpConnectionAttributesBrowserSSO(false),
			},
			{
				// Browser SSO WsFed connection minimal
				Config: testAccSpConnectionBrowserSso(spConnectionId, true),
				Check:  testAccCheckExpectedSpConnectionAttributesBrowserSSO(true),
			},
			{
				// WS Trust connection, minimal
				Config: testAccSpConnectionWsTrust(spConnectionId),
				Check:  testAccCheckExpectedSpConnectionAttributesWsTrust(),
			},
			{
				// Complete connection with all three types
				Config: testAccSpConnectionComplete(spConnectionId),
				Check:  testAccCheckExpectedSpConnectionAttributesAll(),
			},
			{
				// Test importing the resource
				Config:            testAccSpConnectionComplete(spConnectionId),
				ResourceName:      "pingfederate_idp_sp_connection." + spConnectionId,
				ImportStateId:     spConnectionId,
				ImportState:       true,
				ImportStateVerify: true,
				// These attributes have "_all" versions where values will be imported instead
				ImportStateVerifyIgnore: []string{
					"outbound_provision.channels.0.attribute_mapping",
					"outbound_provision.target_settings",
				},
			},
			{
				// Back to Outbound Provision connection, minimal
				Config: testAccSpConnectionOutboundProvision(spConnectionId),
				Check:  testAccCheckExpectedSpConnectionAttributesOutboundProvision(),
			},
		},
	})
}

func baseHcl(resourceName string) string {
	return fmt.Sprintf(`
	connection_id = "%s"
	entity_id     = "myEntity"
	name          = "mySpConn"
	credentials = {
		certs = []
		signing_settings = {
		  signing_key_pair_ref = {
			id = "419x9yg43rlawqwq9v6az997k"
		  }
		  include_raw_key_in_signature = false
		  include_cert_in_signature    = false
		  algorithm                    = "SHA256withRSA"
		}
	  }
	active                 = false
	contact_info           = {
	  company = "Example Corp"
	  first_name = "Bugs"
	  phone = "5555555"
	  email = "bugsbunny@example.com"
	}
	base_url               = "https://api.pingone.com/v5"
	logging_mode           = "STANDARD"
	virtual_entity_ids     = [
	  "example1",
	  "example2"
	]
	default_virtual_entity_id = "example2"
	connection_target_type = "STANDARD"
	application_name = "MyApp"
	application_icon_url = "https://example.com/icon.png"
	`, resourceName,
	)
}

func outboundProvisionHcl() string {
	return `
  outbound_provision = {
    type = "PingOne"
    target_settings = [
      {
        name  = "PINGONE_ENVIRONMENT"
        value = "example"
      }
    ]
    channels = [
      {
        name        = "Channel1"
        max_threads = 1
        timeout     = 120
        active      = false
        channel_source = {
          base_dn = "dc=example,dc=com"
          data_source = {
            id = "pingdirectory"
          }
          guid_attribute_name = "entry_uuid"
          change_detection_settings = {
            user_object_class         = "inetOrgPerson"
            changed_users_algorithm   = "TIMESTAMP_NO_NEGATION"
            group_object_class        = "groupOfUniqueNames"
            time_stamp_attribute_name = "modifyTimestamp"
          }
          account_management_settings = {
            account_status_algorithm      = "ACCOUNT_STATUS_ALGORITHM_FLAG"
            account_status_attribute_name = "nsaccountlock"
            flag_comparison_value         = "true"
            flag_comparison_status        = true
            default_status                = true
          }
          group_membership_detection = {
            group_member_attribute_name = "uniqueMember"
          }
          guid_binary = false
          user_source_location = {
            filter = "cn=John"

          }
        }
        attribute_mapping = [
          {
            field_name = "username"
            saas_field_info = {
              attribute_names = [
                "uid"
              ]
            }
          },
          {
            field_name = "email"
            saas_field_info = {
              attribute_names = [
                "mail"
              ]
            }
          },
          {
            field_name = "populationID"
            saas_field_info = {
              default_value = "example"
            }
          }
        ]
      }
    ]
  }
  `
}

func wsTrustHcl() string {
	return `
	ws_trust = {
		partner_service_ids = [
		  "myid"
		]
		oauth_assertion_profiles = true
		default_token_type       = "SAML20"
		generate_key             = false
		encrypt_saml2_assertion  = false
		minutes_before           = 5
		minutes_after            = 30
		attribute_contract = {
		  core_attributes = [
			{
			  name = "TOKEN_SUBJECT"
			}
		  ]
		  extended_attributes = []
		}
		token_processor_mappings = [
		  {
			attribute_sources = []
			attribute_contract_fulfillment = {
			  "TOKEN_SUBJECT" : {
				source = {
				  type = "TOKEN"
				}
				value = "username"
			  }
			}
			issuance_criteria = {
			  conditional_criteria = []
			}
			idp_token_processor_ref = {
			  id = "UsernameTokenProcessor"
			}
			restricted_virtual_entity_ids = []
		  }
		]
	  }
	`
}

func spBrowserSSOHcl(authenticationPolicyContractName string) string {
	return `
sp_browser_sso = {
    protocol                      = "SAML20"
    require_signed_authn_requests = false
    sp_saml_identity_mapping      = "STANDARD"
    sign_assertions               = false
    authentication_policy_contract_assertion_mappings = [
      {
        abort_sso_transaction_as_fail_safe = false
        authentication_policy_contract_ref = {
          id = "QGxlec5CX693lBQL"
        }
        restricted_virtual_entity_ids = []
        attribute_contract_fulfillment = {
          "SAML_SUBJECT" = {
            source = {
              type = "AUTHENTICATION_POLICY_CONTRACT"
            }
            value = "subject"
          }
        }
        restrict_virtual_entity_ids = false
        attribute_sources           = []
        issuance_criteria = {
          conditional_criteria = []
        }
      }
    ]
    encryption_policy = {
      encrypt_slo_subject_name_id   = false
      encrypt_assertion             = false
      encrypted_attributes          = []
      slo_subject_name_id_encrypted = false
    }
    enabled_profiles = [
      "IDP_INITIATED_SSO"
    ]
    sign_response_as_required = true
    sso_service_endpoints = [
      {
        is_default = true
        binding    = "POST"
        index      = 0
        url        = "https://httpbin.org/anything"
      }
    ]
    adapter_mappings = []
    assertion_lifetime = {
      minutes_after  = 5
      minutes_before = 5
    }
    attribute_contract = {
      core_attributes = [
        {
          name_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
          name        = "SAML_SUBJECT"
        }
      ]
      extended_attributes = []
    }
  }
`
}

func wsFedSpBrowserSSOHcl(authenticationPolicyContractName string) string {
	return `
sp_browser_sso = {
    protocol                      = "WSFED"
    always_sign_artifact_response = false
    sso_service_endpoints = [
      {
        url = "/sp/prpwrong.wsf"
      }
    ]
    sp_ws_fed_identity_mapping = "EMAIL_ADDRESS"
    assertion_lifetime = {
      minutes_before = 5
      minutes_after = 5
    }
    attribute_contract = {
      core_attributes = [
        {
          name = "SAML_SUBJECT"
        }
      ]
      extended_attributes = []
    }
    adapter_mappings = [
      {
        attribute_sources = []
        attribute_contract_fulfillment = {
          "SAML_SUBJECT" = {
            source = {
              type = "ADAPTER"
            }
            value = "subject"
          }
        }
        issuance_criteria = {
          conditional_criteria = []
        }
        restrict_virtual_entity_ids = false
        restricted_virtual_entity_ids = []
        idp_adapter_ref = {
          id = "OTIdPJava"
        }
        abort_sso_transaction_as_fail_safe = false
      }
    ]
    authentication_policy_contract_assertion_mappings = []
    ws_fed_token_type = "SAML11"
    ws_trust_version = "WSTRUST12"
  }
`
}

func testAccSpConnectionOutboundProvision(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "%[1]s" {
	%s
  %s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		outboundProvisionHcl(),
	)
}

func testAccSpConnectionBrowserSso(resourceName string, useWsFed bool) string {
	var browserHcl string
	if useWsFed {
		browserHcl = wsFedSpBrowserSSOHcl(resourceName)
	} else {
		browserHcl = spBrowserSSOHcl(resourceName)
	}

	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "%[1]s" {
  %s
  %s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		browserHcl,
	)
}

func testAccSpConnectionWsTrust(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "%[1]s" {
  %s
  %s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		wsTrustHcl(),
	)
}

func testAccSpConnectionComplete(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "%[1]s" {
		%s
		%s
		%s
		%s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		outboundProvisionHcl(),
		spBrowserSSOHcl(resourceName),
		wsTrustHcl(),
	)
}

func testCommonExpectedSpConnectionAttributes() (*configurationapi.SpConnection, error) {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	spConn, _, err := testClient.IdpSpConnectionsAPI.GetSpConnection(ctx, spConnectionId).Execute()

	if err != nil {
		return nil, err
	}

	// Entity id
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "entity_id", "myEntity", spConn.EntityId)
	if err != nil {
		return spConn, err
	}

	// Name
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "name", "mySpConn", spConn.Name)
	if err != nil {
		return spConn, err
	}

	// Contact info
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "contact_info.company_name", "Example Corp", *spConn.ContactInfo.Company)
	if err != nil {
		return spConn, err
	}

	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "contact_info.email", "bugsbunny@example.com", *spConn.ContactInfo.Email)
	if err != nil {
		return spConn, err
	}

	// Virtual entity ids
	err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(spConnectionId), "virtual_entity_ids",
		[]string{"example1", "example2"}, spConn.VirtualEntityIds)
	if err != nil {
		return spConn, err
	}

	return spConn, nil
}

func testExpectedSpConnectionOutboundProvisionAttributes(response *configurationapi.SpConnection) error {
	// Type
	err := acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "outbound_provision.type", "PingOne", response.OutboundProvision.Type)
	if err != nil {
		return err
	}

	// Target settings
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "outbound_provision.target_settings[0].name",
		"PINGONE_ENVIRONMENT", response.OutboundProvision.TargetSettings[0].Name)
	if err != nil {
		return err
	}

	// Channels
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "outbound_provision.channels[0].name",
		"Channel1", response.OutboundProvision.Channels[0].Name)
	if err != nil {
		return err
	}

	// base dn of channel source
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "outbound_provision.channels[0].channel_source.base_dn",
		"dc=example,dc=com", response.OutboundProvision.Channels[0].ChannelSource.BaseDn)
	if err != nil {
		return err
	}

	return nil
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSpConnectionAttributesOutboundProvision() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		response, err := testCommonExpectedSpConnectionAttributes()
		if err != nil {
			return err
		}

		return testExpectedSpConnectionOutboundProvisionAttributes(response)
	}
}

func testExpectedSpConnectionBrowserSSOAttributes(response *configurationapi.SpConnection, useWsFed bool) error {
	// protocol
	var expected string
	if useWsFed {
		expected = "WSFED"
	} else {
		expected = "SAML20"
	}
	err := acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "sp_browser_sso.protocol",
		expected, response.SpBrowserSso.Protocol)
	if err != nil {
		return err
	}

	if useWsFed {
		// ws trust version
		err := acctest.TestAttributesMatchStringPointer(resourceType, pointers.String(spConnectionId), "sp_browser_sso.ws_trust_version",
			"WSTRUST12", response.SpBrowserSso.WsTrustVersion)
		if err != nil {
			return err
		}
	} else {
		// enabled profiles
		err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(spConnectionId), "sp_browser_sso.enabled_profiles",
			[]string{"IDP_INITIATED_SSO"}, response.SpBrowserSso.EnabledProfiles)
		if err != nil {
			return err
		}
	}

	// attribute contract
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "sp_browser_sso.attribute_contract.core_attributes[0].name",
		"SAML_SUBJECT", response.SpBrowserSso.AttributeContract.CoreAttributes[0].Name)
	if err != nil {
		return err
	}

	return nil
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSpConnectionAttributesBrowserSSO(useWsFed bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		response, err := testCommonExpectedSpConnectionAttributes()
		if err != nil {
			return err
		}

		return testExpectedSpConnectionBrowserSSOAttributes(response, useWsFed)
	}
}

func testExpectedSpConnectionWsTrustAttributes(response *configurationapi.SpConnection) error {
	// minutes_before
	err := acctest.TestAttributesMatchInt(resourceType, pointers.String(spConnectionId), "ws_trust.minutes_before",
		5, *response.WsTrust.MinutesBefore)
	if err != nil {
		return err
	}

	// core attributes
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "ws_trust.attribute_contract.core_attributes[0].name",
		"TOKEN_SUBJECT", response.WsTrust.AttributeContract.CoreAttributes[0].Name)
	if err != nil {
		return err
	}

	// attribute contract fulfillment
	err = acctest.TestAttributesMatchString(resourceType, pointers.String(spConnectionId), "ws_trust.token_processor_mappings.attribute_contract_fulfillment[\"TOKEN_SUBJECT\"].value",
		"username", response.WsTrust.TokenProcessorMappings[0].AttributeContractFulfillment["TOKEN_SUBJECT"].Value)
	if err != nil {
		return err
	}

	return nil
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSpConnectionAttributesWsTrust() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		response, err := testCommonExpectedSpConnectionAttributes()
		if err != nil {
			return err
		}

		return testExpectedSpConnectionWsTrustAttributes(response)
	}
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSpConnectionAttributesAll() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		response, err := testCommonExpectedSpConnectionAttributes()
		if err != nil {
			return err
		}

		err = testExpectedSpConnectionOutboundProvisionAttributes(response)
		if err != nil {
			return err
		}

		err = testExpectedSpConnectionBrowserSSOAttributes(response, false)
		if err != nil {
			return err
		}

		return testExpectedSpConnectionWsTrustAttributes(response)
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSpConnectionDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(ctx, spConnectionId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("IdP SP Connection", spConnectionId)
	}
	return nil
}
