resource "pingfederate_sp_target_url_mappings" "spTargetUrlMappings" {
  items = [
    {
      ref = {
        id = "myspadapter"
      }
      type = "SP_ADAPTER"
      url  = "*"
    }
  ]
}
