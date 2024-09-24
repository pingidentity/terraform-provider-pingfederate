---
layout: ""
page_title: "Version 1 Upgrade Guide (from version 0)"
description: |-
  Version 1.0.0 of the PingFederate Terraform provider is a major release that introduces breaking changes to existing HCL.  This guide describes the changes that are required to upgrade v0.* PingFederate Terraform provider releases to v1.0.0 onwards.
---

# PingFederate Terraform Provider Version 1 Upgrade Guide (from version 0)

Version `1.0` of the PingFederate Terraform provider is a major release that introduces breaking changes to existing HCL. This guide describes the changes that are required to upgrade `v0.*` PingFederate Terraform provider releases to `v1.*`.

## Why have schemas changed?

As part of ensuring the ongoing maintainability of the Terraform provider integration and to solve functional issues, some resource schemas have changed going to version `1` from version `0`.

The schemas may have changed in the following ways in this release:

* Removal of previously deprecated fields
* Removal of previously deprecated resources/data sources

The following sections detail the rationale for the above changes, and whether the changes are routine for a major version upgrade or one off changes that aren't expected in future major version changes.

### Removal of previously deprecated fields

Removal of deprecated fields are expected on each major release going forward.  Ping maintains a deprecation and release strategy according to [Terraform provider creation best practices](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning) and support of PingFederate product versions aligns with the [Ping Identity End of Life Policy](https://www.pingidentity.com/en/legal/end-of-life-policy.html).

### Removal of previously deprecated resources/data sources

Removal of deprecated resources / data sources are expected on each major release going forward.  Ping maintains a deprecation and release strategy according to [Terraform provider creation best practices](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning) and support of PingFederate product versions aligns with the [Ping Identity End of Life Policy](https://www.pingidentity.com/en/legal/end-of-life-policy.html).

### Removal of support for previous PingFederate versions

The support of PingFederate product versions aligns with the [Ping Identity End of Life Policy](https://www.pingidentity.com/en/legal/end-of-life-policy.html).  In this major release, supported PingFederate versions remain unchanged, versions `11.2` through `12.1` of PingFederate (supported in the `v0.*` release) are still supported.

## Provider Configuration Changes

### Major Version Change

Customers can keep operating existing `v0.*` releases until ready to upgrade to `v1.*`.  Remaining on the latest `v0.*` release can be achieved using the following syntax:

```terraform
terraform {
  required_providers {
    pingfederate = {
      source  = "pingidentity/pingfederate"
      version = "~> 0.16"
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
  product_version                     = "12.1"
}
```

It is highly recommended to go through the guide and make updates to each impacted resource before changing the version, as there are backward-incompatible changes.  Once ready to upgrade, the version can be incremented as follows:

```terraform
terraform {
  required_providers {
    pingfederate = {
      source  = "pingidentity/pingfederate"
      version = "~> 1.0"
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
  product_version                     = "12.1"
}
```

Ping recommends using [Provider version control](https://terraform.pingidentity.com/best-practices/#use-provider-version-control), detailed in the [Terraform best practices guide](https://terraform.pingidentity.com/best-practices/).

## Resource: pingfederate_administrative_account

### `encrypted_password` computed attribute removed

The `encrypted_password` computed attribute has been removed as it is no longer used.

## Resource: pingfederate_authentication_api_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_authentication_policies_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_extended_properties

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_idp_default_urls

This resource has been previously deprecated and has now been removed. Use the `pingfederate_default_urls` resource going forward.

## Resource: pingfederate_idp_sp_connection

### `type` computed attribute removed

The unnecessary `type` computed attribute has been removed.

### `credentials.inbound_back_channel_auth.type` computed attribute removed

The unnecessary `credentials.inbound_back_channel_auth.type` computed attribute has been removed.

### `credentials.inbound_back_channel_auth.http_basic_credentials.encrypted_password` optional parameter removed

The unnecessary `credentials.inbound_back_channel_auth.http_basic_credentials.encrypted_password` optional parameter has been removed.

### `credentials.outbound_back_channel_auth.type` computed attribute removed

The unnecessary `credentials.outbound_back_channel_auth.type` computed attribute has been removed.

### `credentials.outbound_back_channel_auth.http_basic_credentials.encrypted_password` optional parameter removed

The unnecessary `credentials.outbound_back_channel_auth.http_basic_credentials.encrypted_password` optional parameter has been removed.

## Resource: pingfederate_incoming_proxy_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_key_pair_signing_import

This resource has been previously deprecated and has now been removed. Use the `pingfederate_keypairs_signing_key` resource going forward.

## Resource: pingfederate_key_pair_ssl_server_import

This resource has been previously deprecated and has now been removed. Use the `pingfederate_keypairs_ssl_server_key` resource going forward.

## Resource: pingfederate_license_agreement

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_license

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_local_identity_identity_profile

This resource has been previously deprecated and has now been removed. Use the `pingfederate_local_identity_profile` resource going forward.

## Resource: pingfederate_notification_publisher_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_notification_publishers_settings

This resource has been previously deprecated and has now been removed. Use the `pingfederate_notification_publisher_settings` resource going forward.

## Resource: pingfederate_oauth_auth_server_settings_scopes_common_scope

This resource has been previously deprecated and has now been removed. Use the `pingfederate_oauth_server_settings` resource going forward.

## Resource: pingfederate_oauth_auth_server_settings_scopes_exclusive_scope

This resource has been previously deprecated and has now been removed. Use the `pingfederate_oauth_server_settings` resource going forward.

## Resource: pingfederate_oauth_auth_server_settings

This resource has been previously deprecated and has now been removed. Use the `pingfederate_oauth_server_settings` resource going forward.

## Resource: pingfederate_oauth_ciba_server_policy_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_oauth_open_id_connect_policy

This resource has been previously deprecated and has now been removed. Use the `pingfederate_openid_connect_policy` resource going forward.

## Resource: pingfederate_oauth_server_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_oauth_token_exchange_generator_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_open_id_connect_settings

This resource has been previously deprecated and has now been removed. Use the `pingfederate_openid_connect_settings` resource going forward.

## Resource: pingfederate_openid_connect_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

### `session_settings` optional parameter removed

The `session_settings` optional parameter has been removed.  Use the `pingfederate_session_settings` resource going forward.

## Resource: pingfederate_ping_one_connection

This resource has been previously deprecated and has now been removed. Use the `pingfederate_pingone_connection` resource going forward.

## Resource: pingfederate_protocol_metadata_lifetime_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_redirect_validation

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_server_settings_general_settings

This resource has been previously deprecated and has now been removed. Use the `pingfederate_server_settings_general` resource going forward.

## Resource: pingfederate_server_settings_general

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_server_settings_log_settings

This resource has been previously deprecated and has now been removed. Use the `pingfederate_server_settings_logging` resource going forward.

## Resource: pingfederate_server_settings_logging

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_server_settings_system_keys

This resource has been previously deprecated and has now been removed. Use the `pingfederate_server_settings_system_keys_rotate` resource going forward.

## Resource: pingfederate_server_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

### `captcha_settings` optional parameter removed

The `captcha_settings` optional parameter has been removed.  Use the `pingfederate_captcha_provider` resource going forward.

### `email_server` optional parameter removed

The `email_server` optional parameter has been removed.  Use the `pingfederate_notification_publisher` resource going forward.

### `federation_info.auto_connect_entity_id` computed attribute removed

The `federation_info.auto_connect_entity_id` computed attribute parameter has been removed as it is no longer used.

### `roles_and_protocols.idp_role.saml_2_0_profile.enable_auto_connect` computed attribute removed

The `roles_and_protocols.idp_role.saml_2_0_profile.enable_auto_connect` computed attribute parameter has been removed as it is no longer used.

### `roles_and_protocols.sp_role.saml_2_0_profile.enable_auto_connect` computed attribute removed

The `roles_and_protocols.sp_role.saml_2_0_profile.enable_auto_connect` computed attribute parameter has been removed as it is no longer used.

## Resource: pingfederate_session_application_policy

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_session_application_session_policy

This resource has been previously deprecated and has now been removed. Use the `pingfederate_session_application_policy` resource going forward.

## Resource: pingfederate_session_authentication_policies_global

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_session_authentication_session_policies_global

This resource has been previously deprecated and has now been removed. Use the `pingfederate_session_authentication_policies_global` resource going forward.

## Resource: pingfederate_session_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Resource: pingfederate_virtual_host_names

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Data Source: pingfederate_key_pair_signing_import

This data source has been previously deprecated and has now been removed. Use the `pingfederate_keypairs_signing_key` data source going forward.

## Data Source: pingfederate_key_pair_ssl_server_import

This data source has been previously deprecated and has now been removed. Use the `pingfederate_keypairs_ssl_server_key` data source going forward.

## Data Source: pingfederate_oauth_auth_server_settings_scopes_common_scope

This data source has been previously deprecated and has now been removed. Use the `pingfederate_oauth_server_settings` data source going forward.

## Data Source: pingfederate_oauth_auth_server_settings_scopes_exclusive_scope

This data source has been previously deprecated and has now been removed. Use the `pingfederate_oauth_server_settings` data source going forward.

## Data Source: pingfederate_server_settings_system_keys

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Data Source: pingfederate_server_settings

### `federation_info.auto_connect_entity_id` computed attribute removed

The `federation_info.auto_connect_entity_id` computed attribute has been removed as it is no longer used.

### `roles_and_protocols.idp_role.saml_2_0_profile.enable_auto_connect` computed attribute removed

The `roles_and_protocols.idp_role.saml_2_0_profile.enable_auto_connect` computed attribute has been removed as it is no longer used.

### `roles_and_protocols.sp_role.saml_2_0_profile.enable_auto_connect` computed attribute removed

The `roles_and_protocols.sp_role.saml_2_0_profile.enable_auto_connect` computed attribute has been removed as it is no longer used.

## Data Source: pingfederate_session_settings

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.

## Data Source: pingfederate_virtual_host_names

### `id` computed attribute removed

The unnecessary `id` computed attribute has previously been deprecated and has now been removed.
