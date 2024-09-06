---
page_title: "pingfederate_oauth_out_of_band_auth_plugin Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to create and manage an Out of Band authenticator plugin.
---

# pingfederate_oauth_out_of_band_auth_plugin (Resource)

Resource to create and manage an Out of Band authenticator plugin.

## Example Usage

```terraform
resource "pingfederate_oauth_out_of_band_auth_plugin" "authPlugin" {
  plugin_id = "CIBAAuthenticator"
  configuration = {
    tables = [
      {
        name = "PingOne Template Variables"
        rows = []
      }
    ]
    fields = [
      {
        name  = "PingOne Environment"
        value = var.pingone_environment
      },
      {
        name  = "Application"
        value = var.pingone_application
      },
      {
        name  = "PingOne Authentication Policy"
        value = "Standalone_MFA"
      },
      {
        name  = "Test Username"
        value = "user.0"
      },
      {
        name  = "PingOne Template Name"
        value = "transaction"
      },
      {
        name  = "PingOne Template Variant"
        value = ""
      },
      {
        name  = "Client Context"
        value = "Example"
      },
      {
        name  = "Messages Files"
        value = "pingone-mfa-messages"
      },
      {
        name  = "API Request Timeout"
        value = "5000"
      },
      {
        name  = "Proxy Settings"
        value = "System Defaults"
      }
    ]
  }
  name = "CIBA Authenticator (PingOne MFA)"
  plugin_descriptor_ref = {
    id = "com.pingidentity.oobauth.pingone.mfa.PingOneMfaCibaAuthenticator"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--configuration))
- `name` (String) The plugin instance name. The name can be modified once the instance is created.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. (see [below for nested schema](#nestedatt--plugin_descriptor_ref))
- `plugin_id` (String) The ID of the plugin instance. The ID cannot be modified once the instance is created.

### Optional

- `attribute_contract` (Attributes) A set of attributes exposed by an out of band authenticator plugin instance. (see [below for nested schema](#nestedatt--attribute_contract))
- `parent_ref` (Attributes) The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. (see [below for nested schema](#nestedatt--parent_ref))

### Read-Only

- `id` (String) The ID of this resource.

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


<a id="nestedatt--attribute_contract"></a>
### Nested Schema for `attribute_contract`

Optional:

- `extended_attributes` (Attributes Set) A list of additional attributes that can be returned by the out of band authenticator plugin instance. The extended attributes are only used if the plugin supports them. (see [below for nested schema](#nestedatt--attribute_contract--extended_attributes))

Read-Only:

- `core_attributes` (Attributes Set) A list of out of band authenticator attributes. (see [below for nested schema](#nestedatt--attribute_contract--core_attributes))

<a id="nestedatt--attribute_contract--extended_attributes"></a>
### Nested Schema for `attribute_contract.extended_attributes`

Required:

- `name` (String) The name of this attribute.


<a id="nestedatt--attribute_contract--core_attributes"></a>
### Nested Schema for `attribute_contract.core_attributes`

Required:

- `name` (String) The name of this attribute.



<a id="nestedatt--parent_ref"></a>
### Nested Schema for `parent_ref`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "pluginId" should be the id of the OAuth Out of Band Authenticator Plugin to be imported

```shell
terraform import pingfederate_oauth_out_of_band_auth_plugin.authPlugin pluginId
```