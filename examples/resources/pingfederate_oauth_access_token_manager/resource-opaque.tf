resource "pingfederate_oauth_access_token_manager" "internally_managed_example" {
  manager_id = "internallyManagedReferenceOatm"
  name       = "internallyManagedReferenceExample"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = []
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
    coreAttributes = []
    extended_attributes = [
      {
        name         = "extended_contract"
        multi_valued = true
      }
    ]
  }
  selection_settings = {
    resource_uris = []
  }
  access_control_settings = {
    restrict_clients = false
    allowedClients   = []
  }
  session_validation_settings = {
    check_valid_authn_session       = false
    check_session_revocation_status = false
    update_authn_session_activity   = false
    include_session_id              = false
  }
}
