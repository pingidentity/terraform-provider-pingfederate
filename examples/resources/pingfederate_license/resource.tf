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
# this resource does not support import
resource "pingfederate_license" "licenseExample" {
  # this property needs to contain base64 encoded value of your license.
	file_data = ""
}