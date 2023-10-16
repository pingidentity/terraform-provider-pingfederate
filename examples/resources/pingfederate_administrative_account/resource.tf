resource "pingfederate_administrative_account" "myAdministrativeAccount" {
  username = "example"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}