// Example of using the time provider to control regular export of CSR
resource "time_rotating" "csr_export" {
  rotation_days = 30
}

resource "pingfederate_keypairs_ssl_client_csr_export" "example" {
  keypair_id = "sslclientkeypair"
  export_trigger_values = {
    "export_rfc3339" : time_rotating.csr_export.rotation_rfc3339,
  }
}