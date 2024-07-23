---
page_title: "pingfederate_kerberos_realm_settings Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to manage the Kerberos realm settings.
---

# pingfederate_kerberos_realm_settings (Resource)

Resource to manage the Kerberos realm settings.

## Example Usage

```terraform
resource "pingfederate_kerberos_realm_settings" "kerberosRealmSettings" {
  debug_log_output              = false
  force_tcp                     = false
  kdc_retries                   = 3
  kdc_timeout                   = 3
  key_set_retention_period_mins = 610
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `kdc_retries` (Number) Reference to the default Key Distribution Center Retries.
- `kdc_timeout` (Number) Reference to the default Key Distribution Center Timeout (in seconds).

### Optional

- `debug_log_output` (Boolean) Reference to the default logging. Default value is `false`
- `force_tcp` (Boolean) Reference to the default security. Default value is `false`
- `key_set_retention_period_mins` (Number) The key set retention period in minutes. When 'retain_previous_keys_on_password_change' is set to `true` for a realm, this setting determines how long keys will be retained after a password change occurs. Default value is `610`

## Import

Import is supported using the following syntax:

~> This resource is singleton, so the value of "id" doesn't matter - it is just a placeholder, and required by Terraform

```shell
terraform import pingfederate_kerberos_realm_settings.kerberosRealmSettings id
```