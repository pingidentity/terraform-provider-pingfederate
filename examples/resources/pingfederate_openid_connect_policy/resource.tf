resource "pingfederate_oauth_access_token_manager" "internally_managed_example" {
  manager_id = "internallyManagedReferenceOATM"
  name       = "Internally Managed Token Manager"

  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }

  configuration = {
    fields = [
      {
        name  = "Token Length"
        value = "56"
      },
      {
        name  = "Token Lifetime"
        value = "240"
      },
      {
        name  = "Lifetime Extension Policy"
        value = "NONE"
      },
      {
        name  = "Maximum Token Lifetime"
        value = ""
      },
      {
        name  = "Lifetime Extension Threshold Percentage"
        value = "30"
      },
      {
        name  = "Mode for Synchronous RPC"
        value = "3"
      },
      {
        name  = "RPC Timeout"
        value = "500"
      },
      {
        name  = "Expand Scope Groups"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name         = "givenName"
        multi_valued = false
      },
      {
        name         = "familyName"
        multi_valued = false
      },
      {
        name         = "email"
        multi_valued = false
      },
      {
        name         = "groups"
        multi_valued = true
      }
    ]
  }
  access_control_settings = {
    restrict_clients = false
  }
  session_validation_settings = {
    check_valid_authn_session       = false
    check_session_revocation_status = false
    update_authn_session_activity   = false
    include_session_id              = false
  }
}


resource "pingfederate_openid_connect_policy" "OIDCPolicy" {
  policy_id = "exampleOIDCPolicy"
  name      = "Example OpenID Connect Policy"

  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.internally_managed_example.id
  }

  attribute_contract = {
    extended_attributes = [
      { name = "customData1" },
      { name = "customData2", multi_valued = true },
      { name = "email" },
      { name = "firstName" },
      { name = "lastName" },
    ]
  }

  attribute_mapping = {
    attribute_contract_fulfillment = {
      "customData1" = {
        source = {
          type = "NO_MAPPING"
        }
      },
      "customData2" = {
        source = {
          type = "NO_MAPPING"
        }
      },
      "email" = {
        source = {
          type = "TOKEN"
        }
        value = "email"
      },
      "firstName" = {
        source = {
          type = "TOKEN"
        }
        value = "first_name"
      },
      "lastName" = {
        source = {
          type = "TOKEN"
        }
        value = "last_name"
      },
      "sub" = {
        source = {
          type = "TOKEN"
        }
        value = "directory_id"
      },
    }
  }

  scope_attribute_mappings = {
    "email" = {
      values = ["email"]
    },
    "my_custom_user_data" = {
      values = ["customData1", "customData2"]
    },
    "profile" = {
      values = ["firstName", "lastName"]
    },
  }

  return_id_token_on_refresh_grant = false
  include_sri_in_id_token          = true
  include_s_hash_in_id_token       = false
  include_user_info_in_id_token    = true
  reissue_id_token_in_hybrid_flow  = false
  id_token_lifetime                = 5
}
