terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username = "administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:9999"
}
# this resource does not support import as the PF API only supports PUT Method
resource "pingfederate_authenticationapi_settings" "authenticationApiSettingsExample" {
	api_enabled = true
  enable_api_descriptions = false
  restrict_access_to_redirectless_mode = false
  include_request_context = true
  # To remove a previously added default application ref, change id and location values to empty strings
  default_application_ref = {
    id = ""
    location = ""
  }
}