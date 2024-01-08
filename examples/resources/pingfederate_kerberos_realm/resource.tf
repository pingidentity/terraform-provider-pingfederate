# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #

resource "pingfederate_kerberos_realm" "myKerberosRealmExample" {
  realm_id            = "myKerberosRealm"
  kerberos_realm_name = "myKerberosRealm"
  kerberos_username   = "myKerberosUsername"
  kerberos_password   = "myKerberosPassword"
}
