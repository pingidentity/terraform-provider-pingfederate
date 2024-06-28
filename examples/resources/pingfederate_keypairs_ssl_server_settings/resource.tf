resource "pingfederate_keypairs_ssl_server_settings" "sslServerSettings" {
  active_admin_console_certs = [
    {
      id = "sslservercert"
    }
  ]
  active_runtime_server_certs = [
    {
      id = "sslservercert"
    }
  ]
  admin_console_cert_ref = {
    id = "sslservercert"
  }
  runtime_server_cert_ref = {
    id = "sslservercert"
  }
}