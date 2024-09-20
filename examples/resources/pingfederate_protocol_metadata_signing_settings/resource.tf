resource "pingfederate_keypairs_signing_key" "signingKey" {
  file_data = filebase64("./assets/signingkey.p12")
  password  = var.signing_key_password
  format    = "PKCS12"
}

resource "pingfederate_protocol_metadata_signing_settings" "signingSettings" {
  signature_algorithm = "SHA256withRSA"
  signing_key_ref = {
    id = pingfederate_keypairs_signing_key.signingKey.key_id
  }
}
