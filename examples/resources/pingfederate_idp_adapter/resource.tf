resource "pingfederate_idp_adapter" "idpAdapterExample" {
  adapter_id = "HTMLFormPD"
  name       = "HTMLFormPD"
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.htmlform.idp.HtmlFormIdpAuthnAdapter"
  }
  attribute_mapping = {
    attribute_contract_fulfillment = {
      "entryUUID" = {
        source = {
          type = "ADAPTER"
        },
        value = "entryUUID"
      }
      "policy.action" = {
        source = {
          type = "ADAPTER"
        },
        value = "policy.action"
      },
      "username" = {
        source = {
          type = "ADAPTER"
        },
        value = "username"
      }
    },
    attribute_sources = [],
    issuance_criteria = {
      conditional_criteria = []
    }
  }
  configuration = {
    tables = [
      {
        name = "Credential Validators"
        rows = [
          {
            default_row = false
            fields = [
              {
                name  = "Password Credential Validator Instance"
                value = "pingdirectory"
              }
            ]
          }
        ]
      }
    ]
    fields = [
      {
        name  = "Challenge Retries"
        value = 3
      },
      {
        name  = "Session State",
        value = "None"
      },
      {
        name  = "Session Timeout",
        value = "60"
      },
      {
        name  = "Session Max Timeout",
        value = "480"
      },
      {
        name  = "Allow Password Changes",
        value = "false"
      },
      {
        name  = "Password Management System",
        value = ""
      },
      {
        name  = "Enable 'Remember My Username'",
        value = "false"
      },
      {
        name  = "Enable 'This is My Device'",
        value = "false"
      },
      {
        name  = "Change Password Email Notification",
        value = "false"
      },
      {
        name  = "Show Password Expiring Warning",
        value = "false"
      },
      {
        name  = "Password Reset Type",
        value = "NONE"
      },
      {
        name  = "Account Unlock",
        value = "false"
      },
      {
        name  = "Local Identity Profile",
        value = "example"
      },
      {
        name  = "Enable Username Recovery",
        value = "false"
      },
      {
        name  = "Login Template",
        value = "html.form.login.template.html"
      },
      {
        name  = "Logout Path",
        value = ""
      },
      {
        name  = "Logout Redirect",
        value = ""
      },
      {
        name  = "Logout Template",
        value = "idp.logout.success.page.template.html"
      },
      {
        name  = "Change Password Template",
        value = "html.form.change.password.template.html"
      },
      {
        name  = "Change Password Message Template",
        value = "html.form.message.template.html"
      },
      {
        name  = "Password Management System Message Template",
        value = "html.form.message.template.html"
      },
      {
        name  = "Change Password Email Template",
        value = "message-template-end-user-password-change.html"
      },
      {
        name  = "Expiring Password Warning Template",
        value = "html.form.password.expiring.notification.template.html"
      },
      {
        name  = "Threshold for Expiring Password Warning",
        value = "7"
      },
      {
        name  = "Snooze Interval for Expiring Password Warning",
        value = "24"
      },
      {
        name  = "Login Challenge Template",
        value = "html.form.login.challenge.template.html"
      },
      {
        name  = "'Remember My Username' Lifetime",
        value = "30"
      },
      {
        name  = "'This is My Device' Lifetime",
        value = "30"
      },
      {
        name  = "Allow Username Edits During Chaining",
        value = "false"
      },
      {
        name  = "Track Authentication Time",
        value = "true"
      },
      {
        name  = "Post-Password Change Re-Authentication Delay",
        value = "0"
      },
      {
        name  = "Password Reset Username Template",
        value = "forgot-password.html"
      },
      {
        name  = "Password Reset Code Template",
        value = "forgot-password-resume.html"
      },
      {
        name  = "Password Reset Template",
        value = "forgot-password-change.html"
      },
      {
        name  = "Password Reset Error Template",
        value = "forgot-password-error.html"
      },
      {
        name  = "Password Reset Success Template",
        value = "forgot-password-success.html"
      },
      {
        name  = "Account Unlock Template",
        value = "account-unlock.html"
      },
      {
        name  = "OTP Length",
        value = "8"
      },
      {
        name  = "OTP Time to Live",
        value = "10"
      },
      {
        name  = "PingID Properties",
        value = ""
      },
      {
        name  = "Require Verified Email",
        value = "false"
      },
      {
        name  = "Username Recovery Template",
        value = "username.recovery.template.html"
      },
      {
        name  = "Username Recovery Info Template",
        value = "username.recovery.info.template.html"
      },
      {
        name  = "Username Recovery Email Template",
        value = "message-template-username-recovery.html"
      },
      {
        name  = "CAPTCHA for Authentication",
        value = "false"
      },
      {
        name  = "CAPTCHA for Password change",
        value = "false"
      },
      {
        name  = "CAPTCHA for Password Reset",
        value = "false"
      },
      {
        name  = "CAPTCHA for Username recovery",
        value = "false"
      }
    ]
  }
  attribute_contract = {
    mask_ognl_values = false
    core_attributes = [
      {
        masked    = false
        name      = "policy.action"
        pseudonym = false
      },
      {
        masked    = false
        name      = "username"
        pseudonym = true
      }
    ]
    extended_attributes = [
      {
        masked    = false
        name      = "entryUUID"
        pseudonym = false
      }
    ]
  }
}
