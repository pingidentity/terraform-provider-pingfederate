resource "pingfederate_authentication_api_application" "authenticationApiApplicationExample" {
  name        = "My Example Application"
  description = "My example application that has the authentication API widget embedded, or implements the authentication API directly."

  url = "https://bxretail.org"
  additional_allowed_origins = [
    "https://bxretail.org",
    "https://bxretail.org/*",
    "https://bxretail.org/cb/*",
    "https://auth.bxretail.org/*",
  ]
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
