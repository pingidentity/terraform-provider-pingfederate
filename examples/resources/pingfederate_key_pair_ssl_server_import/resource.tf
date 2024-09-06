resource "pingfederate_key_pair_ssl_server_import" "keyPairsSslServerImport" {
  import_id = "keyPairSSLServerImport"
  file_data = "example"
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}
