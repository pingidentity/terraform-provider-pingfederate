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
        name  = "Audience",
        value = "myaudience"
      }
    ]
  }
  name = "My SAML2 token processor"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.processor.saml.Saml20TokenProcessor"
  }
}

resource "pingfederate_oauth_token_exchange_processor_policy" "processorPolicy" {
  policy_id = "mypolicy"
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