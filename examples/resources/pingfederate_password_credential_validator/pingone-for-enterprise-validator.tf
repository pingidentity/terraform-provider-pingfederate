resource "pingfederate_password_credential_validator" "pingOneForEnterpriseDirectoryPasswordCredentialValidatorExample" {
  validator_id = "pingOneForEnterpriseDirectoryPCV"
  name         = "pingOneForEnterpriseDirectoryPasswordCredentialValidatorExample"
  plugin_descriptor_ref = {
    id = "com.pingconnect.alexandria.pingfed.pcv.PingOnePasswordValidator"
  }
  configuration = {
    fields = [
      {
        name  = "Client Id"
        value = "ping_federate_client_id"
      },
      {
        name  = "Client Secret"
        value = var.pcv_client_secret
      },
      {
        name  = "PingOne URL"
        value = "https://directory-api.pingone.com/api"
      },
      {
        name  = "Authenticate by Subject URL"
        value = "/directory/users/authenticate?by=subject"
      },
      {
        name  = "Reset Password URL"
        value = "/directory/users/password-reset"
      },
      {
        name  = "SCIM User URL"
        value = "/directory/user"
      },
      {
        name  = "Connection Pool Size"
        value = "100"
      },
      {
        name  = "Connection Pool Idle Timeout"
        value = "4000"
      }
    ]
  }
}
