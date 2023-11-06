# v0.2.0 (Unreleased)
### BUG FIXES
* **Administrative Account Resource:** Resolved issue where updating a managed `pingfederate_administrative_account` resource forces replacment ([#86](https://github.com/pingidentity/terraform-provider-pingfederate/pull/86)) 
### FEATURES
#### Resources
* **New Resource:** `pingfederate_idp_adapter` ([#64](https://github.com/pingidentity/terraform-provider-pingfederate/pull/64))
* **New Resource:** `pingfederate_token_processor_to_token_generator_mapping` ([#65](https://github.com/pingidentity/terraform-provider-pingfederate/pull/65))
* **New Resource:** `pingfederate_oauth_token_exchange_processor_policy_token_generator_mapping` ([#68](https://github.com/pingidentity/terraform-provider-pingfederate/pull/68))
* **New Resource:** `pingfederate_sp_authentication_policy_contract_mapping` ([#75](https://github.com/pingidentity/terraform-provider-pingfederate/pull/75))

​
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
* **New Data Source:** `pingfederate_oauth_access_token_manager` ([#76](https://github.com/pingidentity/terraform-provider-pingfederate/pull/76))
* **New Data Source:** `pingfederate_oauth_auth_server_settings` ([#78](https://github.com/pingidentity/terraform-provider-pingfederate/pull/78))
* **New Data Source:** `pingfederate_oauth_auth_server_settings_scopes_common_scope` ([#85](https://github.com/pingidentity/terraform-provider-pingfederate/pull/85))
* **New Data Source:** `pingfederate_oauth_auth_server_settings_scopes_exclusive_scope` ([#95](https://github.com/pingidentity/terraform-provider-pingfederate/pull/95))

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