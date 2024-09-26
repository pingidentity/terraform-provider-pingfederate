resource "pingfederate_keypairs_ssl_client_key" "sslClientKey" {
  file_data = filebase64("./assets/sslclientkey.p12")
  password  = var.ssl_client_key_password
  format    = "PKCS12"
}

// Example of using the time provider to control regular export of CSR
resource "time_rotating" "csr_export" {
  rotation_days = 30
}

resource "pingfederate_keypairs_ssl_client_csr_export" "example" {
  keypair_id = pingfederate_keypairs_ssl_client_key.sslClientKey.id

  export_trigger_values = {
    "export_rfc3339" : time_rotating.csr_export.rotation_rfc3339,
  }
}
