# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_server_settings_system_keys" "serverSettingsSystemKeysExample" {
  current = {
    encrypted_key_data = ""
  }
  pending = {
    encrypted_key_data = ""
  }
}