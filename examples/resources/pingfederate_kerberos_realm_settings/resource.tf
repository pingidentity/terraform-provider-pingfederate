resource "pingfederate_kerberos_realm_settings" "kerberosRealmSettings" {
  debug_log_output              = false
  force_tcp                     = false
  kdc_retries                   = 3
  kdc_timeout                   = 3
  key_set_retention_period_mins = 610
}