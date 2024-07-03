data "pingfederate_connection_metadata_export" "metadataExport" {
  connection_type = "SP"
  connection_id   = "mySpConnection"
  signing_settings = {
    signing_key_pair_ref = {
      id = "mySigningKeyId"
    }
    algorithm = "SHA512withRSA"
  }
}