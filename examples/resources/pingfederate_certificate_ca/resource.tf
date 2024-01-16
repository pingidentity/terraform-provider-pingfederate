resource "pingfederate_certificate_ca" "certificateCa" {
  ca_id = "certificateCA"
  # this property needs to contain base64 encoded value of your pem certificate.
  file_data = ""
}
