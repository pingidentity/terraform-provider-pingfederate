resource "pingfederate_keypairs_ssl_client_csr" "example" {
  keypair_id = "sslclientkeypair"
  file_data  = filebase64("./path/to/csr_response.pem")
}