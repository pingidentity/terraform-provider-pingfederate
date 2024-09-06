resource "pingfederate_idp_adapter" "myIdpAdapter" {
  adapter_id = "myIdPAdapter"
  name       = "My IdP Adapter"

  attribute_contract = {
    core_attributes = [
      {
        masked    = false
        name      = "policy.action"
        pseudonym = false
      },
      {
        masked    = false
        name      = "username"
        pseudonym = true
      }
    ]

    extended_attributes = [
      {
        name = "givenName_idp"
      },
      {
        name = "familyName_idp"
      },
      {
        name = "email_idp"
      }
    ]
  }

  # ... other required configuration parameters
}

resource "pingfederate_sp_adapter" "mySPAdapter" {
  adapter_id = "mySPAdapter"
  name       = "My SP Adapter"

  attribute_contract = {
    extended_attributes = [
      {
        name = "givenName_sp"
      },
      {
        name = "familyName_sp"
      },
      {
        name = "email_sp"
      }
    ]
  }

  # ... other required configuration parameters
}

resource "pingfederate_idp_to_sp_adapter_mapping" "idpToSpAdapterMapping" {
  source_id = pingfederate_idp_adapter.myIdpAdapter.adapter_id
  target_id = pingfederate_sp_adapter.mySPAdapter.adapter_id

  attribute_contract_fulfillment = {
    "subject" = {
      source = {
        type = "ADAPTER"
      }
      value = "username"
    }
    "givenName_sp" = {
      source = {
        type = "ADAPTER"
      }
      value = "givenName_idp"
    }
    "familyName_sp" = {
      source = {
        type = "ADAPTER"
      }
      value = "familyName_idp"
    }
    "email_sp" = {
      source = {
        type = "ADAPTER"
      }
      value = "email_idp"
    }
  }

  application_name = "My Awesome Application"
}