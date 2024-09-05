resource "pingfederate_password_credential_validator" "pingOnePasswordCredentialValidatorExample" {
  validator_id = "pingOnePCV"
  name         = "pingOnePasswordCredentialValidatorExample"
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
        value = "mydatastore"
      },
      {
        name  = "Case-Sensitive Matching"
        value = "true"
      }
    ]
  }
}