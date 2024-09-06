resource "pingfederate_key_pair_signing_import" "keyPairsSigningImport" {
  import_id = "signingImportId"
  file_data = "example"
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}
