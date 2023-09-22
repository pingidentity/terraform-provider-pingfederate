resource "pingfederate_session_authentication_session_policies_global" "sessionAuthenticationSessionPoliciesGlobalExample" {
  enable_sessions                = true
  persistent_sessions            = false
  hash_unique_user_key_attribute = true
  idle_timeout_mins              = 60
  idle_timeout_display_unit      = "MINUTES"
  max_timeout_mins               = 90
  max_timeout_display_unit       = "MINUTES"
}
