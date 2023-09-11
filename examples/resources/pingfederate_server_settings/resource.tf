resource "pingfederate_server_settings" "serverSettingsExample" {
  contact_info = {
    company    = "example company"
    email      = "adminemail@company.com"
    first_name = "Jane"
    last_name  = "Admin"
    phone      = "555-555-1222"
  }

  notifications = {
    license_events = {
      email_address = "license-events-email@company.com"
      notification_publisher_ref = {
        id = "<uiInstanceID>"
      }
    }
    certificate_expirations = {
      email_address          = "cert-expire-notifications@company.com"
      initial_warning_period = 45
      final_warning_period   = 7
      notification_publisher_ref = {
        id = "<uiInstanceID>"
      }
    }
    notify_admin_user_password_changes = true
    account_changes_notification_publisher_ref = {
      id = "<uiInstanceID>"
    }
    metadata_notification_settings = {
      email_address = "metadata-notification@company.com"
      notification_publisher_ref = {
        id = "<uiInstanceID>"
      }
    }


    federation_info = {
      // base_url must be standard URL format: http(s)://<company-or-hostname> with optional domain and port
      base_url = "https://localhost:9999"
      // SAML entities have to be defined first
      saml2_entity_id  = "urn:auth0:example:myserverconnection2"
      saml1x_issuer_id = "pingidentity2.com"
      //saml1x_source_id should be a hex if supplied.  Value can be empty string or not set at all.
      saml1x_source_id = ""
      wsfed_realm      = "myrealm2"
    }

    email_server = {
      source_addr  = "emailServerAdmin2@company.com"
      email_server = "myemailserver2.company.com"
      use_ssl      = true
      // cannot set both TLS and SSL at the same time.  SSL has priority
      //use_tls= true
      verify_hostname             = true
      enable_utf8_message_headers = true
      use_debugging               = false
      username                    = "emailServerAdmin"
      password                    = "emailServerAdminPassword"
    }

    captcha_settings = {
      site_key   = "mySiteKey"
      secret_key = "mySiteKeySecret"
    }
  }
}