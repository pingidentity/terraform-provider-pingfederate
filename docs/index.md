---
page_title: "Provider: PingFederate"
description: |-
  The PingFederate provider is used to manage the configuration of a PingFederate server.
---

# PingFederate Provider

The PingFederate provider manages the configuration of a PingFederate server via the management API.

## PingFederate Server Support

This PingFederate Terraform provider version supports versions `11.3` through `12.3` of PingFederate.

Further information on PingFederate server version support and released version compatibility can be found in the [PingFederate Server Support](guides/server-support) guide.

## Getting Started

### Configure PingFederate for Terraform access

For detailed instructions on how to prepare PingFederate for Terraform access, see the [PingFederate getting started guide](https://terraform.pingidentity.com/getting-started/pingfederate/) at [terraform.pingidentity.com](https://terraform.pingidentity.com).

### PingFederate Server Documentation

Detailed documentation on the PingFederate server can be found in the [online documentation](https://docs.pingidentity.com/r/en-us/pingfederate-121/pf_pingfederate_landing_page)

## Provider Authentication

### Simple Example using basic authentication with a resource
```terraform
terraform {
  required_version = ">=1.4"
  required_providers {
    pingfederate = {
      version = "~> 1.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username                            = "administrator"
  password                            = "2FederateM0re"
  https_host                          = "https://localhost:9999"
  admin_api_path                      = "/pf-admin-api/v1"
  insecure_trust_all_tls              = true
  x_bypass_external_validation_header = true
  product_version                     = "12.3"
}

resource "pingfederate_administrative_account" "administrativeAccount" {
  username    = "example"
  description = "description"
  password    = "2FederateM0re"
  roles       = ["USER_ADMINISTRATOR"]
}
```

### When using basic authentication, `username` and `password` must be defined in the `provider` block
```terraform
provider "pingfederate" {
  username        = "administrator"
  password        = "2FederateM0re"
  https_host      = "https://localhost:9999"
  admin_api_path  = "/pf-admin-api/v1"
  product_version = "12.3"
}
```

### When using OAuth2 Client Credentials flow authentication, `client_id`, `client_secret`, and `token_url` are required, while `scopes` is optional in the `provider` block
```terraform
provider "pingfederate" {
  client_id       = "clientid"
  client_secret   = "clientsecret"
  scopes          = ["scope"]
  token_url       = "https://localhost:9031/as/token.oauth2"
  https_host      = "https://localhost:9999"
  admin_api_path  = "/pf-admin-api/v1"
  product_version = "12.3"
}
```

### When using Access Token authentication, `access_token` is required in the `provider` block
```terraform
provider "pingfederate" {
  access_token    = "accesstokenvaluefromclient"
  https_host      = "https://localhost:9999"
  admin_api_path  = "/pf-admin-api/v1"
  product_version = "12.3"
}
```

## Custom User Agent information

The PingFederate provider allows custom information to be appended to the default user agent string (that includes Terraform provider version information) by setting the `PINGFEDERATE_TF_APPEND_USER_AGENT` environment variable.  This can be useful when troubleshooting issues with Ping Identity Support, or adding context to HTTP requests.

```shell
export PINGFEDERATE_TF_APPEND_USER_AGENT="Jenkins/2.426.2"
```

## Schema

### Required

- `https_host` (String) URI for PingFederate HTTPS port. Default value can be set with the `PINGFEDERATE_PROVIDER_HTTPS_HOST` environment variable.
- `product_version` (String) Version of the PingFederate server being configured. Default value can be set with the `PINGFEDERATE_PROVIDER_PRODUCT_VERSION` environment variable.

### Optional

- `access_token` (String, Sensitive) Access token for PingFederate Admin API. Cannot be used in conjunction with username and password, or oauth. Default value can be set with the `PINGFEDERATE_PROVIDER_ACCESS_TOKEN` environment variable.
- `admin_api_path` (String) Path for PingFederate Admin API. Default value can be set with the `PINGFEDERATE_PROVIDER_ADMIN_API_PATH` environment variable. If no value is supplied, the value used will be `/pf-admin-api/v1`.
- `ca_certificate_pem_files` (Set of String) Paths to files containing PEM-encoded certificates to be trusted as root CAs when connecting to the PingFederate server over HTTPS. If not set, the host's root CA set will be used. Default value can be set with the `PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES` environment variable, using commas to delimit multiple PEM files if necessary.
- `client_id` (String) OAuth client ID for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID` environment variable.
- `client_secret` (String, Sensitive) OAuth client secret for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET` environment variable.
- `insecure_trust_all_tls` (Boolean) Set to true to trust any certificate when connecting to the PingFederate server. This is insecure and should not be enabled outside of testing. Default value can be set with the `PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS` environment variable.
- `password` (String, Sensitive) Password for PingFederate Admin user. Must only be set with username. Cannot be used in conjunction with access_token, or oauth.  Default value can be set with the `PINGFEDERATE_PROVIDER_PASSWORD` environment variable.
- `scopes` (List of String) OAuth scopes for access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_SCOPES` environment variable.
- `token_url` (String) OAuth token URL for requesting access token. Default value can be set with the `PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL` environment variable.
- `username` (String) Username for PingFederate Admin user. Must only be set with password. Cannot be used in conjunction with access_token, or oauth. Default value can be set with the `PINGFEDERATE_PROVIDER_USERNAME` environment variable.
- `x_bypass_external_validation_header` (Boolean) Header value in request for PingFederate. When set to `true`, connectivity checks for resources such as `pingfederate_data_store` will be skipped. Default value can be set with the `PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER` environment variable.
