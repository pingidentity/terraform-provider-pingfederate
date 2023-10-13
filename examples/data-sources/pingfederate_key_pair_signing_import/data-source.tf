resource "pingfederate_key_pair_signing_import" "myKeyPairsSigningImport" {
  custom_id = "id"
  file_data = ""
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}

data "pingfederate_key_pair_signing_import" "myKeyPairsSigningImport" {
  id = pingfederate_key_pair_signing_import.myKeyPairsSigningImport.custom_id
}
resource "pingfederate_key_pair_signing_import" "keyPairsSigningImportExample" {
  custom_id = "${data.pingfederate_key_pair_signing_import.myKeyPairsSigningImport.id}2"
  file_data = ""
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}