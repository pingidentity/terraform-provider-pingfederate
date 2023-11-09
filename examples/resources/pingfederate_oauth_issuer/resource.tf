resource "pingfederate_oauth_issuer" "example" {
  oauth_issuer_id = "MyOauthIssuer"
  description     = "example description"
  host            = "example"
  name            = "example"
  path            = "/example"
}
