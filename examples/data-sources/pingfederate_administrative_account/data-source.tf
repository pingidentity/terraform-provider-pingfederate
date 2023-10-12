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

resource "pingfederate_administrative_account" "myAdministrativeAccount" {
  id       = "example"
  username = "data-source-example"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}

data "pingfederate_administrative_account" "myAdministrativeAccount" {
  id       = "example"
  username = pingfederate_administrative_account.myAdministrativeAccount.username
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}
resource "pingfederate_administrative_account" "administrativeAccountExample" {
  id       = "example"
  username = "${data.pingfederate_administrative_account.myAdministrativeAccount.username}2"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}