resource "pingfederate_incoming_proxy_settings" "incomingProxySettingsExample" {
  forwarded_ip_address_header_name  = "X-Forwarded-For"
  forwarded_ip_address_header_index = "FIRST"
  forwarded_host_header_name        = "X-Forwarded-Host"
  forwarded_host_header_index       = "LAST"
  client_cert_ssl_header_name       = "X-Client-Cert"
  client_cert_chain_ssl_header_name = "X-Client-Cert-Chain"
  proxy_terminates_https_conns      = true
}