---
page_title: "pingfederate_data_store Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to create and manage data stores.
---

# pingfederate_data_store (Resource)

Resource to create and manage data stores.

## Example Usage - LDAP data store

```terraform
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
```

## Example Usage - PingOne LDAP gateway data store

```terraform
resource "pingfederate_pingone_connection" "example" {
  name       = "My PingOne Tenant"
  credential = var.pingone_connection_credential
}

resource "pingfederate_data_store" "pingOneLdapDataStore" {
  ping_one_ldap_gateway_data_store = {
    name      = "PingOne LDAP Gateway"
    ldap_type = "PING_DIRECTORY"

    ping_one_connection_ref = {
      id = pingfederate_pingone_connection.example.id
    }

    ping_one_environment_id  = var.pingone_environment_id
    ping_one_ldap_gateway_id = var.pingone_gateway_id

    use_ssl = true
  }
}
```

## Example Usage - JDBC data store

```terraform
resource "pingfederate_data_store" "jdbcDataStore" {
  jdbc_data_store = {
    name = "JDBC Data Store"

    connection_url = "jdbc:sqlserver://localhost;encrypt=true;integratedSecurity=true;"
    driver_class   = "org.hsqldb.jdbcDriver"

    user_name = var.jdbc_data_store_username
    password  = var.jdbc_data_store_password

    allow_multi_value_attributes = false

    connection_url_tags = [
      {
        connection_url = "jdbc:sqlserver://localhost;encrypt=true;integratedSecurity=true;"
        default_source = true
      }
    ]

    min_pool_size    = 10
    max_pool_size    = 100
    blocking_timeout = 5000
    idle_timeout     = 5
  }
}
```

## Example Usage - PingOne Directory data store

```terraform
resource "pingfederate_pingone_connection" "example" {
  name       = "My PingOne Tenant"
  credential = var.pingone_connection_credential
}

resource "pingfederate_data_store" "pingOneDataStore" {
  custom_data_store = {
    name = format("PingOne Data Store (%s)", var.pingone_environment_name)

    plugin_descriptor_ref = {
      id = "com.pingidentity.plugins.datastore.p14c.PingOneForCustomersDataStore"
    }

    configuration = {
      tables = [
        {
          name = "Custom Attributes Details",
          rows = [
            {
              fields = [
                {
                  name  = "Local Attribute",
                  value = "local_attribute"
                },
                {
                  name  = "PingOne for Customers Attribute",
                  value = "/pingone_attribute"
                }
              ],
              defaultRow = false
            }
          ]
        }
      ],
      fields = [
        {
          name  = "PingOne Environment",
          value = format("%s|%s", pingfederate_pingone_connection.example.id, var.pingone_environment_id)
        },
        {
          name  = "Connection Timeout",
          value = "10000"
        },
        {
          name  = "Retry Request",
          value = "true"
        },
        {
          name  = "Maximum Retries Limit",
          value = "5"
        },
        {
          name  = "Retry Error Codes",
          value = "429"
        },
        {
          name  = "Proxy Settings",
          value = "System Defaults"
        },
        {
          name  = "Custom Proxy Host",
          value = ""
        },
        {
          name  = "Custom Proxy Port",
          value = ""
        }
      ]
    }
    mask_attribute_values = false
  }
}
```

## Example Usage - Custom REST API data store

```terraform
resource "pingfederate_data_store" "customDataStore" {
  custom_data_store = {
    name = "Custom REST Data Store"

    plugin_descriptor_ref = {
      id = "com.pingidentity.pf.datastore.other.RestDataSourceDriver"
    }

    configuration = {
      tables = [
        {
          name = "Base URLs and Tags"
          rows = [
            {
              fields = [
                {
                  name  = "Base URL"
                  value = "https://my_rest_datasource.bxretail.org/api/v1/users"
                },
                {
                  name  = "Tags"
                  value = "production"
                }
              ],
              default_row = true
            }
          ]
        },
        {
          name = "HTTP Request Headers"
          rows = [
            {
              fields = [
                {
                  name  = "Header Name"
                  value = "header"
                },
                {
                  name  = "Header Value"
                  value = "header_value"
                }
              ],
              default_row = false
            }
          ]
        },
        {
          name = "Attributes"
          rows = [
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "givenName"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/givenName"
                }
              ],
              default_row = false
            },
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "familyName"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/familyName"
                }
              ],
              default_row = false
            },
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "email"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/email"
                }
              ],
              default_row = false
            },
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "password"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "/password"
                }
              ],
              default_row = false
            }
          ]
        }
      ],
      fields = [
        {
          name  = "Authentication Method"
          value = "OAuth 2.0 Bearer Token"
        },
        {
          name  = "HTTP Method"
          value = "GET"
        },
        {
          name  = "Username"
          value = var.rest_data_store_basic_auth_username
        },
        {
          name  = "Password Reference"
          value = ""
        },
        {
          name  = "OAuth Token Endpoint"
          value = "https://authservices.bxretail.org/as/token"
        },
        {
          name  = "OAuth Scope"
          value = "restapiscope"
        },
        {
          name  = "Client ID"
          value = var.rest_data_store_oauth2_client_id
        },
        {
          name  = "Client Secret Reference"
          value = ""
        },
        {
          name  = "Enable HTTPS Hostname Verification"
          value = "true"
        },
        {
          name  = "Read Timeout (ms)"
          value = "10000"
        },
        {
          name  = "Connection Timeout (ms)"
          value = "10000"
        },
        {
          name  = "Max Payload Size (KB)"
          value = "1024"
        },
        {
          name  = "Retry Request"
          value = "true"
        },
        {
          name  = "Maximum Retries Limit"
          value = "5"
        },
        {
          name  = "Retry Error Codes"
          value = "429"
        },
        {
          name  = "Test Connection URL"
          value = "https://my_rest_datasource.bxretail.org/api/v1/connectiontest"
        },
        {
          name  = "Test Connection Body"
          value = "{\"foo\":\"bar\"}"
        }
      ]
      sensitive_fields = [
        {
          name  = "Password"
          value = var.rest_data_store_basic_auth_password
        },
        {
          name  = "Client Secret"
          value = var.rest_data_store_oauth2_client_secret
        }
      ]
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `custom_data_store` (Attributes) A custom data store. (see [below for nested schema](#nestedatt--custom_data_store))
- `data_store_id` (String) The persistent, unique ID for the data store. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified.
- `jdbc_data_store` (Attributes) A JDBC data store. (see [below for nested schema](#nestedatt--jdbc_data_store))
- `ldap_data_store` (Attributes) An LDAP Data Store (see [below for nested schema](#nestedatt--ldap_data_store))
- `mask_attribute_values` (Boolean) Whether attribute values should be masked in the log. Default value is `false`.
- `ping_one_ldap_gateway_data_store` (Attributes) A PingOne LDAP Gateway data store. (see [below for nested schema](#nestedatt--ping_one_ldap_gateway_data_store))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--custom_data_store"></a>
### Nested Schema for `custom_data_store`

Required:

- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--custom_data_store--configuration))
- `name` (String) The plugin instance name.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. (see [below for nested schema](#nestedatt--custom_data_store--plugin_descriptor_ref))

Optional:

- `parent_ref` (Attributes) The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides). Supported prior to PingFederate `12.0`. (see [below for nested schema](#nestedatt--custom_data_store--parent_ref))

Read-Only:

- `type` (String) The data store type.

<a id="nestedatt--custom_data_store--configuration"></a>
### Nested Schema for `custom_data_store.configuration`

Optional:

- `fields` (Attributes Set) List of configuration fields. (see [below for nested schema](#nestedatt--custom_data_store--configuration--fields))
- `sensitive_fields` (Attributes Set) List of sensitive configuration fields. (see [below for nested schema](#nestedatt--custom_data_store--configuration--sensitive_fields))
- `tables` (Attributes List) List of configuration tables. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables))

Read-Only:

- `fields_all` (Attributes Set) List of configuration fields. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--custom_data_store--configuration--fields_all))
- `tables_all` (Attributes List) List of configuration tables. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables_all))

<a id="nestedatt--custom_data_store--configuration--fields"></a>
### Nested Schema for `custom_data_store.configuration.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The value for the configuration field.


<a id="nestedatt--custom_data_store--configuration--sensitive_fields"></a>
### Nested Schema for `custom_data_store.configuration.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The sensitive value for the configuration field.


<a id="nestedatt--custom_data_store--configuration--tables"></a>
### Nested Schema for `custom_data_store.configuration.tables`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables--rows))

<a id="nestedatt--custom_data_store--configuration--tables--rows"></a>
### Nested Schema for `custom_data_store.configuration.tables.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables--rows--fields))
- `sensitive_fields` (Attributes Set) The sensitive configuration fields in the row. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables--rows--sensitive_fields))

<a id="nestedatt--custom_data_store--configuration--tables--rows--fields"></a>
### Nested Schema for `custom_data_store.configuration.tables.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The value for the configuration field.


<a id="nestedatt--custom_data_store--configuration--tables--rows--sensitive_fields"></a>
### Nested Schema for `custom_data_store.configuration.tables.rows.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The sensitive value for the configuration field.




<a id="nestedatt--custom_data_store--configuration--fields_all"></a>
### Nested Schema for `custom_data_store.configuration.fields_all`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The value for the configuration field.


<a id="nestedatt--custom_data_store--configuration--tables_all"></a>
### Nested Schema for `custom_data_store.configuration.tables_all`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables_all--rows))

<a id="nestedatt--custom_data_store--configuration--tables_all--rows"></a>
### Nested Schema for `custom_data_store.configuration.tables_all.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables_all--rows--fields))

<a id="nestedatt--custom_data_store--configuration--tables_all--rows--fields"></a>
### Nested Schema for `custom_data_store.configuration.tables_all.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The value for the configuration field.





<a id="nestedatt--custom_data_store--plugin_descriptor_ref"></a>
### Nested Schema for `custom_data_store.plugin_descriptor_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--custom_data_store--parent_ref"></a>
### Nested Schema for `custom_data_store.parent_ref`

Required:

- `id` (String) The ID of the resource.



<a id="nestedatt--jdbc_data_store"></a>
### Nested Schema for `jdbc_data_store`

Required:

- `driver_class` (String) The name of the driver class used to communicate with the source database.

Optional:

- `allow_multi_value_attributes` (Boolean) Indicates that this data store can select more than one record from a column and return the results as a multi-value attribute. Default value is `false`.
- `blocking_timeout` (Number) The amount of time in milliseconds a request waits to get a connection from the connection pool before it fails. The default value is `5000` milliseconds.
- `connection_url` (String) The default location of the JDBC database. This field is required if `connection_url_tags` is not specified.
- `connection_url_tags` (Attributes Set) The set of connection URLs and associated tags for this JDBC data store. This is required if 'connection_url' is not provided. (see [below for nested schema](#nestedatt--jdbc_data_store--connection_url_tags))
- `idle_timeout` (Number) The length of time in minutes the connection can be idle in the pool before it is closed. The default value is `5` minutes.
- `max_pool_size` (Number) The largest number of database connections in the connection pool for the given data store. The default value is `100`.
- `min_pool_size` (Number) The smallest number of database connections in the connection pool for the given data store. The default value is `10`.
- `name` (String) The data store name with a unique value across all data sources. Defaults to a combination of the `connection_url` and `username`.
- `password` (String, Sensitive) The password needed to access the database.
- `user_name` (String) The name that identifies the user when connecting to the database.
- `validate_connection_sql` (String) A simple SQL statement used by PingFederate at runtime to verify that the database connection is still active and to reconnect if needed.

Read-Only:

- `type` (String) The data store type.

<a id="nestedatt--jdbc_data_store--connection_url_tags"></a>
### Nested Schema for `jdbc_data_store.connection_url_tags`

Required:

- `connection_url` (String) The location of the JDBC database.

Optional:

- `default_source` (Boolean) Whether this is the default connection. Default value is `false`.
- `tags` (String) Tags associated with the `connection_url`. At runtime, nodes will use the first `connection_url_tags` element that has a tag that matches with `node.tags` in the run.properties file.



<a id="nestedatt--ldap_data_store"></a>
### Nested Schema for `ldap_data_store`

Required:

- `ldap_type` (String) A type that allows PingFederate to configure many provisioning settings automatically. The `UNBOUNDID_DS` type has been deprecated, please use the `PING_DIRECTORY` type instead. Supported values are `ACTIVE_DIRECTORY`, `ORACLE_DIRECTORY_SERVER`, `ORACLE_UNIFIED_DIRECTORY`, `PING_DIRECTORY`, `GENERIC`.

Optional:

- `binary_attributes` (Set of String) A list of LDAP attributes to be handled as binary data.
- `bind_anonymously` (Boolean) Whether username and password are required. If `true`, then `user_dn` and `client_tls_certificate_ref` cannot be set. The default value is `false`.
- `client_tls_certificate_ref` (Attributes) The client TLS certificate used to access the data store. If specified, authentication to the data store will be done using mutual TLS. See '/keyPairs/sslClient' to manage certificates. Supported in PF version `11.3` or later. In order to use this authentication method, you must set either `use_start_tls` or `use_ssl` to `true`. Mutually exclusive with `bind_anonymously` and `user_dn` (see [below for nested schema](#nestedatt--ldap_data_store--client_tls_certificate_ref))
- `connection_timeout` (Number) The maximum number of milliseconds that a connection attempt should be allowed to continue before returning an error. A value of `-1` causes the pool to wait indefinitely. Defaults to `0`.
- `create_if_necessary` (Boolean) Indicates whether temporary connections can be created when the Maximum Connections threshold is reached. Default value is `false`.
- `dns_ttl` (Number) The maximum time in milliseconds that DNS information are cached. Defaults to `0`.
- `follow_ldap_referrals` (Boolean) Follow LDAP Referrals in the domain tree. The default value is `false`. This property does not apply to PingDirectory as this functionality is configured in PingDirectory.
- `hostnames` (List of String) The default LDAP host names. This field is required if `hostnames_tags` is not specified. Failover can be configured by providing multiple host names.
- `hostnames_tags` (Attributes Set) The set of host names and associated tags for this LDAP data store. This is required if 'hostnames' is not provided. (see [below for nested schema](#nestedatt--ldap_data_store--hostnames_tags))
- `ldap_dns_srv_prefix` (String) The prefix value used to discover LDAP DNS SRV record. Defaults to `_ldap._tcp`.
- `max_connections` (Number) The largest number of active connections that can remain in each pool without releasing extra ones. Defaults to `100`.
- `max_wait` (Number) The maximum number of milliseconds the pool waits for a connection to become available when trying to obtain a connection from the pool. Setting a value of `-1` causes the pool not to wait at all and to either create a new connection or produce an error (when no connections are available). Defaults to `-1`.
- `min_connections` (Number) The smallest number of connections that can remain in each pool, without creating extra ones. Defaults to `10`.
- `name` (String) The data store name with a unique value across all data sources. Defaults to a combination of the values of `hostnames` and `user_dn`.
- `password` (String, Sensitive) The password credential required to access the data store. Requires `user_dn` to be set.
- `read_timeout` (Number) The maximum number of milliseconds a connection waits for a response to be returned before producing an error. A value of `-1` causes the connection to wait indefinitely. Defaults to `0`.
- `retry_failed_operations` (Boolean) Indicates whether failed operations should be retried. The default is `false`. Supported in PF version `11.3` or later.
- `test_on_borrow` (Boolean) Indicates whether objects are validated before being borrowed from the pool. Default value is `false`.
- `test_on_return` (Boolean) Indicates whether objects are validated before being returned to the pool. Default value is `false`.
- `time_between_evictions` (Number) The frequency, in milliseconds, that the evictor cleans up the connections in the pool. A value of `-1` disables the evictor. Defaults to `0`.
- `use_dns_srv_records` (Boolean) Use DNS SRV Records to discover LDAP server information. The default value is `false`.
- `use_ssl` (Boolean) Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS). The default value is `false`.
- `use_start_tls` (Boolean) Connects to the LDAP data store using secure StartTLS encryption. The default value is `false`.
- `user_dn` (String) The username credential required to access the data store. Mutually exclusive with `bind_anonymously` and `client_tls_certificate_ref`. `password` must also be set to use this attribute.
- `verify_host` (Boolean) Verifies that the presented server certificate includes the address to which the client intended to establish a connection. Defaults to `true`.

Read-Only:

- `type` (String) The data store type.

<a id="nestedatt--ldap_data_store--client_tls_certificate_ref"></a>
### Nested Schema for `ldap_data_store.client_tls_certificate_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--ldap_data_store--hostnames_tags"></a>
### Nested Schema for `ldap_data_store.hostnames_tags`

Required:

- `hostnames` (List of String) The LDAP host names. Failover can be configured by providing multiple host names.

Optional:

- `default_source` (Boolean) Whether this is the default connection. Defaults to `false`.
- `tags` (String) Tags associated with the host names. At runtime, nodes will use the first `hostnames_tags` element that has a tag that matches with node.tags in the run.properties file.



<a id="nestedatt--ping_one_ldap_gateway_data_store"></a>
### Nested Schema for `ping_one_ldap_gateway_data_store`

Required:

- `ldap_type` (String) A type that allows PingFederate to configure many provisioning settings automatically. The value is validated against the LDAP gateway configuration in PingOne unless the provider setting 'x_bypass_external_validation_header' is set to `true`. Supported values are `ACTIVE_DIRECTORY`, `ORACLE_DIRECTORY_SERVER`, `ORACLE_UNIFIED_DIRECTORY`, `UNBOUNDID_DS`, `PING_DIRECTORY`, and `GENERIC`.
- `ping_one_connection_ref` (Attributes) Reference to the PingOne connection this gateway uses. (see [below for nested schema](#nestedatt--ping_one_ldap_gateway_data_store--ping_one_connection_ref))
- `ping_one_environment_id` (String) The environment ID to which the gateway belongs.
- `ping_one_ldap_gateway_id` (String) The ID of the PingOne LDAP Gateway this data store uses.

Optional:

- `binary_attributes` (Set of String) A list of LDAP attributes to be handled as binary data.
- `name` (String) The data store name with a unique value across all data sources. Defaults to `ping_one_connection_ref.id` plus `ping_one_environment_id` plus `ping_one_ldap_gateway_id`, each separated by `:`.
- `use_ssl` (Boolean) Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS). The default value is `false`. The value is validated against the LDAP gateway configuration in PingOne unless the provider setting 'x_bypass_external_validation_header' is set to `true`.
- `use_start_tls` (Boolean) Connects to the LDAP data store using StartTLS. The default value is `false`. The value is validated against the LDAP gateway configuration in PingOne unless the provider setting 'x_bypass_external_validation_header' is set to `true`. Supported in PingFederate `12.1` and later.

Read-Only:

- `type` (String) The data store type.

<a id="nestedatt--ping_one_ldap_gateway_data_store--ping_one_connection_ref"></a>
### Nested Schema for `ping_one_ldap_gateway_data_store.ping_one_connection_ref`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "data-store-id" should be the id of the Data Store to be imported

```shell
terraform import pingfederate_data_store.dataStore data-store-id
```