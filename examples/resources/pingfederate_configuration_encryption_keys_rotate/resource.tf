// Example of using the time provider to control regular rotaion of encryption keys
resource "time_rotating" "encryption_key_rotation" {
  rotation_days = 30
}

resource "pingfederate_configuration_encryption_keys_rotate" "encryptionKeysRotate" {
  rotation_trigger_values = {
    "rotation_rfc3339" : time_rotating.encryption_key_rotation.rotation_rfc3339,
  }
}
