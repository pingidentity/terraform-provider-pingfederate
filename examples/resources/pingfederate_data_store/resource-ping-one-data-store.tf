resource "pingfederate_data_store" "pingOneDataStore" {
  data_store_id = "pingOneDataStore"
  custom_data_store = {
    name = "pingOneDataStore"
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
