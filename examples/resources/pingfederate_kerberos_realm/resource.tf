# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_kerberos_realm" "kerberosRealmExample" {
  realm_id            = "kerberosRealm"
  kerberos_realm_name = "kerberosRealm"
  kerberos_username   = "kerberosUsername"
  kerberos_password   = "kerberosPassword"
}
