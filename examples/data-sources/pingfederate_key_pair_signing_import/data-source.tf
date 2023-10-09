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

resource "pingfederate_key_pair_signing_import" "myKeyPairsSigningImport" {
  file_data = "example"
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}

data "pingfederate_key_pair_signing_import" "myKeyPairsSigningImport" {
  file_data = "pingfederate_key_pair_signing_import.myKeyPairsSigningImport.file_data"
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}
resource "pingfederate_key_pair_signing_import" "keyPairsSigningImportExample" {
  file_data = "${data.pingfederate_key_pair_signing_import.myKeyPairsSigningImport.file_data}2"
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}