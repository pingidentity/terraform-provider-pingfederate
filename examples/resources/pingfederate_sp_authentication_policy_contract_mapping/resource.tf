resource "pingfederate_authentication_policy_contract" "example" {
  name = "User"
  extended_attributes = [
    { name = "email" },
  ]
}

resource "pingfederate_sp_adapter" "opentoken_example" {
  adapter_id = "myOpenTokenAdapter"
  name       = "My OpenToken Adapter"

  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.opentoken.SpAuthnAdapter"
  }

  attribute_contract = {
    extended_attributes = [
      { name = "firstName" },
      { name = "lastName" },
      { name = "email" },
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Token Name"
        value = "opentoken"
      },
    ]
    sensitive_fields = [
      {
        name  = "Confirm Password"
        value = var.sp_adapter_opentoken_password
      },
      {
        name  = "Password"
        value = var.sp_adapter_opentoken_password
      },
    ]
  }
}

resource "pingfederate_sp_authentication_policy_contract_mapping" "spAuthenticationPolicyContractMappingExample" {
  source_id = pingfederate_authentication_policy_contract.example.id
  target_id = pingfederate_sp_adapter.opentoken_example.id

  attribute_sources = [
    {
      ldap_attribute_source = {
        base_dn = "ou=Users,dc=bxretail,dc=org"
        data_store_ref = {
          id = pingfederate_data_store.ldap_data_store.id
        }
        description            = "Directory"
        id                     = "Directory"
        member_of_nested_group = false
        search_attributes = [
          "Subject DN",
          "givenName",
          "sn",
          "mail",
          "type",
        ]
        search_filter = "cn=$${subject}"
        search_scope  = "SUBTREE"
        type          = "LDAP"
      }
    },
  ]
  attribute_contract_fulfillment = {
    "subject" = {
      source = {
        type = "AUTHENTICATION_POLICY_CONTRACT"
      },
      value = "subject"
    }
    "firstName" = {
      source = {
        id   = "Directory"
        type = "LDAP_DATA_STORE"
      },
      value = "givenName"
    }
    "lastName" = {
      source = {
        id   = "Directory"
        type = "LDAP_DATA_STORE"
      },
      value = "sn"
    }
    "email" = {
      source = {
        type = "AUTHENTICATION_POLICY_CONTRACT"
      },
      value = "email"
    }
  }
  issuance_criteria = {
    conditional_criteria = [
      {
        error_result = "Cannot issue token"
        source = {
          id   = "Directory"
          type = "LDAP_DATA_STORE"
        }
        attribute_name = "type"
        condition      = "EQUALS_CASE_INSENSITIVE"
        value          = "deleted"
      }
    ]
  }
}
