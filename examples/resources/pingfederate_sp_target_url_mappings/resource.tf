resource "pingfederate_sp_target_url_mappings" "spTargetUrlMappings" {
  items = [
    {
      ref = {
        id = pingfederate_sp_adapter.reference_id.id
      }
      type = "SP_ADAPTER"
      url  = "https://www.bxretail.org/acct101/"
    },
    {
      ref = {
        id = pingfederate_sp_adapter.opentoken.id
      }
      type = "SP_ADAPTER"
      url  = "https://www.bxretail.org/*"
    },
  ]
}