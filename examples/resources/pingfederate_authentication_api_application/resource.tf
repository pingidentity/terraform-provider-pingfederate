resource "pingfederate_authentication_api_application" "myAuthenticationApiApplicationExample" {
  application_id             = "example"
  name                       = "example"
  url                        = "https://example.com"
  description                = "example"
  additional_allowed_origins = ["https://example1.com"]
  client_for_redirectless_mode_ref = {
    id = pingfederate_oauth_client.myOauthClientExample.id
  }
}
