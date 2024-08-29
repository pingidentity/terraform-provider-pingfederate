resource "pingfederate_oauth_ciba_server_policy_settings" "policySettings" {
  default_request_policy_ref = {
    id = pingfederate_oauth_ciba_server_policy.examplePolicy.policy_id
  }
}
