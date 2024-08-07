resource "pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate" "issuerCert" {
  certificate_id = "mycertid"
  # Include base64-encoded cert data here
  file_data = ""
}