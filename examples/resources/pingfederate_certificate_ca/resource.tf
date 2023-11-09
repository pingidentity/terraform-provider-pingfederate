resource "pingfederate_certificate_ca" "myCertificateCa" {
  certificate_ca_id = "MyCertificateCA"
  # this property needs to contain base64 encoded value of your pem certificate.
  file_data = ""
}
