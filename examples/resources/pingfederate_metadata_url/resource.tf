resource "pingfederate_metadata_url" "metadataUrl" {
  url_id             = "myUrlId"
  name               = "My Metadata Url"
  url                = "https://bxretail.org/metadata"
  validate_signature = true
  x509_file = {
    file_data = filebase64("./assets/my-certificate.pem")
  }
}