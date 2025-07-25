---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_data_store Data Source - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Describes a data store.
---

# pingfederate_data_store (Data Source)

Describes a data store.

## Example Usage

```terraform
data "pingfederate_data_store" "myDataStore" {
  data_store_id = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `data_store_id` (String) Unique ID for the data store.

### Read-Only

- `custom_data_store` (Attributes) A custom data store. (see [below for nested schema](#nestedatt--custom_data_store))
- `id` (String) ID of this resource.
- `jdbc_data_store` (Attributes) A JDBC data store. (see [below for nested schema](#nestedatt--jdbc_data_store))
- `ldap_data_store` (Attributes) An LDAP Data Store (see [below for nested schema](#nestedatt--ldap_data_store))
- `mask_attribute_values` (Boolean) Whether attribute values should be masked in the log.
- `ping_one_ldap_gateway_data_store` (Attributes) A PingOne LDAP Gateway data store. (see [below for nested schema](#nestedatt--ping_one_ldap_gateway_data_store))

<a id="nestedatt--custom_data_store"></a>
### Nested Schema for `custom_data_store`

Read-Only:

- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--custom_data_store--configuration))
- `name` (String) The plugin instance name.
- `parent_ref` (Attributes) The reference to this plugin's parent instance. Supported prior to PingFederate `12.0`. (see [below for nested schema](#nestedatt--custom_data_store--parent_ref))
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. (see [below for nested schema](#nestedatt--custom_data_store--plugin_descriptor_ref))
- `type` (String) The data store type.

<a id="nestedatt--custom_data_store--configuration"></a>
### Nested Schema for `custom_data_store.configuration`

Read-Only:

- `fields` (Attributes List) List of configuration fields. (see [below for nested schema](#nestedatt--custom_data_store--configuration--fields))
- `tables` (Attributes List) List of configuration tables. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables))

<a id="nestedatt--custom_data_store--configuration--fields"></a>
### Nested Schema for `custom_data_store.configuration.fields`

Read-Only:

- `encrypted_value` (String) For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined.
- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute.


<a id="nestedatt--custom_data_store--configuration--tables"></a>
### Nested Schema for `custom_data_store.configuration.tables`

Read-Only:

- `name` (String) The name of the table.
- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables--rows))

<a id="nestedatt--custom_data_store--configuration--tables--rows"></a>
### Nested Schema for `custom_data_store.configuration.tables.rows`

Read-Only:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes List) The configuration fields in the row. (see [below for nested schema](#nestedatt--custom_data_store--configuration--tables--rows--fields))

<a id="nestedatt--custom_data_store--configuration--tables--rows--fields"></a>
### Nested Schema for `custom_data_store.configuration.tables.rows.fields`

Read-Only:

- `encrypted_value` (String) For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined.
- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute.





<a id="nestedatt--custom_data_store--parent_ref"></a>
### Nested Schema for `custom_data_store.parent_ref`

Read-Only:

- `id` (String) The ID of the resource.


<a id="nestedatt--custom_data_store--plugin_descriptor_ref"></a>
### Nested Schema for `custom_data_store.plugin_descriptor_ref`

Read-Only:

- `id` (String) The ID of the resource.



<a id="nestedatt--jdbc_data_store"></a>
### Nested Schema for `jdbc_data_store`

Read-Only:

- `allow_multi_value_attributes` (Boolean) Indicates that this data store can select more than one record from a column and return the results as a multi-value attribute.
- `blocking_timeout` (Number) The amount of time in milliseconds a request waits to get a connection from the connection pool before it fails.
- `connection_url` (String) The default location of the JDBC database.
- `connection_url_tags` (Attributes Set) The set of connection URLs and associated tags for this JDBC data store. (see [below for nested schema](#nestedatt--jdbc_data_store--connection_url_tags))
- `driver_class` (String) The name of the driver class used to communicate with the source database.
- `encrypted_password` (String) The encrypted password needed to access the database.
- `idle_timeout` (Number) The length of time in minutes the connection can be idle in the pool before it is closed.
- `max_pool_size` (Number) The largest number of database connections in the connection pool for the given data store.
- `min_pool_size` (Number) The smallest number of database connections in the connection pool for the given data store.
- `name` (String) The data store name with a unique value across all data sources.
- `type` (String) The data store type.
- `user_name` (String) The name that identifies the user when connecting to the database.
- `validate_connection_sql` (String) A simple SQL statement used by PingFederate at runtime to verify that the database connection is still active and to reconnect if needed.

<a id="nestedatt--jdbc_data_store--connection_url_tags"></a>
### Nested Schema for `jdbc_data_store.connection_url_tags`

Read-Only:

- `connection_url` (String) The location of the JDBC database.
- `default_source` (Boolean) Whether this is the default connection.
- `tags` (String) Tags associated with the connection URL. At runtime, nodes will use the first `connection_url_tags` element that has a tag that matches with node.tags in the run.properties file.



<a id="nestedatt--ldap_data_store"></a>
### Nested Schema for `ldap_data_store`

Optional:

- `name` (String) The data store name with a unique value across all data sources.

Read-Only:

- `binary_attributes` (Set of String) A list of LDAP attributes to be handled as binary data.
- `bind_anonymously` (Boolean) Whether username and password are required.
- `client_tls_certificate_ref` (Attributes) The client TLS certificate used to access the data store. If specified, authentication to the data store will be done using mutual TLS. See '/keyPairs/sslClient' to manage certificates. (see [below for nested schema](#nestedatt--ldap_data_store--client_tls_certificate_ref))
- `connection_timeout` (Number) The maximum number of milliseconds that a connection attempt should be allowed to continue before returning an error.
- `create_if_necessary` (Boolean) Indicates whether temporary connections can be created when the Maximum Connections threshold is reached.
- `dns_ttl` (Number) The maximum time in milliseconds that DNS information are cached.
- `encrypted_password` (String) The encrypted password credential required to access the data store.
- `follow_ldap_referrals` (Boolean) Follow LDAP Referrals in the domain tree.
- `hostnames` (List of String) The default LDAP host names. Failover can be configured by providing multiple host names.
- `hostnames_tags` (Attributes Set) The set of host names and associated tags for this LDAP data store. (see [below for nested schema](#nestedatt--ldap_data_store--hostnames_tags))
- `ldap_dns_srv_prefix` (String) The prefix value used to discover LDAP DNS SRV record.
- `ldap_type` (String) A type that allows PingFederate to configure many provisioning settings automatically.
- `ldaps_dns_srv_prefix` (String) The prefix value used to discover LDAPS DNS SRV record.
- `max_connections` (Number) The largest number of active connections that can remain in each pool without releasing extra ones.
- `max_wait` (Number) The maximum number of milliseconds the pool waits for a connection to become available when trying to obtain a connection from the pool.
- `min_connections` (Number) The smallest number of connections that can remain in each pool, without creating extra ones.
- `read_timeout` (Number) The maximum number of milliseconds a connection waits for a response to be returned before producing an error.
- `retry_failed_operations` (Boolean) Indicates whether failed operations should be retried.
- `test_on_borrow` (Boolean) Indicates whether objects are validated before being borrowed from the pool.
- `test_on_return` (Boolean) Indicates whether objects are validated before being returned to the pool.
- `time_between_evictions` (Number) The frequency, in milliseconds, that the evictor cleans up the connections in the pool.
- `type` (String) The data store type.
- `use_dns_srv_records` (Boolean) Use DNS SRV Records to discover LDAP server information.
- `use_ssl` (Boolean) Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS).
- `use_start_tls` (Boolean) Connects to the LDAP data store using secure StartTLS encryption.
- `user_dn` (String) The username credential required to access the data store.
- `verify_host` (Boolean) Verifies that the presented server certificate includes the address to which the client intended to establish a connection.

<a id="nestedatt--ldap_data_store--client_tls_certificate_ref"></a>
### Nested Schema for `ldap_data_store.client_tls_certificate_ref`

Read-Only:

- `id` (String) The ID of the resource.


<a id="nestedatt--ldap_data_store--hostnames_tags"></a>
### Nested Schema for `ldap_data_store.hostnames_tags`

Read-Only:

- `default_source` (Boolean) Whether this is the default connection.
- `hostnames` (List of String) The LDAP host names. Failover can be configured by providing multiple host names.
- `tags` (String) Tags associated with the host names. At runtime, nodes will use the first `hostname_tags` element that has a tag that matches with node.tags in the run.properties file.



<a id="nestedatt--ping_one_ldap_gateway_data_store"></a>
### Nested Schema for `ping_one_ldap_gateway_data_store`

Read-Only:

- `binary_attributes` (Set of String) A list of LDAP attributes to be handled as binary data.
- `ldap_type` (String) A type that allows PingFederate to configure many provisioning settings automatically.
- `name` (String) The data store name with a unique value across all data sources.
- `ping_one_connection_ref` (Attributes) Reference to the PingOne connection this gateway uses. (see [below for nested schema](#nestedatt--ping_one_ldap_gateway_data_store--ping_one_connection_ref))
- `ping_one_environment_id` (String) The environment ID to which the gateway belongs.
- `ping_one_ldap_gateway_id` (String) The ID of the PingOne LDAP Gateway this data store uses.
- `type` (String) The data store type.
- `use_ssl` (Boolean) Connects to the LDAP data store using secure SSL/TLS encryption (LDAPS).
- `use_start_tls` (Boolean) Connects to the LDAP data store using StartTLS.

<a id="nestedatt--ping_one_ldap_gateway_data_store--ping_one_connection_ref"></a>
### Nested Schema for `ping_one_ldap_gateway_data_store.ping_one_connection_ref`

Read-Only:

- `id` (String) The ID of the resource.
