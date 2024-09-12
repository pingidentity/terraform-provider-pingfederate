resource "pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate" "issuerCert" {
  file_data = filebase64("path/to/my/issuercert.pem")
}

resource "pingfederate_server_settings_ws_trust_sts_settings" "wsTrustStstSettings" {
  basic_authn_enabled       = true
  client_cert_authn_enabled = true
  restrict_by_subject_dn    = true
  restrict_by_issuer_cert   = true
  subject_dns = [
    "cn=my-restricted-issuer1",
    "cn=my-restricted-issuer2",
  ]
  users = [
    {
      username = "basic_auth_user_1"
      password = var.basic_auth_user_1_password
    },
    {
      username = "basic_auth_user_2"
      password = var.basic_auth_user_2_password
    },
  ]
  issuer_certs = [
    {
      id = pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate.issuerCert.id
    }
  ]
}
