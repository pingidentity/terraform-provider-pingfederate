resource "pingfederate_local_identity_identity_profile" "myLocalIdentityIdentityProfile" {
  name = "yourIdentityProfileName"
  #id   = "yourid"
  apc_id = {
    id = "apcid"
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
    captcha_enabled = true
    captcha_provider_ref = {
      id = "testCaptchaAndRisk"
    }
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
    email_verification_enabled = true
    verify_email_template_name = "message-template-email-ownership-verification.html"
    /* email_verification_sent_template_name = "local.identity.email.verification.sent.html"  */
    email_verification_success_template_name = "local.identity.email.verification.success.html"
    email_verification_error_template_name   = "local.identity.email.verification.error.html"
    /* TO ENABLE OTL as the email verification type, remove verification OTP template,otp_character_set, otp_retry attempts, otp_length attribute  and uncomment otl_time_to_live,email_verification_sent and require_verified email template*/
    email_verification_type              = "OTP"
    allowed_otp_character_set            = "23456789BCDFGHJKMNPQRSTVWXZbcdfghjkmnpqrstvwxz"
    email_verification_otp_template_name = "message-template-email-ownership-verification.html"
    otp_length                           = 8
    otp_retry_attempts                   = 3
    otp_time_to_live                     = 3
    /* otl_time_to_live = 1440  */
    field_for_email_to_verify         = "mail"
    field_storing_verification_status = "entryUUID"
    notification_publisher_ref = {
      id = "testnp",
    }
    require_verified_email = true
    /* require_verified_email_template_name = "local.identity.email.verification.required.html" */
  }
  // Local Identity profile only support Directory as DataStore
  data_store_config = {
    type = "LDAP"
    data_store_ref = {
      id = "directoryid"
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
