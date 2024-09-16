resource "pingfederate_authentication_api_application" "authenticationApiApplicationExample" {
  name        = "My Example Application"
  description = "My example application that has the authentication API widget embedded, or implements the authentication API directly."

  url = "https://bxretail.org"
  additional_allowed_origins = [
    "https://bxretail.org",
    "https://bxretail.org/*",
    "https://bxretail.org/cb/*",
    "https://auth.bxretail.org/*",
  ]

  client_for_redirectless_mode_ref = {
    id = pingfederate_oauth_client.oauthClientExample.id
  }
}
