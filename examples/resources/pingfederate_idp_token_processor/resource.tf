resource "pingfederate_idp_token_processor" "idpTokenProcessor" {
  processor_id = "myProcessor"
  name         = "My token processor"

  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.processor.saml.Saml20TokenProcessor"
  }

  attribute_contract = {
    core_attributes = [
      {
        name = "SAML_SUBJECT"
      }
    ]

    extended_attributes = [
      {
        masked = false
        name   = "lastName"
      },
      {
        masked = true
        name   = "firstName"
      },
      {
        masked = false
        name   = "email"
      },
    ]
  }

  configuration = {
    fields = [
      {
        name  = "Audience",
        value = "myAudience"
      }
    ]

    tables = [
      {
        name = "Valid Certificate Issuer DNs",
        rows = [
          {
            fields = [
              {
                name  = "Issuer DN",
                value = "cn=issuer1"
              }
            ]
          },
          {
            fields = [
              {
                name  = "Issuer DN",
                value = "cn=issuer2"
              }
            ]
          }
        ]
      },
      {
        name = "Valid Certificate Subject DNs",
        rows = [
          {
            fields = [
              {
                name  = "Subject DN",
                value = "cn=validcert1"
              }
            ]
          },
          {
            fields = [
              {
                name  = "Subject DN",
                value = "cn=validcert2"
              }
            ]
          }
        ]
      }
    ]
  }
}
