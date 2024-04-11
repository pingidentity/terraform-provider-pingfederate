package acctest_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccAuthenticationPoliciesPolicy(t *testing.T) {
	resourceName := "myAuthenticationPoliciesPolicy"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test simple policy
			{
				Config: testAccAuthenticationPoliciesPolicySimple(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(resourceName, false),
			},
			// Test a more complex policy
			{
				Config: testAccAuthenticationPoliciesPolicyComplex(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(resourceName, true),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationPoliciesPolicyComplex(resourceName),
				ResourceName:      "pingfederate_authentication_policies_policy." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccAuthenticationPoliciesPolicySimple(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(resourceName, false),
			},
		},
	})
}

func testAccAuthenticationPoliciesPolicySimple(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_policy" "%[1]s" {
  policy_id = "%[1]s"
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
}

data "pingfederate_authentication_policies_policy" "%[1]s" {
  policy_id = pingfederate_authentication_policies_policy.%[1]s.policy_id
}
`, resourceName,
	)
}

func testAccAuthenticationPoliciesPolicyComplex(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "myAuthenticationPolicyContract" {
  contract_id = "myAuthenticationPolicyContract"
  name = "example"
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
  
resource "pingfederate_authentication_policies_policy" "%[1]s" {
  policy_id = "%[1]s"
  name        = "%[1]s"
  description = "Registration with PingOne Verify (GovID + Selfie)"
  root_node = {
    action = {
      authn_source_policy_action = {
        authentication_source = {
          type = "IDP_ADAPTER"
          source_ref = {
            id = "OTIdPJava"
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
              id = pingfederate_authentication_policy_contract.myAuthenticationPolicyContract.contract_id
            }
            attribute_mapping = {
              attribute_sources = []
              attribute_contract_fulfillment = {
                "firstName" : {
                  source = {
                    type = "ADAPTER",
                    id   = "OTIdPJava"
                  }
                  value = "firstName"
                }
                "lastName" : {
                  source = {
                    type = "ADAPTER",
                    id   = "OTIdPJava"
                  }
                  value = "lastName"
                }
                "subject" : {
                  source = {
                    type = "ADAPTER",
                    id   = "OTIdPJava"
                  }
                  value = "subject"
                }
                "fullName" : {
                  source = {
                    type = "EXPRESSION"
                  },
                  value = "fullName"
                }
                "photo" : {
                  source = {
                    type = "ADAPTER",
                    id   = "OTIdPJava"
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
}

data "pingfederate_authentication_policies_policy" "%[1]s" {
  policy_id = pingfederate_authentication_policies_policy.%[1]s.policy_id
}

`, resourceName,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(id string, isComplex bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPoliciesPolicy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationPoliciesAPI.GetPolicy(ctx, id).Execute()

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
		if response.RootNode == nil || response.RootNode.Action.AuthnSourcePolicyAction == nil {
			return errors.New("Expected root node with AUTHN_SOURCE policy action")
		}

		// Check for the expected children
		if len(response.RootNode.Children) != 2 {
			return errors.New("Expected root node with two child nodes")
		}

		if response.RootNode.Children[0].Action.DonePolicyAction == nil {
			return errors.New("Expected root node with first child being a DONE policy action")
		}

		if !isComplex && response.RootNode.Children[1].Action.DonePolicyAction == nil {
			return errors.New("Expected root node with second child being a DONE policy action")
		}

		if isComplex && response.RootNode.Children[1].Action.ApcMappingPolicyAction == nil {
			return errors.New("Expected root node with second child being an APC_MAPPING policy action")
		}

		return nil
	}
}
