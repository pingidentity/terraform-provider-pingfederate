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
  username = "data-source-example"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}

data "pingfederate_administrative_account" "myAdministrativeAccount" {
  username = pingfederate_administrative_account.myAdministrativeAccount.username
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}
resource "pingfederate_administrative_account" "administrativeAccountExample" {
  username = "${data.pingfederate_administrative_account.myAdministrativeAccount.username}2"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}