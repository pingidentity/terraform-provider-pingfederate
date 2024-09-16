resource "pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate" "issuerCert" {
  file_data = filebase64("path/to/my/issuercert.pem")
}
