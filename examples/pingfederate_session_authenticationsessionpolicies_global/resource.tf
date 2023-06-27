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

resource "pingfederate_session_authenticationsessionpolicies_global" "sessionAuthenticationSessionPoliciesGlobalExample" {
	enable_sessions = true
  persistent_sessions = false
  hash_unique_user_key_attribute = true
  idle_timeout_mins = 60
  idle_timeout_display_unit = "MINUTES"
  max_timeout_mins = 90
  max_timeout_display_unit = "MINUTES"
}