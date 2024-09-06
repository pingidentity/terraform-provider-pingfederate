resource "pingfederate_authentication_api_application" "authenticationApiApplicationExample" {
  application_id             = "example"
  name                       = "example"
  url                        = "https://example.com"
  description                = "example"
  additional_allowed_origins = ["https://example.com"]
}

resource "pingfederate_authentication_api_settings" "apiSettings" {
  api_enabled                          = true
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
  default_application_ref = {
    id = pingfederate_authentication_api_application.authenticationApiApplicationExample.application_id
  }
}
