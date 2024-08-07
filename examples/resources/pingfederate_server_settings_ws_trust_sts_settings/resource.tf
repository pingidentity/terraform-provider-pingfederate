resource "pingfederate_server_settings_ws_trust_sts_settings" "wsTrustStstSettings" {
  basic_authn_enabled       = false
  client_cert_authn_enabled = false
  restrict_by_subject_dn    = false
  restrict_by_issuer_cert   = false
  subject_dns               = []
  users                     = []
  issuer_certs              = []
}