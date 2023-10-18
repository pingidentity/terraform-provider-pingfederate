resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractExample" {
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
}

resource "pingfederate_sp_authentication_policy_contract_mapping" "spAuthenticationPolicyContractMappingExample" {
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
    "subject" = {
      source = {
        type = "TEXT"
      },
      value = "test"
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
  source_id = pingfederate_authentication_policy_contract.authenticationPolicyContractExample.id
  target_id = "spadapter"
}