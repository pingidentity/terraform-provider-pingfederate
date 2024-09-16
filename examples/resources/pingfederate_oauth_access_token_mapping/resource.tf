resource "pingfederate_oauth_access_token_manager" "devices" {
  manager_id = "devicesATM"
  name       = "Device Token Manager"

  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }

  configuration = {
    fields = [
      {
        name  = "Token Length"
        value = "56"
      },
      {
        name  = "Token Lifetime"
        value = "240"
      },
    ]
  }
  attribute_contract = {
    extended_attributes = [
      { name = "directory_id" },
      { name = "device_id" },
      { name = "device_owner_user_id" },
      { name = "device_custodian_user_id" },
    ]
  }
}

resource "pingfederate_oauth_access_token_mapping" "example_device" {
  context = {
    type = "DEFAULT"
  }

  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.devices.id
  }

  attribute_sources = [
    {
      ldap_attribute_source = {
        base_dn = "ou=Devices,dc=bxretail,dc=org"
        data_store_ref = {
          id = pingfederate_data_store.ldap_data_store.id
        }
        description            = "Directory"
        id                     = "Directory"
        member_of_nested_group = false
        search_attributes = [
          "Subject DN",
          "serialNumber",
          "owner",
          "seeAlso",
          "uid",
        ]
        search_filter = "cn=$${USER_KEY}"
        search_scope  = "SUBTREE"
        type          = "LDAP"
      }
    },
  ]

  attribute_contract_fulfillment = {
    "directory_id" = {
      source = {
        id   = "Directory"
        type = "LDAP_DATA_STORE"
      }
      value = "uid"
    },
    "device_id" = {
      source = {
        id   = "Directory"
        type = "LDAP_DATA_STORE"
      }
      value = "serialNumber"
    },
    "device_owner_user_id" = {
      source = {
        id   = "Directory"
        type = "LDAP_DATA_STORE"
      }
      value = "owner"
    },
    "device_custodian_user_id" = {
      source = {
        id   = "Directory"
        type = "LDAP_DATA_STORE"
      }
      value = "seeAlso"
    },
  }

  issuance_criteria = {
    conditional_criteria = [
      {
        attribute_name = "Device Type"
        condition      = "NOT_EQUAL_CASE_INSENSITIVE"
        error_result   = "Cannot issue access token"
        source = {
          type = "EXTENDED_PROPERTIES"
        }
        value = "User Device"
      },
    ]
  }
}
