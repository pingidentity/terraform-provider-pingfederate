resource "pingfederate_data_store" "jdbcDataStore" {
  data_store_id         = "jdbcDataStore"
  mask_attribute_values = false
  jdbc_data_store = {
    name                         = "jdbcDataStore"
    connection_url               = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false"
    driver_class                 = "org.hsqldb.jdbcDriver"
    user_name                    = "sa"
    password                     = "2FederateM0re"
    allow_multi_value_attributes = false
    connection_url_tags = [
      {
        connection_url = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false",
        default_source = true
      }
    ]
    min_pool_size    = 10
    max_pool_size    = 100
    blocking_timeout = 5000
    idle_timeout     = 5
  }
}

resource "pingfederate_server_settings_outbound_provisioning" "outboundProvisioningSettings" {
  data_store_ref = {
    id = pingfederate_data_store.jdbcDataStore.data_store_id
  }
  synchronization_frequency = 60
}