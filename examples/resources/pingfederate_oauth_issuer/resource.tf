resource "pingfederate_oauth_issuer" "example" {
  issuer_id   = "MyOauthIssuer"
  description = "example description"
  host        = "example"
  name        = "example"
  path        = "/example"
}
