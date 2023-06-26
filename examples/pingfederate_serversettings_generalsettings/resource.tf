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

resource "pingfederate_serversettings_generalsettings" "serverSettingsGeneralSettingsExample" {
  datastore_validation_interval_secs = 300
	disable_automatic_connection_validation = false
  idp_connection_transaction_logging_override = "NONE"
  request_header_for_correlation_id = "example"
  sp_connection_transaction_logging_override = "FULL"
}
