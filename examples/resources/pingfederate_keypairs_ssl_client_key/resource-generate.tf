resource "pingfederate_keypairs_ssl_client_key" "sslClientKey" {
  key_id                    = "sslclientkey"
  city                      = "Austin"
  common_name               = "example"
  country                   = "US"
  key_algorithm             = "RSA"
  key_size                  = 2048
  organization              = "BXRetail"
  organization_unit         = "Auth Services"
  signature_algorithm       = "SHA256withRSA"
  state                     = "Texas"
  subject_alternative_names = ["bxretail.org"]
  valid_days                = 365
}