resource "pingfederate_password_credential_validator" "radiusUsernamePasswordCredentialValidatorExample" {
  validator_id = "radiusUnPwPCV"
  name         = "RADIUS Username Password Credential Validator"

  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.RadiusUsernamePasswordCredentialValidator"
  }

  configuration = {
    tables = [
      {
        name = "RADIUS Servers"
        rows = [
          {
            fields = [
              {
                name  = "Hostname"
                value = "localhost"
              },
              {
                name  = "Authentication Port"
                value = "1812"
              },
              {
                name  = "Authentication Protocol"
                value = "PAP"
              },
            ]
            sensitive_fields = [
              {
                name  = "Shared Secret"
                value = var.pcv_shared_secret
              }
            ]
            default_row = false
          }
        ]
      }
    ],
    fields = [
      {
        name  = "NAS Identifier"
        value = "PingFederate"
      },
      {
        name  = "Timeout"
        value = "3000"
      },
      {
        name  = "Retry Count"
        value = "3"
      },
      {
        name  = "Allow Challenge Retries after Access-Reject"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "contract"
      }
    ]
  }
}
