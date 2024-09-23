resource "pingfederate_oauth_client" "spa_oic_client" {
  client_id = "spa_oic_client"
  name      = "OpenID Connect Authorization Code with PKCE"
  enabled   = true

  client_auth = {
    type = "NONE"
  }

  require_proof_key_for_code_exchange = true
  bypass_approval_page                = true

  default_access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.jwt_example.id
  }

  grant_types = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]

  redirect_uris = ["https://www.bxretail.org/oidc/callback"]

  oidc_policy = {
    policy_group = {
      id = pingfederate_openid_connect_policy.OIDCPolicy.id
    }
  }
}
