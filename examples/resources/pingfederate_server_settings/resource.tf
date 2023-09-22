resource "pingfederate_server_settings" "serverSettingsExample" {
  contact_info = {
    company    = "example company"
    email      = "adminemail@example.com"
    first_name = "Jane"
    last_name  = "Admin"
    phone      = "555-555-1222"
  }

  notifications = {
    license_events = {
      email_address = "license-events-email@example.com"
      notification_publisher_ref = {
        id = "<uiInstanceID>"
      }
    }
    certificate_expirations = {
      email_address          = "cert-expire-notifications@example.com"
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
      email_address = "metadata-notification@example.com"
      notification_publisher_ref = {
        id = "<uiInstanceID>"
      }
    }


    federation_info = {
      // base_url must be standard URL format: http(s)://<company-or-hostname> with optional domain and port
      base_url = "https://localhost:9999"
      // SAML entities have to be defined first
      saml_2_entity_id  = "urn:auth0:example:myserverconnection"
      saml_1x_issuer_id = "example.com"
      //saml_1x_source_id should be a hex if supplied.  Value can be empty string or not set at all.
      saml_1x_source_id = ""
      wsfed_realm       = "myrealm"
    }

    email_server = {
      source_addr                 = "emailServerAdmin@example.com"
      email_server                = "myemailserver.example.com"
      use_ssl                     = true
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
