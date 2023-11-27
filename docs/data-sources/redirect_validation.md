---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_redirect_validation Data Source - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Describes the Redirect Validation Settings.
---

# pingfederate_redirect_validation (Data Source)

Describes the Redirect Validation Settings.

## Example Usage

```terraform
data "pingfederate_redirect_validation" "myRedirectValidationExample" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `redirect_validation_local_settings` (Attributes) Settings for local redirect validation. (see [below for nested schema](#nestedatt--redirect_validation_local_settings))
- `redirect_validation_partner_settings` (Attributes) Settings for partner redirect validation. (see [below for nested schema](#nestedatt--redirect_validation_partner_settings))

<a id="nestedatt--redirect_validation_local_settings"></a>
### Nested Schema for `redirect_validation_local_settings`

Read-Only:

- `enable_in_error_resource_validation` (Boolean) Enable validation for error resource.
- `enable_target_resource_validation_for_idp_discovery` (Boolean) Enable target resource validation for IdP discovery.
- `enable_target_resource_validation_for_slo` (Boolean) Enable target resource validation for SLO.
- `enable_target_resource_validation_for_sso` (Boolean) Enable target resource validation for SSO.
- `white_list` (Attributes List) List of URLs that are designated as valid target resources. (see [below for nested schema](#nestedatt--redirect_validation_local_settings--white_list))

<a id="nestedatt--redirect_validation_local_settings--white_list"></a>
### Nested Schema for `redirect_validation_local_settings.white_list`

Read-Only:

- `allow_query_and_fragment` (Boolean) Allow any query parameters and fragment in the resource.
- `idp_discovery` (Boolean) Enable this target resource for IdP discovery validation.
- `in_error_resource` (Boolean) Enable this target resource for in error resource validation.
- `require_https` (Boolean) Require HTTPS for accessing this resource.
- `target_resource_slo` (Boolean) Enable this target resource for SLO redirect validation.
- `target_resource_sso` (Boolean) Enable this target resource for SSO redirect validation.
- `valid_domain` (String) Domain of a valid resource.
- `valid_path` (String) Path of a valid resource.



<a id="nestedatt--redirect_validation_partner_settings"></a>
### Nested Schema for `redirect_validation_partner_settings`

Read-Only:

- `enable_wreply_validation_slo` (Boolean) Enable wreply validation for SLO.