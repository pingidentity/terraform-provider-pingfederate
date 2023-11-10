# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_administrative_account" "myAdministrativeAccount" {
  username    = "example"
  description = "description"
  password    = "2FederateM0re"
  roles       = ["USER_ADMINISTRATOR"]
}
