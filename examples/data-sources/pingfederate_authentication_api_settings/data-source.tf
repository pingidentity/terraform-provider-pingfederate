terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username               = "administrator"
  password               = "2FederateM0re"
  https_host             = "https://localhost:9999"
  insecure_trust_all_tls = true
}

resource "pingfederate_authentication_api_settings" "myAuthenticationApiSettings" {
  api_enabled                          = true
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
}

data "pingfederate_authentication_api_settings" "myAuthenticationApiSettings" {
  api_enabled                          = false
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
}
resource "pingfederate_authentication_api_settings" "authenticationApiSettingsExample" {
  api_enabled                          = data.pingfederate_authentication_api_settings.myAuthenticationApiSettings.api_enabled
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
}