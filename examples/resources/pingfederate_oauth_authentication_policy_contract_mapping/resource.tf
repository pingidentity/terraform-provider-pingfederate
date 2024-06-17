resource "pingfederate_oauth_authentication_policy_contract_mapping" "oauthAuthenticationPolicyContractMapping" {
  attribute_contract_fulfillment = {
    "USER_NAME" = {
      source = {
        type = "AUTHENTICATION_POLICY_CONTRACT"
      }
      value = "subject"
    }
    "USER_KEY" = {
      source = {
        type = "AUTHENTICATION_POLICY_CONTRACT"
      }
      value = "ImmutableID"
    }
  }
  authentication_policy_contract_ref = {
    id = "authPolicyContractId"
  }
}