resource "pingfederate_session_authentication_session_policies_global" "sessionApplicationSessionPolicyExample" {
  idle_timeout_mins = 60
  max_timeout_mins  = 60
}
