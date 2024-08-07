---
page_title: "pingfederate_default_urls Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to manage IdP and SP default URL settings.
---

# pingfederate_default_urls (Resource)

Resource to manage IdP and SP default URL settings.

## Example Usage

```terraform
resource "pingfederate_default_urls" "defaultUrlsExample" {
  confirm_sp_slo      = true
  sp_slo_success_url  = "https://example.com/slo_success_url"
  sp_sso_success_url  = "https://example.com/sso_success_url"
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "https://example.com/idp_slo_success_url"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `idp_error_msg` (String) IdP setting for the error text displayed in a user's browser when an SSO operation fails.

### Optional

- `confirm_idp_slo` (Boolean) IdP setting to prompt user to confirm Single Logout (SLO). The default value is `false`.
- `confirm_sp_slo` (Boolean) SP setting to prompt user to confirm Single Logout (SLO). The default is `false`.
- `idp_slo_success_url` (String) Idp setting for the default URL you would like to send the user to when Single Logout has succeeded.
- `sp_slo_success_url` (String) SP setting for the default URL you would like to send the user to when Single Logout (SLO) has succeeded.
- `sp_sso_success_url` (String) SP setting for the default URL you would like to send the user to when Single Sign On (SSO) has succeeded.

## Import

Import is supported using the following syntax:

~> This resource is singleton, so the value of "id" doesn't matter - it is just a placeholder, and required by Terraform

```shell
terraform import pingfederate_default_urls.defaultUrlsExample id
```
