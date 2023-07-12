terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username = "administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:9999"
}
# Please do not delete any log categories from the resource file. 
# PingFederate API does not currently support accessing logCategories by ID.
# Please change the enabled field to desired configuration.
resource "pingfederate_server_settings_log_settings" "serverSettingsLogSettingsExample" {
  log_categories = [
    {
      id = "core"
      name = "Core"
      description = "Debug logging for core components."
      enabled = false
    },
    {
      id = "policytree"
      name = "Policy Tree"
      description = "Policy tree debug logging."
      enabled = false
    },
    {
      id = "trustedcas"
      name = "Trusted CAs"
      description = "Log PingFederate and JRE trusted CAs when they are loaded."
      enabled = false
    },
    {
      id = "xmlsig"
      name = "XML Signatures"
      description = "Debug logging for XML signature operations."
      enabled = false
    },
    {
      id = "requestheaders"
      name = "HTTP Request Headers"
      description = "Log HTTP request headers. Sensitive information, such as passwords, may be logged when this category is enabled."
      enabled = false
    },
    {
      id = "requestparams"
      name = "HTTP Request Parameters"
      description = "Log HTTP GET request parameters. Sensitive information, such as passwords, may be logged when this category is enabled."
      enabled = false
    },
    {
      id = "restdatastore"
      name = "REST Data Store Requests and Responses"
      description = "Log REST datastore requests and responses. Sensitive information, such as passwords, may be logged when this category is enabled."
      enabled = false
    },
  ]
}