# this resource does not support import
resource "pingfederate_license" "myLicense" {
  # this property needs to contain base64 encoded value of your license.
  file_data = ""
}
