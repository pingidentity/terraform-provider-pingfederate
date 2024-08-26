resource "pingfederate_server_settings_system_keys" "systemKeys" {
  current = {
    encrypted_key_data = var.current_server_encrypted_key_data
  }
  pending = {
    encrypted_key_data = var.pending_server_encrypted_key_data
  }
}