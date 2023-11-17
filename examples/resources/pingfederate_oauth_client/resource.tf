terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.1.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username                            = "administrator"
  password                            = "2FederateM0re"
  https_host                          = "https://localhost:9999"
  insecure_trust_all_tls              = true
  x_bypass_external_validation_header = true
}

# resource "pingfederate_oauth_client" "myOauthClient" {
#   client_id            = "myOauthClient"
#   grant_types          = ["DEVICE_CODE"]
#   name                 = "myOauthClient"
#   enabled              = true
#   bypass_approval_page = false
#   description          = "myOauthClient"
#   logo_url             = "https://example.com"
# }

resource "pingfederate_oauth_client" "myOtherOauthClient" {
  client_id = "myOtherOauthClient"
  enabled   = true
  redirect_uris = [
    "https://example.com"
  ]
  grant_types = [
    "IMPLICIT",
    "AUTHORIZATION_CODE",
    "RESOURCE_OWNER_CREDENTIALS",
    # "CLIENT_CREDENTIALS",
    "REFRESH_TOKEN",
    "EXTENSION",
    "DEVICE_CODE",
    "ACCESS_TOKEN_VALIDATION",
    # "CIBA",
    "TOKEN_EXCHANGE"
  ]
  name = "testing"
  # logo_url = "https://example.com"
  # allow_authentication_api_init = false
  # bypass_approval_page = true
  # require_pushed_authorization_requests           = false
  require_jwt_secured_authorization_response_mode = false
  restrict_scopes                                 = false
  # restricted_scopes = [
  #   "openid"
  # ]
  # exclusive_scopes = [
  #   "scope"
  # ]
  restricted_response_types = [
    "code",
    "code id_token",
    "code id_token token",
    "code token",
    "id_token",
    "id_token token",
    "token"
  ]
  # restrict_to_default_access_token_manager = true
  # validate_using_all_eligible_atms = false
  # oidc_policy = {
  # id_token_signing_algorithm = "HS256"
  # grant_access_session_revocation_api         = false
  # grant_access_session_session_management_api = true
  # ping_access_logout_capable    = false
  # pairwise_identifier_user_type = true
  # sector_identifier_uri         = "https://example.com"
  # id_token_encryption_algorithm               = "A192GCMKW"
  # id_token_content_encryption_algorithm = "AES_128_CBC_HMAC_SHA_256"
  # }

  client_auth = {
    type                   = "CERTIFICATE"
    secret                 = "mySecretValue"
    client_cert_issuer_dn  = "EMAILADDRESS=test@gmail.com, CN=terraformtest, OU=Devops, O=ping Identity, L=san Jose, ST=SJC, C=US"
    client_cert_subject_dn = "dn=subject"
    # secondary_secrets = [
    #   {
    #     secret      = "myOtherSecretValue"
    #     expiry_time = "2025-01-02T15:24:00Z"
    #   }
    # ]
  }

  #  test this
  jwks_settings = {
    jwks_url = "https://example.com"
  }
  # require_proof_key_for_code_exchange = false
  # ciba_delivery_mode           = "PING"
  # ciba_polling_interval        = 1
  # ciba_require_signed_requests = false
  # ciba_user_code_supported     = false
  # ciba_notification_endpoint   = "https://example.com"
  # jwt_secured_authorization_response_mode_encryption_algorithm         = "RSA_OAEP"
  # jwt_secured_authorization_response_mode_content_encryption_algorithm = "AES_128_CBC_HMAC_SHA_256"


  #  either of these two require jwks_settings

  #  this cannot be defined
  # token_introspection_signing_algorithm = "RS256"

  default_access_token_manager_ref = {
    id = "test"
  }

  # token_introspection_encryption_algorithm         = "DIR"
  token_introspection_content_encryption_algorithm = "AES_128_CBC_HMAC_SHA_256"

  require_signed_requests = false


  # token_exchange_processor_policy_ref = {
  #   id = "tokenexchangeprocessorpolicy"
  # }
}
