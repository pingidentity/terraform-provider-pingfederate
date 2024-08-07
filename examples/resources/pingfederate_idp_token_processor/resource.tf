resource "pingfederate_idp_token_processor" "idpTokenProcessor" {
  processor_id = "myProcessor"
  attribute_contract = {
    core_attributes = [
      {
        name = "SAML_SUBJECT"
      }
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Audience",
        value = "myAudience"
      }
    ]
  }
  name = "My token processor"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.processor.saml.Saml20TokenProcessor"
  }
}