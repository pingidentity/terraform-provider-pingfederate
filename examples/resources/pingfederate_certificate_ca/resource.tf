resource "pingfederate_certificate_ca" "myCertificateCa" {
  # this property needs to contain base64 encoded value of your pem certificate.
  # when importing this resource, a subsequent apply is needed to store file_data into state for the future management of the resource
  file_data = ""
}
