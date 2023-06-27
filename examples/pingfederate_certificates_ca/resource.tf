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

resource "pingfederate_certificates_ca" "example" {
  # this property needs to contain base64 encoded value of your pem certificate.
  # when importing this resource, a subsequent apply is needed to store file_data into state for the future management of the resource
  file_data = ""
}