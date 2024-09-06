resource "pingfederate_session_application_policy" "sessionApplicationPolicy" {
  idle_timeout_mins = 60
  max_timeout_mins  = 60
}
