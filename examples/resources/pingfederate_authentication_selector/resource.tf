resource "pingfederate_authentication_selector" "samlAuthnContextExample" {
  selector_id = "samlAuthnContextExample"
  name        = "samlAuthnContextExample"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.selectors.saml.SamlAuthnContextAdapterSelector"
  }
  configuration = {
    tables = []
    fields = [
      {
        name  = "Add or Update AuthN Context Attribute"
        value = "true"
      },
      {
        name  = "Override AuthN Context for Flow"
        value = "true"
      },
      {
        name  = "Enable 'No Match' Result Value"
        value = "false"
      },
      {
        name  = "Enable 'Not in Request' Result Value"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "result_value2"
      }
    ]
  }
}

