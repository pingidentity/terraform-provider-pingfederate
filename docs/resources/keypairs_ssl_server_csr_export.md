---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_keypairs_ssl_server_csr_export Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to export CSRs for SSL server key pairs.
---

# pingfederate_keypairs_ssl_server_csr_export (Resource)

Resource to export CSRs for SSL server key pairs.

## Example Usage

```terraform
resource "pingfederate_keypairs_ssl_server_key" "sslServerKey" {
  file_data = filebase64("./assets/sslserverkey.p12")
  password  = var.ssl_server_key_password
  format    = "PKCS12"
}

// Example of using the time provider to control regular export of CSR
resource "time_rotating" "csr_export" {
  rotation_days = 30
}

resource "pingfederate_keypairs_ssl_server_csr_export" "example" {
  keypair_id = pingfederate_keypairs_ssl_server_key.sslServerKey.id

  export_trigger_values = {
    "export_rfc3339" : time_rotating.csr_export.rotation_rfc3339,
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `keypair_id` (String) The id of the key pair.

### Optional

- `export_trigger_values` (Map of String) A meta-argument map of values that, if any values are changed, will force export of a new CSR. Adding values to and removing values from the map will not trigger an export. This parameter can be used to control time-based exports using Terraform.

### Read-Only

- `exported_csr` (String) The exported PEM-encoded certificate signing request.
- `id` (String) The ID of this resource.