resource "pingfederate_kerberos_realm" "kerberos_realm" {
  realm_id            = "myKerberosRealm"
  kerberos_realm_name = "My Kerberos Realm"
  kerberos_username   = var.kerberos_realm_username
  kerberos_password   = var.kerberos_realm_password
  connection_type     = "DIRECT"
}

resource "pingfederate_kerberos_realm_settings" "kerberosRealmSettings" {
  depends_on = [pingfederate_kerberos_realm.kerberos_realm]

  debug_log_output              = false
  force_tcp                     = false
  kdc_retries                   = 3
  kdc_timeout                   = 3
  key_set_retention_period_mins = 610
}
