resource "pingfederate_server_settings_general" "generalSettings" {
  datastore_validation_interval_secs          = 300
  disable_automatic_connection_validation     = false
  idp_connection_transaction_logging_override = "NONE"
  request_header_for_correlation_id           = "example"
  sp_connection_transaction_logging_override  = "FULL"
}
