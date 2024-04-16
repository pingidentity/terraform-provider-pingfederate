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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(resourceName, true),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "enabled", "true"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "handle_failures_locally", "false"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.action.authn_source_policy_action.authentication_source.source_ref.id", "OTIdPJava"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.0.action.done_policy_action.context", "Fail"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.action.authn_source_policy_action.context", "Success"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.0.action.authn_source_policy_action.context", "Fail"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.0.children.0.action.done_policy_action.context", "Fail"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.0.children.1.action.done_policy_action.context", "Success"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.context", "Success"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.authentication_policy_contract_ref.id", "QGxlec5CX693lBQL"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.attribute_mapping.attribute_sources.0.jdbc_attribute_source.data_store_ref.id", "ProvisionerDS"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.attribute_mapping.attribute_sources.0.jdbc_attribute_source.id", "test"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.attribute_mapping.attribute_sources.0.jdbc_attribute_source.description", "test"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.attribute_mapping.attribute_sources.0.jdbc_attribute_source.schema", "INFORMATION_SCHEMA"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.attribute_mapping.attribute_sources.0.jdbc_attribute_source.table", "ADMINISTRABLE_ROLE_AUTHORIZATIONS"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies_policy.%s", resourceName), "root_node.children.1.children.1.action.apc_mapping_policy_action.attribute_mapping.attribute_sources.0.jdbc_attribute_source.filter", "filter"),
				),
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
  name      = "%[1]s"
  policy_id = "%[1]s"
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
resource "pingfederate_authentication_policies_policy" "%[1]s" {
  name                    = "%[1]s"
  policy_id               = "%[1]s"
  enabled                 = true
  handle_failures_locally = false
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
          authn_source_policy_action = {
            context = "Success"
            authentication_source = {
              type = "IDP_ADAPTER"
              source_ref = {
                id = "OTIdPJava"
              }
            }
          }
        }
        children = [
          {
            action = {
              authn_source_policy_action = {
                context = "Fail"
                authentication_source = {
                  type = "IDP_ADAPTER"
                  source_ref = {
                    id = "OTIdPJava"
                  }
                }
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
                  done_policy_action = {
                    context = "Success"
                  }
                }
              }
            ]
          },
          {
            action = {
              apc_mapping_policy_action = {
                context = "Success"
                authentication_policy_contract_ref = {
                  id = "QGxlec5CX693lBQL"
                }
                attribute_mapping = {
                  attribute_sources = [
                    {
                      jdbc_attribute_source = {
                        data_store_ref = {
                          id = "ProvisionerDS"
                        }
                        id          = "test"
                        description = "test"
                        schema      = "INFORMATION_SCHEMA"
                        table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
                        filter      = "filter"
                        column_names = [
                          "GRANTEE",
                          "IS_GRANTABLE",
                          "ROLE_NAME"
                        ]
                      }
                    }
                  ]
                  attribute_contract_fulfillment = {
                    subject = {
                      source = {
                        type = "ADAPTER"
                        id   = "OTIdPJava"
                      }
                      value = "subject"
                    }
                  }
                  issuance_criteria = {
                    conditional_criteria = [
                      {
                        error_result = "error"
                        source = {
                          type = "MAPPED_ATTRIBUTES"
                        }
                        attribute_name = "subject"
                        condition      = "EQUALS"
                        value          = "value"
                      }
                    ]
                  }
                }
              }
            }
          }
        ]
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

		if isComplex && response.RootNode.Children[1].Children[1].Action.ApcMappingPolicyAction == nil {
			return errors.New("Expected root node with second child being an APC_MAPPING policy action")
		}

		return nil
	}
}
