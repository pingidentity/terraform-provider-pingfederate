resource "pingfederate_data_store" "pingOneLdapDataStore" {
  mask_attribute_values = false
  ping_one_ldap_gateway_data_store = {
    name      = "PingOne Gateway Data Store"
    ldap_type = "PING_DIRECTORY"

    ping_one_connection_ref = {
      id = var.pingone_connection_id
    }

    ping_one_environment_id  = var.pingone_environment_id
    ping_one_ldap_gateway_id = var.pingone_ldap_gateway_id

    use_ssl = true
  }
}
