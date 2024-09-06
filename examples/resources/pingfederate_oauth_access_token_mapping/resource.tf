resource "pingfederate_oauth_access_token_manager" "jwt_example" {
  manager_id = "jsonWebTokenOatm"
  name       = "jsonWebTokenExample"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.access.token.management.plugins.JwtBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = [
      {
        name = "Symmetric Keys"
        rows = [
          {
            fields = [
              {
                name  = "Key ID"
                value = "keyidentifier"
              },
              {
                name  = "Key"
                value = var.access_token_manager_key
              },
              {
                name  = "Encoding"
                value = "b64u"
              }
            ]
            default_row = false
          }
        ]
      },
      {
        name = "Certificates"
        rows = []
      }
    ]
    fields = [
      {
        name  = "Token Lifetime"
        value = "120"
      },
      {
        name  = "Use Centralized Signing Key"
        value = "false"
      },
      {
        name  = "JWS Algorithm"
        value = ""
      },
      {
        name  = "Active Symmetric Key ID"
        value = "keyidentifier"
      },
      {
        name  = "Active Signing Certificate Key ID"
        value = ""
      },
      {
        name  = "JWE Algorithm"
        value = "dir"
      },
      {
        name  = "JWE Content Encryption Algorithm"
        value = "A192CBC-HS384"
      },
      {
        name  = "Active Symmetric Encryption Key ID"
        value = "keyidentifier"
      },
      {
        name  = "Asymmetric Encryption Key"
        value = ""
      },
      {
        name  = "Asymmetric Encryption JWKS URL"
        value = ""
      },
      {
        name  = "Enable Token Revocation"
        value = "false"
      },
      {
        name  = "Include Key ID Header Parameter"
        value = "true"
      },
      {
        name  = "Default JWKS URL Cache Duration"
        value = "720"
      },
      {
        name  = "Include JWE Key ID Header Parameter"
        value = "true"
      },
      {
        name  = "Client ID Claim Name"
        value = "client_id"
      },
      {
        name  = "Scope Claim Name"
        value = "scope"
      },
      {
        name  = "Space Delimit Scope Values"
        value = "true"
      },
      {
        name  = "Authorization Details Claim Name"
        value = "authorization_details"
      },
      {
        name  = "Issuer Claim Value"
        value = ""
      },
      {
        name  = "Audience Claim Value"
        value = ""
      },
      {
        name  = "JWT ID Claim Length"
        value = "22"
      },
      {
        name  = "Access Grant GUID Claim Name"
        value = ""
      },
      {
        name  = "JWKS Endpoint Path"
        value = ""
      },
      {
        name  = "JWKS Endpoint Cache Duration"
        value = "720"
      },
      {
        name  = "Expand Scope Groups"
        value = "false"
      },
      {
        name  = "Type Header Value"
        value = ""
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name         = "contract"
        multi_valued = false
      }
    ]
  }
  selection_settings = {
    resource_uris = []
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

resource "pingfederate_oauth_access_token_mapping" "oauthAccessTokenMappingExample" {
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.jwt_example.manager_id
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
