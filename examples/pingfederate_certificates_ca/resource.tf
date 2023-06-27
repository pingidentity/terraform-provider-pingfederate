terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.0"
      source = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username = "Administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:9999"
}
# this resource does not support import
resource "pingfederate_certificates_ca" "example" {
  # this property needs to contain base64 encode value of your pem certificate.
  file_data = ""
}