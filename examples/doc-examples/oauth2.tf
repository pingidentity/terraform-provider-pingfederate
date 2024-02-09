provider "pingfederate" {
  client_id                           = "clientid"
  client_secret                       = "clientsecret"
  scopes                              = ["scope"]
  token_url                           = "https://localhost:9031/as/token.oauth2"
  https_host                          = "https://localhost:9999"
  admin_api_path                      = "/pf-admin-api/v1"
  insecure_trust_all_tls              = true
  x_bypass_external_validation_header = true
  product_version                     = "12.0"
}
