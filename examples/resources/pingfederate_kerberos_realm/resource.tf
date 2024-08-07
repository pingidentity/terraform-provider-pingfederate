resource "pingfederate_kerberos_realm" "kerberos_realm" {
  realm_id            = "myKerberosRealm"
  kerberos_realm_name = "My Kerberos Realm"
  kerberos_username   = var.kerberos_realm_username
  kerberos_password   = var.kerberos_realm_password
  connection_type     = "DIRECT"
}
