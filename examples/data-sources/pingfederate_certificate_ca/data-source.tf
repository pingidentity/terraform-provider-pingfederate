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

resource "pingfederate_certificate_ca" "myCertificateCa" {
  custom_id = "example"
}

data "pingfederate_certificate_ca" "myCertificateCa" {
  id = "pingfederate_certificate_ca.myCertificateCa.custom_id"
  depends_on = [
    pingfederate_certificate_ca.myCertificateCa
  ]
}
resource "pingfederate_certificate_ca" "certificateCaExample" {
  custom_id = "${data.pingfederate_certificate_ca.myCertificateCa.id}2"
}