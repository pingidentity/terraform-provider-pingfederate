resource "pingfederate_open_id_connect_settings" "openIdConnectSettingsExample" {
  default_policy_ref = {
    id = "oauth_example_policy"
  }
  session_settings = {
    track_user_sessions_for_logout = true
    revoke_user_session_on_logout  = false
    session_revocation_lifetime    = 180
  }
}