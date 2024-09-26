resource "pingfederate_keypairs_ssl_client_key" "sslClientKey" {
  key_id    = "sslclientkey"
  file_data = filebase64("./assets/sslclientkey.p12")
  password  = var.ssl_client_key_password
  format    = "PKCS12"
}