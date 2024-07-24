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
  }

  issuance_criteria = {
    conditional_criteria = [
      {
        attribute_name = "OAuthAuthorizationDetails"
        condition      = "EQUALS"
        error_result   = "Invalid Authorization Details"
        source = {
          type = "CONTEXT"
        }
        value = "Auth Details"
      },
    ]
  }
}