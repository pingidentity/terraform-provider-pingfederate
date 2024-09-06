resource "pingfederate_pingone_connection" "example" {
  name       = "My PingOne Tenant"
  credential = var.pingone_connection_credential
}

resource "pingfederate_data_store" "pingOneDataStore" {
  custom_data_store = {
    name = format("PingOne Data Store (%s)", var.pingone_environment_name)

    plugin_descriptor_ref = {
      id = "com.pingidentity.plugins.datastore.p14c.PingOneForCustomersDataStore"
    }

    configuration = {
      fields = [
        {
          name  = "PingOne Environment",
          value = format("%s|%s", pingfederate_pingone_connection.example.id, var.pingone_environment_id)
        }
      ]
    }
    mask_attribute_values = false
  }
}

resource "pingfederate_password_credential_validator" "pingOnePasswordCredentialValidatorExample" {
  validator_id = "pingOnePCV"
  name         = "PingOne Directory Password Credential Validator"

  plugin_descriptor_ref = {
    id = "com.pingidentity.plugins.pcvs.p14c.PingOneForCustomersPCV"
  }

  configuration = {
    tables = [
      {
        name = "Authentication Error Overrides"
        rows = []
      }
    ],
    fields = [
      {
        name  = "PingOne For Customers Datastore"
        value = pingfederate_data_store.pingOneDataStore.id
      },
      {
        name  = "Case-Sensitive Matching"
        value = "true"
      }
    ]
  }
}