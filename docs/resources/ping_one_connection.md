---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_ping_one_connection Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages Ping One Connection
---

# pingfederate_ping_one_connection (Resource)

Manages Ping One Connection

## Example Usage

```terraform
resource "pingfederate_ping_one_connection" "pingOneConnectionExample" {
  name = "pingOneConnection"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the PingOne Connection

### Optional

- `active` (Boolean) Whether the PingOne Connection is active. Defaults to true.
- `connection_id` (String) The persistent, unique ID of the connection.
- `credential` (String, Sensitive) The credential for the PingOne connection.
- `description` (String) The description of the PingOne Connection

### Read-Only

- `creation_date` (String) The creation date of the PingOne connection. This field is read only.
- `credential_id` (String) The ID of the PingOne credential. This field is read only.
- `environment_id` (String) The ID of the environment of the PingOne credential. This field is read only.
- `id` (String) The ID of this resource.
- `organization_name` (String) The name of the organization associated with this PingOne connection. This field is read only.
- `ping_one_authentication_api_endpoint` (String) The PingOne Authentication API endpoint. This field is read only.
- `ping_one_connection_id` (String) The ID of the PingOne connection. This field is read only.
- `ping_one_management_api_endpoint` (String) The PingOne Management API endpoint. This field is read only.
- `region` (String) The region of the PingOne connection. This field is read only.

## Import

Import is supported using the following syntax:

```shell
# "connectionId" should be the id of the PingOne Connection to be imported
terraform import pingfederate_ping_one_connection.pingOneConnection connectionId
```