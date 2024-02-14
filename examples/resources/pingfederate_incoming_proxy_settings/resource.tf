resource "pingfederate_incoming_proxy_settings" "incomingProxySettingsExample" {
  forwarded_ip_address_header_name  = "X-Forwarded-For"
  forwarded_host_header_name        = "X-Forwarded-Host"
  forwarded_host_header_index       = "LAST"

  proxy_terminates_https_conns      = false
}