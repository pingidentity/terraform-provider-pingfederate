terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source = "pingidentity/pingfederate"
    }
  }
}
provider "pingfederate" {
  username = "administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:9999"
  insecure_trust_all_tls = true
}

resource "pingfederate_oauth_access_token_managers" "jsonWebTokenOauthAccessTokenManagersExample2" {
  id = "test2"
  name = "test2"
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
                name = "Key ID"
                value = "keyidentifier"
              },
              {
                name = "Key"
                value = "+d5OB5b+I4dqn1Mjp8YE/M/QFWvDX7Nxz3gC8mAEwRLqL67SrHcwRyMtGvZKxvIn"
              },
              {
                name = "Encoding"
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
        name = "Token Lifetime"
        value = "120"
      },
      {
        name = "Use Centralized Signing Key"
        value = "false"
      },
      {
        name = "JWS Algorithm"
        value = ""
      },
      {
        name = "Active Symmetric Key ID"
        value = "keyidentifier"
      },
      {
        name = "Active Signing Certificate Key ID"
        value = ""
      },
      {
        name = "JWE Algorithm"
        value = "dir"
      },
      {
        name = "JWE Content Encryption Algorithm"
        value = "A192CBC-HS384"
      },
      {
        name = "Active Symmetric Encryption Key ID"
        value = "keyidentifier"
      },
      {
        name = "Asymmetric Encryption Key"
        value = ""
      },
      {
        name = "Asymmetric Encryption JWKS URL"
        value = ""
      },
      {
        name = "Enable Token Revocation"
        value = "false"
      },
      {
        name = "Include Key ID Header Parameter"
        value = "true"
      },
      {
        name = "Default JWKS URL Cache Duration"
        value = "720"
      },
      {
        name = "Include JWE Key ID Header Parameter"
        value = "true"
      },
      {
        name = "Client ID Claim Name"
        value = "client_id"
      },
      {
        name = "Scope Claim Name"
        value = "scope"
      },
      {
        name = "Space Delimit Scope Values"
        value = "true"
      },
      {
        name = "Authorization Details Claim Name"
        value = "authorization_details"
      },
      {
        name = "Issuer Claim Value"
        value = ""
      },
      {
        name = "Audience Claim Value"
        value = ""
      },
      {
        name = "JWT ID Claim Length"
        value = "22"
      },
      {
        name = "Access Grant GUID Claim Name"
        value = ""
      },
      {
        name = "JWKS Endpoint Path"
        value = ""
      },
      {
        name = "JWKS Endpoint Cache Duration"
        value = "720"
      },
      {
        name = "Expand Scope Groups"
        value = "false"
      },
      {
        name = "Type Header Value"
        value = ""
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "contract"
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
    check_valid_authn_session = false
    check_session_revocation_status = false
    update_authn_session_activity = false
    include_session_id = false
  }
}

resource "pingfederate_oauth_access_token_managers" "internallyManagedReferenceOauthAccessTokenManagersExample" {
  id = "test4"
  name = "test4"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = []
    fields = [
      {
        name = "Token Length"
        value = "28"
      },
      {
        name = "Token Lifetime"
        value = "120"
      },
      {
        name = "Lifetime Extension Policy"
        value = "NONE"
      },
      {
        name = "Maximum Token Lifetime"
        value = ""
      },
      {
        name = "Lifetime Extension Threshold Percentage"
        value = "30"
      },
      {
        name = "Mode for Synchronous RPC"
        value = "3"
      },
      {
        name = "RPC Timeout"
        value = "500"
      },
      {
        name = "Expand Scope Groups"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    coreAttributes = []
    extended_attributes = [
      {
        name = "extended_contract"
        multi_valued = true
      }
    ]
  }
  selection_settings = {
    resource_uris = []
  }
  access_control_settings = {
    restrict_clients = false
    allowedClients = []
  }
  session_validation_settings = {
    check_valid_authn_session = false
    check_session_revocation_status = false
    update_authn_session_activity = false
    include_session_id = false
  }
}