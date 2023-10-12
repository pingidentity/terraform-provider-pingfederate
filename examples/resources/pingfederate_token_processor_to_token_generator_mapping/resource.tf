resource "pingfederate_token_processor_to_token_generator_mapping" "tokenProcessorToTokenGeneratorMappingExample" {
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
  source_id = "tokenprocessor"
  target_id = "tokengenerator"
}
