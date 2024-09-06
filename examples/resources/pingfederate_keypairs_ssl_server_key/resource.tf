resource "pingfederate_keypairs_ssl_server_key" "sslServerKey" {
  key_id    = "sslserverkey"
  file_data = filebase64("./assets/sslserverkey.p12")
  password  = var.ssl_server_key_password
  format    = "PKCS12"
}