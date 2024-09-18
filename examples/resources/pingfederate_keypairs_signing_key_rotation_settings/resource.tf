resource "pingfederate_keypairs_signing_key" "signingKey" {
  key_id    = "signingkey"
  file_data = filebase64("./assets/signingkey.p12")
  password  = var.signing_key_password
  format    = "PKCS12"
}

resource "pingfederate_keypairs_signing_key_rotation_settings" "keyRotationSettings" {
  key_pair_id            = pingfederate_keypairs_signing_key.signingKey.key_id
  activation_buffer_days = 90
  creation_buffer_days   = 180
}