resource "pingfederate_pingone_connection" "example" {
  name       = "My PingOne Tenant"
  credential = var.pingone_connection_credential
}

resource "pingfederate_data_store" "pingOneLdapDataStore" {
  ping_one_ldap_gateway_data_store = {
    name      = "PingOne LDAP Gateway"
    ldap_type = "PING_DIRECTORY"

    ping_one_connection_ref = {
      id = pingfederate_pingone_connection.example.id
    }

    ping_one_environment_id  = var.pingone_environment_id
    ping_one_ldap_gateway_id = var.pingone_gateway_id

    use_ssl = true
  }
}