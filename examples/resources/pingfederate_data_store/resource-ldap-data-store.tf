resource "pingfederate_data_store" "pingDirectoryLdapDataStore" {
  data_store_id = "pingDirectoryLdapDataStore"
  ldap_data_store = {
    ldap_type          = "PING_DIRECTORY"
    bind_anonymously   = false
    user_dn            = var.ping_directory_user_dn
    password           = var.ping_directory_password
    use_ssl            = false
    use_dns_srv_record = false
    name               = "pingDirectoryLdapDataStore"
    hostnames = [
      "pingdirectory:1389"
    ]
    hostnames_tags = [
      {
        hostnames = [
          "pingdirectory:1389"
        ]
        default_source = true
      }
    ]
    test_on_borrow         = true
    test_on_return         = false
    create_if_necessary    = true
    verify_host            = true
    min_connections        = 10
    max_connections        = 100
    max_wait               = -1
    time_between_evictions = 6000
    read_timeout           = 300
    connection_timeout     = 300
    dns_ttl                = 6000
  }
}
