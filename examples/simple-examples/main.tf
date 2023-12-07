terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.4.0"
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

resource "pingfederate_administrative_account" "myAdministrativeAccount" {
  username    = "example"
  description = "description"
  password    = "2FederateM0re"
  roles       = ["USER_ADMINISTRATOR"]
}
