resource "pingfederate_server_settings" "serverSettings" {
  contact_info = {
    company    = "BXRetail.org"
    email      = "authadmin@bxretail.org"
    first_name = "Jane"
    last_name  = "Admin"
    phone      = "555-555-1222"
  }

  notifications = {
    license_events = {
      email_address = "license-events-email@bxretail.org"
      notification_publisher_ref = {
        id = pingfederate_notification_publisher.license_publisher.publisher_id
      }
    }
    certificate_expirations = {
      email_address          = "cert-expire-notifications@bxretail.org"
      initial_warning_period = 45
      final_warning_period   = 7
      notification_publisher_ref = {
        id = pingfederate_notification_publisher.cert_publisher.publisher_id
      }
    }
    notify_admin_user_password_changes = true
    account_changes_notification_publisher_ref = {
      id = pingfederate_notification_publisher.account_changes_publisher.publisher_id
    }
    metadata_notification_settings = {
      email_address = "metadata-notification@bxretail.org"
      notification_publisher_ref = {
        id = pingfederate_notification_publisher.metadata_publisher.publisher_id
      }
    }
  }

  federation_info = {
    base_url = "https://auth.bxretail.org"

    // SAML entities have to be defined first
    saml_2_entity_id  = "org:bxretail:auth"
    saml_1x_issuer_id = "auth.bxretail.org"

    wsfed_realm = "realm"
  }
}
