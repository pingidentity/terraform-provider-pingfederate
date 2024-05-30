package authenticationpolicies_test

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

func TestAccAuthenticationPolicies(t *testing.T) {
	resourceName := "myAuthenticationPolicies"

	simpleDescription := "This is an OTIdPJava Authentication Policy"
	complexDescription := "This Sample Policy uses the Extended Properties Selector on the Application to allow easy switching between: Single_Factor (First_Factor Fragment) Internal (Employee HTML Form)"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test simple policy
			{
				Config: testAccAuthenticationPoliciesSimple(resourceName, simpleDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAuthenticationPoliciesAttributes(resourceName, simpleDescription, false),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies.%s", resourceName), "authn_selection_trees.0.enabled", "true"),
				),
			}, // Test a more complex policy
			{
				Config: testAccAuthenticationPoliciesComplex(resourceName, complexDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAuthenticationPoliciesAttributes(resourceName, complexDescription, true),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_policies.%s", resourceName), "fail_if_no_selection", "true"),
				),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationPoliciesComplex(resourceName, complexDescription),
				ResourceName:      "pingfederate_authentication_policies." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccAuthenticationPoliciesSimple(resourceName, simpleDescription),
				Check:  testAccCheckExpectedAuthenticationPoliciesAttributes(resourceName, simpleDescription, false),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.AuthenticationPoliciesAPI.DeletePolicy(ctx, resourceName).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Back to minimal model
				Config: testAccAuthenticationPoliciesSimple(resourceName, simpleDescription),
				Check:  testAccCheckExpectedAuthenticationPoliciesAttributes(resourceName, simpleDescription, false),
			},
		},
	})
}

func testAccAuthenticationPoliciesSimple(resourceName, description string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies" "%[1]s" {
  authn_selection_trees = [
    {
      name        = "%[1]s"
      description = "%[2]s"
      id          = "%[1]s"
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
  ]
}
`, resourceName, description,
	)
}

func testAccAuthenticationPoliciesComplex(resourceName, description string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies" "%[1]s" {
  fail_if_no_selection = true
  authn_selection_trees = [
    {
      root_node = {
        action = {
          authn_selector_policy_action = {
            authentication_selector_ref = {
              id = "authnExp"
            }
          }
        }
        children = [
          {
            action = {
              fragment_policy_action = {
                context = "Internal"
                fragment = {
                  id = "InternalAuthN"
                }
                fragment_mapping = {
                  attribute_sources = []
                  attribute_contract_fulfillment = {
                    "subject" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    }
                  }
                  issuance_criteria = {
                    conditional_criteria = []
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
              fragment_policy_action = {
                context = "Single_Factor"
                fragment = {
                  id = "FirstFactor"
                }
                fragment_mapping = {
                  attribute_sources = []
                  attribute_contract_fulfillment = {
                    "subject" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    }
                  }
                  issuance_criteria = {
                    conditional_criteria = []
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
          }
        ]
      }
      name        = "%[1]s"
      description = "%[2]s"
      id          = "%[1]s"
      enabled     = true
    },
    {
      root_node = {
        action = {
          fragment_policy_action = {
            fragment = {
              id = "FirstFactor"
            }
            fragment_mapping = {
              attribute_sources = []
              attribute_contract_fulfillment = {
                "subject" = {
                  source = {
                    type = "NO_MAPPING"
                  }
                }
              }
              issuance_criteria = {
                conditional_criteria = []
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
      }
      name        = "Fallback Policy"
      description = "Used to perform authentication if other Policies are not invoked"
      enabled     = true
    },
    {
      root_node = {
        action = {
          fragment_policy_action = {
            fragment = {
              id = "FirstFactor"
            },
            fragment_mapping = {
              attribute_sources = []
              attribute_contract_fulfillment = {
                "subject" = {
                  source = {
                    type = "NO_MAPPING"
                  }
                }
              }
              issuance_criteria = {
                conditional_criteria = []
              }
            }
          }
        },
        children = [
          {
            action = {
              authn_selector_policy_action = {
                context = "Fail",
                authentication_selector_ref = {
                  id = "authnExp",
                }
              }
            }
            children = [
              {
                action = {
                  restart_policy_action = {
                    context = "Internal"
                  }
                }
              },
              {
                action = {
                  fragment_policy_action = {
                    context = "Single_Factor"
                    fragment = {
                      id = "FirstFactor"
                    },
                    fragment_mapping = {
                      attribute_sources = []
                      attribute_contract_fulfillment = {
                        "subject" = {
                          source = {
                            type = "NO_MAPPING"
                          }
                        }
                      }
                      issuance_criteria = {
                        conditional_criteria = []
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
                      apc_mapping_policy_action = {
                        context = "Success"
                        authentication_policy_contract_ref = {
                          id = "QGxlec5CX693lBQL"
                        },
                        attribute_mapping = {
                          attribute_sources = []
                          attribute_contract_fulfillment = {
                            "subject" = {
                              source = {
                                type = "NO_MAPPING"
                              }
                            }
                          },
                          issuance_criteria = {
                            conditional_criteria = []
                          }
                        }
                      }
                    }
                  }
                ]
              }
            ]
          },
          {
            action = {
              local_identity_mapping_policy_action = {
                context = "Success",
                local_identity_ref = {
                  id = "adminIdentityProfile",
                },
                inbound_mapping = {
                  attribute_sources = [],
                  attribute_contract_fulfillment = {
                    "pf.local.identity.unique.id" = {
                      source = {
                        type = "TEXT"
                      },
                      value = "test"
                    }
                  },
                  issuance_criteria = {
                    conditional_criteria = []
                  }
                },
                outbound_attribute_mapping = {
                  attribute_sources = [],
                  attribute_contract_fulfillment = {
                    "firstName" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    },
                    "lastName" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    },
                    "ImmutableID" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    },
                    "mail" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    },
                    "subject" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    },
                    "SAML_AUTHN_CTX" = {
                      source = {
                        type = "NO_MAPPING"
                      }
                    }
                  },
                  issuance_criteria = {
                    conditional_criteria = []
                  }
                }
              }
            }
          }
        ]
      }
      name                    = "Coverage For Testing",
      enabled                 = true,
      handle_failures_locally = false
    }
  ]
}
  `, resourceName, description,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPoliciesAttributes(id, description string, isComplex bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPolicies"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationPoliciesAPI.GetDefaultAuthenticationPolicy(ctx).Execute()
		if err != nil {
			return err
		}

		authenticationPolicy := response.AuthnSelectionTrees[0]

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "name",
			id, authenticationPolicy.Name)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "description",
			description, authenticationPolicy.Description)
		if err != nil {
			return err
		}

		// Check for the expected root action
		if !isComplex {
			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"IDP_ADAPTER",
				authenticationPolicy.RootNode.Action.AuthnSourcePolicyAction.AuthenticationSource.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"OTIdPJava",
				authenticationPolicy.RootNode.Action.AuthnSourcePolicyAction.AuthenticationSource.SourceRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Fail",
				authenticationPolicy.RootNode.Children[0].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Success",
				authenticationPolicy.RootNode.Children[1].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}
		} else {
			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"authnExp",
				authenticationPolicy.RootNode.Action.AuthnSelectorPolicyAction.AuthenticationSelectorRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Internal",
				authenticationPolicy.RootNode.Children[0].Action.FragmentPolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"InternalAuthN",
				authenticationPolicy.RootNode.Children[0].Action.FragmentPolicyAction.Fragment.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"NO_MAPPING",
				authenticationPolicy.RootNode.Children[0].Action.FragmentPolicyAction.FragmentMapping.AttributeContractFulfillment["subject"].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Fail",
				authenticationPolicy.RootNode.Children[0].Children[0].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Success",
				authenticationPolicy.RootNode.Children[0].Children[1].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Single_Factor",
				authenticationPolicy.RootNode.Children[1].Action.FragmentPolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"FirstFactor",
				authenticationPolicy.RootNode.Children[1].Action.FragmentPolicyAction.Fragment.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"NO_MAPPING",
				authenticationPolicy.RootNode.Children[1].Action.FragmentPolicyAction.FragmentMapping.AttributeContractFulfillment["subject"].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Fail",
				authenticationPolicy.RootNode.Children[1].Children[0].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Success",
				authenticationPolicy.RootNode.Children[1].Children[1].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			//  check failback policy
			failbackPolicy := response.AuthnSelectionTrees[1]
			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"name",
				"Fallback Policy",
				failbackPolicy.Name)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil,
				"enabled",
				true,
				*failbackPolicy.Enabled)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"description",
				"Used to perform authentication if other Policies are not invoked",
				failbackPolicy.Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"FirstFactor",
				failbackPolicy.RootNode.Action.FragmentPolicyAction.Fragment.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"NO_MAPPING",
				failbackPolicy.RootNode.Action.FragmentPolicyAction.FragmentMapping.AttributeContractFulfillment["subject"].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Fail",
				failbackPolicy.RootNode.Children[0].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Success",
				failbackPolicy.RootNode.Children[1].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			//  check coverage for testing policy
			coverageForTestingPolicy := response.AuthnSelectionTrees[2]
			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"name",
				"Coverage For Testing",
				coverageForTestingPolicy.Name)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil,
				"enabled",
				true,
				*coverageForTestingPolicy.Enabled)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil,
				"handle_failures_locally",
				false,
				*coverageForTestingPolicy.HandleFailuresLocally)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"authnExp",
				coverageForTestingPolicy.RootNode.Children[0].Action.AuthnSelectorPolicyAction.AuthenticationSelectorRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"FirstFactor",
				coverageForTestingPolicy.RootNode.Children[0].Children[1].Action.FragmentPolicyAction.Fragment.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"NO_MAPPING",
				coverageForTestingPolicy.RootNode.Children[0].Children[1].Action.FragmentPolicyAction.FragmentMapping.AttributeContractFulfillment["subject"].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Fail",
				coverageForTestingPolicy.RootNode.Children[0].Children[1].Children[0].Action.DonePolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"QGxlec5CX693lBQL",
				coverageForTestingPolicy.RootNode.Children[0].Children[1].Children[1].Action.ApcMappingPolicyAction.AuthenticationPolicyContractRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"NO_MAPPING",
				coverageForTestingPolicy.RootNode.Children[0].Children[1].Children[1].Action.ApcMappingPolicyAction.AttributeMapping.AttributeContractFulfillment["subject"].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Success",
				coverageForTestingPolicy.RootNode.Children[0].Children[1].Children[1].Action.ApcMappingPolicyAction.Context)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"id",
				"adminIdentityProfile",
				coverageForTestingPolicy.RootNode.Children[1].Action.LocalIdentityMappingPolicyAction.LocalIdentityRef.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"TEXT",
				coverageForTestingPolicy.RootNode.Children[1].Action.LocalIdentityMappingPolicyAction.InboundMapping.AttributeContractFulfillment["pf.local.identity.unique.id"].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"value",
				"test",
				coverageForTestingPolicy.RootNode.Children[1].Action.LocalIdentityMappingPolicyAction.InboundMapping.AttributeContractFulfillment["pf.local.identity.unique.id"].Value)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil,
				"type",
				"NO_MAPPING",
				coverageForTestingPolicy.RootNode.Children[1].Action.LocalIdentityMappingPolicyAction.OutboundAttributeMapping.AttributeContractFulfillment["subject"].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringPointer(resourceType, nil,
				"context",
				"Success",
				coverageForTestingPolicy.RootNode.Children[1].Action.LocalIdentityMappingPolicyAction.Context)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
