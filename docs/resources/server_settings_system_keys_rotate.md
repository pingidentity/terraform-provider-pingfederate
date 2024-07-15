---
page_title: "pingfederate_server_settings_system_keys_rotate Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  
---

# pingfederate_server_settings_system_keys_rotate (Resource)



## Example Usage

```terraform
// Example of using the time provider to control regular rotaion of system keys
resource "time_rotating" "system_key_rotation" {
  rotation_days = 30
}

resource "pingfederate_server_settings_system_keys_rotate" "systemKeysRotate" {
  rotation_trigger_values = {
    "rotation_rfc3339" : time_rotating.system_key_rotation.rotation_rfc3339,
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `rotation_trigger_values` (Map of String) A meta-argument map of values that, if any values are changed, will force rotation of the system keys. Adding values to and removing values from the map will not trigger a key rotation. This parameter can be used to control time-based rotation using Terraform.

### Read-Only

- `current` (Attributes) (see [below for nested schema](#nestedatt--current))
- `pending` (Attributes) (see [below for nested schema](#nestedatt--pending))
- `previous` (Attributes) (see [below for nested schema](#nestedatt--previous))

<a id="nestedatt--current"></a>
### Nested Schema for `current`

Read-Only:

- `creation_date` (String)
- `encrypted_key_data` (String, Sensitive)


<a id="nestedatt--pending"></a>
### Nested Schema for `pending`

Read-Only:

- `creation_date` (String)
- `encrypted_key_data` (String, Sensitive)


<a id="nestedatt--previous"></a>
### Nested Schema for `previous`

Read-Only:

- `creation_date` (String)
- `encrypted_key_data` (String, Sensitive)

## Import

Import is supported using the following syntax:

~> This resource is singleton, so the value of "id" doesn't matter - it is just a placeholder, and required by Terraform

```shell
terraform import pingfederate_server_settings_system_keys_rotate.systemKeysRotate id
```