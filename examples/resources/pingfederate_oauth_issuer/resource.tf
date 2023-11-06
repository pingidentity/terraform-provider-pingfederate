resource "pingfederate_oauth_issuer" "example" {
  custom_id = "MyOauthIssuer"
  description = "example description"
  host        = "example"
  name        = "example"
  path        = "/example"
}
