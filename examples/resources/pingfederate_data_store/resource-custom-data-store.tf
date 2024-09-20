resource "pingfederate_data_store" "customDataStore" {
  custom_data_store = {
    name = "Custom REST Data Store"

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
                  value = "https://my_rest_datasource.bxretail.org/api/v1/users"
                },
                {
                  name  = "Tags"
                  value = "production"
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
                  value = "givenName"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/givenName"
                }
              ],
              default_row = false
            },
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "familyName"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/familyName"
                }
              ],
              default_row = false
            },
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "email"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/email"
                }
              ],
              default_row = false
            },
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "password"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/password"
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
          value = "OAuth 2.0 Bearer Token"
        },
        {
          name  = "HTTP Method"
          value = "GET"
        },
        {
          name  = "Username"
          value = var.rest_data_store_basic_auth_username
        },
        {
          name  = "Password Reference"
          value = ""
        },
        {
          name  = "OAuth Token Endpoint"
          value = "https://authservices.bxretail.org/as/token"
        },
        {
          name  = "OAuth Scope"
          value = "restapiscope"
        },
        {
          name  = "Client ID"
          value = var.rest_data_store_oauth2_client_id
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
          value = "https://my_rest_datasource.bxretail.org/api/v1/connectiontest"
        },
        {
          name  = "Test Connection Body"
          value = "{\"foo\":\"bar\"}"
        }
      ]
      sensitive_fields = [
        {
          name  = "Password"
          value = var.rest_data_store_basic_auth_password
        },
        {
          name  = "Client Secret"
          value = var.rest_data_store_oauth2_client_secret
        }
      ]
    }
  }
}
