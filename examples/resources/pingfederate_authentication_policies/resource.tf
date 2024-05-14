resource "pingfederate_authentication_policies" "authenticationPolicies" {
  fail_if_no_selection = false
  authn_selection_trees = [
    {
      root_node = {
        action = {
          authn_selector_policy_action = {
            authentication_selector_ref = {
              id = pingfederate_authentication_selector.authnexp.id
            }
          }
        }
        children = [
          {
            action = {
              fragment_policy_action = {
                context = "Internal"
                fragment = {
                  id = pingfederate_authentication_policies_fragment.internal.id
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
                  id = pingfederate_authentication_policies_fragment.first_factor.id
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
      name        = "Sample Authentication Experiences"
      description = "This Sample Policy uses the Extended Properties Selector on the Application to allow easy switching between:\\r Single_Factor (First_Factor Fragment)\\r Internal (Employee HTML Form)"
      enabled     = true
    },
    {
      root_node = {
        action = {
          fragment_policy_action = {
            fragment = {
              id = pingfederate_authentication_policies_fragment.first_factor.id
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
    }
  ]
}
