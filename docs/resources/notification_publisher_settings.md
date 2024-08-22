---
page_title: "pingfederate_notification_publisher_settings Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages Notification Publisher Settings
---

# pingfederate_notification_publisher_settings (Resource)

Manages Notification Publisher Settings

## Example Usage

```terraform
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

resource "pingfederate_notification_publisher_settings" "publisherSettings" {
  default_notification_publisher_ref = {
    id = pingfederate_notification_publisher.notificationPublisher.publisher_id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `default_notification_publisher_ref` (Attributes) The default notification publisher reference (see [below for nested schema](#nestedatt--default_notification_publisher_ref))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--default_notification_publisher_ref"></a>
### Nested Schema for `default_notification_publisher_ref`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> This resource is singleton, so the value of "id" doesn't matter - it is just a placeholder, and required by Terraform

```shell
terraform import pingfederate_notification_publisher_settings.publisherSettings id
```