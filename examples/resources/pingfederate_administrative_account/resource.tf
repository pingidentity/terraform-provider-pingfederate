resource "pingfederate_administrative_account" "myAdministrativeAccount" {
  username    = "example"
  description = "description"
  password    = "2FederateM0re"
  roles       = ["USER_ADMINISTRATOR"]
}
