resource "pingfederate_certificates_group" "certGroup" {
  group_name = "MyGroup"
  file_data  = filebase64("path/to/my/certificate.pem")
}