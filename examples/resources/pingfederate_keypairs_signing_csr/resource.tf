resource "pingfederate_keypairs_signing_csr" "csrResponse" {
  keypair_id = "mysigningkeypair"
  file_data  = filebase64("./path/to/csr_response.pem")
}