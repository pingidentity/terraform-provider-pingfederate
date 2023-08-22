resource "pingfederate_oauth_auth_server_settings" "oauthAuthServerSettingsExample" {
  authorization_code_entropy          = 20
  authorization_code_timeout          = 50
  bypass_activation_code_confirmation = false
  default_scope_description           = "example"
  device_polling_interval             = 1
  pending_authorization_timeout       = 550
  refresh_rolling_interval            = 1
  refresh_token_length                = 40
  registered_authorization_path       = "/example"
  scopes = [
    {
      name        = "examplescope",
      description = "example scope",
      dynamic     = false
    }
  ]
  scope_groups = [
    {
      name        = "examplescopegroup",
      description = "example scope group"
      scopes      = ["examplescope"]
    }
  ]
  exclusive_scopes = [
    {
      name        = "exampleexclusivescope",
      description = "example scope",
      dynamic     = false
    }
  ]
  exclusive_scope_groups = [
    {
      name        = "exampleexclusivescopegroup",
      description = "example exclusive scope group"
      scopes      = ["exampleexclusivescope"]
    }
  ]
  disallow_plain_pkce                      = false
  include_issuer_in_authorization_response = false
  persistent_grant_lifetime                = -1
  persistent_grant_lifetime_unit           = "DAYS"
  persistent_grant_idle_timeout            = 30
  persistent_grant_idle_timeout_time_unit  = "DAYS"
  roll_refresh_token_values                = true
  refresh_token_rolling_grace_period       = 0
  persistent_grant_reuse_grant_types       = []
  persistent_grant_contract = {
    extended_attributes = [
      {
        name = "example_extended_attribute"
      }
    ]
  }
  bypass_authorization_for_approved_grants         = false
  allow_unidentified_client_ro_creds               = false
  allow_unidentified_client_extension_grants       = false
  token_endpoint_base_url                          = ""
  user_authorization_url                           = ""
  activation_code_check_mode                       = "BEFORE_AUTHENTICATION"
  user_authorization_consent_page_setting          = "INTERNAL"
  atm_id_for_oauth_grant_management                = ""
  scope_for_oauth_grant_management                 = ""
  allowed_origins                                  = []
  track_user_sessions_for_logout                   = false
  par_reference_timeout                            = 60
  par_reference_length                             = 24
  par_status                                       = "ENABLED"
  client_secret_retention_period                   = 0
  jwt_secured_authorization_response_mode_lifetime = 600
}