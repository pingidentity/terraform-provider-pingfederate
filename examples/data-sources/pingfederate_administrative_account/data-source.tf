resource "pingfederate_administrative_account" "myAdministrativeAccount" {
  username = "data-source-example"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}

data "pingfederate_administrative_account" "myAdministrativeAccount" {
  id = pingfederate_administrative_account.myAdministrativeAccount.username
}
resource "pingfederate_administrative_account" "administrativeAccountExample" {
  username = "${data.pingfederate_administrative_account.myAdministrativeAccount.id}-new"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}