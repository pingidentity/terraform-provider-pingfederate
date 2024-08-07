# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_sp_adapter" "spAdapter" {
  adapter_id = "OTSPJava"
  name       = "OTSPJava"
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.opentoken.SpAuthnAdapter"
  }
  configuration = {
    fields = [
      {
        name  = "Password",
        value = "2FederateM0re"
      },
      {
        name  = "Confirm Password",
        value = "2FederateM0re"
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
        value = "https://localhost:9031/SpSample/?cmd=accountlink"
      },
      {
        name  = "Logout Service",
        value = "https://localhost:9031/SpSample/?cmd=slo"
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
        name = "spadapter_attr3"
      },
      {
        name = "spadapter_attr2"
      },
      {
        name = "spadapter_attr1"
      }
    ]
  }
}