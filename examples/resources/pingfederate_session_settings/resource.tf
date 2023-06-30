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

resource "pingfederate_session_settings" "sessionSettingsExample" {
	track_adapter_sessions_for_logout = false
  revoke_user_session_on_logout = true
  session_revocation_lifetime = 490
}