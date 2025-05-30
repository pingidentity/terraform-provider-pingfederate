---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_config_store Data Source - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Data source to retrieve a bundle of settings.
---

# pingfederate_config_store (Data Source)

Data source to retrieve a bundle of settings.

## Example Usage

```terraform
data "pingfederate_config_store" "example" {
  bundle = "MyBundle"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bundle` (String) This field represents a configuration file that contains a bundle of settings.

### Read-Only

- `items` (Attributes List) List of configuration settings. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `id` (String) The id of the configuration setting.
- `list_value` (List of String) The list of values for the configuration setting. This is used when the setting has a list of string values.
- `map_value` (Map of String) The map of key/value pairs for the configuration setting. This is used when the setting has a map of string keys and values.
- `string_value` (String) The value of the configuration setting. This is used when the setting has a single string value.
- `type` (String) The type of configuration setting. This could be a single string, list of strings, or map of string keys and values.
