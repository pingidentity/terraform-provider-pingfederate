resource "pingfederate_oauth_token_exchange_processor_policy" "my_awesome_processor_policy" {
  policy_id = "myProcessorPolicy"
  name      = "My processor policy"
  processor_mappings = [
    {
      attribute_contract_fulfillment = {
        "subject" = {
          source = {
            type = "TEXT"
          }
          value = "value"
        }
      }
      subject_token_processor = {
        id = pingfederate_idp_token_processor.saml2.processor_id
      }
      subject_token_type = "urn:ietf:params:oauth:token-type:saml2"
    }
  ]
}

resource "pingfederate_idp_token_processor" "saml2" {
  processor_id = "saml2TokenProcessor"
  attribute_contract = {
    core_attributes = [
      {
        name = "SAML_SUBJECT"
      }
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Audience"
        value = "myaudience"
      }
    ]
  }
  name = "My token processor"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.processor.saml.Saml20TokenProcessor"
  }
}

resource "pingfederate_oauth_token_exchange_processor_policy_mapping" "oauthTokenExchangeProcessorPolicyMapping" {
  attribute_contract_fulfillment = {
    "USER_NAME" = {
      source = {
        type = "TOKEN_EXCHANGE_PROCESSOR_POLICY"
      }
      value = "subject"
    }
    "USER_KEY" = {
      source = {
        type = "TOKEN_EXCHANGE_PROCESSOR_POLICY"
      }
      value = "subject"
    }
  }

  processor_policy_ref = {
    id = pingfederate_oauth_token_exchange_processor_policy.my_awesome_processor_policy.policy_id
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
