resource "pingfederate_sp_adapter" "spAdapter" {
  adapter_id = "myOpenTokenAdapter"
  name       = "My OpenToken Adapter"

  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.opentoken.SpAuthnAdapter"
  }

  configuration = {
    fields = [
      {
        name  = "Password",
        value = var.opentoken_sp_adapter_password
      },
      {
        name  = "Confirm Password",
        value = var.opentoken_sp_adapter_password
      },
      {
        name  = "Transport Mode",
        value = "2"
      },
      {
        name  = "Token Name",
        value = "spopentoken"
      },
      {
        name  = "Cipher Suite",
        value = "2"
      },
      {
        name  = "Authentication Service",
        value = ""
      },
      {
        name  = "Account Link Service",
        value = "https://auth.bxretail.org/SpSample/?cmd=accountlink"
      },
      {
        name  = "Logout Service",
        value = "https://auth.bxretail.org/SpSample/?cmd=slo"
      },
      {
        name  = "Cookie Domain",
        value = ""
      },
      {
        name  = "Cookie Path",
        value = "/"
      },
      {
        name  = "Token Lifetime",
        value = "300"
      },
      {
        name  = "Session Lifetime",
        value = "43200"
      },
      {
        name  = "Not Before Tolerance",
        value = "0"
      },
      {
        name  = "Force SunJCE Provider",
        value = "false"
      },
      {
        name  = "Use Verbose Error Messages",
        value = "false"
      },
      {
        name  = "Obfuscate Password",
        value = "true"
      },
      {
        name  = "Session Cookie",
        value = "false"
      },
      {
        name  = "Secure Cookie",
        value = "false"
      },
      {
        name  = "HTTP Only Flag",
        value = "true"
      },
      {
        name  = "Send Subject as Query Parameter",
        value = "false"
      },
      {
        name  = "Subject Query Parameter                 ",
        value = ""
      },
      {
        name  = "Send Extended Attributes",
        value = "0"
      },
      {
        name  = "Skip Trimming of Trailing Backslashes",
        value = "false"
      },
      {
        name  = "SameSite Cookie",
        value = "3"
      },
      {
        name  = "URL Encode Cookie Values",
        value = "true"
      },
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "firstName"
      },
      {
        name = "lastName"
      },
      {
        name = "email"
      }
    ]
  }
}
