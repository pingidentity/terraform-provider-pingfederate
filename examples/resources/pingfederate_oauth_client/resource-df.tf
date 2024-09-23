resource "pingfederate_oauth_client" "df_client" {
  client_id = "df_client"
  name      = "Device Authorization"
  enabled   = true

  client_auth = {
    type   = "SECRET"
    secret = var.df_client_secret
  }

  default_access_token_manager_ref = {
    id = "jwt"
  }

  grant_types = ["DEVICE_CODE"]
}
