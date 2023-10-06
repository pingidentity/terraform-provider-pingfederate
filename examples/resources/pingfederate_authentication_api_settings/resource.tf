# this resource does not support import as the PF API only supports PUT Method
resource "pingfederate_authentication_api_settings" "myAuthenticationApiSettings" {
  api_enabled                          = true
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
  # To remove a previously added default application ref, change id and location values to empty strings
  default_application_ref = {
    id       = ""
    location = ""
  }
}
