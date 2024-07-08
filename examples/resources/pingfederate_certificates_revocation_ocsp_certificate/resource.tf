resource "pingfederate_certificates_revocation_ocsp_certificate" "certificate" {
  certificate_id = "certid"
  # Include base64-encoded cert data here
  file_data = ""
}