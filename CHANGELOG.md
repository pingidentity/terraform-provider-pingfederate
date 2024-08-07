# v0.14.0 (Unreleased)
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
