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

resource "pingfederate_oauth_auth_server_settings_scopes_common_scopes" "oauthAuthServerSettingsScopesCommonScopesExample" {
	dynamic = true
  description = "example"	
	name = "*exampleCommonScope"
}