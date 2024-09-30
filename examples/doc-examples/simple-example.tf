terraform {
  required_version = ">=1.4"
  required_providers {
    pingfederate = {
      version = "~> 1.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username                            = "administrator"
  password                            = "2FederateM0re"
  https_host                          = "https://localhost:9999"
  admin_api_path                      = "/pf-admin-api/v1"
  insecure_trust_all_tls              = true
  x_bypass_external_validation_header = true
  product_version                     = "12.1"
}

resource "pingfederate_administrative_account" "administrativeAccount" {
  username    = "example"
  description = "description"
  password    = "2FederateM0re"
  roles       = ["USER_ADMINISTRATOR"]
}
