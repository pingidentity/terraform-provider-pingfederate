resource "pingfederate_certificates_group" "certGroup" {
  group_name = "MyGroup"
  group_id   = "mygroupid"
  # Include base64-encoded cert data here
  file_data = ""
}