---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_authentication_selector Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages Authentication Selectors
---

# pingfederate_authentication_selector (Resource)

Manages Authentication Selectors

## Example Usage

```terraform
resource "pingfederate_authentication_selector" "samlAuthnContextExample" {
  selector_id = "samlAuthnContextExample"
  name        = "samlAuthnContextExample"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.selectors.saml.SamlAuthnContextAdapterSelector"
  }
  configuration = {
    tables = []
    fields = [
      {
        name  = "Add or Update AuthN Context Attribute"
        value = "true"
      },
      {
        name  = "Override AuthN Context for Flow"
        value = "true"
      },
      {
        name  = "Enable 'No Match' Result Value"
        value = "false"
      },
      {
        name  = "Enable 'Not in Request' Result Value"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "result_value2"
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--configuration))
- `name` (String) The plugin instance name. The name can be modified once the instance is created.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. (see [below for nested schema](#nestedatt--plugin_descriptor_ref))
- `selector_id` (String) The ID of the plugin instance. The ID cannot be modified once the instance is created.

### Optional

- `attribute_contract` (Attributes) The list of attributes that the Authentication Selector provides. (see [below for nested schema](#nestedatt--attribute_contract))

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

- `extended_attributes` (Attributes Set) A set of additional attributes that can be returned by the Authentication Selector. The extended attributes are only used if the Authentication Selector supports them. (see [below for nested schema](#nestedatt--attribute_contract--extended_attributes))

<a id="nestedatt--attribute_contract--extended_attributes"></a>
### Nested Schema for `attribute_contract.extended_attributes`

Required:

- `name` (String) An attribute for the Authentication Selector attribute contract.

## Import

Import is supported using the following syntax:

```shell
# "authenticationSelectorId" should be the id of the Authentication Selector to be imported
terraform import pingfederate_authentication_selector.authenticationSelector authenticationSelectorId
```
