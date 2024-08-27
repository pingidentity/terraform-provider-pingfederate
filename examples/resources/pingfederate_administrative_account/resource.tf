resource "pingfederate_administrative_account" "administrativeAccount" {
  username    = "username"
  description = "description"
  password    = var.pingfederate_administrative_account_password
  roles       = ["USER_ADMINISTRATOR"]
}
