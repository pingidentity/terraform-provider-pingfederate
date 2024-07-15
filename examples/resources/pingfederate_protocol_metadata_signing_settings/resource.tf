resource "pingfederate_protocol_metadata_signing_settings" "signingSettings" {
  signature_algorithm = "SHA256withRSA"
  signing_key_ref = {
    id = "mysigningkey"
  }
}