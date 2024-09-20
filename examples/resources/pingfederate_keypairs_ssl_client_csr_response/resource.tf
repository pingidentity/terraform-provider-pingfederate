resource "pingfederate_keypairs_ssl_client_key" "sslClientKey" {
  file_data = filebase64("./assets/sslclientkey.p12")
  password  = var.ssl_client_key_password
  format    = "PKCS12"
}

resource "pingfederate_keypairs_ssl_client_csr_response" "example" {
  keypair_id = pingfederate_keypairs_ssl_client_key.sslClientKey.id
  file_data  = filebase64("./path/to/csr_response.pem")
}
