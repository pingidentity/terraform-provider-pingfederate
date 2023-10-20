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

# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_data_store" "myCustomDataStore" {
  custom_id = "myCustomDataStore"
  # mask_attribute_value= false
  custom_data_store = {
    name = "custom"
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

# resource "pingfederate_data_store" "myJdbcDataStore" {
#   custom_id             = "myJdbcDataStore"
#   mask_attribute_values = false
#   jdbc_data_store = {
#     name                         = "jdbc"
#     connection_url               = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false"
#     driver_class                 = "org.hsqldb.jdbcDriver"
#     user_name                    = "sa"
#     password                     = "secretpass"
#     allow_multi_value_attributes = false
#     connection_url_tags = [
#       {
#         connection_url = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false",
#         default_source = true
#       }
#     ],
#     min_pool_size    = 10
#     max_pool_size    = 100
#     blocking_timeout = 5000
#     idle_timeout     = 5
#   }
# }
