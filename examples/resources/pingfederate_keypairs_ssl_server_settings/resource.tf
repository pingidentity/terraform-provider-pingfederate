resource "pingfederate_keypairs_ssl_server_key" "admin_console" {
  key_id    = "adminconsole"
  file_data = filebase64("./path/to/admin_console.p12")
  format    = "PKCS12"
  password  = var.runtime_server_cert_admin_console_pkcs12_password
}

resource "pingfederate_keypairs_ssl_server_key" "runtime" {
  key_id    = "runtime"
  file_data = filebase64("./path/to/runtime.p12")
  format    = "PKCS12"
  password  = var.runtime_server_cert_runtime_pkcs12_password
}

resource "pingfederate_keypairs_ssl_server_settings" "sslServerSettings" {
  admin_console_cert_ref = {
    id = pingfederate_keypairs_ssl_server_key.admin_console.id
  }
  active_admin_console_certs = [
    {
      id = pingfederate_keypairs_ssl_server_key.admin_console.id
    }
  ]

  runtime_server_cert_ref = {
    id = pingfederate_keypairs_ssl_server_key.runtime.id
  }
  active_runtime_server_certs = [
    {
      id = pingfederate_keypairs_ssl_server_key.runtime.id
    }
  ]
}