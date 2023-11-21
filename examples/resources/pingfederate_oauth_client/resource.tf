resource "pingfederate_oauth_client" "myOauthClient" {
  client_id = "oauthClientId"
  grant_types = [
    "IMPLICIT",
    "AUTHORIZATION_CODE",
    "RESOURCE_OWNER_CREDENTIALS",
    "REFRESH_TOKEN",
    "EXTENSION",
    "DEVICE_CODE",
    "ACCESS_TOKEN_VALIDATION",
    "CIBA",
    "TOKEN_EXCHANGE"
  ]
  name                          = "myOauthClient"
  allow_authentication_api_init = false
  bypass_approval_page          = true
  ciba_delivery_mode            = "PING"
  ciba_polling_interval         = 1
  ciba_require_signed_requests  = true
  ciba_user_code_supported      = false
  ciba_notification_endpoint    = "https://example.com"
  client_auth = {
    type   = "SECRET"
    secret = "mySecretValue"
    secondary_secrets = [
      {
        secret      = "myOtherSecretValue"
        expiry_time = "2030-01-02T15:24:00Z"
      }
    ]
  }
  enabled = true
  jwks_settings = {
    jwks_url = "https://example.com"
  }
  jwt_secured_authorization_response_mode_encryption_algorithm         = "RSA_OAEP"
  jwt_secured_authorization_response_mode_content_encryption_algorithm = "AES_128_CBC_HMAC_SHA_256"
  logo_url                                                             = "https://example.com"
  oidc_policy = {
    id_token_signing_algorithm                  = "HS256"
    grant_access_session_revocation_api         = false
    grant_access_session_session_management_api = true
    ping_access_logout_capable                  = false
    pairwise_identifier_user_type               = true
    sector_identifier_uri                       = "https://example.com"
    id_token_encryption_algorithm               = "A192GCMKW"
    id_token_content_encryption_algorithm       = "AES_128_CBC_HMAC_SHA_256"
  }
  redirect_uris = [
    "https://example.com"
  ]
  require_jwt_secured_authorization_response_mode = false
  require_pushed_authorization_requests           = false
  require_proof_key_for_code_exchange             = false
  require_signed_requests                         = true
  restrict_scopes                                 = true
  restricted_scopes = [
    "openid"
  ]
  restricted_response_types = [
    "code",
    "code id_token",
    "code id_token token",
    "code token",
    "id_token",
    "id_token token",
    "token"
  ]
  restrict_to_default_access_token_manager         = false
  token_introspection_signing_algorithm            = "RS256"
  token_introspection_encryption_algorithm         = "DIR"
  token_introspection_content_encryption_algorithm = "AES_128_CBC_HMAC_SHA_256"
  validate_using_all_eligible_atms                 = false
}
