resource "pingfederate_keypairs_signing_key" "signingKeyGenerate" {
  common_name               = "AuthSigning"
  subject_alternative_names = ["bxretail.org", "192.168.0.1"]

  organization      = "BXRetail"
  organization_unit = "Auth Services"

  city    = "Austin"
  state   = "Texas"
  country = "US"

  key_algorithm       = "RSA"
  key_size            = 2048
  signature_algorithm = "SHA256withRSA"

  valid_days = 365
}