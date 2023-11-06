resource "pingfederate_server_settings_log_settings" "serverSettingsLogSettingsExample" {
  log_categories = [
    {
      id      = "core"
      enabled = false
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
      enabled = false
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
