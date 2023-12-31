---
page_title: "Provider: PingFederate"
description: |-
  The PingFederate provider is used to manage the configuration of a PingFederate server.
---

# PingFederate Provider

The PingFederate provider manages the configuration of a PingFederate server.

## PingFederate Version Support

The PingFederate provider supports versions `11.2` through `12.0` of PingFederate.

## Documentation
Detailed documentation on PingFederate can be found in the [online docs](https://docs.pingidentity.com/r/en-us/pingfederate-112/pf_pingfederate_landing_page)
### Example Usage of PingFederate Provider
```terraform
terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.4.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username                            = "administrator"
  password                            = "2FederateM0re"
  https_host                          = "https://localhost:9999"
  insecure_trust_all_tls              = true
  x_bypass_external_validation_header = true
}

resource "pingfederate_administrative_account" "myAdministrativeAccount" {
  username    = "example"
  description = "description"
  password    = "2FederateM0re"
  roles       = ["USER_ADMINISTRATOR"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `ca_certificate_pem_files` (Set of String) Paths to files containing PEM-encoded certificates to be trusted as root CAs when connecting to the PingFederate server over HTTPS. If not set, the host's root CA set will be used. Default value can be set with the `PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES` environment variable, using commas to delimit multiple PEM files if necessary.
- `https_host` (String) URI for PingFederate HTTPS port. Default value can be set with the `PINGFEDERATE_PROVIDER_HTTPS_HOST` environment variable.
- `insecure_trust_all_tls` (Boolean) Set to true to trust any certificate when connecting to the PingFederate server. This is insecure and should not be enabled outside of testing. Default value can be set with the `PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS` environment variable.
- `password` (String, Sensitive) Password for PingFederate Admin user. Default value can be set with the `PINGFEDERATE_PROVIDER_PASSWORD` environment variable.
- `username` (String) Username for PingFederate Admin user. Default value can be set with the `PINGFEDERATE_PROVIDER_USERNAME` environment variable.
- `x_bypass_external_validation_header` (Boolean) Header value in request for PingFederate. The connection test will be bypassed when set to true. Default value is false.
