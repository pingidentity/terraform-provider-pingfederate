resource "pingfederate_oauth_issuer" "myOauthIssuer" {
  issuer_id   = "MyOauthIssuer"
  description = "example description"
  host        = "example"
  name        = "example"
  path        = "/example"
}
