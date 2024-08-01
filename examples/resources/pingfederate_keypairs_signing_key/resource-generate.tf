resource "pingfederate_keypairs_signing_key" "signingKey" {
  key_id = "signingKey"
  city = "Austin"
  common_name = "Example"
  country = "US"
  key_algorithm = "RSA"
  key_size = 2048
  organization = "Ping Identity"
  organization_unit = "Engineering"
  signature_algorithm = "SHA256withRSA"
  state = "Texas"
  subject_alternative_names = ["example.com"]
  valid_days = 365
}
