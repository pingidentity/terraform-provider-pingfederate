---
page_title: "pingfederate_idp_adapter Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages an IdP adapter instance.
---

# pingfederate_idp_adapter (Resource)

Manages an IdP adapter instance.

## Example Usage

```terraform
resource "pingfederate_idp_adapter" "idpAdapterExample" {
  adapter_id = "HTMLFormPD"
  name       = "HTMLFormPD"
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.htmlform.idp.HtmlFormIdpAuthnAdapter"
  }
  attribute_mapping = {
    attribute_contract_fulfillment = {
      "entryUUID" = {
        source = {
          type = "ADAPTER"
        },
        value = "entryUUID"
      }
      "policy.action" = {
        source = {
          type = "ADAPTER"
        },
        value = "policy.action"
      },
      "username" = {
        source = {
          type = "ADAPTER"
        },
        value = "username"
      }
    },
    attribute_sources = [],
    issuance_criteria = {
      conditional_criteria = []
    }
  }
  configuration = {
    tables = [
      {
        name = "Credential Validators"
        rows = [
          {
            default_row = false
            fields = [
              {
                name  = "Password Credential Validator Instance"
                value = "pingdirectory"
              }
            ]
          }
        ]
      }
    ]
    fields = [
      {
        name  = "Challenge Retries"
        value = 3
      },
      {
        name  = "Session State",
        value = "None"
      },
      {
        name  = "Session Timeout",
        value = "60"
      },
      {
        name  = "Session Max Timeout",
        value = "480"
      },
      {
        name  = "Allow Password Changes",
        value = "false"
      },
      {
        name  = "Password Management System",
        value = ""
      },
      {
        name  = "Enable 'Remember My Username'",
        value = "false"
      },
      {
        name  = "Enable 'This is My Device'",
        value = "false"
      },
      {
        name  = "Change Password Email Notification",
        value = "false"
      },
      {
        name  = "Show Password Expiring Warning",
        value = "false"
      },
      {
        name  = "Password Reset Type",
        value = "NONE"
      },
      {
        name  = "Account Unlock",
        value = "false"
      },
      {
        name  = "Local Identity Profile",
        value = "example"
      },
      {
        name  = "Enable Username Recovery",
        value = "false"
      },
      {
        name  = "Login Template",
        value = "html.form.login.template.html"
      },
      {
        name  = "Logout Path",
        value = ""
      },
      {
        name  = "Logout Redirect",
        value = ""
      },
      {
        name  = "Logout Template",
        value = "idp.logout.success.page.template.html"
      },
      {
        name  = "Change Password Template",
        value = "html.form.change.password.template.html"
      },
      {
        name  = "Change Password Message Template",
        value = "html.form.message.template.html"
      },
      {
        name  = "Password Management System Message Template",
        value = "html.form.message.template.html"
      },
      {
        name  = "Change Password Email Template",
        value = "message-template-end-user-password-change.html"
      },
      {
        name  = "Expiring Password Warning Template",
        value = "html.form.password.expiring.notification.template.html"
      },
      {
        name  = "Threshold for Expiring Password Warning",
        value = "7"
      },
      {
        name  = "Snooze Interval for Expiring Password Warning",
        value = "24"
      },
      {
        name  = "Login Challenge Template",
        value = "html.form.login.challenge.template.html"
      },
      {
        name  = "'Remember My Username' Lifetime",
        value = "30"
      },
      {
        name  = "'This is My Device' Lifetime",
        value = "30"
      },
      {
        name  = "Allow Username Edits During Chaining",
        value = "false"
      },
      {
        name  = "Track Authentication Time",
        value = "true"
      },
      {
        name  = "Post-Password Change Re-Authentication Delay",
        value = "0"
      },
      {
        name  = "Password Reset Username Template",
        value = "forgot-password.html"
      },
      {
        name  = "Password Reset Code Template",
        value = "forgot-password-resume.html"
      },
      {
        name  = "Password Reset Template",
        value = "forgot-password-change.html"
      },
      {
        name  = "Password Reset Error Template",
        value = "forgot-password-error.html"
      },
      {
        name  = "Password Reset Success Template",
        value = "forgot-password-success.html"
      },
      {
        name  = "Account Unlock Template",
        value = "account-unlock.html"
      },
      {
        name  = "OTP Length",
        value = "8"
      },
      {
        name  = "OTP Time to Live",
        value = "10"
      },
      {
        name  = "PingID Properties",
        value = ""
      },
      {
        name  = "Require Verified Email",
        value = "false"
      },
      {
        name  = "Username Recovery Template",
        value = "username.recovery.template.html"
      },
      {
        name  = "Username Recovery Info Template",
        value = "username.recovery.info.template.html"
      },
      {
        name  = "Username Recovery Email Template",
        value = "message-template-username-recovery.html"
      },
      {
        name  = "CAPTCHA for Authentication",
        value = "false"
      },
      {
        name  = "CAPTCHA for Password change",
        value = "false"
      },
      {
        name  = "CAPTCHA for Password Reset",
        value = "false"
      },
      {
        name  = "CAPTCHA for Username recovery",
        value = "false"
      }
    ]
  }
  attribute_contract = {
    mask_ognl_values = false
    core_attributes = [
      {
        masked    = false
        name      = "policy.action"
        pseudonym = false
      },
      {
        masked    = false
        name      = "username"
        pseudonym = true
      }
    ]
    extended_attributes = [
      {
        masked    = false
        name      = "entryUUID"
        pseudonym = false
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `adapter_id` (String) The ID of the plugin instance. This field is immutable and will trigger a replacement plan if changed.
- `attribute_mapping` (Attributes) The attributes mapping from attribute sources to attribute targets. (see [below for nested schema](#nestedatt--attribute_mapping))
- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--configuration))
- `name` (String) The plugin instance name. The name can be modified once the instance is created.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. (see [below for nested schema](#nestedatt--plugin_descriptor_ref))

### Optional

- `attribute_contract` (Attributes) The list of attributes that the IdP adapter provides. (see [below for nested schema](#nestedatt--attribute_contract))
- `authn_ctx_class_ref` (String) The fixed value that indicates how the user was authenticated.
- `parent_ref` (Attributes) The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides) (see [below for nested schema](#nestedatt--parent_ref))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--attribute_mapping"></a>
### Nested Schema for `attribute_mapping`

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_contract_fulfillment))
- `attribute_sources` (Attributes List) A list of configured data stores to look up attributes from. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources))
- `issuance_criteria` (Attributes) The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled. (see [below for nested schema](#nestedatt--attribute_mapping--issuance_criteria))

<a id="nestedatt--attribute_mapping--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_mapping.attribute_contract_fulfillment`

Optional:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_contract_fulfillment--source))
- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_mapping--attribute_contract_fulfillment--source"></a>
### Nested Schema for `attribute_mapping.attribute_contract_fulfillment.source`

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.
- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.



<a id="nestedatt--attribute_mapping--attribute_sources"></a>
### Nested Schema for `attribute_mapping.attribute_sources`

Optional:

- `custom_attribute_source` (Attributes) The configured settings used to look up attributes from a custom data store. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--custom_attribute_source))
- `jdbc_attribute_source` (Attributes) The configured settings used to look up attributes from a JDBC data store. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source))
- `ldap_attribute_source` (Attributes) The configured settings used to look up attributes from a LDAP data store. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source))

<a id="nestedatt--attribute_mapping--attribute_sources--custom_attribute_source"></a>
### Nested Schema for `attribute_mapping.attribute_sources.custom_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--data_store_ref))

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--attribute_contract_fulfillment))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `filter_fields` (Attributes Set) The list of fields that can be used to filter a request to the custom data store. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--filter_fields))
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--data_store_ref"></a>
### Nested Schema for `attribute_mapping.attribute_sources.custom_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_mapping.attribute_sources.custom_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `attribute_mapping.attribute_sources.custom_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--attribute_mapping--attribute_sources--custom_attribute_source--filter_fields"></a>
### Nested Schema for `attribute_mapping.attribute_sources.custom_attribute_source.filter_fields`

Required:

- `name` (String) The name of this field.

Optional:

- `value` (String) The value of this field. Whether or not the value is required will be determined by plugin validation checks.



<a id="nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source"></a>
### Nested Schema for `attribute_mapping.attribute_sources.jdbc_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source--data_store_ref))
- `filter` (String) The JDBC WHERE clause used to query your data store to locate a user record.
- `table` (String) The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source--attribute_contract_fulfillment))
- `column_names` (Set of String) A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `schema` (String) Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source--data_store_ref"></a>
### Nested Schema for `attribute_mapping.attribute_sources.jdbc_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_mapping.attribute_sources.jdbc_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_mapping--attribute_sources--jdbc_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `attribute_mapping.attribute_sources.jdbc_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.




<a id="nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source"></a>
### Nested Schema for `attribute_mapping.attribute_sources.ldap_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--data_store_ref))
- `search_filter` (String) The LDAP filter that will be used to lookup the objects from the directory.
- `search_scope` (String) Determines the node depth of the query.
- `type` (String) The data store type of this attribute source.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `binary_attribute_settings` (Attributes Map) The advanced settings for binary LDAP attributes. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--binary_attribute_settings))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `member_of_nested_group` (Boolean) Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.
- `search_attributes` (Set of String) A list of LDAP attributes returned from search and available for mapping.

<a id="nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `attribute_mapping.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_mapping.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `attribute_mapping.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--attribute_mapping--attribute_sources--ldap_attribute_source--binary_attribute_settings"></a>
### Nested Schema for `attribute_mapping.attribute_sources.ldap_attribute_source.binary_attribute_settings`

Optional:

- `binary_encoding` (String) Get the encoding type for this attribute. If not specified, the default is BASE64.




<a id="nestedatt--attribute_mapping--issuance_criteria"></a>
### Nested Schema for `attribute_mapping.issuance_criteria`

Optional:

- `conditional_criteria` (Attributes Set) A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue. (see [below for nested schema](#nestedatt--attribute_mapping--issuance_criteria--conditional_criteria))
- `expression_criteria` (Attributes Set) A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria. (see [below for nested schema](#nestedatt--attribute_mapping--issuance_criteria--expression_criteria))

<a id="nestedatt--attribute_mapping--issuance_criteria--conditional_criteria"></a>
### Nested Schema for `attribute_mapping.issuance_criteria.conditional_criteria`

Required:

- `attribute_name` (String) The name of the attribute to use in this issuance criterion.
- `condition` (String) The condition that will be applied to the source attribute's value and the expected value. Options are `EQUALS`, `EQUALS_CASE_INSENSITIVE`, `EQUALS_DN`, `NOT_EQUAL`, `NOT_EQUAL_CASE_INSENSITIVE`, `NOT_EQUAL_DN`, `MULTIVALUE_CONTAINS`, `MULTIVALUE_CONTAINS_CASE_INSENSITIVE`, `MULTIVALUE_CONTAINS_DN`, `MULTIVALUE_DOES_NOT_CONTAIN`, `MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE`, `MULTIVALUE_DOES_NOT_CONTAIN_DN`.
- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_mapping--issuance_criteria--conditional_criteria--source))
- `value` (String) The expected value of this issuance criterion.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.

<a id="nestedatt--attribute_mapping--issuance_criteria--conditional_criteria--source"></a>
### Nested Schema for `attribute_mapping.issuance_criteria.conditional_criteria.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--attribute_mapping--issuance_criteria--expression_criteria"></a>
### Nested Schema for `attribute_mapping.issuance_criteria.expression_criteria`

Required:

- `expression` (String) The OGNL expression to evaluate.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.




<a id="nestedatt--configuration"></a>
### Nested Schema for `configuration`

Optional:

- `fields` (Attributes Set) List of configuration fields. (see [below for nested schema](#nestedatt--configuration--fields))
- `sensitive_fields` (Attributes Set) List of sensitive configuration fields. (see [below for nested schema](#nestedatt--configuration--sensitive_fields))
- `tables` (Attributes List) List of configuration tables. (see [below for nested schema](#nestedatt--configuration--tables))

Read-Only:

- `fields_all` (Attributes Set) List of configuration fields. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--configuration--fields_all))
- `tables_all` (Attributes List) List of configuration tables. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--configuration--tables_all))

<a id="nestedatt--configuration--fields"></a>
### Nested Schema for `configuration.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--configuration--sensitive_fields"></a>
### Nested Schema for `configuration.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.

Optional:

- `encrypted_value` (String) For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. Either this attribute or `value` must be specified.
- `value` (String, Sensitive) The sensitive value for the configuration field. Either this attribute or `encrypted_value` must be specified`.


<a id="nestedatt--configuration--tables"></a>
### Nested Schema for `configuration.tables`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--configuration--tables--rows))

<a id="nestedatt--configuration--tables--rows"></a>
### Nested Schema for `configuration.tables.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--configuration--tables--rows--fields))
- `sensitive_fields` (Attributes Set) The sensitive configuration fields in the row. (see [below for nested schema](#nestedatt--configuration--tables--rows--sensitive_fields))

<a id="nestedatt--configuration--tables--rows--fields"></a>
### Nested Schema for `configuration.tables.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--configuration--tables--rows--sensitive_fields"></a>
### Nested Schema for `configuration.tables.rows.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.

Optional:

- `encrypted_value` (String) For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. Either this attribute or `value` must be specified.
- `value` (String, Sensitive) The sensitive value for the configuration field. Either this attribute or `encrypted_value` must be specified`.




<a id="nestedatt--configuration--fields_all"></a>
### Nested Schema for `configuration.fields_all`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--configuration--tables_all"></a>
### Nested Schema for `configuration.tables_all`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--configuration--tables_all--rows))

<a id="nestedatt--configuration--tables_all--rows"></a>
### Nested Schema for `configuration.tables_all.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--configuration--tables_all--rows--fields))

<a id="nestedatt--configuration--tables_all--rows--fields"></a>
### Nested Schema for `configuration.tables_all.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.





<a id="nestedatt--plugin_descriptor_ref"></a>
### Nested Schema for `plugin_descriptor_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--attribute_contract"></a>
### Nested Schema for `attribute_contract`

Required:

- `core_attributes` (Attributes Set) A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type. (see [below for nested schema](#nestedatt--attribute_contract--core_attributes))

Optional:

- `extended_attributes` (Attributes Set) A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them. (see [below for nested schema](#nestedatt--attribute_contract--extended_attributes))
- `mask_ognl_values` (Boolean) Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to `false`.
- `unique_user_key_attribute` (String) The attribute to use for uniquely identify a user's authentication sessions.

Read-Only:

- `core_attributes_all` (Attributes Set) A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--attribute_contract--core_attributes_all))

<a id="nestedatt--attribute_contract--core_attributes"></a>
### Nested Schema for `attribute_contract.core_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `masked` (Boolean) Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.
- `pseudonym` (Boolean) Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.


<a id="nestedatt--attribute_contract--extended_attributes"></a>
### Nested Schema for `attribute_contract.extended_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `masked` (Boolean) Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.
- `pseudonym` (Boolean) Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.


<a id="nestedatt--attribute_contract--core_attributes_all"></a>
### Nested Schema for `attribute_contract.core_attributes_all`

Read-Only:

- `masked` (Boolean) Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.
- `name` (String) The name of this attribute.
- `pseudonym` (Boolean) Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.



<a id="nestedatt--parent_ref"></a>
### Nested Schema for `parent_ref`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "idpAdapterId" should be the ID of the IdP adapter to be imported

```shell
terraform import pingfederate_idp_adapter.idpAdapter idpAdapterId
```