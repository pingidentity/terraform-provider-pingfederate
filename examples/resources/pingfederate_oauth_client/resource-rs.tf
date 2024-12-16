resource "pingfederate_oauth_client" "rs_client" {
  client_id = "rs_client"
  name      = "Resource Server Client"
  enabled   = true

  client_auth = {
    type   = "SECRET"
    secret = var.rs_client_secret
  }

  grant_types = ["ACCESS_TOKEN_VALIDATION"]

  validate_using_all_eligible_atms = true
}
