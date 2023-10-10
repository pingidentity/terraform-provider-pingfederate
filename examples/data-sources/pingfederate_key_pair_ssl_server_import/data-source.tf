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

resource "pingfederate_key_pair_ssl_server_import" "myKeyPairsSslServerImport" {
  id    = "example"
}

data "pingfederate_key_pair_ssl_server_import" "myKeyPairsSslServerImport" {
  id    = "example"
}
resource "pingfederate_key_pair_ssl_server_import" "keyPairsSslServerImportExample" {
  id    = "example"
}