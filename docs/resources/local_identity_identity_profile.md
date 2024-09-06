---
page_title: "pingfederate_local_identity_identity_profile Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages a configured local identity profile
---

# pingfederate_local_identity_identity_profile (Resource)

Manages a configured local identity profile

!> The `pingfederate_local_identity_identity_profile` resource has been renamed and will be removed in a future release. Use the `pingfederate_local_identity_profile` resource instead.

## Example Usage

```terraform
resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractsExample" {
  contract_id         = "myContract"
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "My Contract"
}

resource "pingfederate_notification_publisher" "notificationPublisher" {
  publisher_id = "EmailSMTPPublisherSettings"
  name         = "Email SMTP Publisher Settings"
  configuration = {
    fields = [
      {
        name  = "Email Server"
        value = "localhost"
      },
      {
        name  = "From Address"
        value = "noreply@bxretail.org"
      },
      {
        name  = "Sender Name"
        value = "BXRetail"
      },
      {
        name  = "SMTP Port"
        value = "25"
      },
      {
        name  = "Encryption Method"
        value = "SSL"
      },
      {
        name  = "SMTPS Port"
        value = "465"
      },
      {
        name  = "Username"
        value = var.email_smtp_server_username
      },
      {
        name  = "Password"
        value = var.email_smtp_server_password
      },
      {
        name  = "Verify Hostname"
        value = "true"
      },
      {
        name  = "UTF-8 Message Header Support"
        value = "false"
      },
      {
        name  = "Connection Timeout"
        value = "30"
      },
      {
        name  = "Retry Attempt"
        value = "2"
      },
      {
        name  = "Retry Delay"
        value = "2"
      },
      {
        name  = "Enable SMTP Debugging Messages"
        value = "true"
      }
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.email.SmtpNotificationPlugin"
  }
}

resource "pingfederate_local_identity_identity_profile" "identityProfileExample" {
  name       = "identityProfileName"
  profile_id = "profileId"
  apc_id = {
    id = pingfederate_authentication_policy_contract.authenticationPolicyContractsExample.contract_id
  }
  auth_sources = [
    {
      source = "test",
    },
    {
      source = "username",
    }
  ]
  auth_source_update_policy = {
    store_attributes  = false
    retain_attributes = false
    update_attributes = false
    update_interval   = 0
  }
  registration_enabled = false
  registration_config = {

    template_name                           = "local.identity.registration.html"
    create_authn_session_after_registration = true
    username_field                          = "cn"
    this_is_my_device_enabled               = false
    registration_workflow = {
      id = "registrationid",
    }
    execute_workflow = "AFTER_ACCOUNT_CREATION"
  }
  profile_config = {
    delete_identity_enabled = true
    template_name           = "local.identity.profile.html"
  }
  field_config = {
    fields = [
      {
        type                    = "EMAIL"
        id                      = "mail"
        label                   = "Email address"
        registration_page_field = true
        profile_page_field      = true
        attributes = {
          "Read-Only"       = false,
          "Required"        = true,
          "Unique ID Field" = true,
          "Mask Log Values" = false,
        }
      },
      {
        type                    = "TEXT"
        id                      = "cn"
        label                   = "First Name"
        registration_page_field = true
        profile_page_field      = true
        attributes = {
          "Read-Only"       = false,
          "Required"        = true,
          "Unique ID Field" = false,
          "Mask Log Values" = false,
        }
      },
      {
        type                    = "HIDDEN",
        id                      = "entryUUID",
        label                   = "entryUUID",
        registration_page_field = true
        profile_page_field      = true
        attributes = {
          "Unique ID Field" = false,
          "Mask Log Values" = false,
        }
      },
    ]
    strip_space_from_unique_field = true
  }
  email_verification_config = {
    email_verification_enabled               = true
    verify_email_template_name               = "message-template-email-ownership-verification.html"
    email_verification_success_template_name = "local.identity.email.verification.success.html"
    email_verification_error_template_name   = "local.identity.email.verification.error.html"
    email_verification_type                  = "OTP"
    allowed_otp_character_set                = "23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz"
    email_verification_otp_template_name     = "message-template-email-ownership-verification.html"
    otp_length                               = 8
    otp_retry_attempts                       = 3
    otp_time_to_live                         = 3
    field_for_email_to_verify                = "mail"
    field_storing_verification_status        = "entryUUID"
    notification_publisher_ref = {
      id = pingfederate_notification_publisher.notificationPublisher.publisher_id
    }
    require_verified_email = true
  }
  data_store_config = {
    type = "LDAP"
    data_store_ref = {
      id = "PDdatastore"
    }
    base_dn        = "ou=people,dc=example,dc=com",
    create_pattern = "uid=$${mail}",
    object_class   = "inetOrgPerson",
    data_store_mapping = {
      "entryUUID" = {
        type     = "LDAP"
        name     = "entryUUID"
        metadata = {}

      }
      "cn" = {
        type     = "LDAP"
        name     = "cn"
        metadata = {}
      },
      "mail" = {
        type     = "LDAP"
        name     = "mail"
        metadata = {}
      },
    }
  }
  profile_enabled = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `apc_id` (Attributes) The reference to the authentication policy contract to use for this local identity profile. (see [below for nested schema](#nestedatt--apc_id))
- `name` (String) The local identity profile name. Name is unique.
- `profile_id` (String) The persistent, unique ID for the local identity profile. It can be any combination of `[a-zA-Z0-9._-]`.

### Optional

- `auth_source_update_policy` (Attributes) The attribute update policy for authentication sources. (see [below for nested schema](#nestedatt--auth_source_update_policy))
- `auth_sources` (Attributes Set) The local identity authentication sources. Sources are unique. (see [below for nested schema](#nestedatt--auth_sources))
- `data_store_config` (Attributes) The local identity profile data store configuration. (see [below for nested schema](#nestedatt--data_store_config))
- `email_verification_config` (Attributes) The local identity email verification configuration. (see [below for nested schema](#nestedatt--email_verification_config))
- `field_config` (Attributes) The local identity profile field configuration. (see [below for nested schema](#nestedatt--field_config))
- `profile_config` (Attributes) The local identity profile management configuration. (see [below for nested schema](#nestedatt--profile_config))
- `profile_enabled` (Boolean) Whether the profile configuration is enabled or not. The default value is `false`.
- `registration_config` (Attributes) The local identity profile registration configuration. (see [below for nested schema](#nestedatt--registration_config))
- `registration_enabled` (Boolean) Whether the registration configuration is enabled or not. The default value is `false`.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--apc_id"></a>
### Nested Schema for `apc_id`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--auth_source_update_policy"></a>
### Nested Schema for `auth_source_update_policy`

Optional:

- `retain_attributes` (Boolean) Whether or not to keep attributes after user disconnects. The default value is `false`.
- `store_attributes` (Boolean) Whether or not to store attributes that came from authentication sources. The default value is `false`.
- `update_attributes` (Boolean) Whether or not to update attributes when users authenticate. The default value is `false`.
- `update_interval` (Number) The minimum number of days between updates. The default value is `0`.


<a id="nestedatt--auth_sources"></a>
### Nested Schema for `auth_sources`

Required:

- `source` (String) The local identity authentication source. Source is unique.

Optional:

- `id` (String) The persistent, unique ID for the local identity authentication source. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified.


<a id="nestedatt--data_store_config"></a>
### Nested Schema for `data_store_config`

Required:

- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `data_store_mapping` (Attributes Map) The data store mapping. (see [below for nested schema](#nestedatt--data_store_config--data_store_mapping))
- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--data_store_config--data_store_ref))
- `type` (String) The data store config type. Supported values are `LDAP`, `PING_ONE_LDAP_GATEWAY`, `JDBC`, and `CUSTOM`.

Optional:

- `auxiliary_object_classes` (Set of String) The Auxiliary Object Classes used by the new objects stored in the LDAP data store.
- `create_pattern` (String) The Relative DN Pattern that will be used to create objects in the directory.
- `object_class` (String) The Object Class used by the new objects stored in the LDAP data store.

<a id="nestedatt--data_store_config--data_store_mapping"></a>
### Nested Schema for `data_store_config.data_store_mapping`

Required:

- `name` (String) The data store attribute name.
- `type` (String) The data store attribute type. Supported values are `LDAP`, `PING_ONE_LDAP_GATEWAY`, `JDBC`, and `CUSTOM`.

Optional:

- `metadata` (Map of String) The data store attribute metadata.


<a id="nestedatt--data_store_config--data_store_ref"></a>
### Nested Schema for `data_store_config.data_store_ref`

Required:

- `id` (String) The ID of the resource.



<a id="nestedatt--email_verification_config"></a>
### Nested Schema for `email_verification_config`

Optional:

- `allowed_otp_character_set` (String) The allowed character set used to generate the OTP. The default is `23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz`. Note: Only applicable if `email_verification_type` is `OTP`.
- `email_verification_enabled` (Boolean) Whether the email ownership verification is enabled. The default value is `false`.
- `email_verification_error_template_name` (String) The template name for email verification error. The default is `local.identity.email.verification.error.html`.
- `email_verification_otp_template_name` (String) The template name for email verification OTP verification. The default is `local.identity.email.verification.otp.html`. Note: Only applicable if `email_verification_type` is `OTP`.
- `email_verification_sent_template_name` (String) The template name for email verification sent. The default is `local.identity.email.verification.sent.html`. Note:Only applicable if `email_verification_type` is `OTL`.
- `email_verification_success_template_name` (String) The template name for email verification success. The default is `local.identity.email.verification.success.html`.
- `email_verification_type` (String) Email Verification Type. Supported values are `OTP` and `OTL`.
- `field_for_email_to_verify` (String) Field used for email ownership verification. Note: Not required when `email_verification_enabled` is set to `false`.
- `field_storing_verification_status` (String) Field used for storing email verification status. Note: Not required when `email_verification_enabled` is set to `false`.
- `notification_publisher_ref` (Attributes) Reference to the associated notification publisher. (see [below for nested schema](#nestedatt--email_verification_config--notification_publisher_ref))
- `otl_time_to_live` (Number) Field used OTL time to live. The default is `1440`. Note: Only applicable if `email_verification_type` is `OTL`.
- `otp_length` (Number) The OTP length generated for email verification. The default is `8`. Note: Only applicable if `email_verification_type` is `OTP`. The value must be between `5` and `100`.
- `otp_retry_attempts` (Number) The number of OTP retry attempts for email verification. The default is `3`. Note: Only applicable if `email_verification_type` is `OTP`.
- `otp_time_to_live` (Number) Field used OTP time to live. The default is `15`. Note: Only applicable if `email_verification_type` is `OTP`.
- `require_verified_email` (Boolean) Whether the user must verify their email address before they can complete a single sign-on transaction. The default is `false`.
- `require_verified_email_template_name` (String) The template to render when the user must verify their email address before they can complete a single sign-on transaction. The default is `local.identity.email.verification.required.html`. Note: Only applicable if `email_verification_type` is OTL and `require_verified_email` is true.
- `verify_email_template_name` (String) The template name for verify email. The default is `message-template-email-ownership-verification.html`.

<a id="nestedatt--email_verification_config--notification_publisher_ref"></a>
### Nested Schema for `email_verification_config.notification_publisher_ref`

Required:

- `id` (String) The ID of the resource.



<a id="nestedatt--field_config"></a>
### Nested Schema for `field_config`

Optional:

- `fields` (Attributes Set) The field configuration for the local identity profile. (see [below for nested schema](#nestedatt--field_config--fields))
- `strip_space_from_unique_field` (Boolean) Strip leading/trailing spaces from unique ID field. The default value is `false`.

<a id="nestedatt--field_config--fields"></a>
### Nested Schema for `field_config.fields`

Required:

- `id` (String) Id of the local identity field.
- `label` (String) Label of the local identity field.
- `type` (String) The type of the local identity field. Supported values are `CHECKBOX`, `CHECKBOX_GROUP`, `DATE`, `DROP_DOWN`, `EMAIL`, `PHONE`, `TEXT`, and `HIDDEN`.

Optional:

- `attributes` (Map of Boolean) Attributes of the local identity field.
- `profile_page_field` (Boolean) Whether this is a profile page field or not. The default value is `false`.
- `registration_page_field` (Boolean) Whether this is a registration page field or not. The default value is `false`.



<a id="nestedatt--profile_config"></a>
### Nested Schema for `profile_config`

Required:

- `template_name` (String) The template name for end-user profile management.

Optional:

- `delete_identity_enabled` (Boolean) Whether the end user is allowed to use delete functionality. The default value is `false`.


<a id="nestedatt--registration_config"></a>
### Nested Schema for `registration_config`

Required:

- `template_name` (String) The template name for the registration configuration.

Optional:

- `captcha_enabled` (Boolean) Whether CAPTCHA is enabled or not in the registration configuration. The default value is `false`.
- `captcha_provider_ref` (Attributes) Reference to the associated CAPTCHA provider. (see [below for nested schema](#nestedatt--registration_config--captcha_provider_ref))
- `create_authn_session_after_registration` (Boolean) Whether to create an Authentication Session when registering a local account. The default value is `true`.
- `execute_workflow` (String) This setting indicates whether PingFederate should execute the workflow before or after account creation. The default is to run the registration workflow after account creation. Supported values are `BEFORE_ACCOUNT_CREATION` and `AFTER_ACCOUNT_CREATION`. Requires that `registration_workflow` is also set.
- `registration_workflow` (Attributes) The policy fragment to be executed as part of the registration workflow. (see [below for nested schema](#nestedatt--registration_config--registration_workflow))
- `this_is_my_device_enabled` (Boolean) Allows users to indicate whether their device is shared or private. In this mode, PingFederate Authentication Sessions will not be stored unless the user indicates the device is private. The default value is `false`.
- `username_field` (String) When creating an Authentication Session after registering a local account, PingFederate will pass the Unique ID field's value as the username. If the Unique ID value is not the username, then override which field's value will be used as the username.

<a id="nestedatt--registration_config--captcha_provider_ref"></a>
### Nested Schema for `registration_config.captcha_provider_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--registration_config--registration_workflow"></a>
### Nested Schema for `registration_config.registration_workflow`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "profileId" should be the id of the Local Identity Profile to be imported

```shell
terraform import pingfederate_local_identity_identity_profile.identityProfileExample profileId
```