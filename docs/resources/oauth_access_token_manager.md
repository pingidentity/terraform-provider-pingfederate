---
page_title: "pingfederate_oauth_access_token_manager Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages an OAuth access token manager plugin instance.
---

# pingfederate_oauth_access_token_manager (Resource)

Manages an OAuth access token manager plugin instance.

## Example Usage - Internally Managed Reference Tokens

```terraform
resource "pingfederate_oauth_access_token_manager" "internally_managed_example" {
  manager_id = "internallyManagedReferenceOatm"
  name       = "internallyManagedReferenceExample"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = []
    fields = [
      {
        name  = "Token Length"
        value = "56"
      },
      {
        name  = "Token Lifetime"
        value = "240"
      },
      {
        name  = "Lifetime Extension Policy"
        value = "NONE"
      },
      {
        name  = "Maximum Token Lifetime"
        value = ""
      },
      {
        name  = "Lifetime Extension Threshold Percentage"
        value = "30"
      },
      {
        name  = "Mode for Synchronous RPC"
        value = "3"
      },
      {
        name  = "RPC Timeout"
        value = "500"
      },
      {
        name  = "Expand Scope Groups"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    coreAttributes = []
    extended_attributes = [
      {
        name         = "extended_contract"
        multi_valued = true
      }
    ]
  }
  selection_settings = {
    resource_uris = []
  }
  access_control_settings = {
    restrict_clients = false
    allowedClients   = []
  }
  session_validation_settings = {
    check_valid_authn_session       = false
    check_session_revocation_status = false
    update_authn_session_activity   = false
    include_session_id              = false
  }
}
```

## Example Usage - JWT Tokens

```terraform
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
                value = "e1oDxOiC3Jboz3um8hBVmW3JRZNo9z7C0DMm/oj2V1gclQRcgi2gKM2DBj9N05G4"
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `attribute_contract` (Attributes) The list of attributes that will be added to an access token. (see [below for nested schema](#nestedatt--attribute_contract))
- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--configuration))
- `manager_id` (String) The ID of the plugin instance. The ID cannot be modified once the instance is created.
- `name` (String) The plugin instance name. The name can be modified once the instance is created.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. (see [below for nested schema](#nestedatt--plugin_descriptor_ref))

### Optional

- `access_control_settings` (Attributes) Settings which determine which clients may access this token manager. (see [below for nested schema](#nestedatt--access_control_settings))
- `parent_ref` (Attributes) The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides) (see [below for nested schema](#nestedatt--parent_ref))
- `selection_settings` (Attributes) Settings which determine how this token manager can be selected for use by an OAuth request. (see [below for nested schema](#nestedatt--selection_settings))
- `session_validation_settings` (Attributes) Settings which determine how the user session is associated with the access token. (see [below for nested schema](#nestedatt--session_validation_settings))

### Read-Only

- `id` (String) The ID of this resource.
- `sequence_number` (Number) Number added to an access token to identify which Access Token Manager issued the token.

<a id="nestedatt--attribute_contract"></a>
### Nested Schema for `attribute_contract`

Required:

- `extended_attributes` (Attributes Set) A list of additional token attributes that are associated with this access token management plugin instance. (see [below for nested schema](#nestedatt--attribute_contract--extended_attributes))

Optional:

- `default_subject_attribute` (String) Default subject attribute to use for audit logging when validating the access token. Blank value means to use USER_KEY attribute value after grant lookup.

Read-Only:

- `core_attributes` (Attributes Set) A list of core token attributes that are associated with the access token management plugin type. This field is read-only and is ignored on POST/PUT. (see [below for nested schema](#nestedatt--attribute_contract--core_attributes))

<a id="nestedatt--attribute_contract--extended_attributes"></a>
### Nested Schema for `attribute_contract.extended_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `multi_valued` (Boolean) Indicates whether attribute value is always returned as an array.


<a id="nestedatt--attribute_contract--core_attributes"></a>
### Nested Schema for `attribute_contract.core_attributes`

Read-Only:

- `multi_valued` (Boolean) Indicates whether attribute value is always returned as an array.
- `name` (String) The name of this attribute.



<a id="nestedatt--configuration"></a>
### Nested Schema for `configuration`

Optional:

- `fields` (Attributes Set) List of configuration fields. (see [below for nested schema](#nestedatt--configuration--fields))
- `tables` (Attributes Set) List of configuration tables. (see [below for nested schema](#nestedatt--configuration--tables))

Read-Only:

- `fields_all` (Attributes Set) List of configuration fields. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--configuration--fields_all))
- `tables_all` (Attributes Set) List of configuration tables. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--configuration--tables_all))

<a id="nestedatt--configuration--fields"></a>
### Nested Schema for `configuration.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.


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

<a id="nestedatt--configuration--tables--rows--fields"></a>
### Nested Schema for `configuration.tables.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.




<a id="nestedatt--configuration--fields_all"></a>
### Nested Schema for `configuration.fields_all`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.


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
- `value` (String) The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.





<a id="nestedatt--plugin_descriptor_ref"></a>
### Nested Schema for `plugin_descriptor_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--access_control_settings"></a>
### Nested Schema for `access_control_settings`

Optional:

- `allowed_clients` (Attributes List) If 'restrictClients' is true, this field defines the list of OAuth clients that are allowed to access the token manager. (see [below for nested schema](#nestedatt--access_control_settings--allowed_clients))
- `restrict_clients` (Boolean) Determines whether access to this token manager is restricted to specific OAuth clients. If false, the 'allowedClients' field is ignored. The default value is false.

<a id="nestedatt--access_control_settings--allowed_clients"></a>
### Nested Schema for `access_control_settings.allowed_clients`

Required:

- `id` (String) The ID of the resource.



<a id="nestedatt--parent_ref"></a>
### Nested Schema for `parent_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--selection_settings"></a>
### Nested Schema for `selection_settings`

Optional:

- `resource_uris` (List of String) The list of base resource URI's which map to this token manager. A resource URI, specified via the 'aud' parameter, can be used to select a specific token manager for an OAuth request.


<a id="nestedatt--session_validation_settings"></a>
### Nested Schema for `session_validation_settings`

Optional:

- `check_session_revocation_status` (Boolean) Check the session revocation status when validating the access token.
- `check_valid_authn_session` (Boolean) Check for a valid authentication session when validating the access token.
- `include_session_id` (Boolean) Include the session identifier in the access token. Note that if any of the session validation features is enabled, the session identifier will already be included in the access tokens.
- `update_authn_session_activity` (Boolean) Update authentication session activity when validating the access token.

## Import

Import is supported using the following syntax:

```shell
# "oauthAccessTokenManagerId" should be the id of the Access Token Manager to be imported
# After importing this resource, a subsequent terraform apply will be needed if plain-text values are used
terraform import pingfederate_oauth_access_token_manager.oauthAccessTokenManager oauthAccessTokenManagerId
```
