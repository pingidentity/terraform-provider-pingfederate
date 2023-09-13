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
  insecure_trust_all_tls = true
}

resource "pingfederate_oauth_access_token_managers" "oauthAccessTokenManagersExample" {
    id = "test_id"
    name = "test_token_manager"
    selection_settings = {
        resource_uris = ["http://resource",]
    }
}