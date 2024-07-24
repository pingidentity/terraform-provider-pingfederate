resource "pingfederate_key_pair_ssl_server_import" "admin_console" {
  import_id = "adminconsole"
  file_data = filebase64("./path/to/admin_console.p12")
  format    = "PKCS12"
  password  = var.runtime_server_cert_admin_console_pkcs12_password
}

resource "pingfederate_key_pair_ssl_server_import" "runtime" {
  import_id = "runtime"
  file_data = filebase64("./path/to/runtime.p12")
  format    = "PKCS12"
  password  = var.runtime_server_cert_runtime_pkcs12_password
}

resource "pingfederate_keypairs_ssl_server_settings" "sslServerSettings" {
  admin_console_cert_ref = {
    id = pingfederate_key_pair_ssl_server_import.admin_console.id
  }
  active_admin_console_certs = [
    {
      id = pingfederate_key_pair_ssl_server_import.admin_console.id
    }
  ]

  runtime_server_cert_ref = {
    id = pingfederate_key_pair_ssl_server_import.runtime.id
  }
  active_runtime_server_certs = [
    {
      id = pingfederate_key_pair_ssl_server_import.runtime.id
    }
  ]
}