resource "pingfederate_license" "myLicense" {
  id = "id"
  # this property needs to contain base64 encoded value of your license.
  file_data = ""
}

data "pingfederate_license" "myLicense" {
  id = pingfederate_license.myLicense.id
  depends_on = [
    pingfederate_license.myLicense
  ]
}
resource "pingfederate_license" "licenseExample" {
  id = "${data.pingfederate_license.myLicense.id}2"
  # this property needs to contain base64 encoded value of your license.
  file_data = ""
}