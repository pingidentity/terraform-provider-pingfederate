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
    sensitive_fields = [
      {
        name  = "Password"
        value = var.email_smtp_server_password
      }
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.email.SmtpNotificationPlugin"
  }
}

resource "pingfederate_notification_publisher_settings" "publisherSettings" {
  default_notification_publisher_ref = {
    id = pingfederate_notification_publisher.notificationPublisher.publisher_id
  }
}