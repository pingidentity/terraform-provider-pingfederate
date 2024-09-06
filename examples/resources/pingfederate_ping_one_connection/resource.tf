resource "pingfederate_ping_one_connection" "example" {
  name       = "My PingOne Environment"
  credential = var.pingone_connection_credential
}