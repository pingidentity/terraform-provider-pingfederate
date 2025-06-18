resource "pingfederate_oauth_client" "cc_secret_client" {
  client_id = "cc_secret_client"
  name      = "Client Credentials (Secret)"
  enabled   = true

  client_auth = {
    type   = "SECRET"
    secret = var.client_credentials_client_secret
  }

  grant_types = ["CLIENT_CREDENTIALS"]
}
