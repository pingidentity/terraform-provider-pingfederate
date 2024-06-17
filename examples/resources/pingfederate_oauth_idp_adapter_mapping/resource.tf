resource "pingfederate_oauth_idp_adapter_mapping" "oauthIdpAdapterMapping" {
  attribute_contract_fulfillment = {
    "USER_NAME" = {
      source = {
        type = "ADAPTER"
      }
      value = "subject"
    }
    "USER_KEY" = {
      source = {
        type = "ADAPTER"
      }
      value = "uid"
    }
  }
  mapping_id = "idpAdapterId"
}