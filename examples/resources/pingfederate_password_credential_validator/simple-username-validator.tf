resource "pingfederate_password_credential_validator" "simpleUsernamePasswordCredentialValidatorExample" {
  validator_id = "simpleUsernamePCV"
  name         = "simpleUsernamePasswordCredentialValidator"
  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator"
  }
  configuration = {
    tables = [
      {
        name = "Users"
        rows = [
          {
            fields = [
              {
                name  = "Username"
                value = "example"
              },
              {
                name  = "Password"
                value = var.pcv_password_user1
              },
              {
                name  = "Confirm Password"
                value = var.pcv_password_user1
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          },
          {
            fields = [
              {
                name  = "Username"
                value = "example2"
              },
              {
                name  = "Password"
                value = var.pcv_password_user2
              },
              {
                name  = "Confirm Password"
                value = var.pcv_password_user2
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          }
        ],
      }
    ]
  }
}
