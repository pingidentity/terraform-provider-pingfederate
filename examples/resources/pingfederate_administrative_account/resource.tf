resource "pingfederate_administrative_account" "administrativeAccountExample" {
  username = "example"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}

data "pingfederate_administrative_account" "account1"{
  username = pingfederate_administrative_account.administrativeAccountExample.username
}

resource "pingfederate_administrative_account" "administrativeAccountExample2" {
  username = "${data.pingfederate_administrative.account1.username}2"
  password = "2FederateM0re"
  roles    = ["USER_ADMINISTRATOR"]
}