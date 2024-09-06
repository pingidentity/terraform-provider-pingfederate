resource "pingfederate_data_store" "pingDirectoryLdapDataStore" {
  ldap_data_store = {
    name      = "PingDirectory LDAP Data Store"
    ldap_type = "PING_DIRECTORY"

    bind_anonymously = false
    user_dn          = var.pingdirectory_bind_dn
    password         = var.pingdirectory_bind_dn_password

    use_ssl       = true
    use_start_tls = false

    use_dns_srv_record = false

    hostnames = [
      "pingdirectory:636"
    ]
    hostnames_tags = [
      {
        hostnames = [
          "pingdirectory:636"
        ]
        default_source = true
      }
    ]

    test_on_borrow         = false
    test_on_return         = false
    create_if_necessary    = true
    verify_host            = true
    min_connections        = 10
    max_connections        = 100
    max_wait               = -1
    time_between_evictions = 60000
    read_timeout           = 3000
    connection_timeout     = 3000
    dns_ttl                = 60000
  }
}