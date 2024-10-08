resource "pingfederate_extended_properties" "example" {
  items = [
    {
      name        = "Attribute 1"
      description = "My single valued extended attribute"
    },
    {
      name         = "Attribute 2"
      description  = "My multi-valued extended attribute"
      multi_valued = true
    },
  ]
}

resource "pingfederate_oauth_client_registration_policy" "registrationPolicy" {
  policy_id = "myRegistrationPolicy"
  name      = "My client registration policy"

  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.client.registration.ResponseTypesConstraintsPlugin"
  }

  configuration = {
    fields = [
      {
        name  = "code"
        value = "true"
      },
      {
        name  = "code id_token"
        value = "true"
      },
      {
        name  = "code id_token token"
        value = "true"
      },
      {
        name  = "code token"
        value = "true"
      },
      {
        name  = "id_token"
        value = "true"
      },
      {
        name  = "id_token token"
        value = "true"
      },
      {
        name  = "token"
        value = "true"
      }
    ]
  }
}

resource "pingfederate_oauth_client_settings" "oauthClientSettings" {
  depends_on = [pingfederate_extended_properties.example]
  dynamic_client_registration = {
    initial_access_token_scope = "urn:pingidentity:register-client"
    restrict_common_scopes     = false
    restricted_common_scopes   = []
    allowed_exclusive_scopes = [
      "urn:pingidentity:directory",
      "urn:pingidentity:scim"
    ]
    enforce_replay_prevention                = false
    require_signed_requests                  = false
    restrict_to_default_access_token_manager = false
    persistent_grant_expiration_type         = "SERVER_DEFAULT"
    persistent_grant_idle_timeout_type       = "SERVER_DEFAULT"
    client_cert_issuer_type                  = "NONE"
    refresh_rolling                          = "SERVER_DEFAULT"
    refresh_token_rolling_interval_type      = "SERVER_DEFAULT"
    policy_refs = [
      {
        id = pingfederate_oauth_client_registration_policy.registrationPolicy.id
      }
    ]
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
