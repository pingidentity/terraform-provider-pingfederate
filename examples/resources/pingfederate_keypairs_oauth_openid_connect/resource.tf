resource "pingfederate_keypairs_oauth_openid_connect" "keypairsOAuthOpenIDConnect" {
  rsa_active_cert_ref = {
    id = "rsaactive"
  }
  rsa_decryption_active_cert_ref = {
    id = "rsadecryptionactive"
  }
  rsa_decryption_previous_cert_ref = {
    id = "rsadecryptionprevious"
  }
  rsa_previous_cert_ref = {
    id = "rsaprevious"
  }
  rsa_publish_x5c_parameter = true
  static_jwks_enabled       = true
}
