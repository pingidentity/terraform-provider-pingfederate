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

resource "pingfederate_oauth_issuer" "myOauthIssuer" {
  custom_id   = "MyOauthIssuer"
  description = "example description"
  host        = "example"
  name        = "example"
  path        = "/example"
}

data "pingfederate_oauth_issuer" "myOauthIssuer" {
  custom_id = pingfederate_oauth_issuer.myOauthIssuer.custom_id
}