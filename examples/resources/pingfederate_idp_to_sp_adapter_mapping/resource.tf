resource "pingfederate_idp_to_sp_adapter_mapping" "idpToSpAdapterMapping" {
  attribute_contract_fulfillment = {
    "subject" = {
      source = {
        type = "ADAPTER"
      }
      value = "subject"
    }
  }
  source_id = "OTIdPJava"
  target_id = "spadapter"
}
