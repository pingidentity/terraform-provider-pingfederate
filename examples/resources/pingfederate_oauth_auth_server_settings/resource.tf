resource "pingfederate_oauth_auth_server_settings" "example" {
  authorization_code_entropy               = 20
  authorization_code_timeout               = 50
  disallow_plain_pkce                      = false
  include_issuer_in_authorization_response = false
  track_user_sessions_for_logout           = false
  client_secret_retention_period           = 0

  # Pushed Authorization Request (PAR) Settings
  par_status            = "ENABLED"
  par_reference_timeout = 60
  par_reference_length  = 24

  # Refresh Token and Persistent Grant Settings
  persistent_grant_lifetime                  = -1
  persistent_grant_lifetime_unit             = "DAYS"
  persistent_grant_idle_timeout              = 30
  persistent_grant_idle_timeout_time_unit    = "DAYS"
  refresh_token_length                       = 40
  roll_refresh_token_values                  = true
  refresh_rolling_interval                   = 1
  refresh_token_rolling_grace_period         = 0
  allow_unidentified_client_ro_creds         = false
  allow_unidentified_client_extension_grants = false

  # Persistent Grant Extended Attributes
  persistent_grant_contract = {
    extended_attributes = [
      {
        name = "Persistent Grant Attribute 1"
      },
      {
        name = "Persistent Grant Attribute 2"
      },
    ]
  }

  # Authorization Consent
  bypass_authorization_for_approved_grants   = false
  bypass_authorization_for_approved_consents = false
  user_authorization_consent_page_setting    = "INTERNAL"

  # Cross-Origin Resource Sharing Settings
  allowed_origins = [
    "https://bxretail.org/path1/*",
    "https://bxretail.org/path/*",
    "https://bxretail.org/*",
  ]

  # Device Authorization Grant Settings
  user_authorization_url              = "https://bxretail.org/device"
  registered_authorization_path       = "/deviceAuthz"
  pending_authorization_timeout       = 550
  device_polling_interval             = 1
  activation_code_check_mode          = "BEFORE_AUTHENTICATION"
  bypass_activation_code_confirmation = false

  # JWT Secured Authorization Response Mode (JARM)
  jwt_secured_authorization_response_mode_lifetime = 600

  # Demonstrating Proof-of-Possession (DPoP)
  dpop_proof_require_nonce             = false
  dpop_proof_enforce_replay_prevention = false

  # Scopes
  scopes = [
    {
      name        = "my_custom_user_data",
      description = "A scope to represent the user's custom data in an application",
    },
    {
      name        = "delegated_custom_user_data",
      description = "A scope to represent the other user's custom data in an application, that the current user has access to",
    },
    {
      name        = "custom_app_data_*",
      description = "A scope to represent dynamic custom application data in an application",
      dynamic     = true
    }
  ]

  scope_groups = [
    {
      name        = "custom_user_data",
      description = "A scope group that represents all custom user data in an application that a user should have access to"
      scopes = [
        "delegated_custom_user_data",
        "my_custom_user_data",
      ]
    }
  ]

  # Exclusive Scopes
  exclusive_scopes = [
    {
      name        = "data_app1_a",
      description = "A scope to represent data that is exclusive to the app1 client",
    },
    {
      name        = "data_app1_b",
      description = "A scope to represent additional data that is exclusive to the app1 client",
    },
  ]

  exclusive_scope_groups = [
    {
      name        = "data_app1",
      description = "A scope group to represent all the scoped data that is exclusive to the app1 client"
      scopes = [
        "data_app1_a",
        "data_app1_b",
      ]
    }
  ]
}