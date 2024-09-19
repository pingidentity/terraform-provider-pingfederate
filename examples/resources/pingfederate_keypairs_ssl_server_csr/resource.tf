resource "pingfederate_keypairs_ssl_server_key" "sslServerKey" {
  file_data = filebase64("./assets/sslserverkey.p12")
  password  = var.ssl_server_key_password
  format    = "PKCS12"
}

resource "pingfederate_keypairs_ssl_server_csr" "example" {
  keypair_id = pingfederate_keypairs_ssl_server_key.sslServerKey.id
  file_data  = filebase64("./path/to/csr_response.pem")
}