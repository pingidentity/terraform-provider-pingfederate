# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_password_credential_validators" "simpleUsernamePasswordCredentialValidatorExample" {
  id   = "simpleUsernamePCV"
  name = "simpleUsernamePasswordCredentialValidator"
  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator"
  }
  configuration = {
    tables = [
      {
        name = "Users"
        rows = [
          {
            fields = [
              {
                name  = "Username"
                value = "example"
              },
              {
                name = "Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name = "Confirm Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          },
          {
            fields = [
              {
                name  = "Username"
                value = "example2"
              },
              {
                name = "Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name = "Confirm Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          }
        ],
      }
    ]
  }
}

resource "pingfederate_password_credential_validators" "simpleUsernamePasswordCredentialValidatorWithParentRefExample" {
  depends_on = [pingfederate_password_credential_validators.simpleUsernamePasswordCredentialValidatorExample]
  id         = "simpleUnPCVParentRefExample"
  name       = "simpleUsernamePasswordCredentialValidatorWithParentRefExample"
  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator"
  }
  parent_ref = {
    id = pingfederate_password_credential_validators.simpleUsernamePasswordCredentialValidatorExample.id
  }
  configuration = {
    tables = [
      {
        name = "Users"
        rows = [
          {
            fields = [
              {
                name  = "Username"
                value = "example"
              },
              {
                name = "Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name = "Confirm Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          },
          {
            fields = [
              {
                name  = "Username"
                value = "example2"
              },
              {
                name = "Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name = "Confirm Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          }
        ]
      }
    ]
  }
}

resource "pingfederate_password_credential_validators" "radiusUsernamePasswordCredentialValidatorExample" {
  id   = "radiusUnPwPCV"
  name = "radiusUsernamePasswordCredentialValidator"
  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.RadiusUsernamePasswordCredentialValidator"
  }
  configuration = {
    tables = [
      {
        name = "RADIUS Servers"
        rows = [
          {
            fields = [
              {
                name  = "Hostname"
                value = "localhost"
              },
              {
                name  = "Authentication Port"
                value = "1812"
              },
              {
                name  = "Authentication Protocol"
                value = "PAP"
              },
              {
                name = "Shared Secret"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              }
            ]
            default_row = false
          }
        ]
      }
    ],
    fields = [
      {
        name  = "NAS Identifier"
        value = "PingFederate"
      },
      {
        name  = "Timeout"
        value = "3000"
      },
      {
        name  = "Retry Count"
        value = "3"
      },
      {
        name  = "Allow Challenge Retries after Access-Reject"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "contract"
      }
    ]
  }
}

resource "pingfederate_password_credential_validators" "ldapUsernamePasswordCredentialValidatorExample" {
  id   = "ldapUnPwPCV"
  name = "ldapUsernamePasswordCredentialValidatorExample"
  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.LDAPUsernamePasswordCredentialValidator"
  }
  configuration = {
    tables = [
      {
        name = "Authentication Error Overrides"
        rows = []
      }
    ],
    fields = [
      {
        name = "LDAP Datastore"
        # ID of LDAP Data Store
        value = ""
      },
      {
        name  = "Search Base"
        value = "cn=Users"
      },
      {
        name = "Search Filter"
        # escape $'s
        value = "sAMAccountName=$${username}"
      },
      {
        name  = "Scope of Search"
        value = "Subtree"
      },
      {
        name  = "Case-Sensitive Matching"
        value = "true"
      },
      {
        name  = "Display Name Attribute"
        value = "displayName"
      },
      {
        name  = "Mail Attribute"
        value = "mail"
      },
      {
        name  = "SMS Attribute"
        value = ""
      },
      {
        name  = "PingID Username Attribute"
        value = ""
      },
      {
        name  = "Mail Search Filter"
        value = ""
      },
      {
        name  = "Username Attribute"
        value = ""
      },
      {
        name  = "Trim Username Spaces For Search"
        value = "true"
      },
      {
        name  = "Mail Verified Attribute"
        value = ""
      },
      {
        name  = "Enable PingDirectory Detailed Password Policy Requirement Messaging"
        value = "true"
      },
      {
        name  = "Expect Password Expired Control"
        value = "false"
      }
    ]
  }
}

resource "pingfederate_password_credential_validators" "pingIdPasswordCredentialValidatorExample" {
  id   = "pingIdPCV"
  name = "pingIdPasswordCredentialValidatorExample"
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
                value = ""
              },
              {
                name = "Client Shared Secret"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "2FederateM0re"
              },
              {
                name  = "Label"
                value = ""
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
        name  = "RADIUS Network Policy Server Port"
        value = ""
      },
      {
        name  = "RADIUS Server Authentication Port"
        value = "1812"
      },
      {
        name  = "Domain postfix"
        value = ""
      },
      {
        name = "PingID Properties File"
        # download file from PingID Client Integration menu, paste contents of file in value property
        # escape newlines, for example \n should be \\n
        value = ""
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
        name = "LDAP Data source"
        # LDAP Data Store ID
        value = ""
      },
      {
        name  = "Create Entry For Devices"
        value = "false"
      },
      {
        name  = "Encryption Key For Devices"
        value = ""
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
        name  = "Distinguished Name Pattern"
        value = ""
      },
      {
        name  = "State Attribute"
        value = ""
      },
      {
        name  = "Server Threads"
        value = ""
      },
      {
        name  = "Enable RADIUS Server"
        value = "true"
      },
      {
        name  = "Default Shared Secret"
        value = ""
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
        name  = "Application icon"
        value = ""
      },
      {
        # use_base64_key value from PingID Properties File
        name  = "State Encryption Key"
        value = ""
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
        name  = "PingID username attribute"
        value = ""
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

resource "pingfederate_password_credential_validators" "pingOneForEnterpriseDirectoryPasswordCredentialValidatorExample" {
  id   = "pingOneForEnterpriseDirectoryPCV"
  name = "pingOneForEnterpriseDirectoryPasswordCredentialValidatorExample"
  plugin_descriptor_ref = {
    id = "com.pingconnect.alexandria.pingfed.pcv.PingOnePasswordValidator"
  }
  configuration = {
    fields = [
      {
        name  = "Client Id"
        value = "ping_federate_client_id"
      },
      {
        name = "Client Secret"
        # This value will be stored into your state file and will not detect any configuration changes made in the UI
        value = "2FederateM0re"
      },
      {
        name  = "PingOne URL"
        value = "https://directory-api.pingone.com/api"
      },
      {
        name  = "Authenticate by Subject URL"
        value = "/directory/users/authenticate?by=subject"
      },
      {
        name  = "Reset Password URL"
        value = "/directory/users/password-reset"
      },
      {
        name  = "SCIM User URL"
        value = "/directory/user"
      },
      {
        name  = "Connection Pool Size"
        value = "100"
      },
      {
        name  = "Connection Pool Idle Timeout"
        value = "4000"
      }
    ],
  }
}

resource "pingfederate_password_credential_validators" "pingOnePasswordCredentialValidatorExample" {
  id   = "pingOnePCV"
  name = "pingOnePasswordCredentialValidatorExample"
  plugin_descriptor_ref = {
    id = "com.pingidentity.plugins.pcvs.p14c.PingOneForCustomersPCV"
  }
  configuration = {
    tables = [
      {
        name = "Authentication Error Overrides"
        rows = []
      }
    ],
    fields = [
      {
        name = "PingOne For Customers Datastore"
        # PingOne Data Store ID
        value = ""
      },
      {
        name  = "Case-Sensitive Matching"
        value = "true"
      },
      {
        name  = "Display Name Attribute"
        value = ""
      },
      {
        name  = "PingID Username Attribute"
        value = ""
      },
      {
        name  = "Username Attribute"
        value = ""
      },
      {
        name  = "Mail Verified Attribute"
        value = ""
      },
      {
        name  = "Mail Attribute"
        value = ""
      },
      {
        name  = "SMS Attribute"
        value = ""
      }
    ]
  }
}