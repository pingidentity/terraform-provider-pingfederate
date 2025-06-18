resource "pingfederate_oauth_ciba_server_policy_request_policy" "requestPolicy" {
  policy_id = "CIBAPolicy"
  name      = "CIBA Policy"
  authenticator_ref = {
    id = "authenticatorId"
  }
  user_code_pcv_ref = {
    id = "PingDirectoryPCV"
  }
  transaction_lifetime = 120
  identity_hint_mapping = {
    attribute_contract_fulfillment = {
      "USER_CODE_USER_NAME" = {
        source = {
          type = "REQUEST"
        }
        value = "IDENTITY_HINT_SUBJECT"
      }
      "subject" = {
        source = {
          type = "REQUEST"
        }
        value = "IDENTITY_HINT_SUBJECT"
      }
      "USER_KEY" = {
        source = {
          type = "REQUEST"
        }
        value = "IDENTITY_HINT_SUBJECT"
      }
    }
  }
}