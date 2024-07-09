resource "pingfederate_notification_publisher" "notificationPublisher" {
  publisher_id = "MigratedSmtpSettings"
  name         = "Migrated SMTP Settings"
  configuration = {
    fields = [
      {
        name  = "From Address"
        value = "EmailVerificationRuntimeTests@fakesmtpserver.com"
      },
      {
        name  = "Email Server"
        value = "localhost"
      },
      {
        name  = "SMTP Port"
        value = "2525"
      },
      {
        name  = "Encryption Method"
        value = "NONE"
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