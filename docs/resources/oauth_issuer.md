---
page_title: "pingfederate_oauth_issuer Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to create and manage a virtual OAuth issuer.
---

# pingfederate_oauth_issuer (Resource)

Resource to create and manage a virtual OAuth issuer.

## Example Usage

```terraform
resource "pingfederate_oauth_issuer" "oauthIssuer" {
  issuer_id   = "oauthIssuer"
  description = "example description"
  host        = "example"
  name        = "example"
  path        = "/example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host` (String) The hostname of this virtual issuer.
- `name` (String) The name of this virtual issuer with a unique value.

### Optional

- `description` (String) The description of this virtual issuer.
- `issuer_id` (String) The persistent, unique ID for the virtual issuer. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.
- `path` (String) The path of this virtual issuer. Path must start with a `/`, but cannot end with `/`.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

~> "oauthIssuerId" should be the id of the OAuth Issuer to be imported

```shell
terraform import pingfederate_oauth_issuer.oauthIssuer oauthIssuerId
```