resource "pingfederate_authentication_api_settings" "myAuthenticationApiSettings" {
  api_enabled                          = true
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
}

data "pingfederate_authentication_api_settings" "myAuthenticationApiSettings" {
  depends_on = [
    pingfederate_authentication_api_settings.myAuthenticationApiSettings
  ]
}
resource "pingfederate_authentication_api_settings" "authenticationApiSettingsExample" {
  api_enabled                          = data.pingfederate_authentication_api_settings.myAuthenticationApiSettings.api_enabled
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
}