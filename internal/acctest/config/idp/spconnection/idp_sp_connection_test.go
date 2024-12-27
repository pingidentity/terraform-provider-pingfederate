package idpspconnection_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const (
	spConnectionId = "spConnId"
	resourceType   = "IdP SP Connection"
)

var outboundEnvironmentSetting = ""

func TestAccIdpSpConnection(t *testing.T) {
	connId := os.Getenv("PF_TF_P1_CONNECTION_ID")
	envId := os.Getenv("PF_TF_P1_CONNECTION_ENV_ID")
	outboundEnvironmentSetting = connId + "|" + envId
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if connId == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ID must be set for the IdP SP Connection acceptance test")
			}
			if envId == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ENV_ID must be set for the IdP SP Connection acceptance test")
			}
		},
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
				// Browser SSO SAML connection reproducing inconsistent result bug
				// https://github.com/pingidentity/terraform-provider-pingfederate/issues/319
				Config: testAccSpConnectionBrowserSsoInconsistentResult(spConnectionId),
			},
			{
				// WS Trust connection, minimal
				Config: testAccSpConnectionWsTrust(spConnectionId),
				Check:  testAccCheckExpectedSpConnectionAttributesWsTrust(),
			},
			{
				// Complete connection with all three types
				Config: testAccSpConnectionComplete(spConnectionId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSpConnectionAttributesAll(),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_idp_sp_connection.%s", spConnectionId), "virtual_entity_ids.0", "example1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_idp_sp_connection.%s", spConnectionId), "virtual_entity_ids.1", "example2"),
				),
			},
			{
				// Test importing the resource
				Config:            testAccSpConnectionComplete(spConnectionId),
				ResourceName:      "pingfederate_idp_sp_connection." + spConnectionId,
				ImportStateId:     spConnectionId,
				ImportState:       true,
				ImportStateVerify: true,
				// These attributes have many extra values not being set in the test used in this HCL, so those extra values
				// will change these attributes on import. file_data also gets formatted by PF so it won't exactly match.
				ImportStateVerifyIgnore: []string{
					"outbound_provision.channels.0.attribute_mapping",
					"outbound_provision.target_settings",
					"credentials.certs.0.x509_file.file_data",
				},
				// Ensure that the both versions of the attributes have values set
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_idp_sp_connection.%s", spConnectionId), "outbound_provision.channels.0.attribute_mapping.#", "38"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_idp_sp_connection.%s", spConnectionId), "outbound_provision.target_settings", "9"),
				),
			},
			{
				// Back to Outbound Provision connection, minimal
				Config: testAccSpConnectionOutboundProvision(spConnectionId),
				Check:  testAccCheckExpectedSpConnectionAttributesOutboundProvision(),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.IdpSpConnectionsAPI.DeleteSpConnection(ctx, spConnectionId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Outbound provision connection, minimal
				Config: testAccSpConnectionOutboundProvision(spConnectionId),
				Check:  testAccCheckExpectedSpConnectionAttributesOutboundProvision(),
			},
		},
	})
}

func baseCredentials() string {
	return `
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
`
}

func fullCredentials() string {
	return `
	credentials = {
    certs = [{
      x509_file = {
        id        = "4qrossmq1vxa4p836kyqzp48h"
        file_data = "MIIDOjCCAiICCQCjbB7XBVkxCzANBgkqhkiG9w0BAQsFADBfMRIwEAYDVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwHhcNMjMwNzE0MDI1NDUzWhcNMjQwNzEzMDI1NDUzWjBfMRIwEAYDVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC5yFrh9VR2wk9IjzMz+Ei80K453g1j1/Gv3EQ/SC9h7HZBI6aV9FaEYhGnaquRT5q87p8lzCphKNXVyeL6T/pDJOW70zXItkl8Ryoc0tIaknRQmj8+YA0Hr9GDdmYev2yrxSoVS7s5Bl8poasn3DljgnWT07vsQz+hw3NY4SPp7IFGP2PpGUBBIIvrOaDWpPGsXeznBxSFtis6Qo+JiEoaVql9b9/XyKZj65wOsVyZhFWeM1nCQITSP9OqOc9FSoDFYQ1AVogm4A2AzUrkMnT1SrN2dCuTmNbeVw7gOMqMrVf0CiTv9hI0cATbO5we1sPAlJxscSkJjsaI+sQfjiAnAgMBAAEwDQYJKoZIhvcNAQELBQADggEBACgwoH1qklPF1nI9+WbIJ4K12Dl9+U3ZMZa2lP4hAk1rMBHk9SHboOU1CHDQKT1Z6uxi0NI4JZHmP1qP8KPNEWTI8Q76ue4Q3aiA53EQguzGb3SEtyp36JGBq05Jor9erEebFftVl83NFvio72Fn0N2xvu8zCnlylf2hpz9x1i01Xnz5UNtZ2ppsf2zzT+4U6w3frH+pkp0RDPuoe9mnBF001AguP31hSBZyZzWcwQltuNELnSRCcgJl4kC2h3mAgaVtYalrFxLRa3tA2XF2BHRHmKgocedVhTq+81xrqj+WQuDmUe06DnrS3Ohmyj3jhsCCluznAolmrBhT/SaDuGg="
      }
      active_verification_cert    = true
      encryption_cert             = true
      primary_verification_cert   = true
      secondary_verification_cert = false
    }]

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
  }
	`
}

func baseHcl(resourceName string) string {
	return fmt.Sprintf(`
	connection_id = "%s"
	entity_id     = "myEntity"
	name          = "mySpConn"
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
	return fmt.Sprintf(`
  outbound_provision = {
    type = "PingOne"
    target_settings = [
      {
        name  = "PINGONE_ENVIRONMENT"
        value = "%s"
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
  `, outboundEnvironmentSetting)
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
		encrypt_saml2_assertion  = true
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
  %s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		baseCredentials(),
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
  %s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		baseCredentials(),
		browserHcl,
	)
}

func testAccSpConnectionBrowserSsoInconsistentResult(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "%[1]s" {
  %s
  %s
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
      "IDP_INITIATED_SSO",
      "SP_INITIATED_SSO",
    ]
    incoming_bindings = [
      "REDIRECT",
      "POST"
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
    }
  }
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		baseCredentials(),
	)
}

func testAccSpConnectionWsTrust(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_sp_connection" "%[1]s" {
  %s
  %s
  %s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		fullCredentials(),
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
		%s
}
data "pingfederate_idp_sp_connection" "%[1]s" {
  connection_id = pingfederate_idp_sp_connection.%[1]s.connection_id
}`, resourceName,
		baseHcl(resourceName),
		fullCredentials(),
		outboundProvisionHcl(),
		spBrowserSSOHcl(resourceName),
		wsTrustHcl(),
	)
}

func testCommonExpectedSpConnectionAttributes(s *terraform.State) (*configurationapi.SpConnection, error) {
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
		response, err := testCommonExpectedSpConnectionAttributes(s)
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
		response, err := testCommonExpectedSpConnectionAttributes(s)
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
		response, err := testCommonExpectedSpConnectionAttributes(s)
		if err != nil {
			return err
		}

		return testExpectedSpConnectionWsTrustAttributes(response)
	}
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSpConnectionAttributesAll() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		response, err := testCommonExpectedSpConnectionAttributes(s)
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
