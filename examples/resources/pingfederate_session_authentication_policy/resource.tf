resource "pingfederate_idp_adapter" "http_basic" {
  adapter_id = "HTTPBasicAdapter"
  name       = "HTTPBasic"
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.httpbasic.idp.HttpBasicIdpAuthnAdapter"
  }

  configuration = {
    fields = [
      {
        name  = "Realm",
        value = "example"
      },
      {
        name  = "Challenge Retries",
        value = "3"
      }
    ]
    tables = [
      {
        name = "Credential Validators"
        rows = [
          {
            fields = [
              {
                name  = "Password Credential Validator Instance"
                value = "simple"
              }
            ]
            defaultRow = false
          }
        ]
      }
    ]
  }

  attribute_contract = {
    core_attributes = [
      {
        name      = "username"
        pseudonym = true
      }
    ]
  }

  attribute_mapping = {
    attribute_contract_fulfillment = {
      username = {
        source = {
          type = "ADAPTER"
        }
        value = "username"
      }
    }
  }
}

resource "pingfederate_session_authentication_policy" "sessionAuthenticationPolicy" {
  policy_id = "MyHttpAdapterPolicy"
  authentication_source = {
    source_ref = {
      id = pingfederate_idp_adapter.http_basic.id
    }
    type = "IDP_ADAPTER"
  }
  authn_context_sensitive = false
  enable_sessions         = false
  persistent              = false
  timeout_display_unit    = "MINUTES"
  user_device_type        = "PRIVATE"
}