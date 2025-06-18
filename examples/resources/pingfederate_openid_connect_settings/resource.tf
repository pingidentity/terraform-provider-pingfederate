resource "pingfederate_session_settings" "sessionSettingsExample" {
  track_adapter_sessions_for_logout = false
  revoke_user_session_on_logout     = true
  session_revocation_lifetime       = 490
}

resource "pingfederate_openid_connect_policy" "oauthOIDCPolicyExample" {
  policy_id = "oidcPolicy"
  name      = "oidcPolicy"
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.example.manager_id
  }
  attribute_contract = {
    extended_attributes = []
  }
  attribute_mapping = {
    attribute_contract_fulfillment = {
      "sub" = {
        source = {
          type = "TOKEN"
        }
        value = "Username"
      }
    }
  }
  return_id_token_on_refresh_grant = false
  include_sri_in_id_token          = false
  include_s_hash_in_id_token       = false
  include_user_info_in_id_token    = false
  reissue_id_token_in_hybrid_flow  = false
  id_token_lifetime                = 5
}

resource "pingfederate_openid_connect_settings" "openIdConnectSettingsExample" {
  depends_on = [pingfederate_session_settings.sessionSettingsExample]
  default_policy_ref = {
    id = pingfederate_openid_connect_policy.oauthOIDCPolicyExample.policy_id
  }
}