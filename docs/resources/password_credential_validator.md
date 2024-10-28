---
page_title: "pingfederate_password_credential_validator Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages a password credential validator plugin instance.
---

# pingfederate_password_credential_validator (Resource)

Manages a password credential validator plugin instance.

## Example Usage - LDAP Username Password Credential Validator

```terraform
resource "pingfederate_data_store" "pingDirectoryLdapDataStore" {
  ldap_data_store = {
    name      = "PingDirectory LDAP Data Store"
    ldap_type = "PING_DIRECTORY"

    user_dn  = var.pingdirectory_bind_dn
    password = var.pingdirectory_bind_dn_password

    use_ssl = true

    hostnames = [
      "pingdirectory:636"
    ]
  }
}

resource "pingfederate_password_credential_validator" "ldapUsernamePasswordCredentialValidatorExample" {
  validator_id = "ldapUnPwPCV"
  name         = "LDAP Username Password Credential Validator"

  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.LDAPUsernamePasswordCredentialValidator"
  }

  attribute_contract = {}

  configuration = {
    tables = [
      {
        name = "Authentication Error Overrides"
        rows = []
      }
    ],
    fields = [
      {
        name  = "LDAP Datastore"
        value = pingfederate_data_store.pingDirectoryLdapDataStore.id
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
        name  = "Trim Username Spaces For Search"
        value = "true"
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
```

## Example Usage - PingID Password Credential Validator

```terraform
resource "pingfederate_password_credential_validator" "pingIdPasswordCredentialValidatorExample" {
  validator_id = "pingIdPCV"
  name         = "PingID Password Credential Validator"

  plugin_descriptor_ref = {
    id = "com.pingidentity.plugins.pcvs.pingid.PingIdPCV"
  }

  attribute_contract = {}

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
              }
            ],
            sensitive_fields = [
              {
                name  = "Client Shared Secret"
                value = var.pcv_shared_secret
              }
            ]
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
    sensitive_fields = [
      {
        # use_base64_key value from PingID Properties File
        name  = "State Encryption Key"
        value = var.pcv_state_encryption_key
      }
    ]
  }
}
```

## Example Usage - PingOne Directory Password Credential Validator

```terraform
resource "pingfederate_pingone_connection" "example" {
  name       = "My PingOne Tenant"
  credential = var.pingone_connection_credential
}

resource "pingfederate_data_store" "pingOneDataStore" {
  custom_data_store = {
    name = format("PingOne Data Store (%s)", var.pingone_environment_name)

    plugin_descriptor_ref = {
      id = "com.pingidentity.plugins.datastore.p14c.PingOneForCustomersDataStore"
    }

    configuration = {
      fields = [
        {
          name  = "PingOne Environment",
          value = format("%s|%s", pingfederate_pingone_connection.example.id, var.pingone_environment_id)
        }
      ]
    }
    mask_attribute_values = false
  }
}

resource "pingfederate_password_credential_validator" "pingOnePasswordCredentialValidatorExample" {
  validator_id = "pingOnePCV"
  name         = "PingOne Directory Password Credential Validator"

  plugin_descriptor_ref = {
    id = "com.pingidentity.plugins.pcvs.p14c.PingOneForCustomersPCV"
  }

  attribute_contract = {}

  configuration = {
    tables = [
      {
        name = "Authentication Error Overrides"
        rows = []
      }
    ],
    fields = [
      {
        name  = "PingOne For Customers Datastore"
        value = pingfederate_data_store.pingOneDataStore.id
      },
      {
        name  = "Case-Sensitive Matching"
        value = "true"
      }
    ]
  }
}
```

## Example Usage - PingOne for Enterprise Password Credential Validator

```terraform
resource "pingfederate_password_credential_validator" "pingOneForEnterpriseDirectoryPasswordCredentialValidatorExample" {
  validator_id = "pingOneForEnterpriseDirectoryPCV"
  name         = "PingOne for Enterprise Directory Password Credential Validator"

  plugin_descriptor_ref = {
    id = "com.pingconnect.alexandria.pingfed.pcv.PingOnePasswordValidator"
  }

  attribute_contract = {}

  configuration = {
    fields = [
      {
        name  = "Client Id"
        value = "ping_federate_client_id"
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
    ]
    sensitive_fields = [
      {
        name  = "Client Secret"
        value = var.pcv_client_secret
      }
    ]
  }
}
```

## Example Usage - Radius Username Password Credential Validator

```terraform
resource "pingfederate_password_credential_validator" "radiusUsernamePasswordCredentialValidatorExample" {
  validator_id = "radiusUnPwPCV"
  name         = "RADIUS Username Password Credential Validator"

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
            ]
            sensitive_fields = [
              {
                name  = "Shared Secret"
                value = var.pcv_shared_secret
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
```

## Example Usage - Simple Username Password Credential Validator

```terraform
resource "pingfederate_password_credential_validator" "simpleUsernamePasswordCredentialValidatorExample" {
  validator_id = "simpleUsernamePCV"
  name         = "Simple Username Password Credential Validator"

  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator"
  }

  attribute_contract = {}

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
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            sensitive_fields = [
              {
                name  = "Password"
                value = var.pcv_password_user1
              },
              {
                name  = "Confirm Password"
                value = var.pcv_password_user1
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
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            sensitive_fields = [
              {
                name  = "Password"
                value = var.pcv_password_user2
              },
              {
                name  = "Confirm Password"
                value = var.pcv_password_user2
              }
            ]
            default_row = false
          }
        ],
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `attribute_contract` (Attributes) The list of attributes that the password credential validator provides. (see [below for nested schema](#nestedatt--attribute_contract))
- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--configuration))
- `name` (String) The plugin instance name. The name can be modified once the instance is created.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. (see [below for nested schema](#nestedatt--plugin_descriptor_ref))
- `validator_id` (String) The ID of the plugin instance. This field is immutable and will trigger a replacement plan if changed. Must be less than 33 characters, contain no spaces, and be alphanumeric.

### Optional

- `parent_ref` (Attributes) The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. (see [below for nested schema](#nestedatt--parent_ref))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--attribute_contract"></a>
### Nested Schema for `attribute_contract`

Optional:

- `extended_attributes` (Attributes Set) A list of additional attributes that can be returned by the password credential validator. The extended attributes are only used if the adapter supports them. (see [below for nested schema](#nestedatt--attribute_contract--extended_attributes))

Read-Only:

- `core_attributes` (Attributes Set) A list of read-only attributes that are automatically populated by the password credential validator descriptor. (see [below for nested schema](#nestedatt--attribute_contract--core_attributes))

<a id="nestedatt--attribute_contract--extended_attributes"></a>
### Nested Schema for `attribute_contract.extended_attributes`

Required:

- `name` (String) The name of this attribute.


<a id="nestedatt--attribute_contract--core_attributes"></a>
### Nested Schema for `attribute_contract.core_attributes`

Read-Only:

- `name` (String) The name of this attribute.



<a id="nestedatt--configuration"></a>
### Nested Schema for `configuration`

Optional:

- `fields` (Attributes Set) List of configuration fields. (see [below for nested schema](#nestedatt--configuration--fields))
- `sensitive_fields` (Attributes Set) List of sensitive configuration fields. (see [below for nested schema](#nestedatt--configuration--sensitive_fields))
- `tables` (Attributes List) List of configuration tables. (see [below for nested schema](#nestedatt--configuration--tables))

Read-Only:

- `fields_all` (Attributes Set) List of configuration fields. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--configuration--fields_all))
- `tables_all` (Attributes List) List of configuration tables. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--configuration--tables_all))

<a id="nestedatt--configuration--fields"></a>
### Nested Schema for `configuration.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--configuration--sensitive_fields"></a>
### Nested Schema for `configuration.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.

Optional:

- `encrypted_value` (String) For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. Either this attribute or `value` must be specified.
- `value` (String, Sensitive) The sensitive value for the configuration field. Either this attribute or `encrypted_value` must be specified`.


<a id="nestedatt--configuration--tables"></a>
### Nested Schema for `configuration.tables`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--configuration--tables--rows))

<a id="nestedatt--configuration--tables--rows"></a>
### Nested Schema for `configuration.tables.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--configuration--tables--rows--fields))
- `sensitive_fields` (Attributes Set) The sensitive configuration fields in the row. (see [below for nested schema](#nestedatt--configuration--tables--rows--sensitive_fields))

<a id="nestedatt--configuration--tables--rows--fields"></a>
### Nested Schema for `configuration.tables.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--configuration--tables--rows--sensitive_fields"></a>
### Nested Schema for `configuration.tables.rows.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.

Optional:

- `encrypted_value` (String) For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. Either this attribute or `value` must be specified.
- `value` (String, Sensitive) The sensitive value for the configuration field. Either this attribute or `encrypted_value` must be specified`.




<a id="nestedatt--configuration--fields_all"></a>
### Nested Schema for `configuration.fields_all`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--configuration--tables_all"></a>
### Nested Schema for `configuration.tables_all`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--configuration--tables_all--rows))

<a id="nestedatt--configuration--tables_all--rows"></a>
### Nested Schema for `configuration.tables_all.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--configuration--tables_all--rows--fields))

<a id="nestedatt--configuration--tables_all--rows--fields"></a>
### Nested Schema for `configuration.tables_all.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.





<a id="nestedatt--plugin_descriptor_ref"></a>
### Nested Schema for `plugin_descriptor_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--parent_ref"></a>
### Nested Schema for `parent_ref`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "passwordCredentialValidatorId" should be the id of the Password Credential Validator to be imported

```shell
terraform import pingfederate_password_credential_validator.passwordCredentialValidator passwordCredentialValidatorId
```