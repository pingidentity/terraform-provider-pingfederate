resource "pingfederate_oauth_auth_server_settings_scopes_common_scope" "oauthAuthServerSettingsScopesCommonScope" {
  dynamic     = true
  description = "example"
  name        = "*exampleCommonScope"
}
