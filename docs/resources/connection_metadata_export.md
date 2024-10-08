---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_connection_metadata_export Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to export a connection's SAML metadata that can be given to a partner.
---

# pingfederate_connection_metadata_export (Resource)

Resource to export a connection's SAML metadata that can be given to a partner.

## Example Usage

```terraform
resource "pingfederate_connection_metadata_export" "metadataExport" {
  connection_type = "SP"
  connection_id   = pingfederate_idp_sp_connection.example_saml.connection_id
  signing_settings = {
    signing_key_pair_ref = {
      id = pingfederate_keypairs_signing_key.rsa_saml_signing_1.id
    }
    algorithm = "SHA256withRSA"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_id` (String) The ID of the connection to export. This field is immutable and will trigger a replacement plan if changed.
- `connection_type` (String) The type of connection to export. Options are `IDP` or `SP`. This field is immutable and will trigger a replacement plan if changed.

### Optional

- `signing_settings` (Attributes) The signing settings to sign the metadata with. If `null`, the metadata will not be signed. This field is immutable and will trigger a replacement plan if changed. (see [below for nested schema](#nestedatt--signing_settings))
- `use_secondary_port_for_soap` (Boolean) If PingFederate's secondary SSL port is configured and you want to use it for the SOAP channel, set to `true`. If client-certificate authentication is configured for the SOAP channel, the secondary port is required and this must be set to `true`. This field is immutable and will trigger a replacement plan if changed.
- `virtual_host_name` (String) The virtual host name to be used as the base url. This field is immutable and will trigger a replacement plan if changed.
- `virtual_server_id` (String) The virtual server ID to export the metadata with. If `null`, the connection's default will be used. This field is immutable and will trigger a replacement plan if changed.

### Read-Only

- `exported_connection_metadata` (String) The exported SAML metadata.

<a id="nestedatt--signing_settings"></a>
### Nested Schema for `signing_settings`

Required:

- `signing_key_pair_ref` (Attributes) The ID of the key pair used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console. This field is immutable and will trigger a replacement plan if changed. (see [below for nested schema](#nestedatt--signing_settings--signing_key_pair_ref))

Optional:

- `algorithm` (String) The algorithm used to sign messages sent to this partner. The default is `SHA1withDSA` for DSA certs, `SHA256withRSA` for RSA certs, and `SHA256withECDSA` for EC certs. For RSA certs, `SHA1withRSA`, `SHA384withRSA`, `SHA512withRSA`, `SHA256withRSAandMGF1`, `SHA384withRSAandMGF1` and `SHA512withRSAandMGF1` are also supported. For EC certs, `SHA384withECDSA` and `SHA512withECDSA` are also supported. If the connection is WS-Federation with JWT token type, then the possible values are `RSA SHA256`, `RSA SHA384`, `RSA SHA512`, `RSASSA-PSS SHA256`, `RSASSA-PSS SHA384`, `RSASSA-PSS SHA512`, `ECDSA SHA256`, `ECDSA SHA384`, `ECDSA SHA512`. This field is immutable and will trigger a replacement plan if changed.
- `include_cert_in_signature` (Boolean) Determines whether the signing certificate is included in the signature element. This field is immutable and will trigger a replacement plan if changed.
- `include_raw_key_in_signature` (Boolean) Determines whether the element with the raw public key is included in the signature element. This field is immutable and will trigger a replacement plan if changed.

<a id="nestedatt--signing_settings--signing_key_pair_ref"></a>
### Nested Schema for `signing_settings.signing_key_pair_ref`

Required:

- `id` (String) The ID of the resource.
