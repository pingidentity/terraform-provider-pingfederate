resource "pingfederate_oauth_ciba_server_policy_settings" "oauthCibaServerPolicySettingsExample" {
  default_request_policy_ref = {
    id = "exampleOauthCibaServerPolicy"
  }
}
