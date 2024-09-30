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

resource "pingfederate_oauth_server_settings" "example" {
  authorization_code_entropy = 20
  authorization_code_timeout = 50
  refresh_token_length       = 40
  refresh_rolling_interval   = 1

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
}

resource "pingfederate_oauth_idp_adapter_mapping" "oauthIdpAdapterMapping" {
  mapping_id = pingfederate_idp_adapter.http_basic.id

  attribute_contract_fulfillment = {
    "USER_NAME" = {
      source = {
        type = "ADAPTER"
      }
      value = "username"
    }
    "USER_KEY" = {
      source = {
        type = "ADAPTER"
      }
      value = "username"
    }
    "Persistent Grant Attribute 1" = {
      source = {
        type = "ADAPTER"
      }
      value = "username"
    }
    "Persistent Grant Attribute 2" = {
      source = {
        type = "ADAPTER"
      }
      value = "username"
    }
  }

  depends_on = [
    pingfederate_oauth_server_settings.example
  ]
}