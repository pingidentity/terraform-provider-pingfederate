resource "pingfederate_key_pair_ssl_server_import" "myKeyPairsSslServerImport" {
  custom_id = "id"
  file_data = ""
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}

data "pingfederate_key_pair_ssl_server_import" "myKeyPairsSslServerImport" {
  id = pingfederate_key_pair_ssl_server_import.myKeyPairsSslServerImport.custom_id
}
resource "pingfederate_key_pair_ssl_server_import" "keyPairsSslServerImportExample" {
  custom_id = "${data.pingfederate_key_pair_ssl_server_import.myKeyPairsSslServerImport.id}2"
  file_data = ""
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}