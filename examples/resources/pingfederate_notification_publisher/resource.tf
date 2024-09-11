- do we just warn, or do we error:
            - if a field in non-sensitive-fields comes back from PF as encrypted
              - warn
            - if a field in sensitive-fields comes back from PF as non-encrypted
              - warn
            - if a field name we expect to be sensitive is placed by the user in non-sensitive-fields
              - leave out for now


















terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 1.0.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username   = "administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:9999"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the server.
  insecure_trust_all_tls = true
  x_bypass_external_validation_header = true
  product_version = "12.1"
}


resource "pingfederate_notification_publisher" "notificationPublisher" {
  publisher_id = "EmailSMTPPublisherSettings"
  name         = "Email SMTP Publisher Settings"
  configuration = {
    sensitive_fields = [
      {
        name  = "Password"
        value = "asdFs"
      },
    ]
    fields = [
      {
        name  = "Email Server"
        value = "localhosts"
      },
      {
        name  = "Username"
        value = "asdf"
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