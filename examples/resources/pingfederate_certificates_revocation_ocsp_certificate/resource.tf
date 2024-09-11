resource "pingfederate_certificates_revocation_ocsp_certificate" "certificate" {
  file_data = filebase64("path/to/my/certificate.pem")
}
