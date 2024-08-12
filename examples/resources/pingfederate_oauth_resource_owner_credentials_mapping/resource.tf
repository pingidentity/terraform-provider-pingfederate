resource "pingfederate_password_credential_validator" "pcv" {
  validator_id = "simpleValidator"
  name         = "pcv"
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
                value = var.password_credential_validator_password
              },
              {
                name  = "Confirm Password"
                value = var.password_credential_validator_password
              },
              {
                name  = "Relax Password Requirements"
                value = "true"
              }
            ]
            default_row = false
          },
        ]
      }
    ]
  }
}

resource "pingfederate_oauth_resource_owner_credentials_mapping" "mapping" {
  mapping_id = pingfederate_password_credential_validator.pcv.validator_id
  attribute_contract_fulfillment = {
    "USER_KEY" = {
      source = {
        type = "CONTEXT"
      }
      value = "ClientId"
    }
  }
}