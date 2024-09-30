resource "pingfederate_connection_metadata_export" "metadataExport" {
  connection_type = "SP"
  connection_id   = pingfederate_idp_sp_connection.example_saml.connection_id
  signing_settings = {
    signing_key_pair_ref = {
      id = pingfederate_keypairs_signing_key.rsa_saml_signing_1.id
    }
    algorithm = "SHA256withRSA"
  }
}