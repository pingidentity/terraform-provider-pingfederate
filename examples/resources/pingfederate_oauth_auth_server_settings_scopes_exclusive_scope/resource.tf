resource "pingfederate_oauth_auth_server_settings_scopes_exclusive_scope" "oauthAuthServerSettingsScopesExclusiveScope" {
  dynamic     = true
  description = "example"
  name        = "*exampleExclusiveScope"
}
