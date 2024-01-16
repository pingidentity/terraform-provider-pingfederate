resource "pingfederate_authentication_api_application" "authenticationApiApplicationExample" {
  application_id             = "example"
  name                       = "example"
  url                        = "https://example.com"
  description                = "example"
  additional_allowed_origins = ["https://example.com"]
  client_for_redirectless_mode_ref = {
    id = pingfederate_oauth_client.oauthClientExample.id
  }
}
