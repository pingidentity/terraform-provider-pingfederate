resource "pingfederate_certificate_ca" "example" {
  file_data = filebase64("myCA.pem")
}
