// Copyright © 2025 Ping Identity Corporation

package authenticationpoliciesfragment_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

var pingOneConnection, pingOneEnvironment, pingOnePopulation string

func TestAccAuthenticationPoliciesFragment(t *testing.T) {
	resourceName := "myAuthenticationPoliciesFragment"

	pingOneConnection = os.Getenv("PF_TF_P1_CONNECTION_ID")
	pingOneEnvironment = os.Getenv("PF_TF_P1_CONNECTION_ENV_ID")
	pingOnePopulation = os.Getenv("PF_TF_P1_POPULATION_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if pingOneConnection == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ID must be set for the TestAccAuthenticationPoliciesFragment acceptance test")
			}
			if pingOneEnvironment == "" {
				t.Fatal("PF_TF_P1_CONNECTION_ENV_ID must be set for the TestAccAuthenticationPoliciesFragment acceptance test")
			}
			if pingOnePopulation == "" {
				t.Fatal("PF_TF_P1_POPULATION_ID must be set for the TestAccAuthenticationPoliciesFragment acceptance test")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test simple policy fragment
			{
				Config: testAccAuthenticationPoliciesFragmentSimple(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(resourceName, false),
			},
			// Test a more complex fragment
			{
				Config: testAccAuthenticationPoliciesFragmentComplex(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(resourceName, true),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationPoliciesFragmentComplex(resourceName),
				ResourceName:      "pingfederate_authentication_policies_fragment." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccAuthenticationPoliciesFragmentSimple(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(resourceName, false),
			},
		},
	})
}

func testAccAuthenticationPoliciesFragmentSimple(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_fragment" "%[1]s" {
  fragment_id = "%[1]s"
  name        = "%[1]s"
  root_node = {
    action = {
      authn_source_policy_action = {
        authentication_source = {
          type = "IDP_ADAPTER"
          source_ref = {
            id = "OTIdPJava"
          }
        }
      }
    },
    children = [
      {
        action = {
          done_policy_action = {
            context = "Fail"
          }
        }
      },
      {
        action = {
          done_policy_action = {
            context = "Success"
          }
        }
      }
    ]
  }
  inputs = {
    id = pingfederate_authentication_policy_contract.mycontract.contract_id
  }
  outputs = {
    id = pingfederate_authentication_policy_contract.mycontract.contract_id
  }
}

data "pingfederate_authentication_policies_fragment" "%[1]s" {
  fragment_id = pingfederate_authentication_policies_fragment.%[1]s.fragment_id
}

%[2]s
`, resourceName,
		dependencyHcl(),
	)
}

func testAccAuthenticationPoliciesFragmentComplex(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_fragment" "%[1]s" {
  fragment_id = "%[1]s"
  name        = "%[1]s"
  description = "Registration with PingOne Verify (GovID + Selfie)"
  root_node = {
    action = {
      authn_source_policy_action = {
        authentication_source = {
          type = "IDP_ADAPTER"
          source_ref = {
            id = pingfederate_idp_adapter.myadapter.adapter_id
          }
        }
        input_user_id_mapping = {
          source = {
            type = "INPUTS"
            id   = "Inputs"
          }
          value = "username"
        }
        user_id_authenticated = true
      }
    }
    children = [
      {
        action = {
          done_policy_action = {
            context = "Fail"
          }
        }
      },
      {
        action = {
          apc_mapping_policy_action = {
            context = "Success"
            authentication_policy_contract_ref = {
              id = pingfederate_authentication_policy_contract.mycontract.contract_id
            }
            attribute_mapping = {
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
                        value = "/users/external"
                      },
                    ]
                    id = "APIStubs"
                  }
                },
                {
                  jdbc_attribute_source = {
                    attribute_contract_fulfillment = null
                    column_names                   = ["GRANTEE"]
                    data_store_ref = {
                      id = "ProvisionerDS"
                    }
                    description = "JDBC"
                    filter      = "subject"
                    id          = "jdbcguy"
                    schema      = "INFORMATION_SCHEMA"
                    table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
                  }
                },
              ],
              attribute_contract_fulfillment = {
                "firstName" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.myadapter.adapter_id
                  }
                  value = "firstName"
                }
                "lastName" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.myadapter.adapter_id
                  }
                  value = "lastName"
                }
                "subject" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.myadapter.adapter_id
                  }
                  value = "subject"
                }
                "fullName" : {
                  source = {
                    type = "EXPRESSION"
                  },
                  value = "'test1|test2|test3'.split(\"\\\\|\")[1]"
                }
                "photo" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.myadapter.adapter_id
                  }
                  value = "photo"
                }
                "username" : {
                  source = {
                    type = "INPUTS",
                    id   = "inputs"
                  }
                  value = "username"
                }
              }
            }
          }
        }
      }
    ]
  }
  inputs = {
    id = pingfederate_authentication_policy_contract.mycontract.contract_id
  }
  outputs = {
    id = pingfederate_authentication_policy_contract.mycontract.contract_id
  }

}

data "pingfederate_authentication_policies_fragment" "%[1]s" {
  fragment_id = pingfederate_authentication_policies_fragment.%[1]s.fragment_id
}

%[2]s
`, resourceName,
		dependencyHcl(),
	)
}

func verifyFields() string {
	if acctest.VersionAtLeast(version.PingFederate1130) {
		return `
    {
      name  = "Email Chained Attribute",
      value = "mail"
    },
    {
      name  = "Phone Chained Attribute",
      value = "mobile"
    },
    {
      name  = "Reference Image Chained Attribute",
      value = "photo"
    },
    {
      name  = "Verification URL Delivery Method",
      value = "User Selection"
    },
    {
      name  = "Verify Policy",
      value = ""
    },
    `
	} else {
		return `
    {
      name  = "Verify App Name",
      value = "myappname"
    },
     `
	}
}

func additionalCoreAttributes() string {
	if acctest.VersionAtLeast(version.PingFederate1130) {
		return `
    {
      name      = "gender",
      masked    = false,
      pseudonym = false
    },
    {
      name      = "weight",
      masked    = false,
      pseudonym = false
    },
    {
      name      = "nationality",
      masked    = false,
      pseudonym = false
    },
    {
      name      = "issuingCountry",
      masked    = false,
      pseudonym = false
    },
    `
	} else {
		return ""
	}
}

func additionalAttributeContractFulfillments() string {
	if acctest.VersionAtLeast(version.PingFederate1130) {
		return `
    "gender" : {
      source = {
        type = "ADAPTER"
      },
      value = "gender"
    },
    "weight" : {
      source = {
        type = "ADAPTER"
      },
      value = "weight"
    },
    "nationality" : {
      source = {
        type = "ADAPTER"
      },
      value = "nationality"
    },
    "issuingCountry" : {
      source = {
        type = "ADAPTER"
      },
      value = "issuingCountry"
    },
    `
	} else {
		return ""
	}
}

func dependencyHcl() string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "mycontract" {
  contract_id = "fragmentVerifyReg"
  name        = "Fragment - Verify - Registration"
  extended_attributes = [
    {
      name = "firstName"
    },
    {
      name = "lastName"
    },
    {
      name = "fullName"
    },
    {
      name = "photo"
    },
    {
      name = "username"
    }
  ]
}

resource "pingfederate_idp_adapter" "myadapter" {
  adapter_id = "PingOneVerify"
  name       = "PingOneVerify (GovID)"
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.pingone.verify.PingOneVerifyAdapter"
  }
  configuration = {
    tables = [
      {
        name = "PingOne Verify Response Mappings (optional)"
        rows = [
          {
            fields = [
              {
                name  = "Local Attribute"
                value = "photo"
              },
              {
                name  = "PingOne Verify API Attribute Mapping"
                value = "/_embedded/verifiedData/0/data/IMAGE"
              }
            ],
            default_row = false
          }
        ]
      }
    ],
    fields = [
      {
        name  = "PingOne Environment",
        value = "%s|%s"
      },
      {
        name  = "PingOne Population",
        value = "%s"
      },
      %s
      {
        name  = "Test Username",
        value = ""
      },
      {
        name  = "HTML Template Prefix",
        value = "pingone-verify"
      },
      {
        name  = "Messages Files",
        value = "pingone-verify-messages"
      },
      {
        name  = "Error Message Key Prefix",
        value = "pingone.verify.error."
      },
      {
        name  = "Provision User",
        value = "true"
      },
      {
        name  = "Allow Verification Retries",
        value = "true"
      },
      {
        name  = "User Not Found Failure Mode",
        value = "Block user"
      },
      {
        name  = "Service Unavailable Failure Mode",
        value = "Bypass authentication"
      },
      {
        name  = "Show Success Screens",
        value = "true"
      },
      {
        name  = "Show Failed Screens",
        value = "true"
      },
      {
        name  = "Show Timeout Screens",
        value = "true"
      },
      {
        name  = "State Timeout",
        value = "1200"
      },
      {
        name  = "API Request Timeout",
        value = "5000"
      },
      {
        name  = "Proxy Settings",
        value = "System Defaults"
      },
      {
        name  = "Custom Proxy Host",
        value = ""
      },
      {
        name  = "Custom Proxy Port",
        value = ""
      }
    ]
  }
  attribute_contract = {
    core_attributes = [
      {
        name      = "verifiedDocuments",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "lastName",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "country",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "addressZip",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "transactionStatus",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "subject",
        masked    = false,
        pseudonym = true
      },
      {
        name      = "addressState",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "idNumber",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "birthDate",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "firstName",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "addressStreet",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "issueDate",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "addressCity",
        masked    = false,
        pseudonym = false
      },
      {
        name      = "expirationDate",
        masked    = false,
        pseudonym = false
      },
      %s
    ],
    extended_attributes = [
      {
        name      = "photo",
        masked    = false,
        pseudonym = false
      }
    ]
    mask_ognl_values = false
  }
  attribute_mapping = {
    attribute_sources = [
      {
        ldap_attribute_source = {
          attribute_contract_fulfillment = null
          base_dn                        = "ou=Applications,ou=Ping,ou=Groups,dc=dm,dc=example,dc=com"
          binary_attribute_settings      = null
          id                             = "ldapguy"
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
              value = "/users/external"
            },
          ]
          id = "APIStubs"
        }
      },
    ],
    attribute_contract_fulfillment = {
      "country" : {
        source = {
          type = "ADAPTER"
        },
        value = "country"
      },
      "lastName" : {
        source = {
          type = "ADAPTER"
        },
        value = "lastName"
      },
      "verifiedDocuments" : {
        source = {
          type = "ADAPTER"
        },
        value = "verifiedDocuments"
      },
      "addressZip" : {
        source = {
          type = "ADAPTER"
        },
        value = "addressZip"
      },
      "transactionStatus" : {
        source = {
          type = "ADAPTER"
        },
        value = "transactionStatus"
      },
      "subject" : {
        source = {
          type = "ADAPTER"
        },
        value = "subject"
      },
      "photo" : {
        source = {
          type = "ADAPTER"
        },
        value = "photo"
      },
      "addressState" : {
        source = {
          type = "ADAPTER"
        },
        value = "addressState"
      },
      "idNumber" : {
        source = {
          type = "ADAPTER"
        },
        value = "idNumber"
      },
      "birthDate" : {
        source = {
          type = "ADAPTER"
        },
        value = "birthDate"
      },
      "firstName" : {
        source = {
          type = "ADAPTER"
        },
        value = "firstName"
      },
      "addressStreet" : {
        source = {
          type = "ADAPTER"
        },
        value = "addressStreet"
      },
      "issueDate" : {
        source = {
          type = "ADAPTER"
        },
        value = "issueDate"
      },
      "addressCity" : {
        source = {
          type = "ADAPTER"
        },
        value = "addressCity"
      },
      "expirationDate" : {
        source = {
          type = "ADAPTER"
        },
        value = "expirationDate"
      },
      %s
    }
  }
}`, pingOneConnection, pingOneEnvironment, pingOnePopulation, verifyFields(), additionalCoreAttributes(), additionalAttributeContractFulfillments())
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(id string, isComplexFragment bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPoliciesFragment"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationPoliciesAPI.GetFragment(ctx, id).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "name",
			id, response.Name)
		if err != nil {
			return err
		}

		// Check for the expected root action
		if response.RootNode.Action.AuthnSourcePolicyAction == nil {
			return errors.New("Expected root node with AUTHN_SOURCE policy action")
		}

		// Check for the expected children
		if len(response.RootNode.Children) != 2 {
			return errors.New("Expected root node with two child nodes")
		}

		if response.RootNode.Children[0].Action.DonePolicyAction == nil {
			return errors.New("Expected root node with first child being a DONE policy action")
		}

		if !isComplexFragment && response.RootNode.Children[1].Action.DonePolicyAction == nil {
			return errors.New("Expected root node with second child being a DONE policy action")
		}

		if isComplexFragment && response.RootNode.Children[1].Action.ApcMappingPolicyAction == nil {
			return errors.New("Expected root node with second child being an APC_MAPPING policy action")
		}

		return nil
	}
}
