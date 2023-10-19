terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.1.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username               = "administrator"
  password               = "2FederateM0re"
  https_host             = "https://localhost:9999"
  insecure_trust_all_tls = true
}

resource "pingfederate_data_store" "myJdbcDataStore" {
  custom_id             = "myJdbcDataStore"
  mask_attribute_values = false
  jdbc_data_store = {
    connection_url               = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false"
    driver_class                 = "org.hsqldb.jdbcDriver"
    user_name                    = "sa"
    password                     = "secretpass"
    allow_multi_value_attributes = false
    name                         = "jdbc"
    connection_url_tags = [
      {
        connection_url = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false",
        default_source = true
      }
    ],
    min_pool_size    = 10
    max_pool_size    = 100
    blocking_timeout = 5000
    idle_timeout     = 5
  }
}
