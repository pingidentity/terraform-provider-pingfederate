resource "pingfederate_certificate_ca" "myCertificateCa" {
  custom_id = "example"
  file_data = ""
}

data "pingfederate_certificate_ca" "myCertificateCa" {
  id = "pingfederate_certificate_ca.myCertificateCa.custom_id"
}

resource "pingfederate_certificate_ca" "certificateCaExample" {
  custom_id = "${data.pingfederate_certificate_ca.myCertificateCa.id}2"
  file_data = ""
}