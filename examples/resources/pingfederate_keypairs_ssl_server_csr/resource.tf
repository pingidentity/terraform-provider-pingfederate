resource "pingfederate_keypairs_ssl_server_csr" "example" {
  keypair_id = "sslserverkeypair"
  file_data  = filebase64("./path/to/csr_response.pem")
}