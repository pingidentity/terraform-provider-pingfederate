resource "pingfederate_administrative_account" "administrativeAccount" {
  username    = "example"
  description = "description"
  password    = "2FederateM0re"
  roles       = ["USER_ADMINISTRATOR"]
}
