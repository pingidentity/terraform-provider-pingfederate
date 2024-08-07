resource "pingfederate_data_store" "customDataStore" {
  data_store_id = "customDataStore"
  custom_data_store = {
    name = "customDataStore"
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
