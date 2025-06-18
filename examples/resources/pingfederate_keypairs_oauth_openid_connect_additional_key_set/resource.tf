resource "pingfederate_keypairs_signing_key" "rsa_oidc_signing_1" {
  key_id    = "oidcsigningkey1"
  file_data = filebase64("./assets/oidc_signing_1.p12")
  format    = "PKCS12"
  password  = var.oidc_signing_pkcs12_password1
}

resource "pingfederate_keypairs_signing_key" "rsa_oidc_signing_2" {
  key_id    = "oidcsigningkey2"
  file_data = filebase64("./assets/oidc_signing_2.p12")
  format    = "PKCS12"
  password  = var.oidc_signing_pkcs12_password2
}

resource "pingfederate_oauth_issuer" "oauthIssuer1" {
  issuer_id   = "example"
  description = "example description"
  host        = "bxretail.org"
  name        = "example"
  path        = "/example"
}

resource "pingfederate_oauth_issuer" "oauthIssuer2" {
  issuer_id   = "example2"
  description = "example description"
  host        = "bxretail2.org"
  name        = "example2"
  path        = "/example"
}

resource "pingfederate_keypairs_oauth_openid_connect_additional_key_set" "keypairsOAuthOpenIDConnectAdditionalKeySet" {
  name = "My Key Set"

  issuers = [
    {
      id = pingfederate_oauth_issuer.oauthIssuer1.id
    },
    {
      id = pingfederate_oauth_issuer.oauthIssuer2.id
    }
  ]
  signing_keys = {
    rsa_active_cert_ref = {
      id = pingfederate_keypairs_signing_key.rsa_oidc_signing_2.id
    }
    rsa_previous_cert_ref = {
      id = pingfederate_keypairs_signing_key.rsa_oidc_signing_1.id
    }
    rsa_publish_x5c_parameter = true
  }
}