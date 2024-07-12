resource "pingfederate_oauth_client_settings" "oauthClientSettings" {
  client_metadata = [
    {
      parameter   = "authNexp"
      description = "Authentication Experience"
      multiValued = false
    },
    {
      parameter   = "useAuthApi"
      description = "Use the AuthN API"
      multiValued = false
    }
  ]
  dynamic_client_registration = {
    initial_access_token_scope = "urn:pingidentity:register-client"
    restrict_common_scopes     = false
    restricted_common_scopes   = []
    allowed_exclusive_scopes = [
      "urn:pingidentity:directory",
      "urn:pingidentity:scim"
    ]
    allowed_authorization_detail_types              = []
    enforce_replay_prevention                       = false
    require_signed_requests                         = false
    restrict_to_default_access_token_manager        = false
    persistent_grant_expiration_type                = "SERVER_DEFAULT"
    persistent_grant_idle_timeout_type              = "SERVER_DEFAULT"
    client_cert_issuer_type                         = "NONE"
    refresh_rolling                                 = "SERVER_DEFAULT"
    refresh_token_rolling_interval_type             = "SERVER_DEFAULT"
    policy_refs                                     = []
    device_flow_setting_type                        = "SERVER_DEFAULT"
    require_proof_key_for_code_exchange             = false
    ciba_require_signed_requests                    = false
    ciba_polling_interval                           = 3
    rotate_registration_access_token                = true
    rotate_client_secret                            = true
    allow_client_delete                             = false
    refresh_token_rolling_grace_period_type         = "SERVER_DEFAULT"
    retain_client_secret                            = false
    client_secret_retention_period_type             = "SERVER_DEFAULT"
    require_jwt_secured_authorization_response_mode = false
  }
}