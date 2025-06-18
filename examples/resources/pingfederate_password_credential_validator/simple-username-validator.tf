resource "pingfederate_password_credential_validator" "simpleUsernamePasswordCredentialValidatorExample" {
  validator_id = "simpleUsernamePCV"
  name         = "Simple Username Password Credential Validator"

  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator"
  }

  attribute_contract = {}

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
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            sensitive_fields = [
              {
                name  = "Password"
                value = var.pcv_password_user1
              },
              {
                name  = "Confirm Password"
                value = var.pcv_password_user1
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
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            sensitive_fields = [
              {
                name  = "Password"
                value = var.pcv_password_user2
              },
              {
                name  = "Confirm Password"
                value = var.pcv_password_user2
              }
            ]
            default_row = false
          }
        ],
      }
    ]
  }
}
