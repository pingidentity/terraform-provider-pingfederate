---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_server_settings Data Source - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Describes the global server configuration settings
---

# pingfederate_server_settings (Data Source)

Describes the global server configuration settings

## Example Usage

```terraform
data "pingfederate_server_settings" "myServerSettingsExample" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `captcha_settings` (Attributes) Captcha Settings. (see [below for nested schema](#nestedatt--captcha_settings))
- `contact_info` (Attributes) Information that identifies the server. (see [below for nested schema](#nestedatt--contact_info))
- `email_server` (Attributes) Email Server Settings. (see [below for nested schema](#nestedatt--email_server))
- `federation_info` (Attributes) Federation Info. (see [below for nested schema](#nestedatt--federation_info))
- `id` (String) ID of this resource.
- `notifications` (Attributes) Notification settings for license and certificate expiration events. (see [below for nested schema](#nestedatt--notifications))
- `roles_and_protocols` (Attributes) Configure roles and protocols. (see [below for nested schema](#nestedatt--roles_and_protocols))

<a id="nestedatt--captcha_settings"></a>
### Nested Schema for `captcha_settings`

Read-Only:

- `secret_key` (String) Secret key for reCAPTCHA. GETs will not return this attribute. To update this field, specify the new value in this attribute.
- `site_key` (String) Site key for reCAPTCHA.


<a id="nestedatt--contact_info"></a>
### Nested Schema for `contact_info`

Read-Only:

- `company` (String) Company name.
- `email` (String) Contact email address.
- `first_name` (String) Contact first name.
- `last_name` (String) Contact last name.
- `phone` (String) Contact phone number.


<a id="nestedatt--email_server"></a>
### Nested Schema for `email_server`

Read-Only:

- `email_server` (String) The IP address or hostname of your email server.
- `enable_utf8_message_headers` (Boolean) Only set this flag to true if the email server supports UTF-8 characters in message headers. Otherwise, this is defaulted to false.
- `password` (String, Sensitive) User password. To update the password, specify the plaintext value in this field. This field will not be populated for GET requests.
- `port` (Number) The SMTP port on your email server. Allowable values: 1 - 65535.
- `retry_attempts` (Number) The number of times PingFederate tries to resend an email upon unsuccessful delivery.
- `retry_delay` (Number) The number of minutes PingFederate waits before the next retry attempt.
- `source_addr` (String) The email address that appears in the 'From' header line in email messages generated by PingFederate. The address must be in valid format but need not be set up on your system.
- `ssl_port` (Number) The secure SMTP port on your email server. This field is not active unless Use SSL is enabled. Allowable values: 1 - 65535.
- `timeout` (Number) The amount of time in seconds that PingFederate will wait before it times out connecting to the SMTP server. Allowable values: 0 - 3600.
- `use_debugging` (Boolean) Turns on detailed error messages for the PingFederate server log to help troubleshoot any problems.
- `use_ssl` (Boolean) Requires the use of SSL/TLS on the port specified by 'sslPort'. If this option is enabled, it overrides the 'useTLS' option.
- `use_tls` (Boolean) Requires the use of the STARTTLS protocol on the port specified by 'port'.
- `username` (String) Authorized email username. Required if the password is provided.
- `verify_hostname` (Boolean) If useSSL or useTLS is enabled, this flag determines whether the email server hostname is verified against the server's SMTPS certificate.


<a id="nestedatt--federation_info"></a>
### Nested Schema for `federation_info`

Read-Only:

- `auto_connect_entity_id` (String) This property has been deprecated and no longer used
- `base_url` (String) The fully qualified host name, port, and path (if applicable) on which the PingFederate server runs.
- `saml_1x_issuer_id` (String) This ID identifies your federation server for SAML 1.x transactions. As with SAML 2.0, it is usually defined as an organization's URL or a DNS address. The SourceID used for artifact resolution is derived from this ID using SHA1.
- `saml_1x_source_id` (String) If supplied, the Source ID value entered here is used for SAML 1.x, instead of being derived from the SAML 1.x Issuer/Audience.
- `saml_2_entity_id` (String) This ID defines your organization as the entity operating the server for SAML 2.0 transactions. It is usually defined as an organization's URL or a DNS address; for example: pingidentity.com. The SAML SourceID used for artifact resolution is derived from this ID using SHA1.
- `wsfed_realm` (String) The URI of the realm associated with the PingFederate server. A realm represents a single unit of security administration or trust.


<a id="nestedatt--notifications"></a>
### Nested Schema for `notifications`

Optional:

- `expired_certificate_administrative_console_warning_days` (Number) Indicates the number of days prior to certificate expiry date, the administrative console warning starts. The default value is 14 days. Supported in PF 12.0 or later.

Read-Only:

- `account_changes_notification_publisher_ref` (Attributes) Reference to the associated notification publisher for admin user account changes. (see [below for nested schema](#nestedatt--notifications--account_changes_notification_publisher_ref))
- `bulkhead_alert_notification_settings` (Attributes) Settings for bulkhead notifications (see [below for nested schema](#nestedatt--notifications--bulkhead_alert_notification_settings))
- `certificate_expirations` (Attributes) Notification settings for certificate expiration events. (see [below for nested schema](#nestedatt--notifications--certificate_expirations))
- `expiring_certificate_administrative_console_warning_days` (Number) Indicates the number of days past the certificate expiry date, the administrative console warning ends. The default value is 14 days. Supported in PF 12.0 or later.
- `license_events` (Attributes) Settings for license event notifications. (see [below for nested schema](#nestedatt--notifications--license_events))
- `metadata_notification_settings` (Attributes) Settings for metadata update event notifications. (see [below for nested schema](#nestedatt--notifications--metadata_notification_settings))
- `notify_admin_user_password_changes` (Boolean) Determines whether admin users are notified through email when their account is changed.
- `thread_pool_exhaustion_notification_settings` (Attributes) Notification settings for thread pool exhaustion events. Supported in PF 12.0 or later. (see [below for nested schema](#nestedatt--notifications--thread_pool_exhaustion_notification_settings))

<a id="nestedatt--notifications--account_changes_notification_publisher_ref"></a>
### Nested Schema for `notifications.account_changes_notification_publisher_ref`

Read-Only:

- `id` (String) The ID of the resource.


<a id="nestedatt--notifications--bulkhead_alert_notification_settings"></a>
### Nested Schema for `notifications.bulkhead_alert_notification_settings`

Read-Only:

- `email_address` (String) Email address where notifications are sent.
- `notification_mode` (String) The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to LOGGING_ONLY.
- `notification_publisher_ref` (Attributes) Reference to the associated notification publisher. (see [below for nested schema](#nestedatt--notifications--bulkhead_alert_notification_settings--notification_publisher_ref))
- `thread_dump_enabled` (Boolean) Generate a thread dump when a bulkhead reaches its warning threshold or is full.

<a id="nestedatt--notifications--bulkhead_alert_notification_settings--notification_publisher_ref"></a>
### Nested Schema for `notifications.bulkhead_alert_notification_settings.notification_publisher_ref`

Read-Only:

- `id` (String) The ID of the resource.



<a id="nestedatt--notifications--certificate_expirations"></a>
### Nested Schema for `notifications.certificate_expirations`

Read-Only:

- `email_address` (String) The email address where notifications are sent.
- `final_warning_period` (Number) Time before certificate expiration when final warning is sent (in days).
- `initial_warning_period` (Number) Time before certificate expiration when initial warning is sent (in days).
- `notification_mode` (String) The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to NOTIFICATION_PUBLISHER. Supported in PF version 11.3 or later.
- `notification_publisher_ref` (Attributes) Reference to the associated notification publisher. (see [below for nested schema](#nestedatt--notifications--certificate_expirations--notification_publisher_ref))

<a id="nestedatt--notifications--certificate_expirations--notification_publisher_ref"></a>
### Nested Schema for `notifications.certificate_expirations.notification_publisher_ref`

Read-Only:

- `id` (String) The ID of the resource.



<a id="nestedatt--notifications--license_events"></a>
### Nested Schema for `notifications.license_events`

Read-Only:

- `email_address` (String) The email address where notifications are sent.
- `notification_publisher_ref` (Attributes) Reference to the associated notification publisher. (see [below for nested schema](#nestedatt--notifications--license_events--notification_publisher_ref))

<a id="nestedatt--notifications--license_events--notification_publisher_ref"></a>
### Nested Schema for `notifications.license_events.notification_publisher_ref`

Read-Only:

- `id` (String) The ID of the resource.



<a id="nestedatt--notifications--metadata_notification_settings"></a>
### Nested Schema for `notifications.metadata_notification_settings`

Read-Only:

- `email_address` (String) The email address where notifications are sent.
- `notification_publisher_ref` (Attributes) Reference to the associated notification publisher. (see [below for nested schema](#nestedatt--notifications--metadata_notification_settings--notification_publisher_ref))

<a id="nestedatt--notifications--metadata_notification_settings--notification_publisher_ref"></a>
### Nested Schema for `notifications.metadata_notification_settings.notification_publisher_ref`

Read-Only:

- `id` (String) The ID of the resource.



<a id="nestedatt--notifications--thread_pool_exhaustion_notification_settings"></a>
### Nested Schema for `notifications.thread_pool_exhaustion_notification_settings`

Read-Only:

- `email_address` (String) Email address where notifications are sent.
- `notification_mode` (String) The mode of notification. Set to NOTIFICATION_PUBLISHER to enable email notifications and server log messages. Set to LOGGING_ONLY to enable server log messages. Defaults to LOGGING_ONLY.
- `notification_publisher_ref` (Attributes) Reference to the associated notification publisher. (see [below for nested schema](#nestedatt--notifications--thread_pool_exhaustion_notification_settings--notification_publisher_ref))
- `thread_dump_enabled` (Boolean) Generate a thread dump when approaching thread pool exhaustion.

<a id="nestedatt--notifications--thread_pool_exhaustion_notification_settings--notification_publisher_ref"></a>
### Nested Schema for `notifications.thread_pool_exhaustion_notification_settings.notification_publisher_ref`

Read-Only:

- `id` (String) The ID of the resource.




<a id="nestedatt--roles_and_protocols"></a>
### Nested Schema for `roles_and_protocols`

Read-Only:

- `enable_idp_discovery` (Boolean) Enable IdP Discovery.
- `idp_role` (Attributes) Identity Provider (IdP) settings. (see [below for nested schema](#nestedatt--roles_and_protocols--idp_role))
- `oauth_role` (Attributes) OAuth role settings. (see [below for nested schema](#nestedatt--roles_and_protocols--oauth_role))
- `sp_role` (Attributes) Service Provider (SP) settings. (see [below for nested schema](#nestedatt--roles_and_protocols--sp_role))

<a id="nestedatt--roles_and_protocols--idp_role"></a>
### Nested Schema for `roles_and_protocols.idp_role`

Read-Only:

- `enable` (Boolean) Enable Identity Provider Role.
- `enable_outbound_provisioning` (Boolean) Enable Outbound Provisioning.
- `enable_saml_1_0` (Boolean) Enable SAML 1.0.
- `enable_saml_1_1` (Boolean) Enable SAML 1.1.
- `enable_ws_fed` (Boolean) Enable WS Federation.
- `enable_ws_trust` (Boolean) Enable WS Trust.
- `saml_2_0_profile` (Attributes) SAML 2.0 Profile settings. (see [below for nested schema](#nestedatt--roles_and_protocols--idp_role--saml_2_0_profile))

<a id="nestedatt--roles_and_protocols--idp_role--saml_2_0_profile"></a>
### Nested Schema for `roles_and_protocols.idp_role.saml_2_0_profile`

Read-Only:

- `enable` (Boolean) Enable SAML2.0 profile.
- `enable_auto_connect` (Boolean) This property has been deprecated and no longer used.



<a id="nestedatt--roles_and_protocols--oauth_role"></a>
### Nested Schema for `roles_and_protocols.oauth_role`

Read-Only:

- `enable_oauth` (Boolean) Enable OAuth 2.0 Authorization Server (AS) Role.
- `enable_open_id_connect` (Boolean) Enable Open ID Connect.


<a id="nestedatt--roles_and_protocols--sp_role"></a>
### Nested Schema for `roles_and_protocols.sp_role`

Read-Only:

- `enable` (Boolean) Enable Service Provider Role.
- `enable_inbound_provisioning` (Boolean) Enable Inbound Provisioning.
- `enable_open_id_connect` (Boolean) Enable OpenID Connect.
- `enable_saml_1_0` (Boolean) Enable SAML 1.0.
- `enable_saml_1_1` (Boolean) Enable SAML 1.1.
- `enable_ws_fed` (Boolean) Enable WS Federation.
- `enable_ws_trust` (Boolean) Enable WS Trust.
- `saml_2_0_profile` (Attributes) SAML 2.0 Profile settings. (see [below for nested schema](#nestedatt--roles_and_protocols--sp_role--saml_2_0_profile))

<a id="nestedatt--roles_and_protocols--sp_role--saml_2_0_profile"></a>
### Nested Schema for `roles_and_protocols.sp_role.saml_2_0_profile`

Read-Only:

- `enable` (Boolean) Enable SAML2.0 profile.
- `enable_auto_connect` (Boolean) This property has been deprecated and no longer used.
- `enable_xasp` (Boolean) Enable Attribute Requester Mapping for X.509 Attribute Sharing Profile (XASP)
