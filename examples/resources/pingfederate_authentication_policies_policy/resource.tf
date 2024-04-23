resource "pingfederate_authentication_policies_policy" "policy" {
  name      = "policy"
  policy_id = "policyid"
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
