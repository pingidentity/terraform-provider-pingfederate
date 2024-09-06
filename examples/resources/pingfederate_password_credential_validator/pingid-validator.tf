resource "pingfederate_password_credential_validator" "pingIdPasswordCredentialValidatorExample" {
  validator_id = "pingIdPCV"
  name         = "PingID Password Credential Validator"

  plugin_descriptor_ref = {
    id = "com.pingidentity.plugins.pcvs.pingid.PingIdPCV"
  }

  configuration = {
    tables = [
      {
        name = "RADIUS Clients"
        rows = [
          {
            fields = [
              {
                name = "Client IP"
                # IP of Client
                value = var.pcv_client_ip
              },
              {
                name  = "Client Shared Secret"
                value = var.pcv_shared_secret
              }
            ],
            default_row = false
          }
        ]
      },
      {
        name = "Delegate PCV's"
        rows = [
          {
            fields = [
              {
                name  = "Delegate PCV"
                value = "radiusUnPwPCV"
              }
            ],
            default_row = false
          }
        ]
      },
      {
        name = "Member Of Groups"
        rows = [
          {
            fields = [
              {
                name  = "LDAP Group Attribute"
                value = "memberOf"
              },
              {
                name  = "LDAP Group Name"
                value = "cn=Groups"
              }
            ],
            default_row = false
          }
        ]
      },
      {
        name = "Bypass Member Of Groups"
        rows = [
          {
            fields = [
              {
                name  = "LDAP Group Name For Bypass"
                value = "cn=BypassGroup"
              }
            ],
            default_row = false
          }
        ]
      },
      {
        name = "RADIUS Vendor-Specific attributes"
        rows = []
      },
      {
        name = "Multiple attributes mapping rules"
        rows = []
      },
      {
        name = "User specific groups to Radius Client"
        rows = []
      }
    ],
    fields = [
      {
        name  = "CHECK GROUPS"
        value = "true"
      },
      {
        name  = "CHECK BYPASS GROUPS"
        value = "false"
      },
      {
        name  = "IF THE USER IS NOT ACTIVATED ON PINGID"
        value = "register"
      },
      {
        name  = "FAIL Login if the user is not member of the ldap group"
        value = "true"
      },
      {
        name  = "Enable RADIUS Remote Network Policy Server"
        value = "false"
      },
      {
        name  = "RADIUS Network Policy Server IP"
        value = ""
      },
      {
        name  = "Domain postfix"
        value = ""
      },
      {
        name = "PingID Properties File"
        # download file from PingID Client Integration menu, paste contents of file in value property
        # escape newlines, for example \n should be \\n
        value = var.pcv_pingid_properties_file
      },
      {
        name  = "Authentication During Errors"
        value = "Bypass User"
      },
      {
        name  = "Users Without a Paired Device"
        value = "Block User"
      },
      {
        name  = "LDAP Data source"
        value = "myldapdatastore"
      },
      {
        name  = "Create Entry For Devices"
        value = "false"
      },
      {
        name  = "Search Base"
        value = "cn=Users"
      },
      {
        name  = "Search Filter"
        value = "sAMAccountName=$${username}"
      },
      {
        name  = "Scope of Search"
        value = "Subtree"
      },
      {
        name  = "Enable RADIUS Server"
        value = "true"
      },
      {
        name  = "PingID service ID"
        value = "vpn"
      },
      {
        name  = "Application name"
        value = "VPN"
      },
      {
        # use_base64_key value from PingID Properties File
        name  = "State Encryption Key"
        value = var.pcv_state_encryption_key
      },
      {
        name  = "State Lifetime"
        value = "300"
      },
      {
        name  = "RADIUS client doesn't support challenge"
        value = "false"
      },
      {
        name  = "OTP in password separator"
        value = "Comma"
      },
      {
        name  = "RADIUS client password validation"
        value = "false"
      },
      {
        name  = "PingId Heartbeat Timeout"
        value = "30"
      },
      {
        name  = "Newline Character"
        value = "None"
      }
    ]
  }
}
