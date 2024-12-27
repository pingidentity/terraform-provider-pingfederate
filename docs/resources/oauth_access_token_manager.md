---
page_title: "pingfederate_oauth_access_token_manager Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to create and manage an OAuth access token manager plugin instance.
---

# pingfederate_oauth_access_token_manager (Resource)

Resource to create and manage an OAuth access token manager plugin instance.

## Example Usage - Internally Managed Reference Tokens

```terraform
resource "pingfederate_oauth_access_token_manager" "internally_managed_example" {
  manager_id = "internallyManagedReferenceOATM"
  name       = "Internally Managed Token Manager"

  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }

  configuration = {
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
    extended_attributes = [
      {
        name         = "givenName"
        multi_valued = false
      },
      {
        name         = "familyName"
        multi_valued = false
      },
      {
        name         = "email"
        multi_valued = false
      },
      {
        name         = "groups"
        multi_valued = true
      }
    ]
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

## Example Usage - JWT Tokens

```terraform
resource "pingfederate_oauth_access_token_manager" "jwt_example" {
  manager_id = "jsonWebTokenOATM"
  name       = "JWT Access Token Manager"

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
                value = "jwtSymmetricKey1"
              },
              {
                name  = "Encoding"
                value = "b64u"
              }
            ]
            sensitive_fields = [
              {
                name  = "Key"
                value = var.jwt_symmetric_key
              },
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
        value = "jwtSymmetricKey1"
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
        value = "jwtSymmetricKey1"
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
        name         = "givenName"
        multi_valued = false
      },
      {
        name         = "familyName"
        multi_valued = false
      },
      {
        name         = "email"
        multi_valued = false
      },
      {
        name         = "groups"
        multi_valued = true
      }
    ]
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
- `manager_id` (String) The ID of the plugin instance. The ID cannot be modified once the instance is created. Must be alphanumeric, contain no spaces, and be less than 33 characters. This field is immutable and will trigger a replacement plan if changed.
- `name` (String) The plugin instance name. The name can be modified once the instance is created.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. This field is immutable and will trigger a replacement plan if changed. (see [below for nested schema](#nestedatt--plugin_descriptor_ref))

### Optional

- `access_control_settings` (Attributes) Settings which determine which clients may access this token manager. (see [below for nested schema](#nestedatt--access_control_settings))
- `parent_ref` (Attributes) The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides) (see [below for nested schema](#nestedatt--parent_ref))
- `selection_settings` (Attributes) Settings which determine how this token manager can be selected for use by an OAuth request. (see [below for nested schema](#nestedatt--selection_settings))
- `session_validation_settings` (Attributes) Settings which determine how the user session is associated with the access token. (see [below for nested schema](#nestedatt--session_validation_settings))
- `token_endpoint_attribute_contract` (Attributes) A set of attributes exposed by an Access Token Manager in a token endpoint response. Supported in PingFederate `12.2.0` and later. (see [below for nested schema](#nestedatt--token_endpoint_attribute_contract))

### Read-Only

- `id` (String) The ID of this resource.
- `sequence_number` (Number) Number added to an access token to identify which Access Token Manager issued the token.

<a id="nestedatt--attribute_contract"></a>
### Nested Schema for `attribute_contract`

Required:

- `extended_attributes` (Attributes Set) A list of additional token attributes that are associated with this access token management plugin instance. (see [below for nested schema](#nestedatt--attribute_contract--extended_attributes))

Optional:

- `default_subject_attribute` (String) Default subject attribute to use for audit logging when validating the access token. Blank value means to use `USER_KEY` attribute value after grant lookup.

Read-Only:

- `core_attributes` (Attributes Set) A list of core token attributes that are associated with the access token management plugin type. This field is read-only. (see [below for nested schema](#nestedatt--attribute_contract--core_attributes))

<a id="nestedatt--attribute_contract--extended_attributes"></a>
### Nested Schema for `attribute_contract.extended_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `multi_valued` (Boolean) Indicates whether attribute value is always returned as an array. The default is `false`.


<a id="nestedatt--attribute_contract--core_attributes"></a>
### Nested Schema for `attribute_contract.core_attributes`

Read-Only:

- `multi_valued` (Boolean) Indicates whether attribute value is always returned as an array.
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

- `id` (String) The ID of the resource. This field is immutable and will trigger a replacement plan if changed.


<a id="nestedatt--access_control_settings"></a>
### Nested Schema for `access_control_settings`

Optional:

- `allowed_clients` (Attributes List) If `restrict_clients` is `true`, this field defines the list of OAuth clients that are allowed to access the token manager. (see [below for nested schema](#nestedatt--access_control_settings--allowed_clients))
- `restrict_clients` (Boolean) Determines whether access to this token manager is restricted to specific OAuth clients. If `false`, the `allowed_clients` field is ignored. The default value is `false`.

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

- `resource_uris` (Set of String) The list of base resource URI's which map to this token manager. A resource URI, specified via the 'aud' parameter, can be used to select a specific token manager for an OAuth request.


<a id="nestedatt--session_validation_settings"></a>
### Nested Schema for `session_validation_settings`

Optional:

- `check_session_revocation_status` (Boolean) Check the session revocation status when validating the access token. The default is `false`.
- `check_valid_authn_session` (Boolean) Check for a valid authentication session when validating the access token. The default is `false`.
- `include_session_id` (Boolean) Include the session identifier in the access token. Note that if any of the session validation features is enabled, the session identifier will already be included in the access tokens. The default is `false`.
- `update_authn_session_activity` (Boolean) Update authentication session activity when validating the access token. The default is `false`.


<a id="nestedatt--token_endpoint_attribute_contract"></a>
### Nested Schema for `token_endpoint_attribute_contract`

Optional:

- `attributes` (Attributes Set) A list of token endpoint response attributes that are associated with this access token management plugin instance. (see [below for nested schema](#nestedatt--token_endpoint_attribute_contract--attributes))

<a id="nestedatt--token_endpoint_attribute_contract--attributes"></a>
### Nested Schema for `token_endpoint_attribute_contract.attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `mapped_scopes` (Set of String) List of scopes that will trigger this attribute to be included in the token endpoint response.
- `multi_valued` (Boolean) Indicates whether attribute value is always returned as an array.

## Import

Import is supported using the following syntax:

~> "oauthAccessTokenManagerId" should be the id of the Access Token Manager to be imported

```shell
terraform import pingfederate_oauth_access_token_manager.oauthAccessTokenManager oauthAccessTokenManagerId
```
