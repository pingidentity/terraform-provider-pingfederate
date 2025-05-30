---
page_title: "pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to create and manage certificates for WS-Trust STS settings.
---

# pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate (Resource)

Resource to create and manage certificates for WS-Trust STS settings.

## Example Usage

```terraform
resource "pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate" "issuerCert" {
  file_data = filebase64("path/to/my/issuercert.pem")
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `file_data` (String) The certificate data in PEM format. New line characters should be omitted or encoded in this value. This field is immutable and will trigger a replacement plan if changed.

### Optional

- `certificate_id` (String) The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.
- `crypto_provider` (String) Cryptographic Provider. This is only applicable if Hybrid HSM mode is `true`. Options are `LOCAL` or `HSM`. This field is immutable and will trigger a replacement plan if changed.

### Read-Only

- `active` (Boolean) Indicates whether this an active certificate or not.
- `expires` (String) The end date up until which the item is valid, in ISO 8601 format (UTC).
- `formatted_file_data` (String) The certificate data in PEM format, formatted by PingFederate. This attribute is read-only.
- `id` (String) The ID of this resource.
- `issuer_dn` (String) The issuer's distinguished name.
- `key_algorithm` (String) The public key algorithm.
- `key_size` (Number) The public key size.
- `serial_number` (String) The serial number assigned by the CA.
- `sha1_fingerprint` (String) SHA-1 fingerprint in Hex encoding.
- `sha256_fingerprint` (String) SHA-256 fingerprint in Hex encoding.
- `signature_algorithm` (String) The signature algorithm.
- `status` (String) Status of the item.
- `subject_alternative_names` (List of String) The subject alternative names (SAN).
- `subject_dn` (String) The subject's distinguished name.
- `valid_from` (String) The start date from which the item is valid, in ISO 8601 format (UTC).
- `version` (Number) The X.509 version to which the item conforms.

## Import

Import is supported using the following syntax:

~> "certificateId" should be the id of the issuer certificate to be imported.

```shell
terraform import pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate.issuerCert certificateId
```