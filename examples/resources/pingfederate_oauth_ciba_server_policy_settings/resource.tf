resource "pingfederate_oauth_ciba_server_policy_settings" "myOauthCibaServerPolicySettingsExample" {
  default_request_policy_ref = {
    id = "myExampleOauthCibaServerPolicy"
  }
}
