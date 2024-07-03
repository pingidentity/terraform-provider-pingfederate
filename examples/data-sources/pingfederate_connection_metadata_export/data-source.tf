data "pingfederate_connection_metadata_export" "metadataExport" {
  connection_type = "SP"
  connection_id   = "mySpConnection"
  signing_settings = {
    signing_key_pair_ref = {
      id = "419x9yg43rlawqwq9v6az997k"
    }
    algorithm = "SHA512withRSA"
  }
}