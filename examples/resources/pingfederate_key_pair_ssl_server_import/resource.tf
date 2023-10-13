# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #
resource "pingfederate_key_pair_ssl_server_import" "myKeyPairsSslServerImport" {
  file_data = "example"
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}
