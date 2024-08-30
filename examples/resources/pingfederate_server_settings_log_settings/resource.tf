resource "pingfederate_server_settings_log_settings" "logSettings" {
  log_categories = [
    {
      id      = "core"
      enabled = true
    },
    {
      id      = "policytree"
      enabled = false
    },
    {
      id      = "trustedcas"
      enabled = false
    },
    {
      id      = "xmlsig"
      enabled = false
    },
    {
      id      = "requestheaders"
      enabled = true
    },
    {
      id      = "requestparams"
      enabled = false
    },
    {
      id      = "restdatastore"
      enabled = false
    },
  ]
}
