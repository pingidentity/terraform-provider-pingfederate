// Example of using the time provider to control regular export of CSR
resource "time_rotating" "csr_export" {
  rotation_days = 30
}

resource "pingfederate_keypairs_signing_csr_export" "signingCsr" {
  keypair_id = "mysigningkeypair"
  export_trigger_values = {
    "export_rfc3339" : time_rotating.csr_export.rotation_rfc3339,
  }
}