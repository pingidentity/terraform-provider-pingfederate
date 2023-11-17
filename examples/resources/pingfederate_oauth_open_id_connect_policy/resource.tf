resource "pingfederate_oauth_open_id_connect_policy" "oauthOIDCPolicyExample" {
  policy_id = "myOIDCPolicy"
  name      = "myOIDCPolicy"
  access_token_manager_ref = {
    id = "myATM"
  }
  attribute_contract = {
    extended_attributes = []
  }
  attribute_mapping = {
    attribute_contract_fulfillment = {
      "sub" = {
        source = {
          type = "TOKEN"
        }
        value = "Username"
      }
    }
    attribute_sources = []
    issuance_criteria = {
      conditional_criteria = []
    }
  }
  scope_attribute_mappings         = {}
  return_id_token_on_refresh_grant = false
  include_sri_in_id_token          = false
  include_s_hash_in_id_token       = false
  include_user_info_in_id_token    = false
  id_token_lifetime                = 5
}
