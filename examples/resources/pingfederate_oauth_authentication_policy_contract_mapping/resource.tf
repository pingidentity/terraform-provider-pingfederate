resource "pingfederate_authentication_policy_contract" "my_awesome_policy_contract" {
  name = "Users"
  extended_attributes = [
    { name = "email" },
    { name = "given_name" },
    { name = "family_name" },
    { name = "directory_id" },
  ]
}

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
      value = "directory_id"
    }
  }

  authentication_policy_contract_ref = {
    id = pingfederate_authentication_policy_contract.my_awesome_policy_contract.id
  }

  issuance_criteria = {
    conditional_criteria = [
      {
        attribute_name = "OAuthAuthorizationDetails"
        condition      = "EQUALS"
        error_result   = "Invalid Authorization Details"
        source = {
          type = "CONTEXT"
        }
        value = "Auth Details"
      },
    ]
  }
}