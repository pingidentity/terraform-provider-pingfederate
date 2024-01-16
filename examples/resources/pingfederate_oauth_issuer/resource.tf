resource "pingfederate_oauth_issuer" "oauthIssuer" {
  issuer_id   = "oauthIssuer"
  description = "example description"
  host        = "example"
  name        = "example"
  path        = "/example"
}
