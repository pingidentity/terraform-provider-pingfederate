# v1.6.0 July 23, 2025
### Enhancements
* Added support for PingFederate `12.3.0` and implemented new attributes for the new version. Added support for latest PF patch releases to `11.3`, `12.0`, `12.1`, and `12.2`. ([#528]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/528)))

### Breaking changes
* Removed support for PingFederate `11.2.x`, in accordance with Ping's [end of life policy](https://support.pingidentity.com/s/article/Ping-Identity-EOL-Tracker). ([#529]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/529)))

### Resources
* **New Resource:** `pingfederate_oauth_out_of_band_auth_plugin` ([#331](https://github.com/pingidentity/terraform-provider-pingfederate/pull/331))
* **New Resource:** `pingfederate_oauth_resource_owner_credentials_mapping` ([#327](https://github.com/pingidentity/terraform-provider-pingfederate/pull/327))

### Bug fixes
* Fixed an issue where the `idp_browser_sso.sso_application_endpoint` attribute in `pingfederate_sp_idp_connection` could cause repeated plans after a successful `terraform apply`. ([#527](https://github.com/pingidentity/terraform-provider-pingfederate/pull/527))
* Warn rather than fail when the provider is unable to delete a `pingfederate_keypairs_ssl_server_key` resource due to the key being in use ([#523](https://github.com/pingidentity/terraform-provider-pingfederate/pull/523))

# v1.5.0 May 28, 2025
### Enhancements
* Added `formatted_file_data` field to `pingfederate_certificate_ca` and `pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate` resources, to handle drift detection when the certificate is changed outside of terraform. ([#502](https://github.com/pingidentity/terraform-provider-pingfederate/pull/502))

### Bug fixes
* Fixed the required `file_data` field not being written to state on import for the `pingfederate_certificate_ca`, `pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate`, and `pingfederate_metadata_url` resources. ([#502](https://github.com/pingidentity/terraform-provider-pingfederate/pull/502))
* Fixed validation for the `idp_browser_sso.adapter_mappings.attribute_sources` attribute in the `pingfederate_sp_idp_connection` resource. The attribute sources are now limited to a maximum size of `1`, and the `id` attribute for the individual attribute sources is removed, as it is not supported in IdP connection adapter mappings. ([#503](https://github.com/pingidentity/terraform-provider-pingfederate/pull/503))
* Fix potential panic due to missing creation timestamp values in `pingfederate_idp_sp_connection`. ([#509](https://github.com/pingidentity/terraform-provider-pingfederate/pull/509))

# v1.4.5 April 30, 2025
### Bug fixes
* Fixed an inconsistent result error that would occur when configuring certificates with no `id` value in the `pingfederate_idp_sp_connection` and `pingfederate_sp_idp_connection` resources. Also fixed a related inconsistent result failure that would occur when modifying the `id` of a certificate. ([#487](https://github.com/pingidentity/terraform-provider-pingfederate/pull/487))
* Fixed plan validation logic in various resources that did not correctly handle unknown values, such as values that depend on the output of another resource. ([#488](https://github.com/pingidentity/terraform-provider-pingfederate/pull/488))
* Fixed missing `false` default for `pingfederate_incoming_proxy_settings.enable_client_cert_header_auth`. ([#494](https://github.com/pingidentity/terraform-provider-pingfederate/pull/494))
* Fixed the `encrypted_value` fields of sensitive configuration fields not being correctly written to state. ([#497](https://github.com/pingidentity/terraform-provider-pingfederate/pull/497))
* Fixed potential inconsistent result errors when using certain escaped OGNL expressions in resource configuration. ([#499](https://github.com/pingidentity/terraform-provider-pingfederate/pull/499))
* Fixed the `pingfederate_authentication_policies_settings` resource not correctly resetting to default on destroy. ([#501](https://github.com/pingidentity/terraform-provider-pingfederate/pull/501))

# v1.4.4 April 8, 2025
### Bug fixes
* Updated the `pingfederate_authentication_policies` resource to set no default for custom attribute source `filter_fields.value` attributes, and to validate that the configured filter field value string has length at least 1. This will prevent inconsistent result errors when using custom attribute sources in the authentication policies. Related to a known terraform-plugin-framework bug with defaults in nested sets: [#867](https://github.com/hashicorp/terraform-plugin-framework/issues/867). ([#484](https://github.com/pingidentity/terraform-provider-pingfederate/pull/484))

### Notes
* bump Go 1.23.5 => 1.24.1 ([#478](https://github.com/pingidentity/terraform-provider-pingfederate/pull/478))

# v1.4.3 March 6, 2025
### Bug fixes
* Marked the `default_request_policy_ref` field as optional in the `pingfederate_oauth_ciba_server_policy_settings` resource, to allow setting the ref to null when no request policies are defined ([#463](https://github.com/pingidentity/terraform-provider-pingfederate/pull/463))
* Marked the `default_generator_group_ref` field as optional in the `pingfederate_oauth_token_exchange_generator_settings` resource, to allow setting the ref to null when no generator groups are defined ([#463](https://github.com/pingidentity/terraform-provider-pingfederate/pull/463))
* Marked the `default_captcha_provider_ref` field as optional in the `pingfederate_captcha_provider_settings` resource, to allow setting the ref to null when no captcha providers are defined ([#464](https://github.com/pingidentity/terraform-provider-pingfederate/pull/464))
* Marked the `default_notification_publisher_ref` field as optional in the `pingfederate_notification_publisher_settings` resource, to allow setting the ref to null when no notification publishers are defined ([#464](https://github.com/pingidentity/terraform-provider-pingfederate/pull/464))
* Marked the `default_access_token_manager_ref` field as optional in the `pingfederate_oauth_access_token_manager_settings` resource, to allow setting the ref to null when no access token managers are defined ([#464](https://github.com/pingidentity/terraform-provider-pingfederate/pull/464))
* Fixed an issue in config validation for `pingfederate_openid_connect_policy` that reported an invalid attribute configuration when using a variable value in `attribute_contract_fulfillment` ([#473](https://github.com/pingidentity/terraform-provider-pingfederate/pull/473))

# v1.4.2 February 5, 2025
### Bug fixes
* Fixed a panic that could occur as a result of HTTP connection errors when using the `pingfederate_local_identity_profile` resource ([#456](https://github.com/pingidentity/terraform-provider-pingfederate/pull/456))

### Notes
* bump Go 1.22.2 => 1.23.5 ([#460](https://github.com/pingidentity/terraform-provider-pingfederate/pull/460))
* bump `github.com/hashicorp/terraform-plugin-go` 0.25.0 => 0.26.0 ([#454](https://github.com/pingidentity/terraform-provider-pingfederate/pull/454))

# v1.4.1 January 23, 2025
### Bug fixes
* Fixed inconsistent result errors when `error_result` is set to `null` within `issuance_criteria` attributes. ([#452]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/452)))

# v1.4.0 January 22, 2025
### Breaking changes
* `sp_idp_connection` resource
  * Marked `idp_oauth_grant_attribute_mapping.access_token_manager_mappings` as required and with minimum size of 1. Previously not including this attribute would have been allowed by the provider, but rejected by PingFederate. Now the provider itself will require at least one access token manager mapping when setting `idp_oauth_grant_attribute_mapping`.
  * Marked `idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract` as required. Previously not including this attribute would have been allowed by the provider, but rejected by PingFederate. Now the provider itself will require this attribute when setting `idp_oauth_grant_attribute_mapping`.
  * Changed `idp_browser_sso.adapter_mappings` from a Set to a List to work around terraform set identity issues.

### Enhancements
* Added missing `ldap_data_store.ldaps_dns_srv_prefix` attribute to the `pingfederate_data_store` resource and data source.
* `sp_idp_connection` resource: Added additional validation that `inbound_provisioning` group attributes do not conflict.

### Bug fixes
* Added missing empty default for some `certs` attributes.
* Fixed incorrect length validation for certain resource link `id` attributes.
* `sp_idp_connection` and `idp_sp_connection` resources
  * Fixed `password` incorrectly defaulting to an empty string for `credentials.inbound_back_channel_auth.http_basic_credentials.password` and `credentials.outbound_back_channel_auth.http_basic_credentials.password`.
  * Fixed JSON marshal errors when setting `attribute_sources` values.
* `idp_sp_connection` resource
  * Fixed `sp_browser_sso.artifact.lifetime` attribute incorrectly being marked as required.
  * Fixed invalid validation requiring either `sp_browser_sso.sign_response_as_required` or `sp_browser_sso.sign_assertions` to be set to `true`.
  * Fixed some booleans being left out of state when set to `false`.
  * Added missing `false` default for `attribute_query.policy.encrypt_assertion`, `attribute_query.policy.require_encrypted_name_id`, `attribute_query.policy.require_signed_attribute_query`, `attribute_query.policy.sign_assertion`, and `attribute_query.policy.sign_response`.
  * Added missing `false` default for `credentials.signing_settings.include_raw_key_in_signature`.
  * Added missing `true` default for `outbound_provision.channels.#.channel_source.account_management_settings.default_status` and `outbound_provision.channels.#.channel_source.account_management_settings.flag_comparison_status`.
  * Added missing `false` default for `sp_browser_sso.adapter_mappings.#.restrict_virtual_entity_ids` and `sp_browser_sso.authentication_policy_contract_assertion_mappings.#.restrict_virtual_entity_ids`.
  * Added missing `false` default for `sp_browser_sso.always_sign_artifact_response`, `sp_browser_sso.require_signed_authn_requests`, and `sp_browser_sso.sign_assertions`.
  * Added missing `false` default for `sp_browser_sso.encryption_policy.encrypt_assertion`, `sp_browser_sso.encryption_policy.encrypt_slo_subject_name_id`, and `sp_browser_sso.encryption_policy.slo_subject_name_id_encrypted`
  * Added missing `false` default for `ws_trust.encrypt_saml2_assertion`, `ws_trust.generate_key`, and `ws_trust.oauth_assertion_profiles`.

* `sp_idp_connection` resource
  * Fixed unexpected update plans that could occur when setting `credentials.certs`.
  * Added missing `oidc_client_credentials.encrypted_secret` attribute, to be used as an alternative to `oidc_client_credentials.client_secret`, and marked `oidc_client_credentials.client_secret` as sensitive.
  * Added missing `false` default for `credentials.inbound_back_channel_auth.digital_signature` and `credentials.outbound_back_channel_auth.digital_signature`.
  * Added missing `false` default for `credentials.inbound_back_channel_auth.require_ssl`.
  * Added missing default for `error_page_msg_id`. Defaults to `errorDetail.spSsoFailure` for browser SSO connections, null otherwise.
  * Fixed `idp_browser_sso.adapter_mappings.#.adapter_override_settings.plugin_descriptor_ref`, `idp_browser_sso.adapter_mappings.#.adapter_override_settings.name`, and `idp_browser_sso.adapter_mappings.#.sp_adapter_ref` incorrectly being marked as required.
  * Added missing empty default for `idp_browser_sso.adapter_mappings.#.adapter_override_settings.target_application_info`.
  * Added missing defaults for `idp_browser_sso.adapter_mappings.#.issuance_criteria`, `idp_browser_sso.adapter_mappings.#.restrict_virtual_entity_ids`, and `idp_browser_sso.adapter_mappings.#.restricted_virtual_entity_ids`.
  * Added missing default for `idp_browser_sso.assertions_signed`, `idp_browser_sso.authentication_policy_contract_mappings`, and `idp_browser_sso.sign_authn_requests`.
  * Added missing `false` defaults for each boolean attribute in `idp_browser_sso.decryption_policy`.
  * Fixed provider failures when configuring `idp_browser_sso.jit_provisioning`.
  * Fixed unexpected new value errors for `oidc_provider_settings.back_channel_logout_uri`, `oidc_provider_settings.front_channel_logout_uri`, `oidc_provider_settings.post_logout_redirect_uri`, and `oidc_provider_settings.redirect_uri`.
  * Added missing `false` default for `oidc_provider_settings.enable_pkce` and `oidc_provider_settings.track_user_sessions_for_logout`.
  * Fixed unexpected new value error for `oidc_provider_settings.jwt_secured_authorization_response_mode_type`.
  * Fixed incorrectly required `oidc_provider_settings.request_parameters.#.attribute_value.value` attribute.
  * Fixed `idp_browser_sso.sso_service_endpoints.#.binding` incorrectly being marked as required.
  * Updated to allow blank path for `idp_browser_sso.url_whitelist_entries.#.valid_path`, and defaulted to an empty string.
  * Added missing empty default for `idp_oauth_grant_attribute_mapping.idp_oauth_attribute_contract.extended_attributes`.
  * Fixed failures when configuring `inbound_provisioning.users` and `inbound_provisioning.groups`.
  * Fixed failures related to `core_attributes` when configuring `inbound_provisioning.user_repository`.
  * Fixed unexpected new value errors for `inbound_provisioning.groups.read_groups.attribute_contract.core_attributes` and `inbound_provisioning.users.read_users.attribute_contract.core_attributes`.
  * Added missing empty default for `inbound_provisioning.groups.read_groups.attribute_fulfillment.#.source.value` and `inbound_provisioning.groups.write_groups.attribute_fulfillment.#.source.value` and marked as optional rather than required.
  * Added empty default for `inbound_provisioning.user_repository.ldap.unique_group_id_filter` and marked as optional rather than required.
  * Added missing empty default for `inbound_provisioning.users.read_users.attribute_contract.extended_attributes`.
  * Added missing empty default for `inbound_provisioning.users.read_users.attribute_fulfillment.#.source.value` and `inbound_provisioning.users.write_users.attribute_fulfillment.#.source.value` and marked as optional rather than required.
  * Added empty or null default for `virtual_entity_ids`, depending on the type of connection
  * Added null default for `jwt_secured_authorization_response_mode_type` for PingFederate versions prior to `12.1`.
  * Fixed errors with `false` values being returned as `null` for `credentials.signing_settings.include_cert_in_signature`, `idp_browser_sso.sign_authn_request`, and `idp_browser_sso.assertions_signed`.

# v1.3.0 January 9, 2025
### Enhancements
* Added support for PingFederate `12.2.0` and implemented new attributes for the new version. Added support for latest PF patch releases to `11.2`, `11.3`, `12.0`, and `12.1`. This will be the last release with support for PingFederate `11.2` in accordance with Ping's [end of life policy](https://support.pingidentity.com/s/article/Ping-Identity-EOL-Tracker). ([#440]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/440)))

### Bug Fixes
* Fix terraform import failure for certain `pingfederate_sp_idp_connection` configurations ([#442](https://github.com/pingidentity/terraform-provider-pingfederate/pull/442))
* Fix URL config validator where some asterisks in value returned "Invalid URL Format" ([#445](https://github.com/pingidentity/terraform-provider-pingfederate/pull/445))

### Notes
* bump `golang.org/x/net` 0.31.0 => 0.33.0 ([#441](https://github.com/pingidentity/terraform-provider-pingfederate/pull/441))

# v1.2.0 December 16, 2024
### Enhancements
* Added missing `ldap_data_store.ldaps_dns_srv_prefix` attribute to the `pingfederate_data_store` resource and data source. ([#428](https://github.com/pingidentity/terraform-provider-pingfederate/pull/428))

### Notes
* bump `github.com/hashicorp/terraform-plugin-framework` 1.11.0 => 1.13.0 ([#434](https://github.com/pingidentity/terraform-provider-pingfederate/pull/434))
* bump `github.com/hashicorp/terraform-plugin-framework-validators` 0.12.0 => 0.16.0 ([#434](https://github.com/pingidentity/terraform-provider-pingfederate/pull/434))
* bump `github.com/hashicorp/terraform-plugin-go` 0.23.0 => 0.25.0 ([#434](https://github.com/pingidentity/terraform-provider-pingfederate/pull/434))
* bump `golang.org/x/crypto` 0.29.0 => 0.31.0 ([#434](https://github.com/pingidentity/terraform-provider-pingfederate/pull/434))

# v1.1.0 October 31, 2024
### Enhancements
* Added `encrypted_` attributes for sensitive attributes. The `encrypted_` versions of these attributes can be used as an alternative to the original attribute when importing configuration into Terraform from an existing PingFederate. ([#419](https://github.com/pingidentity/terraform-provider-pingfederate/pull/419))

### Resources
* **New Resource:** `pingfederate_config_store` ([#420](https://github.com/pingidentity/terraform-provider-pingfederate/pull/420))

### Data Sources
* **New Data Source:** `pingfederate_config_store` ([#420](https://github.com/pingidentity/terraform-provider-pingfederate/pull/420))

# v1.0.0 October 1, 2024
### Breaking changes
As this is the first major release of the provider, there are breaking changes from `0.x` versions. The primary breaking changes are removal of previously deprecated attributes. For more information, see the Version 1 Upgrade Guide in the registry documentation. ([#413]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/413)))

### Resources
* **New Resource:** `pingfederate_certificates_group` ([#278](https://github.com/pingidentity/terraform-provider-pingfederate/pull/278))
* **New Resource:** `pingfederate_identity_store_provisioner` ([#303](https://github.com/pingidentity/terraform-provider-pingfederate/pull/303))
* **New Resource:** `pingfederate_idp_token_processor` ([#277](https://github.com/pingidentity/terraform-provider-pingfederate/pull/277))
* **New Resource:** `pingfederate_keypairs_oauth_openid_connect` ([#267](https://github.com/pingidentity/terraform-provider-pingfederate/pull/267))
* **New Resource:** `pingfederate_keypairs_signing_key_rotation_settings` ([#334](https://github.com/pingidentity/terraform-provider-pingfederate/pull/334))
* **New Resource:** `pingfederate_keypairs_ssl_client_csr_export` ([#335](https://github.com/pingidentity/terraform-provider-pingfederate/pull/335))
* **New Resource:** `pingfederate_keypairs_ssl_client_csr_response` ([#335](https://github.com/pingidentity/terraform-provider-pingfederate/pull/335))
* **New Resource:** `pingfederate_keypairs_ssl_client_key` ([#337](https://github.com/pingidentity/terraform-provider-pingfederate/pull/337))
* **New Resource:** `pingfederate_keypairs_ssl_server_csr_export` ([#336](https://github.com/pingidentity/terraform-provider-pingfederate/pull/336))
* **New Resource:** `pingfederate_keypairs_ssl_server_csr_response` ([#336](https://github.com/pingidentity/terraform-provider-pingfederate/pull/336))
* **New Resource:** `pingfederate_oauth_ciba_server_policy_request_policy` ([#285]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/285)))
* **New Resource:** `pingfederate_oauth_client_settings` ([#286](https://github.com/pingidentity/terraform-provider-pingfederate/pull/286))
* **New Resource:** `pingfederate_protocol_metadata_signing_settings` ([#290](https://github.com/pingidentity/terraform-provider-pingfederate/pull/290))
* **New Resource:** `pingfederate_server_settings_system_keys_rotate` ([#292](https://github.com/pingidentity/terraform-provider-pingfederate/pull/292))
* **New Resource:** `pingfederate_server_settings_ws_trust_sts_settings` ([#294](https://github.com/pingidentity/terraform-provider-pingfederate/pull/294))
* **New Resource:** `pingfederate_server_settings_ws_trust_sts_settings_issuer_certificate` ([#293](https://github.com/pingidentity/terraform-provider-pingfederate/pull/293))
* **New Resource:** `pingfederate_service_authentication` ([#295](https://github.com/pingidentity/terraform-provider-pingfederate/pull/295))
* **New Resource:** `pingfederate_sp_adapter` ([#265](https://github.com/pingidentity/terraform-provider-pingfederate/pull/265))
* **New Resource:** `pingfederate_sp_idp_connection` ([#342](https://github.com/pingidentity/terraform-provider-pingfederate/pull/342))

### Data sources
* **New Data Source:** `pingfederate_certificates_ca_export` ([#296](https://github.com/pingidentity/terraform-provider-pingfederate/pull/296))
* **New Data Source:** `pingfederate_keypairs_signing_certificate` ([#305](https://github.com/pingidentity/terraform-provider-pingfederate/pull/305))
* **New Data Source:** `pingfederate_keypairs_ssl_client_key` ([#337]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/337)))
* **New Data Source:** `pingfederate_keypairs_ssl_client_certificate` ([#338]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/338)))
* **New Data Source:** `pingfederate_keypairs_ssl_server_certificate` ([#338]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/338)))

### Enhancements
* Added new `configuration.sensitive_fields` and `configuration.tables.#.rows.#.sensitive_fields` attributes to plugin configuration across the provider. Use these fields when specifying sensitive fields in plugin configuration, such as secrets and passwords. Values specified in these sets will be marked as Sensitive to Terraform and hidden in the CLI and UI output. ([#383]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/383)))
* Added support for PingFederate patches through `11.2.10`, `11.3.8`, `12.0.5`, `12.1.3`. ([#406]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/406)))
* Set defaults for many computed and optional fields that previously did not have defaults.

### Bug Fixes
* Prevent unexpected plans for `pingfederate_oauth_client_settings` when used in conjunction with `pingfederate_extended_properties`. `pingfederate_oauth_client_settings` should always be marked to `depends_on` `pingfederate_extended_properties` when these resources are used together. ([#409](https://github.com/pingidentity/terraform-provider-pingfederate/pull/409))
* Prevent unexpected plans for `pingfederate_openid_connect_settings` when used in conjunction with `pingfederate_session_settings`. `pingfederate_openid_connect_settings` should always be marked to `depends_on` `pingfederate_session_settings` when these resources are used together. ([#413](https://github.com/pingidentity/terraform-provider-pingfederate/pull/413))
* Prevent create failures when creating multiple key pairs simultaneously.
* Fix various config errors when importing configuration with `-generate-config-out`.

# v0.16.0 September 12, 2024
### Resources
* **New Resource:** `pingfederate_certificates_revocation_ocsp_certificate` ([#279]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/279)))
* **New Resource:** `pingfederate_configuration_encryption_keys_rotate` ([#289](https://github.com/pingidentity/terraform-provider-pingfederate/pull/289))
* **New Resource:** `pingfederate_oauth_client_registration_policy` ([#330](https://github.com/pingidentity/terraform-provider-pingfederate/pull/330))
* **New Resource:** `pingfederate_keypairs_signing_key` ([#313]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/313)))
* **New Resource:** `pingfederate_keypairs_ssl_server_key` ([#314]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/314)))
* **New Resource:** `pingfederate_idp_to_sp_adapter_mapping` ([#264]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/264)))

### Data sources
* **New Data source:** `pingfederate_keypairs_signing_key` ([#313]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/313)))
* **New Data source:** `pingfederate_keypairs_ssl_server_key` ([#314]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/314)))

### Enhancements
* Updated resources using read-only `_all` attributes to read to the non-`_all` attribute on import. This includes common attributes such as `configuration.fields` and `configuration.tables`. ([#365]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/365)))
* Various documentations and schema improvements for resources released in previous versions.
* Improved consistency and readability of error and warning messages.

### Bug fixes
* Fixed terraform plan error after generating import HCL for `pingfederate_oauth_client` due to `persistent_grant_expiration_time` and `persistent_grant_expiration_time_unit` attributes. ([#365]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/365)))

### Breaking changes
* The `pingfederate_server_settings_system_keys` attribute `previous` has been changed to read-only. That resource has also been deprecated (see below). ([#362]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/362)))

### Deprecated
* The `id` attribute that was formerly deprecated in non-singleton-resources is no longer deprecated. ([#367]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/367)))
* The `pingfederate_key_pair_ssl_server_import` resource and data source have been renamed. Use `pingfederate_keypairs_ssl_server_key` instead. `pingfederate_key_pair_ssl_server_import` will be removed in a future release. ([#314]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/314)))
* The `pingfederate_key_pair_signing_import` resource and data source have been renamed. Use `pingfederate_keypairs_signing_key` instead. `pingfederate_key_pair_signing_import` will be removed in a future release. ([#313]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/313)))
* The `pingfederate_server_settings_system_keys` resource is deprecated. Use `pingfederate_server_settings_system_keys_rotate` instead. `pingfederate_server_settings_system_keys` will be removed in a future release. ([#382]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/382)))

# v0.15.0 August 29, 2024
### Resources
* **New Resource:** `pingfederate_secret_manager` ([#291]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/291)))
* **New Resource:** `pingfederate_cluster_settings` ([#339]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/339)))
* **New Resource:** `pingfederate_metadata_url` ([#282]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/282)))
* **New Resource:** `pingfederate_captcha_provider_settings` ([#283](https://github.com/pingidentity/terraform-provider-pingfederate/pull/283))
* **New Resource:** `pingfederate_certificates_revocation_settings` ([#280]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/280)))
* **New Resource:** `pingfederate_idp_sts_request_parameters_contract` ([#281]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/281)))

### Data Sources
* **New Data Source:** `pingfederate_cluster_status` ([#297](https://github.com/pingidentity/terraform-provider-pingfederate/pull/297))

### Bug Fixes
* Fixed "inconsistent result after apply" issues with resources using `attribute_sources.ldap_attribute_source.binary_attribute_settings`. ([#318](https://github.com/pingidentity/terraform-provider-pingfederate/issues/318))
* Fixed "inconsistent result after apply" issues with the `pingfederate_idp_sp_connection` resource, caused by usage of the `sp_browser_sso.incoming_bindings`, `sp_browser_sso.enabled_profiles`, and `extended_attributes` attributes. ([#319](https://github.com/pingidentity/terraform-provider-pingfederate/issues/319))
* Fixed panics when using OAuth in the provider configuration while configuring certain resources ([#352](https://github.com/pingidentity/terraform-provider-pingfederate/pull/352))

### Enhancements
* Various documentations and schema improvements for resources released in previous versions.

### Deprecated
* The `pingfederate_session_application_session_policy` resource and data source have been renamed. Use `pingfederate_session_application_policy` instead. `pingfederate_session_application_session_policy` will be removed in a future release. ([#343]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/343)))
* The `pingfederate_local_identity_identity_profile` resource and data source have been renamed. Use `pingfederate_local_identity_profile` instead. `pingfederate_local_identity_identity_profile` will be removed in a future release. ([#346]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/346)))
* The `pingfederate_notification_publishers_settings` resource and data source have been renamed. Use `pingfederate_notification_publisher_settings` instead. `pingfederate_notification_publishers_settings` will be removed in a future release. ([#346]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/346)))
* The `pingfederate_open_id_connect_settings` resource and data source have been renamed. Use `pingfederate_openid_connect_settings` instead. `pingfederate_open_id_connect_settings` will be removed in a future release. ([#346]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/346)))
* The `pingfederate_oauth_auth_server_settings` resource and data source have been renamed. Use `pingfederate_oauth_server_settings` instead. `pingfederate_oauth_auth_server_settings` will be removed in a future release. ([#346]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/346)))
* The `pingfederate_server_settings_general_settings` resource and data source have been renamed. Use `pingfederate_server_settings_general` instead. `pingfederate_server_settings_general_settings` will be removed in a future release. ([#346]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/346)))
* The `pingfederate_session_authentication_session_policies_global` resource and data source have been renamed. Use `pingfederate_session_authentication_policies_global` instead. `pingfederate_session_authentication_session_policies_global` will be removed in a future release. ([#347]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/347)))
* The `pingfederate_server_settings_log_settings` resource and data source have been renamed. Use `pingfederate_server_settings_logging` instead. `pingfederate_server_settings_log_settings` will be removed in a future release. ([#347]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/347)))
* The `pingfederate_oauth_open_id_connect_policy` resource and data source have been renamed. Use `pingfederate_openid_connect_policy` instead. `pingfederate_oauth_open_id_connect_policy` will be removed in a future release. ([#347]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/347)))
* The `pingfederate_ping_one_connection` resource has been renamed. Use `pingfederate_pingone_connection` instead. `pingfederate_ping_one_connection` will be removed in a future release. ([#347]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/347)))

# v0.14.0 August 15, 2024
### Resources
* **New Resource:** `pingfederate_keypairs_oauth_openid_connect_additional_key_set` ([#271]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/271)))
* **New Resource:** `pingfederate_captcha_provider` ([#275]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/275)))
* **New Resource:** `pingfederate_oauth_access_token_manager_settings` ([#274]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/274)))
* **New Resource:** `pingfederate_notification_publisher` ([#284]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/284)))
* **New Resource:** `pingfederate_connection_metadata_export` ([#276](https://github.com/pingidentity/terraform-provider-pingfederate/pull/276))

### Bug Fixes
* Fixed inability to configure mutliple `hostnames_tags` in the `pingfederate_data_store` resource.
* Fixed config validation error when using variables for password and username in the `pingfederate_kerberos_realm` resource.

# v0.13.0 August 1, 2024
### Resources
* **New Resource:** `pingfederate_default_urls` ([#260]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/260)))
* **New Resource:** `pingfederate_sp_target_url_mappings` ([#273]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/273)))
* **New Resource:** `pingfederate_keypairs_ssl_server_settings` ([#272]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/272)))
* **New Resource:** `pingfederate_oauth_authentication_policy_contract_mapping` ([#262]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/262)))
* **New Resource:** `pingfederate_session_authentication_policy` ([#261]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/261)))
* **New Resource:** `pingfederate_kerberos_realm_settings` ([#266]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/266)))
* **New Resource:** `pingfederate_oauth_idp_adapter_mapping` ([#263]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/263)))

### Deprecated
* The `pingfederate_idp_urls` resource has been deprecated. Use `pingfederate_default_urls` instead. `pingfederate_idp_urls` will be removed in a future release. ([#260]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/260)))

### Enhancements
* Improved examples, schema documentation, and attribute validation for `pingfederate_oauth_issuer`, `pingfederate_authentication_policy_contract`, `pingfederate_certificate_ca`, `pingfederate_extended_properties`, `pingfederate_incoming_proxy_settings`.

# v0.12.0 June 27th, 2024
### BREAKING CHANGES
* Removed deprecated `inherited` attribute from various resources and data sources. ([#268](https://github.com/pingidentity/terraform-provider-pingfederate/pull/268))

### ENHANCEMENTS
* Added support for PingFederate `12.1.0` and implemented new attributes for the new version. Added support for latest PF patch releases to `11.2`, `11.3`, and `12.0`. ([#268]([https](https://github.com/pingidentity/terraform-provider-pingfederate/pull/268)))

# v0.11.0 May 30th, 2024
### DEPRECATED
* `location` property in resource reference object types for all Resources and DataSources removed ([#249](https://github.com/pingidentity/terraform-provider-pingfederate/pull/249))

### BUG FIXES
* `pingfederate_authentication_selector` resource `extended_attributes` unexpected new value ([#249](https://github.com/pingidentity/terraform-provider-pingfederate/pull/249))
* `pingfederate_local_identity_identity_profile` resource `data_store_config` attributes unexpected new value ([#249](https://github.com/pingidentity/terraform-provider-pingfederate/pull/249))
* `pingfederate_ping_one_connection` resource correctly requires `credential` value ([#233](https://github.com/pingidentity/terraform-provider-pingfederate/pull/233))

### Resources
* **New Resource:** `pingfederate_authentication_policies` ([#249](https://github.com/pingidentity/terraform-provider-pingfederate/pull/249))

# v0.10.0 April 11th, 2024
### Resources
* **New Resource:** `pingfederate_ping_one_connection` ([#231](https://github.com/pingidentity/terraform-provider-pingfederate/pull/231))

# v0.9.0 March 29th, 2024
### FEATURES
* Resources that no longer incorrectly require a `<resource_type>_id` are listed here:
  - `pingfederate_authentication_api_application`
  - `pingfederate_authentication_policies_fragment`
  - `pingfederate_authentication_policy_contract`
  - `pingfederate_authentication_selector`
  - `pingfederate_certificate_ca`
  - `pingfederate_data_store`
  - `pingfederate_idp_adapter`
  - `pingfederate_idp_sp_connection`
  - `pingfederate_kerberos_realm`
  - `pingfederate_key_pair_signing_import`
  - `pingfederate_key_pair_ssl_server_import`
  - `pingfederate_local_identity_identity_profile`
  - `pingfederate_oauth_access_token_manager`
  - `pingfederate_oauth_issuer`
  - `pingfederate_oauth_open_id_connect_policy`
  - `pingfederate_password_credential_validator`

### ENHANCEMENTS
* Allow `product_version` values that are not explicitly supported as long as the major-minor version is supported. For example, version `11.3.10` would be allowed, but version `30.0.0` would not be allowed. (#223)
* Added support for PingFederate patch versions `11.3.5` and `12.0.1` (#226)

### BUG FIXES
* `pingfederate_authentication_policy_contract` resource sends in required default value for `core_attributes` ([#224](https://github.com/pingidentity/terraform-provider-pingfederate/pull/224))
* `pingfederate_oauth_client` resource has the following bugs resolved ([#221](https://github.com/pingidentity/terraform-provider-pingfederate/pull/221)):
  - `client_id` now correctly forces the resource to be replaced when value is modified after creation.
  - `client_auth.secret` corrected config validation when using a variable value
  - `restrict_scopes` no longer returns an incorrect value error after apply
  - `oidc_policy` child string attribute values no longer return incorrect value errors after apply

# v0.8.0 March 14th, 2024
### BUG FIXES
* `pingfederate_oauth_access_token_mapping` resource, resolved "produced an unexpected new value: .attribute_contract_fulfillment["username"].value: was null, but now cty.StringVal("")" error([#215](https://github.com/pingidentity/terraform-provider-pingfederate/pull/215))
* `pingfederate_oauth_auth_server_settings` resource, removed unnecessary requirement for `bypass_activation_code_confirmation`, `default_scope_description`, `device_polling_interval`, `pending_authorization_timeout`, `registered_authorization_path` properties([#216](https://github.com/pingidentity/terraform-provider-pingfederate/pull/216))
* `pingfederate_oauth_client` resource, "Provider produced inconsistent result after apply" error when applying empty `extended_parameters` map value, require `values` property within to match product behavior. Fixed `logo_url` "Invalid URL Format" when empty string supplied. ([#204](https://github.com/pingidentity/terraform-provider-pingfederate/pull/204))([#214](https://github.com/pingidentity/terraform-provider-pingfederate/pull/214))([#217](https://github.com/pingidentity/terraform-provider-pingfederate/pull/217))

# v0.7.1 February 29th, 2024
### BUG FIXES
* `pingfederate_oauth_access_token_manager` resource `attribute_contract.extended_attributes` Provider produced inconsistent result after apply ([#202](https://github.com/pingidentity/terraform-provider-pingfederate/pull/202))

# v0.7.0 February 28th, 2024
### DEPRECATED
* `location` property in `resourceLink` object types for all Resources and DataSources. This property will be removed in a future release. ([#195](https://github.com/pingidentity/terraform-provider-pingfederate/pull/195))

### Resources
* **New Resource:** `pingfederate_authentication_selector ` ([#199](https://github.com/pingidentity/terraform-provider-pingfederate/pull/199))
* **New Resource:** `pingfederate_incoming_proxy_settings` ([#190](https://github.com/pingidentity/terraform-provider-pingfederate/pull/190))
* **New Resource:** `pingfederate_notification_publishers_settings` ([#187](https://github.com/pingidentity/terraform-provider-pingfederate/pull/187))
* **New Resource:** `pingfederate_oauth_access_token_mapping` ([#195](https://github.com/pingidentity/terraform-provider-pingfederate/pull/195))
* **New Resource:** `pingfederate_open_id_connect_settings` ([#196](https://github.com/pingidentity/terraform-provider-pingfederate/pull/196))

# v0.6.0 February 9th, 2024
### DEPRECATED
* `inherited` property for all Resources and DataSources. This property will be removed in a future release.

### FEATURES
* Added support for PingFederate version `11.2.8` ([#168](https://github.com/pingidentity/terraform-provider-pingfederate/pull/168))
* `admin_url_path` added to `provider` variables. This variable will allow end-users to set their own admin URL location, rather than being forced to use the hard-coded `/pf-admin-api/v1` suffix value. This variable value will default to `/pf-admin-api/v1` if no value is supplied. ([#170](https://github.com/pingidentity/terraform-provider-pingfederate/pull/170))
* Added support for the `PINGFEDERATE_TF_APPEND_USER_AGENT` environment variable, used to append a custom suffix to the User-Agent header used by the provider when making HTTP requests. ([#171](https://github.com/pingidentity/terraform-provider-pingfederate/pull/171))
* Added support for automatically retrying certain HTTP error codes. ([#171](https://github.com/pingidentity/terraform-provider-pingfederate/pull/171))
* Added support for OAuth2 Client Credentials flow authentication, as well as supporting a supplied Access Token ([#183](https://github.com/pingidentity/terraform-provider-pingfederate/pull/183))

#### Resources
* **New Resource:** `pingfederate_extended_properties` ([#182](https://github.com/pingidentity/terraform-provider-pingfederate/pull/182))

### BUG FIXES
* Fixed provider not correctly comparing versions for PingFederate `11.2.6` and `11.2.7` ([#168](https://github.com/pingidentity/terraform-provider-pingfederate/pull/168))
* Fixed provider errors when using strings with escaped quotes in certain attributes ([#178](https://github.com/pingidentity/terraform-provider-pingfederate/pull/178))
* Fixed `pingfederate_oauth_client` resource incorrectly validating supplied `grant_types` ([#183](https://github.com/pingidentity/terraform-provider-pingfederate/pull/183))

# v0.5.0 January 16th, 2024
### BREAKING CHANGES
* New *required* `product_version` provider attribute. This attribute is used to indicate the version of PingFederate that the provider is targeting. The attribute can also be specified by the `PINGFEDERATE_PROVIDER_PRODUCT_VERSION` environment variable. When upgrading to this provider version, this attribute will need to be configured.

### BUG FIXES
* **Server Settings Log Settings Resource:** Fixed provider errors when not specifying all log categories in the server log settings ([#164](https://github.com/pingidentity/terraform-provider-pingfederate/pull/164))
* **OAuth Client Resource:** Resolved issue where some property default values on the `pingfederate_oauth_client` resource resulted in an invalid apply ([#156](https://github.com/pingidentity/terraform-provider-pingfederate/pull/156))
* **Resource Link Update:** Resolved issue where updating a ResourceLink type in resources resulted in an invalid apply ([#159](https://github.com/pingidentity/terraform-provider-pingfederate/pull/159))

### FEATURES
* Added support for PingFederate version `11.3` and full support for all `11.2.x` versions ([#149](https://github.com/pingidentity/terraform-provider-pingfederate/pull/149))
* Added support for PingFederate version `12.0` ([#153](https://github.com/pingidentity/terraform-provider-pingfederate/pull/153))

#### Resources
* **New Resource:** `pingfederate_authentication_api_application` ([#156](https://github.com/pingidentity/terraform-provider-pingfederate/pull/156))
* **New Resource:** `pingfederate_authentication_policies_fragment` ([#161](https://github.com/pingidentity/terraform-provider-pingfederate/pull/161))
* **New Resource:** `pingfederate_authentication_policies_settings` ([#150](https://github.com/pingidentity/terraform-provider-pingfederate/pull/150))
* **New Resource:** `pingfederate_kerberos_realm` ([#158](https://github.com/pingidentity/terraform-provider-pingfederate/pull/158))
* **New Resource:** `pingfederate_oauth_ciba_server_policy_settings` ([#159](https://github.com/pingidentity/terraform-provider-pingfederate/pull/159))
* **New Resource:** `pingfederate_oauth_token_exchange_generator_settings` ([#163](https://github.com/pingidentity/terraform-provider-pingfederate/pull/163))

#### Data Sources
* **New Data Source:** `pingfederate_authentication_api_application` ([#156](https://github.com/pingidentity/terraform-provider-pingfederate/pull/156))
* **New Data Source:** `pingfederate_authentication_policies_settings` ([#150](https://github.com/pingidentity/terraform-provider-pingfederate/pull/150))

# v0.4.0 December 13, 2023 => MVP
### FEATURES
#### Resources
* **New Resource:** `pingfederate_oauth_client` ([#111](https://github.com/pingidentity/terraform-provider-pingfederate/pull/111))
* **Resource Update:** `pingfederate_oauth_token_exchange_processor_policy_token_generator_mapping` corrected to `pingfederate_oauth_token_exchange_token_generator_mapping`

#### Data Sources
* **New Data Source:** `pingfederate_authentication_policy_contract` ([#118](https://github.com/pingidentity/terraform-provider-pingfederate/pull/118))
* **New Data Source:** `pingfederate_data_store` ([#123](https://github.com/pingidentity/terraform-provider-pingfederate/pull/123))
* **New Data Source:** `pingfederate_idp_adapter` ([#125](https://github.com/pingidentity/terraform-provider-pingfederate/pull/125))
* **New Data Source:** `pingfederate_oauth_client` ([#111](https://github.com/pingidentity/terraform-provider-pingfederate/pull/111))
* **New Data Source:** `pingfederate_oauth_token_exchange_token_generator_mapping` ([#122](https://github.com/pingidentity/terraform-provider-pingfederate/pull/122))
* **New Data Source:** `pingfederate_password_credential_validator` ([#119](https://github.com/pingidentity/terraform-provider-pingfederate/pull/119))
* **New Data Source:** `pingfederate_server_settings` ([#127](https://github.com/pingidentity/terraform-provider-pingfederate/pull/127))
* **New Data Source:** `pingfederate_session_policies_global` ([#117](https://github.com/pingidentity/terraform-provider-pingfederate/pull/117))
* **New Data Source:** `pingfederate_sp_authentication_policy_contract_mapping` ([#124](https://github.com/pingidentity/terraform-provider-pingfederate/pull/124))
* **New Data Source:** `pingfederate_token_processor_to_token_generator_mapping` ([#120](https://github.com/pingidentity/terraform-provider-pingfederate/pull/120))
* **New Data Source:** `pingfederate_oauth_open_id_connect_policy` ([#121](https://github.com/pingidentity/terraform-provider-pingfederate/pull/121))

### ENHANCEMENTS
* Added new `tables_all` attribute to plugin `configuration` attributes, to include any tables that are generated by PingFederate when not specified by the user ([#133](https://github.com/pingidentity/terraform-provider-pingfederate/pull/133))
* Added computed metadata attributes to `pingfederate_license`, `pingfederate_key_pair_signing_import`, and `pingfederate_key_pair_ssl_server_import` resource ([#138](https://github.com/pingidentity/terraform-provider-pingfederate/pull/138))

# v0.3.0 November 22, 2023
### FEATURES
#### Resources
* **New Resource:** `pingfederate_oauth_open_id_connect_policy` ([#105](https://github.com/pingidentity/terraform-provider-pingfederate/pull/105))
* **New Resource:** `pingfederate_idp_sp_connection` ([#128](https://github.com/pingidentity/terraform-provider-pingfederate/pull/128))

#### Data Sources
* **New Data Source:** `pingfederate_protocol_metadata_lifetime_settings` ([#100](https://github.com/pingidentity/terraform-provider-pingfederate/pull/100))
* **New Data Source:** `pingfederate_redirect_validation` ([#116](https://github.com/pingidentity/terraform-provider-pingfederate/pull/116))
* **New Data Source:** `pingfederate_server_settings_general_settings` ([#101](https://github.com/pingidentity/terraform-provider-pingfederate/pull/101))
* **New Data Source:** `pingfederate_server_settings_log_settings` ([#104](https://github.com/pingidentity/terraform-provider-pingfederate/pull/104))
* **New Data Source:** `pingfederate_server_settings_general_settings` ([#105](https://github.com/pingidentity/terraform-provider-pingfederate/pull/105))
* **New Data Source:** `pingfederate_server_settings_system_keys` ([#112](https://github.com/pingidentity/terraform-provider-pingfederate/pull/112))
* **New Data Source:** `pingfederate_session_settings` ([#115](https://github.com/pingidentity/terraform-provider-pingfederate/pull/115))


# v0.2.1 November 10, 2023
### BUG FIXES
* **Administrative Account Resource:** Resolved issue where updating a managed `pingfederate_administrative_account` resource forces replacement ([#86](https://github.com/pingidentity/terraform-provider-pingfederate/pull/86))
### ENHANCEMENTS
* Include values computed from PingFederate in provider state ([#72](https://github.com/pingidentity/terraform-provider-pingfederate/pull/72))
* Use lists instead of sets in most cases ([#81](https://github.com/pingidentity/terraform-provider-pingfederate/pull/81))
* Make the `id` field read-only for resources and data sources, and use a `custom_id` field for setting id values rather than letting them be generated by PingFederate ([#59](https://github.com/pingidentity/terraform-provider-pingfederate/pull/59))
* **`custom_id`** changed to resource type id. For example, the `custom_id` in the `pingfederate_authentication_policy_contract` resource is now `contract_id` ([#98](https://github.com/pingidentity/terraform-provider-pingfederate/pull/98))

### FEATURES
#### Resources
* **New Resource:** `pingfederate_data_store` ([#88](https://github.com/pingidentity/terraform-provider-pingfederate/pull/88))
* **New Resource:** `pingfederate_idp_adapter` ([#64](https://github.com/pingidentity/terraform-provider-pingfederate/pull/64))
* **New Resource:** `pingfederate_token_processor_to_token_generator_mapping` ([#65](https://github.com/pingidentity/terraform-provider-pingfederate/pull/65))
* **New Resource:** `pingfederate_oauth_token_exchange_processor_policy_token_generator_mapping` ([#68](https://github.com/pingidentity/terraform-provider-pingfederate/pull/68))
* **New Resource:** `pingfederate_sp_authentication_policy_contract_mapping` ([#75](https://github.com/pingidentity/terraform-provider-pingfederate/pull/75))

#### Data Sources
* **New Data Source:** `pingfederate_administrative_account` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_authentication_api_settings` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_certificate_ca` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_idp_default_urls` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_key_pair_signing_import` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_key_pair_ssl_server_import` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_license_agreement` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_license` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_local_identity_identity_profile` ([#70](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_oauth_access_token_manager` ([#76](https://github.com/pingidentity/terraform-provider-pingfederate/pull/70))
* **New Data Source:** `pingfederate_oauth_auth_server_settings` ([#78](https://github.com/pingidentity/terraform-provider-pingfederate/pull/78))
* **New Data Source:** `pingfederate_oauth_auth_server_settings_scopes_common_scope` ([#85](https://github.com/pingidentity/terraform-provider-pingfederate/pull/85))
* **New Data Source:** `pingfederate_oauth_auth_server_settings_scopes_exclusive_scope` ([#95](https://github.com/pingidentity/terraform-provider-pingfederate/pull/95))
* **New Data Source:** `pingfederate_oauth_issuer` ([#96](https://github.com/pingidentity/terraform-provider-pingfederate/pull/96))
* **New Data Source:** `pingfederate_virtual_host_names` ([#87](https://github.com/pingidentity/terraform-provider-pingfederate/pull/87))
* **New Data Source:** `pingfederate_session_application_session_policy` ([#94](https://github.com/pingidentity/terraform-provider-pingfederate/pull/94))

​
# v0.1.0 September 28, 2023
​
### FEATURES
* **New Resource:** `pingfederate_administrative_account` => [Initial Commit](https://github.com/pingidentity/terraform-provider-pingfederate/commit/fe35b53aac7146d2a75eeb70f4e21aaf52995a96)
* **New Resource:** `pingfederate_authentication_api_settings` ([#7](https://github.com/pingidentity/terraform-provider-pingfederate/pull/7))
* **New Resource:** `pingfederate_authentication_policy_contract` => [Initial Commit](https://github.com/pingidentity/terraform-provider-pingfederate/commit/fe35b53aac7146d2a75eeb70f4e21aaf52995a96)
* **New Resource:** `pingfederate_certificate_ca` ([#13](https://github.com/pingidentity/terraform-provider-pingfederate/pull/13))
* **New Resource:** `pingfederate_idp_default_urls` ([#8](https://github.com/pingidentity/terraform-provider-pingfederate/pull/8))
* **New Resource:** `pingfederate_key_pair_signing_import` => [Initial Commit](https://github.com/pingidentity/terraform-provider-pingfederate/commit/fe35b53aac7146d2a75eeb70f4e21aaf52995a96)
* **New Resource:** `pingfederate_key_pair_ssl_server_import` => [Initial Commit](https://github.com/pingidentity/terraform-provider-pingfederate/commit/fe35b53aac7146d2a75eeb70f4e21aaf52995a96)
* **New Resource:** `pingfederate_license` ([#13](https://github.com/pingidentity/terraform-provider-pingfederate/pull/13))
* **New Resource:** `pingfederate_license_agreement` ([#15](https://github.com/pingidentity/terraform-provider-pingfederate/pull/15))
* **New Resource:** `pingfederate_local_identity_profile` ([#38](https://github.com/pingidentity/terraform-provider-pingfederate/pull/38))
* **New Resource:** `pingfederate_oauth_access_token_manager` ([#55](https://github.com/pingidentity/terraform-provider-pingfederate/pull/55))
* **New Resource:** `pingfederate_oauth_auth_server_settings` ([#34](https://github.com/pingidentity/terraform-provider-pingfederate/pull/34))
* **New Resource:** `pingfederate_oauth_auth_server_settings_scopes_common_scope` => [Initial Commit](https://github.com/pingidentity/terraform-provider-pingfederate/commit/fe35b53aac7146d2a75eeb70f4e21aaf52995a96)
* **New Resource:** `pingfederate_oauth_auth_server_settings_scopes_exclusive_scope` => [Initial Commit](https://github.com/pingidentity/terraform-provider-pingfederate/commit/fe35b53aac7146d2a75eeb70f4e21aaf52995a96)
* **New Resource:** `pingfederate_oauth_issuer` => [Initial Commit](https://github.com/pingidentity/terraform-provider-pingfederate/commit/fe35b53aac7146d2a75eeb70f4e21aaf52995a96)
* **New Resource:** `pingfederate_password_credential_validator` ([#39](https://github.com/pingidentity/terraform-provider-pingfederate/pull/39))
* **New Resource:** `pingfederate_protocol_metadata_lifetime_settings` ([#9)](https://github.com/pingidentity/terraform-provider-pingfederate/pull/9)
* **New Resource:** `pingfederate_redirect_validation` ([#17](https://github.com/pingidentity/terraform-provider-pingfederate/pull/17))
* **New Resource:** `pingfederate_server_settings` ([#47](https://github.com/pingidentity/terraform-provider-pingfederate/pull/47))
* **New Resource:** `pingfederate_server_settings_general_settings` ([#10](https://github.com/pingidentity/terraform-provider-pingfederate/pull/10))
* **New Resource:** `pingfederate_server_settings_log_settings` ([#23](https://github.com/pingidentity/terraform-provider-pingfederate/pull/23))
* **New Resource:** `pingfederate_server_settings_system_keys` ([#30](https://github.com/pingidentity/terraform-provider-pingfederate/pull/30))
* **New Resource:** `pingfederate_session_application_session_policy` ([#11](https://github.com/pingidentity/terraform-provider-pingfederate/pull/11))
* **New Resource:** `pingfederate_session_authentication_session_policies_global` ([#12](https://github.com/pingidentity/terraform-provider-pingfederate/pull/12))
* **New Resource:** `pingfederate_session_settings` ([#16](https://github.com/pingidentity/terraform-provider-pingfederate/pull/16))
* **New Resource:** `pingfederate_virtual_host_names` ([#14](https://github.com/pingidentity/terraform-provider-pingfederate/pull/14))
