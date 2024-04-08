resource "pingfederate_oauth_access_token_mapping" "oauthAccessTokenMappingExample" {
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.myOauthAccessTokenManagerExample.id
  }
  attribute_contract_fulfillment = {
    "extended_contract" = {
      source = {
        type = "TEXT"
      }
      value = "Administrator"
    },
  }
  attribute_sources = [
    {
      jdbc_attribute_source = {
        column_names = [
          "GRANTEE",
          "IS_GRANTABLE",
          "ROLE_NAME",
        ]
        data_store_ref = {
          id = "ProvisionerDS"
        }
        description = "description"
        filter      = "$${client_id}"
        id          = "test"
        schema      = "INFORMATION_SCHEMA"
        table       = "ADMINISTRABLE_ROLE_AUTHORIZATIONS"
      }
    },
  ]
  context = {
    type = "CLIENT_CREDENTIALS"
  }

  issuance_criteria = {
    conditional_criteria = [
      {
        attribute_name = "ClientId"
        condition      = "EQUALS_CASE_INSENSITIVE"
        error_result   = "error"
        source = {
          type = "CONTEXT"
        }
        value = "text"
      },
    ]
  }
}
