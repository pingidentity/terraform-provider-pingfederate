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
  registration_enabled = true
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
