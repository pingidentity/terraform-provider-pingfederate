resource "pingfederate_data_store" "jdbcDataStore" {
  jdbc_data_store = {
    name = "JDBC Data Store"

    connection_url = "jdbc:sqlserver://localhost;encrypt=true;integratedSecurity=true;"
    driver_class   = "org.hsqldb.jdbcDriver"

    user_name = var.jdbc_data_store_username
    password  = var.jdbc_data_store_password

    allow_multi_value_attributes = false

    connection_url_tags = [
      {
        connection_url = "jdbc:sqlserver://localhost;encrypt=true;integratedSecurity=true;"
        default_source = true
      }
    ]

    min_pool_size    = 10
    max_pool_size    = 100
    blocking_timeout = 5000
    idle_timeout     = 5
  }
}
