resource "pingfederate_oauth_access_token_manager_settings" "oauthTokenManagersSettings" {
  default_access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.example.id
  }
}