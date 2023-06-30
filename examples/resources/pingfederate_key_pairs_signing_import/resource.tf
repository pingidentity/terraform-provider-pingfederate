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

# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #
resource "pingfederate_key_pairs_signing_import" "keyPairsSigningImportExample" {
  file_data = "example"	
  format = "PKCS12"
  # This value will be stored into your state file 
	password = "example"
}