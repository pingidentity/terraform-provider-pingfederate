resource "pingfederate_metadata_url" "metadataUrl" {
  url_id             = "myUrlId"
  name               = "My Url"
  url                = "https://example.com"
  validate_signature = true
  x509_file = {
    # Insert base64-encoded cert data here
    file_data = ""
  }
}