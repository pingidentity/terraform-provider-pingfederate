resource "pingfederate_keypairs_oauth_openid_connect_additional_key_set" "keypairsOAuthOpenIDConnectAdditionalKeySet" {
  set_id = "my-key-set"
  issuers = [
    {
      id = "issuer-id"
    }
  ]
  name = "My Key Set"
  signing_keys = {
    rsa_active_cert_ref = {
      id = "rsaactive"
    }
    rsa_previous_cert_ref = {
      id = "rsaprevious"
    }
    rsa_publish_x5c_parameter = true
  }
}