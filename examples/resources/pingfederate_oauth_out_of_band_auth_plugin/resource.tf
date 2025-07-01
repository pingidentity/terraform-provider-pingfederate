resource "pingfederate_oauth_out_of_band_auth_plugin" "authPlugin" {
  plugin_id = "CIBAAuthenticator"
  name      = "CIBA Authenticator (PingOne MFA)"

  configuration = {
    tables = [
      {
        name = "PingOne Template Variables"
      }
    ]
    fields = [
      {
        name  = "PingOne Environment"
        value = format("%s|%s", pingfederate_pingone_connection.example.id, var.pingone_environment_id)
      },
      {
        name  = "Application"
        value = var.pingone_mfa_application_id
      },
      {
        name  = "PingOne Authentication Policy"
        value = var.pingone_sign_on_policy_name
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
  plugin_descriptor_ref = {
    id = "com.pingidentity.oobauth.pingone.mfa.PingOneMfaCibaAuthenticator"
  }
}