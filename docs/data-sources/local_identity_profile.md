---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_local_identity_profile Data Source - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Describes a configured local identity profile.
---

# pingfederate_local_identity_profile (Data Source)

Describes a configured local identity profile.

## Example Usage

```terraform
data "pingfederate_local_identity_profile" "myLocalIdentityProfile" {
  profile_id = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `profile_id` (String) Unique ID for the local identity profile

### Read-Only

- `apc_id` (Attributes) The reference to the authentication policy contract to use for this local identity profile. (see [below for nested schema](#nestedatt--apc_id))
- `auth_source_update_policy` (Attributes) The attribute update policy for authentication sources. (see [below for nested schema](#nestedatt--auth_source_update_policy))
- `auth_sources` (Attributes Set) The local identity authentication sources. Sources are unique. (see [below for nested schema](#nestedatt--auth_sources))
- `data_store_config` (Attributes) The local identity profile data store configuration. (see [below for nested schema](#nestedatt--data_store_config))
- `email_verification_config` (Attributes) The local identity email verification configuration. (see [below for nested schema](#nestedatt--email_verification_config))
- `field_config` (Attributes) The local identity profile field configuration. (see [below for nested schema](#nestedatt--field_config))
- `id` (String) ID of this resource.
- `name` (String) The local identity profile name. Name is unique.
- `profile_config` (Attributes) The local identity profile management configuration. (see [below for nested schema](#nestedatt--profile_config))
- `profile_enabled` (Boolean) Whether the profile configuration is enabled or not.
- `registration_config` (Attributes) The local identity profile registration configuration. (see [below for nested schema](#nestedatt--registration_config))
- `registration_enabled` (Boolean) Whether the registration configuration is enabled or not.

<a id="nestedatt--apc_id"></a>
### Nested Schema for `apc_id`

Read-Only:

- `id` (String) The ID of the resource.


<a id="nestedatt--auth_source_update_policy"></a>
### Nested Schema for `auth_source_update_policy`

Read-Only:

- `retain_attributes` (Boolean) Whether or not to keep attributes after user disconnects.
- `store_attributes` (Boolean) Whether or not to store attributes that came from authentication sources.
- `update_attributes` (Boolean) Whether or not to update attributes when users authenticate.
- `update_interval` (Number) The minimum number of days between updates.


<a id="nestedatt--auth_sources"></a>
### Nested Schema for `auth_sources`

Read-Only:

- `id` (String) The persistent, unique ID for the local identity authentication source. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.
- `source` (String) The local identity authentication source. Source is unique.


<a id="nestedatt--data_store_config"></a>
### Nested Schema for `data_store_config`

Read-Only:

- `auxiliary_object_classes` (Set of String) The Auxiliary Object Classes used by the new objects stored in the LDAP data store.
- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `create_pattern` (String) The Relative DN Pattern that will be used to create objects in the directory.
- `data_store_mapping` (Attributes Map) The data store mapping. (see [below for nested schema](#nestedatt--data_store_config--data_store_mapping))
- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--data_store_config--data_store_ref))
- `object_class` (String) The Object Class used by the new objects stored in the LDAP data store.
- `type` (String) The data store config type.

<a id="nestedatt--data_store_config--data_store_mapping"></a>
### Nested Schema for `data_store_config.data_store_mapping`

Read-Only:

- `metadata` (Map of String) The data store attribute metadata.
- `name` (String) The data store attribute name.
- `type` (String) The data store attribute type.


<a id="nestedatt--data_store_config--data_store_ref"></a>
### Nested Schema for `data_store_config.data_store_ref`

Read-Only:

- `id` (String) The ID of the resource.



<a id="nestedatt--email_verification_config"></a>
### Nested Schema for `email_verification_config`

Read-Only:

- `allowed_otp_character_set` (String) The allowed character set used to generate the OTP. The default is 23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz. Note: Only applicable if EmailVerificationType is OTP.
- `email_verification_enabled` (Boolean) Whether the email ownership verification is enabled.
- `email_verification_error_template_name` (String) The template name for email verification error. The default is local.identity.email.verification.error.html.
- `email_verification_otp_template_name` (String) The template name for email verification OTP verification. The default is local.identity.email.verification.otp.html. Note: Only applicable if EmailVerificationType is OTP.
- `email_verification_sent_template_name` (String) The template name for email verification sent. The default is local.identity.email.verification.sent.html. Note:Only applicable if EmailVerificationType is OTL.
- `email_verification_success_template_name` (String) The template name for email verification success. The default is local.identity.email.verification.success.html.
- `email_verification_type` (String) Email Verification Type.
- `field_for_email_to_verify` (String) Field used for email ownership verification. Note: Not required when emailVerificationEnabled is set to false.
- `field_storing_verification_status` (String) Field used for storing email verification status. Note: Not required when emailVerificationEnabled is set to false.
- `notification_publisher_ref` (Attributes) Reference to the associated notification publisher. (see [below for nested schema](#nestedatt--email_verification_config--notification_publisher_ref))
- `otl_time_to_live` (Number) Field used OTL time to live. The default is 1440. Note: Only applicable if EmailVerificationType is OTL.
- `otp_length` (Number) The OTP length generated for email verification. The default is 8. Note: Only applicable if EmailVerificationType is OTP.
- `otp_retry_attempts` (Number) The number of OTP retry attempts for email verification. The default is 3. Note: Only applicable if EmailVerificationType is OTP.
- `otp_time_to_live` (Number) Field used OTP time to live. The default is 15. Note: Only applicable if EmailVerificationType is OTP.
- `require_verified_email` (Boolean) Whether the user must verify their email address before they can complete a single sign-on transaction. The default is false.
- `require_verified_email_template_name` (String) The template to render when the user must verify their email address before they can complete a single sign-on transaction. The default is local.identity.email.verification.required.html. Note:Only applicable if EmailVerificationType is OTL and requireVerifiedEmail is true.
- `verify_email_template_name` (String) The template name for verify email. The default is message-template-email-ownership-verification.html.

<a id="nestedatt--email_verification_config--notification_publisher_ref"></a>
### Nested Schema for `email_verification_config.notification_publisher_ref`

Read-Only:

- `id` (String) The ID of the resource.



<a id="nestedatt--field_config"></a>
### Nested Schema for `field_config`

Read-Only:

- `fields` (Attributes Set) The field configuration for the local identity profile. (see [below for nested schema](#nestedatt--field_config--fields))
- `strip_space_from_unique_field` (Boolean) Strip leading/trailing spaces from unique ID field. Default is true.

<a id="nestedatt--field_config--fields"></a>
### Nested Schema for `field_config.fields`

Read-Only:

- `attributes` (Map of Boolean) Attributes of the local identity field.
- `default_value` (String) The default value for this field.
- `id` (String) Id of the local identity field.
- `label` (String) Label of the local identity field.
- `options` (Set of String) The list of options for this selection field.
- `profile_page_field` (Boolean) Whether this is a profile page field or not.
- `registration_page_field` (Boolean) Whether this is a registration page field or not.
- `type` (String) The type of the local identity field.



<a id="nestedatt--profile_config"></a>
### Nested Schema for `profile_config`

Read-Only:

- `delete_identity_enabled` (Boolean) Whether the end user is allowed to use delete functionality.
- `template_name` (String) The template name for end-user profile management.


<a id="nestedatt--registration_config"></a>
### Nested Schema for `registration_config`

Read-Only:

- `captcha_enabled` (Boolean) Whether CAPTCHA is enabled or not in the registration configuration.
- `captcha_provider_ref` (Attributes) Reference to the associated CAPTCHA provider. (see [below for nested schema](#nestedatt--registration_config--captcha_provider_ref))
- `create_authn_session_after_registration` (Boolean) Whether to create an Authentication Session when registering a local account. Default is true.
- `execute_workflow` (String) This setting indicates whether PingFederate should execute the workflow before or after account creation. The default is to run the registration workflow after account creation.
- `registration_workflow` (Attributes) The policy fragment to be executed as part of the registration workflow. (see [below for nested schema](#nestedatt--registration_config--registration_workflow))
- `template_name` (String) The template name for the registration configuration.
- `this_is_my_device_enabled` (Boolean) Allows users to indicate whether their device is shared or private. In this mode, PingFederate Authentication Sessions will not be stored unless the user indicates the device is private.
- `username_field` (String) When creating an Authentication Session after registering a local account, PingFederate will pass the Unique ID field's value as the username. If the Unique ID value is not the username, then override which field's value will be used as the username.

<a id="nestedatt--registration_config--captcha_provider_ref"></a>
### Nested Schema for `registration_config.captcha_provider_ref`

Read-Only:

- `id` (String) The ID of the resource.


<a id="nestedatt--registration_config--registration_workflow"></a>
### Nested Schema for `registration_config.registration_workflow`

Read-Only:

- `id` (String) The ID of the resource.
