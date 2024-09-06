resource "pingfederate_idp_adapter" "idpAdapterExample" {
  adapter_id = "HTMLForm"
  name       = "HTML Form Adapter Example"

  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.htmlform.idp.HtmlFormIdpAuthnAdapter"
  }

  # ... other required fields
}

resource "pingfederate_authentication_policy_contract" "registration" {
  name = "User Registration"
  extended_attributes = [
    { name = "email" },
    { name = "given_name" },
    { name = "family_name" }
  ]
}

resource "pingfederate_authentication_policies_fragment" "policyFragment" {
  name        = "Registration"
  description = "Sample Registration"

  inputs = {
    id = pingfederate_authentication_policy_contract.registration.id
  }
  outputs = {
    id = pingfederate_authentication_policy_contract.registration.id
  }

  root_node = {
    action = {
      authn_source_policy_action = {
        authentication_source = {
          type = "IDP_ADAPTER"
          source_ref = {
            id = pingfederate_idp_adapter.idpAdapterExample.id
          }
        }
        input_user_id_mapping = {
          source = {
            type = "INPUTS"
            id   = "Inputs"
          }
          value = "username"
        }
        user_id_authenticated = true
      }
    }
    children = [
      {
        action = {
          done_policy_action = {
            context = "Fail"
          }
        }
      },
      {
        action = {
          apc_mapping_policy_action = {
            context = "Success"
            authentication_policy_contract_ref = {
              id = pingfederate_authentication_policy_contract.registration.id
            }
            attribute_mapping = {
              attribute_sources = []
              attribute_contract_fulfillment = {
                "firstName" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.idpAdapterExample.id
                  }
                  value = "firstName"
                }
                "lastName" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.idpAdapterExample.id
                  }
                  value = "lastName"
                }
                "subject" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.idpAdapterExample.id
                  }
                  value = "subject"
                }
                "fullName" : {
                  source = {
                    type = "EXPRESSION"
                  },
                  value = "fullName"
                }
                "photo" : {
                  source = {
                    type = "ADAPTER",
                    id   = pingfederate_idp_adapter.idpAdapterExample.id
                  }
                  value = "photo"
                }
                "username" : {
                  source = {
                    type = "INPUTS",
                    id   = "inputs"
                  }
                  value = "username"
                }
              }
            }
          }
        }
      }
    ]
  }
}
