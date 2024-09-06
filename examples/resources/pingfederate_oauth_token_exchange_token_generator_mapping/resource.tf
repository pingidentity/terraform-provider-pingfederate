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
resource "pingfederate_sp_token_generator" "tokenGenerator" {
  generator_id = "myGenerator"
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
        name  = "Minutes Before"
        value = "60"
      },
      {
        name  = "Minutes After"
        value = "60"
      },
      {
        name  = "Issuer"
        value = "issuer"
      },
      {
        name  = "Signing Certificate"
        value = pingfederate_keypairs_signing_key.signingKey.key_id
      },
      {
        name  = "Signing Algorithm"
        value = "SHA1"
      },
      {
        name  = "Include Certificate in KeyInfo"
        value = "false"
      },
      {
        name  = "Include Raw Key in KeyValue"
        value = "false"
      },
      {
        name  = "Audience"
        value = "audience"
      },
      {
        name  = "Confirmation Method"
        value = "urn:oasis:names:tc:SAML:2.0:cm:sender-vouches"
      }
    ]
    tables = []
  }
  name = "My token generator"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.generator.saml.Saml20TokenGenerator"
  }
}

resource "pingfederate_oauth_token_exchange_token_generator_mapping" "exchangeGeneratorMapping" {
  attribute_sources = [
    {
      jdbc_attribute_source = {
        data_store_ref = {
          id = "ProvisionerDS"
        }
        id           = "attributesourceid"
        description  = "description"
        schema       = "INFORMATION_SCHEMA"
        table        = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
        filter       = "CONDITION"
        column_names = ["GRANTEE", "IS_GRANTABLE", "ROLE_NAME"]
      }
    }
  ]
  attribute_contract_fulfillment = {
    "SAML_SUBJECT" = {
      source = {
        type = "TEXT"
      },
      value = "value"
    }
  }
  issuance_criteria = {
    conditional_criteria = [
      {
        error_result = "error"
        source = {
          type = "CONTEXT"
        }
        attribute_name = "ClientIp"
        condition      = "EQUALS"
        value          = "value"
      }
    ]
  }
  source_id = pingfederate_oauth_token_exchange_processor_policy.processorPolicy.policy_id
  target_id = pingfederate_sp_token_generator.tokenGenerator.generator_id
}
