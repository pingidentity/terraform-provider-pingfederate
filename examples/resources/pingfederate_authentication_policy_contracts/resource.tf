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

resource "pingfederate_authentication_policy_contracts" "authenticationPolicyContractsExample" {
  core_attributes = [{name = "subject"}]
  extended_attributes = [{name = "extended_attribute"},{name = "extended_attribute2"}]
  name = "example"
}