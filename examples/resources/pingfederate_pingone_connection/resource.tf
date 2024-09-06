resource "pingfederate_pingone_connection" "example" {
  name       = "My PingOne Environment"
  credential = var.pingone_connection_credential
}