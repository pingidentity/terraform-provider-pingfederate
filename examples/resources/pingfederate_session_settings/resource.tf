resource "pingfederate_session_settings" "sessionSettingsExample" {
  track_adapter_sessions_for_logout = false
  revoke_user_session_on_logout     = true
  session_revocation_lifetime       = 490
}
