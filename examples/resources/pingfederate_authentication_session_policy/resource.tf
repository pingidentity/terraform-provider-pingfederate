terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 1.0.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username   = "administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:9999"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  insecure_trust_all_tls = true
  x_bypass_external_validation_header = true
  product_version = "12.0.0"
}

resource "pingfederate_authentication_session_policy" "example" {
  policy_id = "example"
  enable_sessions = false
  authentication_source = {
    source_ref = {
      id = "OTIdPJava"
    }
    type = "IDP_ADAPTER"
  }
}

