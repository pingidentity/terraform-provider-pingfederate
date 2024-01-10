resource "pingfederate_authentication_policies_fragment" "myPolicyFragment" {
  fragment_id = "myFragment"
  name        = "Verify_Register"
  description = "Sample Registration"
  root_node = {
    action = {
      authn_source_policy_action = {
        authentication_source = {
          type = "IDP_ADAPTER"
          source_ref = {
            id = "MyAdapter"
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
              id = "MyPolicyContract"
            }
            attribute_mapping = {
              attribute_sources = []
              attribute_contract_fulfillment = {
                "firstName" : {
                  source = {
                    type = "ADAPTER",
                    id   = "MyAdapter"
                  }
                  value = "firstName"
                }
                "lastName" : {
                  source = {
                    type = "ADAPTER",
                    id   = "MyAdapter"
                  }
                  value = "lastName"
                }
                "subject" : {
                  source = {
                    type = "ADAPTER",
                    id   = "MyAdapter"
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
                    id   = "MyAdapter"
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
    id = "MyPolicyContract"
  }
  outputs = {
    id = "MyPolicyContract"
  }

}
