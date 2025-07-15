package accesstokenmanager

import "fmt"

func AccessTokenManagerTestHCL(id string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager" "%[1]s" {
  access_control_settings = {
    allowed_clients = [
    ]
    restrict_clients = false
  }
  attribute_contract = {
    default_subject_attribute = "OrgName"
    extended_attributes = [
      {
        multi_valued = false
        name         = "OrgName"
      },
      {
        multi_valued = false
        name         = "Username"
      },
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Access Grant GUID Claim Name"
        value = "agid"
      },
      {
        name  = "Active Signing Certificate Key ID"
        value = "k1"
      },
      {
        name  = "Authorization Details Claim Name"
        value = "authorization_details"
      },
      {
        name  = "Client ID Claim Name"
        value = "client_id_name"
      },
      {
        name  = "Default JWKS URL Cache Duration"
        value = "720"
      },
      {
        name  = "Enable Token Revocation"
        value = "false"
      },
      {
        name  = "Expand Scope Groups"
        value = "false"
      },
      {
        name  = "Include Issued At Claim"
        value = "false"
      },
      {
        name  = "Include JWE Key ID Header Parameter"
        value = "true"
      },
      {
        name  = "Include JWE X.509 Thumbprint Header Parameter"
        value = "false"
      },
      {
        name  = "Include Key ID Header Parameter"
        value = "true"
      },
      {
        name  = "Include X.509 Thumbprint Header Parameter"
        value = "false"
      },
      {
        name  = "JWKS Endpoint Cache Duration"
        value = "720"
      },
      {
        name  = "JWS Algorithm"
        value = "RS256"
      },
      {
        name  = "JWT ID Claim Length"
        value = "22"
      },
      {
        name  = "Publish Key ID X.509 URL"
        value = "false"
      },
      {
        name  = "Publish Thumbprint X.509 URL"
        value = "false"
      },
      {
        name  = "Scope Claim Name"
        value = "scope"
      },
      {
        name  = "Space Delimit Scope Values"
        value = "false"
      },
      {
        name  = "Token Lifetime"
        value = "120"
      },
      {
        name  = "Type Header Value"
        value = ""
      },
      {
        name  = "Use Centralized Signing Key"
        value = "false"
      },
    ]
    tables = [
      {
        name = "Symmetric Keys"
        rows = null
      },
      {
        name = "Certificates"
        rows = [
          {
            default_row = false
            fields = [
              {
                name  = "Certificate"
                value = "419x9yg43rlawqwq9v6az997k"
              },
              {
                name  = "Key ID"
                value = "k1"
              },
            ]
          },
        ]
      },
    ]
  }
  manager_id = "%[1]s"
  name       = "%[1]s"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.access.token.management.plugins.JwtBearerAccessTokenManagementPlugin"
  }
  selection_settings = {
    resource_uris = []
  }
  session_validation_settings = {
    check_session_revocation_status = false
    check_valid_authn_session       = false
    include_session_id              = false
    update_authn_session_activity   = false
  }
}`, id)
}
