resource "pingfederate_session_authentication_policy" "sessionAuthenticationPolicy" {
  policy_id = "PingOneProtect"
  authentication_source = {
    source_ref = {
        id = "PingOneProtect"
    }
    type = "IDP_ADAPTER"
  }
  authn_context_sensitive = false
  enable_sessions = false
  persistent = false
  timeout_display_unit = "MINUTES"
  user_device_type = "PRIVATE"
}