terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.1.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username                            = "administrator"
  password                            = "2FederateM0re"
  https_host                          = "https://localhost:9999"
  insecure_trust_all_tls              = true
  x_bypass_external_validation_header = true
}

resource "pingfederate_oauth_client" "myOauthClientExample" {
  client_id                     = "myOauthClientExample"
  name                          = "myOauthClientExample"
  grant_types                   = ["EXTENSION"]
  allow_authentication_api_init = true
}

resource "pingfederate_authentication_api_application" "myAuthenticationApiApplicationExample" {
  depends_on                 = [pingfederate_oauth_client.myOauthClientExample]
  application_id             = "example"
  name                       = "example"
  url                        = "https://example.com"
  description                = "example"
  additional_allowed_origins = ["https://example1.com"]
  client_for_redirectless_mode_ref = {
    id = pingfederate_oauth_client.myOauthClientExample.id
  }
}
