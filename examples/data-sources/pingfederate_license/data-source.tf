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

resource "pingfederate_license" "myLicense" {
  # this property needs to contain base64 encoded value of your license.
  file_data = ""
}

data "pingfederate_license" "myLicense" {
  # this property needs to contain base64 encoded value of your license.
  file_data = pingfederate_license.myLicense.file_data
}
resource "pingfederate_license" "licenseExample" {
  # this property needs to contain base64 encoded value of your license.
  file_data = "${data.pingfederate_license.myLicense.file_data}2"
}