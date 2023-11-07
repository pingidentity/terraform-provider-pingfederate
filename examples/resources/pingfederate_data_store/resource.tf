# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

# x_bypass_external_validation_header may be used in the provider block to bypass the connection test.
resource "pingfederate_data_store" "myCustomDataStore" {
  custom_id = "myCustomDataStore"
  custom_data_store = {
    name = "myCustomDataStore"
    plugin_descriptor_ref = {
      id = "com.pingidentity.pf.datastore.other.RestDataSourceDriver"
    }
    configuration = {
      tables = [
        {
          name = "Base URLs and Tags"
          rows = [
            {
              fields = [
                {
                  name  = "Base URL"
                  value = "http://localhost"
                },
                {
                  name  = "Tags"
                  value = "tag"
                }
              ],
              default_row = true
            }
          ]
        },
        {
          name = "HTTP Request Headers"
          rows = [
            {
              fields = [
                {
                  name  = "Header Name"
                  value = "header"
                },
                {
                  name  = "Header Value"
                  value = "header_value"
                }
              ],
              default_row = false
            }
          ]
        },
        {
          name = "Attributes"
          rows = [
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "attribute"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/json_response_attr_path"
                }
              ],
              default_row = false
            }
          ]
        }
      ],
      fields = [
        {
          name  = "Authentication Method"
          value = "Basic Authentication"
        },
        {
          name  = "HTTP Method"
          value = "GET"
        },
        {
          name  = "Username"
          value = "Administrator"
        },
        {
          name  = "Password"
          value = "2FederateM0re"
        },
        {
          name  = "Password Reference"
          value = ""
        },
        {
          name  = "OAuth Token Endpoint"
          value = "https://example.com"
        },
        {
          name  = "OAuth Scope"
          value = "scope"
        },
        {
          name  = "Client ID"
          value = "client_id"
        },
        {
          name  = "Client Secret"
          value = "2FederateM0re"
        },
        {
          name  = "Client Secret Reference"
          value = ""
        },
        {
          name  = "Enable HTTPS Hostname Verification"
          value = "true"
        },
        {
          name  = "Read Timeout (ms)"
          value = "10000"
        },
        {
          name  = "Connection Timeout (ms)"
          value = "10000"
        },
        {
          name  = "Max Payload Size (KB)"
          value = "1024"
        },
        {
          name  = "Retry Request"
          value = "true"
        },
        {
          name  = "Maximum Retries Limit"
          value = "5"
        },
        {
          name  = "Retry Error Codes"
          value = "429"
        },
        {
          name  = "Test Connection URL"
          value = "https://example.com"
        },
        {
          name  = "Test Connection Body"
          value = "body"
        }
      ]
    }
  }
}

resource "pingfederate_data_store" "myJdbcDataStore" {
  custom_id             = "myJdbcDataStore"
  mask_attribute_values = false
  jdbc_data_store = {
    name                         = "myJdbcDataStore"
    connection_url               = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false"
    driver_class                 = "org.hsqldb.jdbcDriver"
    user_name                    = "sa"
    password                     = "secretpass"
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

resource "pingfederate_data_store" "myPingDirectoryLdapDataStore" {
  custom_id = "myPingDirectoryLdapDataStore"
  ldap_data_store = {
    ldap_type          = "PING_DIRECTORY"
    bind_anonymously   = false
    user_dn            = "cn=pingfederate"
    password           = "2FederateM0re"
    use_ssl            = false
    use_dns_srv_record = false
    name               = "myPingDirectoryLdapDataStore"
    hostnames = [
      "pingdirectory:1389"
    ]
    hostnames_tags = [
      {
        hostnames = [
          "pingdirectory:1389"
        ]
        default_source = true
      }
    ]
    test_on_borrow         = true
    test_on_return         = false
    create_if_necessary    = true
    verify_host            = true
    min_connections        = 10
    max_connections        = 100
    max_wait               = -1
    time_between_evictions = 6000
    read_timeout           = 300
    connection_timeout     = 300
    binary_attributes      = []
    dns_ttl                = 6000
  }
}

resource "pingfederate_data_store" "myPingOneDataStore" {
  custom_id = "myPingOneDataStore"
  custom_data_store = {
    name = "myPingOneDataStore"
    plugin_descriptor_ref = {
      id = "com.pingidentity.plugins.datastore.p14c.PingOneForCustomersDataStore"
    }
    configuration = {
      tables = [
        {
          name = "Custom Attributes Details",
          rows = [
            {
              fields = [
                {
                  name  = "Local Attribute",
                  value = "local_attribute"
                },
                {
                  name  = "PingOne for Customers Attribute",
                  value = "/pingone_attribute"
                }
              ],
              defaultRow = false
            }
          ]
        }
      ],
      fields = [
        {
          name  = "PingOne Environment",
          value = ""
        },
        {
          name  = "Connection Timeout",
          value = "10000"
        },
        {
          name  = "Retry Request",
          value = "true"
        },
        {
          name  = "Maximum Retries Limit",
          value = "5"
        },
        {
          name  = "Retry Error Codes",
          value = "429"
        },
        {
          name  = "Proxy Settings",
          value = "System Defaults"
        },
        {
          name  = "Custom Proxy Host",
          value = ""
        },
        {
          name  = "Custom Proxy Port",
          value = ""
        }
      ]
    }
    mask_attribute_values = false
  }
}

resource "pingfederate_data_store" "myPingOneLdapDataStore" {
  custom_id             = "myPingOneLdapDataStore"
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
