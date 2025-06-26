resource "pingfederate_oauth_out_of_band_auth_plugin" "authPlugin" {
  plugin_id = "CIBAAuthenticator"
  configuration = {
    tables = [
      {
        name = "PingOne Template Variables"
        rows = []
      }
    ]
    fields = [
      {
        name  = "PingOne Environment"
        value = var.pingone_environment
      },
      {
        name  = "Application"
        value = var.pingone_application
      },
      {
        name  = "PingOne Authentication Policy"
        value = "Standalone_MFA"
      },
      {
        name  = "Test Username"
        value = "user.0"
      },
      {
        name  = "PingOne Template Name"
        value = "transaction"
      },
      {
        name  = "PingOne Template Variant"
        value = ""
      },
      {
        name  = "Client Context"
        value = "Example"
      },
      {
        name  = "Messages Files"
        value = "pingone-mfa-messages"
      },
      {
        name  = "API Request Timeout"
        value = "5000"
      },
      {
        name  = "Proxy Settings"
        value = "System Defaults"
      }
    ]
  }
  name = "CIBA Authenticator (PingOne MFA)"
  plugin_descriptor_ref = {
    id = "com.pingidentity.oobauth.pingone.mfa.PingOneMfaCibaAuthenticator"
  }
}