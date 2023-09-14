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
  name = "test_token_manager"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = [],
    fields = [
      {
        name = "Token Length",
        value = "28"
      }
    ]
  }
  selection_settings = {
    inherited: true
    resource_uris = [
      "http://resource",
    ]
  }
}