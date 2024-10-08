---
page_title: "pingfederate_authentication_policy_contract Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages an authentication policy contract.
---

# pingfederate_authentication_policy_contract (Resource)

Manages an authentication policy contract.

## Example Usage

```terraform
resource "pingfederate_authentication_policy_contract" "example" {
  name = "User"
  extended_attributes = [
    { name = "email" },
    { name = "given_name" },
    { name = "family_name" }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The Authentication Policy contract name. Name is unique.

### Optional

- `contract_id` (String) The persistent, unique ID for the authentication policy contract. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.
- `extended_attributes` (Attributes Set) A list of additional attributes as needed. (see [below for nested schema](#nestedatt--extended_attributes))

### Read-Only

- `core_attributes` (Attributes List) A list of read-only assertion attributes (for example, subject) that are automatically populated by PingFederate. (see [below for nested schema](#nestedatt--core_attributes))
- `id` (String) The ID of this resource.

<a id="nestedatt--extended_attributes"></a>
### Nested Schema for `extended_attributes`

Required:

- `name` (String) The name of this attribute.


<a id="nestedatt--core_attributes"></a>
### Nested Schema for `core_attributes`

Required:

- `name` (String) The name of this attribute.

## Import

Import is supported using the following syntax:

~> "contract_id" should be the id of the Authentication Policy Contract to be imported

```shell
terraform import pingfederate_authentication_policy_contract.authenticationPolicyContract contract_id
```