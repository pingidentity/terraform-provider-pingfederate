resource "pingfederate_data_store" "myPingOneLdapDataStore" {
  data_store_id         = "myPingOneLdapDataStore"
  mask_attribute_values = false
  ping_one_ldap_gateway_data_store = {
    ldap_type = "PING_DIRECTORY"
    name      = "myPingOneLdapDataStore"
    ping_one_connection_ref = {
      id = ""
    },
    ping_one_environment_id  = ""
    ping_one_ldap_gateway_id = ""
    use_ssl                  = true
  }
}