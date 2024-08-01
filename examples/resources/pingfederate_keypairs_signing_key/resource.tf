resource "pingfederate_keypairs_signing_key" "signingKey" {
  key_id = "signingKey"
  file_data = filebase64("./assets/signingkey.p12")
  password = var.signing_key_password
  format = "PKCS12"
}
