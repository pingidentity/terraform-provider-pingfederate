// Example of using the time provider to control regular rotaion of system keys
resource "time_rotating" "system_key_rotation" {
  rotation_days = 30
}

resource "pingfederate_server_settings_system_keys_rotate" "systemKeysRotate" {
  rotation_trigger_values = {
    "rotation_rfc3339" : time_rotating.system_key_rotation.rotation_rfc3339,
  }
}