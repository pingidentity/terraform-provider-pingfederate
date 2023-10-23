---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_idp_default_urls Data Source - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages a IdpDefaultUrls.
---

# pingfederate_idp_default_urls (Data Source)

Manages a IdpDefaultUrls.

## Example Usage

```terraform
data "pingfederate_idp_default_urls" "myIdpDefaultUrl" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `confirm_idp_slo` (Boolean) Prompt user to confirm Single Logout (SLO).
- `id` (String) The ID of this resource.
- `idp_error_msg` (String) Provide the error text displayed in a user's browser when an SSO operation fails.
- `idp_slo_success_url` (String) Provide the default URL you would like to send the user to when Single Logout has succeeded.