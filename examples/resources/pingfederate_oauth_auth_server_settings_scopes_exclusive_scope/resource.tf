resource "pingfederate_oauth_auth_server_settings_scopes_exclusive_scope" "myOauthAuthServerSettingsScopesExclusiveScope" {
  dynamic     = true
  description = "example"
  name        = "*exampleExclusiveScope"
}
