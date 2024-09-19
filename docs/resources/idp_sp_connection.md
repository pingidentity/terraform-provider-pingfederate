---
page_title: "pingfederate_idp_sp_connection Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages an IdP SP Connection
---

# pingfederate_idp_sp_connection (Resource)

Manages an IdP SP Connection

## Example Usage - WS Fed SP Browser SSO

```terraform
resource "pingfederate_idp_sp_connection" "wsFedSpBrowserSSOExample" {
  connection_id = "connectionId"
  name          = "wsfedspconn1"
  entity_id     = "wsfed1"
  active        = false
  contact_info = {
    company = "Example Corp"
  }
  base_url               = "https://localhost:9031"
  logging_mode           = "STANDARD"
  virtual_entity_ids     = []
  connection_target_type = "STANDARD"
  credentials = {
    certs = []
    signing_settings = {
      signing_key_pair_ref = {
        id = "exampleKeyId"
      }
      include_raw_key_in_signature = false
      include_cert_in_signature    = false
      algorithm                    = "SHA256withRSA"
    }
  }
  sp_browser_sso = {
    protocol                      = "WSFED"
    always_sign_artifact_response = false
    sso_service_endpoints = [
      {
        url = "/sp/prpwrong.wsf"
      }
    ]
    sp_ws_fed_identity_mapping = "EMAIL_ADDRESS"
    assertion_lifetime = {
      minutes_before = 5
      minutes_after  = 5
    }
    attribute_contract = {
      core_attributes = [
        {
          name = "SAML_SUBJECT"
        }
      ]
      extended_attributes = []
    }
    adapter_mappings = [
      {
        attribute_sources = []
        attribute_contract_fulfillment = {
          "SAML_SUBJECT" = {
            source = {
              type = "ADAPTER"
            }
            value = "subject"
          }
        }
        issuance_criteria = {
          conditional_criteria = []
        }
        restrict_virtual_entity_ids   = false
        restricted_virtual_entity_ids = []
        idp_adapter_ref = {
          id = "OTIdPJava"
        }
        abort_sso_transaction_as_fail_safe = false
      }
    ]
    authentication_policy_contract_assertion_mappings = []
    ws_fed_token_type                                 = "SAML11"
    ws_trust_version                                  = "WSTRUST12"
  }
}
```

## Example Usage - SAML SP Browser SSO

```terraform
resource "pingfederate_idp_sp_connection" "samlSpBrowserSSOExample" {
  connection_id      = "connection"
  name               = "connection"
  entity_id          = "entity"
  active             = true
  contact_info       = {}
  base_url           = "https://localhost:9032"
  logging_mode       = "STANDARD"
  virtual_entity_ids = []
  credentials = {
    certs = []
    signing_settings = {
      signing_key_pair_ref = {
        id = "signingKey"
      }
      include_raw_key_in_signature = false
      include_cert_in_signature    = false
      algorithm                    = "SHA256withRSA"
    }
  }
  sp_browser_sso = {
    protocol                      = "SAML20"
    require_signed_authn_requests = false
    sp_saml_identity_mapping      = "STANDARD"
    sign_assertions               = false
    authentication_policy_contract_assertion_mappings = [
      {
        abort_sso_transaction_as_fail_safe = false
        authentication_policy_contract_ref = {
          id = "contractId"
        }
        restricted_virtual_entity_ids = []
        attribute_contract_fulfillment = {
          "SAML_SUBJECT" = {
            source = {
              type = "AUTHENTICATION_POLICY_CONTRACT"
            }
            value = "subject"
          }
        }
        restrict_virtual_entity_ids = false
        attribute_sources           = []
        issuance_criteria = {
          conditional_criteria = []
        }
      }
    ]
    encryption_policy = {
      encrypt_slo_subject_name_id   = false
      encrypt_assertion             = false
      encrypted_attributes          = []
      slo_subject_name_id_encrypted = false
    }
    enabled_profiles = [
      "IDP_INITIATED_SSO"
    ]
    sign_response_as_required = true
    sso_service_endpoints = [
      {
        is_default = true
        binding    = "POST"
        index      = 0
        url        = "https://httpbin.org/anything"
      }
    ]
    adapter_mappings = []
    assertion_lifetime = {
      minutes_after  = 5
      minutes_before = 5
    }
    attribute_contract = {
      core_attributes = [
        {
          name_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
          name        = "SAML_SUBJECT"
        }
      ]
      extended_attributes = []
    }
  }
}
```

## Example Usage - Outbound Provision

```terraform
resource "pingfederate_idp_sp_connection" "outboundProvisionExample" {
  connection_id = "connectionId"
  name          = "PingOne Connector"
  entity_id     = "entity"
  active        = false
  contact_info = {
    company = "Example Corp"
  }
  base_url               = "https://api.pingone.com/v5"
  logging_mode           = "STANDARD"
  connection_target_type = "STANDARD"
  credentials = {
    certs = []
    signing_settings = {
      signing_key_pair_ref = {
        id = ""
      }
      include_raw_key_in_signature = false
      include_cert_in_signature    = false
      algorithm                    = "SHA256withRSA"
    }
  }
  outbound_provision = {
    type = "PingOne"
    target_settings = [
      {
        name  = "PINGONE_ENVIRONMENT"
        value = "example"
      }
    ]
    channels = [
      {
        name        = "Channel1"
        max_threads = 1
        timeout     = 120
        active      = false
        channel_source = {
          base_dn = "dc=example,dc=com"
          data_source = {
            id = "pingdirectory"
          }
          guid_attribute_name = "entry_uuid"
          change_detection_settings = {
            user_object_class         = "inetOrgPerson"
            changed_users_algorithm   = "TIMESTAMP_NO_NEGATION"
            group_object_class        = "groupOfUniqueNames"
            time_stamp_attribute_name = "modifyTimestamp"
          }
          account_management_settings = {
            account_status_algorithm      = "ACCOUNT_STATUS_ALGORITHM_FLAG"
            account_status_attribute_name = "nsaccountlock"
            flag_comparison_value         = "true"
            flag_comparison_status        = true
            default_status                = true
          }
          group_membership_detection = {
            group_member_attribute_name = "uniqueMember"
          }
          guid_binary = false
          user_source_location = {
            filter = "cn=John"

          }
        }
        attribute_mapping = [
          {
            field_name = "username"
            saas_field_info = {
              attribute_names = [
                "uid"
              ]
            }
          },
          {
            field_name = "email"
            saas_field_info = {
              attribute_names = [
                "mail"
              ]
            }
          },
          {
            field_name = "populationID"
            saas_field_info = {
              default_value = "example"
            }
          }
        ]
      }
    ]
  }
}
```

## Example Usage - WS Trust

```terraform
resource "pingfederate_idp_sp_connection" "wsTrustExample" {
  connection_id      = "connection"
  name               = "connection"
  entity_id          = "entity"
  active             = true
  contact_info       = {}
  base_url           = "https://localhost:9031"
  logging_mode       = "STANDARD"
  virtual_entity_ids = []
  credentials = {
    certs = []
    signing_settings = {
      signing_key_pair_ref = {
        id = "signingKey"
      }
      include_raw_key_in_signature = false
      include_cert_in_signature    = false
      algorithm                    = "SHA256withRSA"
    }
  }
  ws_trust = {
    partner_service_ids = [
      "id"
    ]
    oauth_assertion_profiles = true
    default_token_type       = "SAML20"
    generate_key             = false
    encrypt_saml2_assertion  = false
    minutes_before           = 5
    minutes_after            = 30
    attribute_contract = {
      core_attributes = [
        {
          name = "TOKEN_SUBJECT"
        }
      ]
      extended_attributes = []
    }
    token_processor_mappings = [
      {
        attribute_sources = []
        attribute_contract_fulfillment = {
          "TOKEN_SUBJECT" : {
            source = {
              type = "TOKEN"
            }
            value = "username"
          }
        }
        issuance_criteria = {
          conditional_criteria = []
        }
        idp_token_processor_ref = {
          id = "UsernameTokenProcessor"
        }
        restricted_virtual_entity_ids = []
      }
    ]
  }
  connection_target_type = "STANDARD"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_id` (String) The persistent, unique ID for the connection. It can be any combination of `[a-zA-Z0-9._-]`.
- `entity_id` (String) The partner's entity ID (connection ID) or issuer value (for OIDC Connections).
- `name` (String) The connection name.

### Optional

- `active` (Boolean) Specifies whether the connection is active and ready to process incoming requests. The default value is `false`.
- `additional_allowed_entities_configuration` (Attributes) Additional allowed entities or issuers configuration. Currently only used in OIDC IdP (RP) connection. (see [below for nested schema](#nestedatt--additional_allowed_entities_configuration))
- `application_icon_url` (String) The application icon url.
- `application_name` (String) The application name.
- `attribute_query` (Attributes) The attribute query profile supports SPs in requesting user attributes. (see [below for nested schema](#nestedatt--attribute_query))
- `base_url` (String) The fully-qualified hostname and port on which your partner's federation deployment runs.
- `connection_target_type` (String) The connection target type. This field is intended for bulk import/export usage. Changing its value may result in unexpected behavior. The default value is `STANDARD`. Options are `STANDARD`, `SALESFORCE`, `SALESFORCE_CP`, `SALESFORCE_PP`, `PINGONE_SCIM11`.
- `contact_info` (Attributes) Contact information. (see [below for nested schema](#nestedatt--contact_info))
- `credentials` (Attributes) The certificates and settings for encryption, signing, and signature verification. (see [below for nested schema](#nestedatt--credentials))
- `default_virtual_entity_id` (String) The default alternate entity ID that identifies the local server to this partner. It is required when `virtual_entity_ids` is not empty and must be included in that list.
- `extended_properties` (Attributes Map) Extended Properties allows to store additional information for IdP/SP Connections. The names of these extended properties should be defined in /extendedProperties. (see [below for nested schema](#nestedatt--extended_properties))
- `license_connection_group` (String) The license connection group. If your PingFederate license is based on connection groups, each connection must be assigned to a group before it can be used.
- `logging_mode` (String) The level of transaction logging applicable for this connection. Default is `STANDARD`. Options are `NONE`, `STANDARD`, `ENHANCED`, `FULL`. If the `sp_connection_transaction_logging_override` attribute is set to anything other than `DONT_OVERRIDE` in the `server_settings_general` resource, then this attribute must be set to the same value.
- `metadata_reload_settings` (Attributes) Configuration settings to enable automatic reload of partner's metadata. (see [below for nested schema](#nestedatt--metadata_reload_settings))
- `outbound_provision` (Attributes) Outbound Provisioning allows an IdP to create and maintain user accounts at standards-based partner sites using SCIM as well as select-proprietary provisioning partner sites that are protocol-enabled. (see [below for nested schema](#nestedatt--outbound_provision))
- `sp_browser_sso` (Attributes) The SAML settings used to enable secure browser-based SSO to resources at your partner's site. (see [below for nested schema](#nestedatt--sp_browser_sso))
- `virtual_entity_ids` (Set of String) List of alternate entity IDs that identifies the local server to this partner.
- `ws_trust` (Attributes) Ws-Trust STS provides security-token validation and creation to extend SSO access to identity-enabled Web Services (see [below for nested schema](#nestedatt--ws_trust))

### Read-Only

- `creation_date` (String) The time at which the connection was created. This property is read only.
- `id` (String) The ID of this resource.
- `type` (String, Deprecated) The type of this connection.

<a id="nestedatt--additional_allowed_entities_configuration"></a>
### Nested Schema for `additional_allowed_entities_configuration`

Optional:

- `additional_allowed_entities` (Attributes Set) An array of additional allowed entities or issuers to be accepted during entity or issuer validation. (see [below for nested schema](#nestedatt--additional_allowed_entities_configuration--additional_allowed_entities))
- `allow_additional_entities` (Boolean) Set to true to configure additional entities or issuers to be accepted during entity or issuer validation.
- `allow_all_entities` (Boolean) Set to true to accept any entity or issuer during entity or issuer validation. (Not Recommended)

<a id="nestedatt--additional_allowed_entities_configuration--additional_allowed_entities"></a>
### Nested Schema for `additional_allowed_entities_configuration.additional_allowed_entities`

Optional:

- `entity_description` (String) Entity description.
- `entity_id` (String) Unique entity identifier.



<a id="nestedatt--attribute_query"></a>
### Nested Schema for `attribute_query`

Required:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_query--attribute_contract_fulfillment))
- `attributes` (Set of String) The list of attributes that may be returned to the SP in the response to an attribute request.

Optional:

- `attribute_sources` (Attributes Set) A list of configured data stores to look up attributes from. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources))
- `issuance_criteria` (Attributes) The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled. (see [below for nested schema](#nestedatt--attribute_query--issuance_criteria))
- `policy` (Attributes) The attribute query profile's security policy. (see [below for nested schema](#nestedatt--attribute_query--policy))

<a id="nestedatt--attribute_query--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_query.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_query--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_query--attribute_contract_fulfillment--source"></a>
### Nested Schema for `attribute_query.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--attribute_query--attribute_sources"></a>
### Nested Schema for `attribute_query.attribute_sources`

Optional:

- `custom_attribute_source` (Attributes) The configured settings used to look up attributes from a custom data store. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--custom_attribute_source))
- `jdbc_attribute_source` (Attributes) The configured settings used to look up attributes from a JDBC data store. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--jdbc_attribute_source))
- `ldap_attribute_source` (Attributes) The configured settings used to look up attributes from a LDAP data store. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--ldap_attribute_source))

<a id="nestedatt--attribute_query--attribute_sources--custom_attribute_source"></a>
### Nested Schema for `attribute_query.attribute_sources.custom_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--custom_attribute_source--data_store_ref))

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--custom_attribute_source--attribute_contract_fulfillment))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `filter_fields` (Attributes Set) The list of fields that can be used to filter a request to the custom data store. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--custom_attribute_source--filter_fields))
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--attribute_query--attribute_sources--custom_attribute_source--data_store_ref"></a>
### Nested Schema for `attribute_query.attribute_sources.custom_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--attribute_query--attribute_sources--custom_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_query.attribute_sources.custom_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--custom_attribute_source--type--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_query--attribute_sources--custom_attribute_source--type--source"></a>
### Nested Schema for `attribute_query.attribute_sources.custom_attribute_source.type.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--attribute_query--attribute_sources--custom_attribute_source--filter_fields"></a>
### Nested Schema for `attribute_query.attribute_sources.custom_attribute_source.filter_fields`

Required:

- `name` (String) The name of this field.

Optional:

- `value` (String) The value of this field. Whether or not the value is required will be determined by plugin validation checks.



<a id="nestedatt--attribute_query--attribute_sources--jdbc_attribute_source"></a>
### Nested Schema for `attribute_query.attribute_sources.jdbc_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--jdbc_attribute_source--data_store_ref))
- `filter` (String) The JDBC WHERE clause used to query your data store to locate a user record.
- `table` (String) The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--jdbc_attribute_source--attribute_contract_fulfillment))
- `column_names` (List of String) A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `schema` (String) Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--attribute_query--attribute_sources--jdbc_attribute_source--data_store_ref"></a>
### Nested Schema for `attribute_query.attribute_sources.jdbc_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--attribute_query--attribute_sources--jdbc_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_query.attribute_sources.jdbc_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--jdbc_attribute_source--type--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_query--attribute_sources--jdbc_attribute_source--type--source"></a>
### Nested Schema for `attribute_query.attribute_sources.jdbc_attribute_source.type.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.




<a id="nestedatt--attribute_query--attribute_sources--ldap_attribute_source"></a>
### Nested Schema for `attribute_query.attribute_sources.ldap_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--ldap_attribute_source--data_store_ref))
- `search_filter` (String) The LDAP filter that will be used to lookup the objects from the directory.
- `search_scope` (String) Determines the node depth of the query.
- `type` (String) The data store type of this attribute source.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `binary_attribute_settings` (Attributes Map) The advanced settings for binary LDAP attributes. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--ldap_attribute_source--binary_attribute_settings))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `member_of_nested_group` (Boolean) Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.
- `search_attributes` (Set of String) A list of LDAP attributes returned from search and available for mapping.

<a id="nestedatt--attribute_query--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `attribute_query.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--attribute_query--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `attribute_query.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_query--attribute_sources--ldap_attribute_source--search_attributes--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--attribute_query--attribute_sources--ldap_attribute_source--search_attributes--source"></a>
### Nested Schema for `attribute_query.attribute_sources.ldap_attribute_source.search_attributes.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--attribute_query--attribute_sources--ldap_attribute_source--binary_attribute_settings"></a>
### Nested Schema for `attribute_query.attribute_sources.ldap_attribute_source.binary_attribute_settings`

Optional:

- `binary_encoding` (String) Get the encoding type for this attribute. If not specified, the default is BASE64.




<a id="nestedatt--attribute_query--issuance_criteria"></a>
### Nested Schema for `attribute_query.issuance_criteria`

Optional:

- `conditional_criteria` (Attributes Set) A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue. (see [below for nested schema](#nestedatt--attribute_query--issuance_criteria--conditional_criteria))
- `expression_criteria` (Attributes Set) A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria. (see [below for nested schema](#nestedatt--attribute_query--issuance_criteria--expression_criteria))

<a id="nestedatt--attribute_query--issuance_criteria--conditional_criteria"></a>
### Nested Schema for `attribute_query.issuance_criteria.conditional_criteria`

Required:

- `attribute_name` (String) The name of the attribute to use in this issuance criterion.
- `condition` (String) The condition that will be applied to the source attribute's value and the expected value. Options are `EQUALS`, `EQUALS_CASE_INSENSITIVE`, `EQUALS_DN`, `NOT_EQUAL`, `NOT_EQUAL_CASE_INSENSITIVE`, `NOT_EQUAL_DN`, `MULTIVALUE_CONTAINS`, `MULTIVALUE_CONTAINS_CASE_INSENSITIVE`, `MULTIVALUE_CONTAINS_DN`, `MULTIVALUE_DOES_NOT_CONTAIN`, `MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE`, `MULTIVALUE_DOES_NOT_CONTAIN_DN`.
- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--attribute_query--issuance_criteria--conditional_criteria--source))
- `value` (String) The expected value of this issuance criterion.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.

<a id="nestedatt--attribute_query--issuance_criteria--conditional_criteria--source"></a>
### Nested Schema for `attribute_query.issuance_criteria.conditional_criteria.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--attribute_query--issuance_criteria--expression_criteria"></a>
### Nested Schema for `attribute_query.issuance_criteria.expression_criteria`

Required:

- `expression` (String) The OGNL expression to evaluate.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.



<a id="nestedatt--attribute_query--policy"></a>
### Nested Schema for `attribute_query.policy`

Optional:

- `encrypt_assertion` (Boolean) Encrypt the assertion.
- `require_encrypted_name_id` (Boolean) Require an encrypted name identifier.
- `require_signed_attribute_query` (Boolean) Require signed attribute query.
- `sign_assertion` (Boolean) Sign the assertion.
- `sign_response` (Boolean) Sign the response.



<a id="nestedatt--contact_info"></a>
### Nested Schema for `contact_info`

Optional:

- `company` (String) Company name.
- `email` (String) Contact email address.
- `first_name` (String) Contact first name.
- `last_name` (String) Contact last name.
- `phone` (String) Contact phone number.


<a id="nestedatt--credentials"></a>
### Nested Schema for `credentials`

Optional:

- `block_encryption_algorithm` (String) The algorithm used to encrypt assertions sent to this partner. `AES_128`, `AES_256`, `AES_128_GCM`, `AES_192_GCM`, `AES_256_GCM` and `Triple_DES` are supported.
- `certs` (Attributes Set) The certificates used for signature verification and XML encryption. (see [below for nested schema](#nestedatt--credentials--certs))
- `decryption_key_pair_ref` (Attributes) A reference to a resource. (see [below for nested schema](#nestedatt--credentials--decryption_key_pair_ref))
- `inbound_back_channel_auth` (Attributes) The SOAP authentication methods when sending or receiving a message using SOAP back channel. (see [below for nested schema](#nestedatt--credentials--inbound_back_channel_auth))
- `key_transport_algorithm` (String) The algorithm used to transport keys to this partner. `RSA_OAEP`, `RSA_OAEP_256` and `RSA_v15` are supported.
- `outbound_back_channel_auth` (Attributes) The SOAP authentication methods when sending or receiving a message using SOAP back channel. (see [below for nested schema](#nestedatt--credentials--outbound_back_channel_auth))
- `secondary_decryption_key_pair_ref` (Attributes) A reference to a resource. (see [below for nested schema](#nestedatt--credentials--secondary_decryption_key_pair_ref))
- `signing_settings` (Attributes) Settings related to signing messages sent to this partner. (see [below for nested schema](#nestedatt--credentials--signing_settings))
- `verification_issuer_dn` (String) If a verification Subject DN is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.
- `verification_subject_dn` (String) If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.

<a id="nestedatt--credentials--certs"></a>
### Nested Schema for `credentials.certs`

Required:

- `x509file` (Attributes) Encoded certificate data. (see [below for nested schema](#nestedatt--credentials--certs--x509file))

Optional:

- `active_verification_cert` (Boolean) Indicates whether this is an active signature verification certificate.
- `encryption_cert` (Boolean) Indicates whether to use this cert to encrypt outgoing assertions. Only one certificate in the collection can have this flag set.
- `primary_verification_cert` (Boolean) Indicates whether this is the primary signature verification certificate. Only one certificate in the collection can have this flag set.
- `secondary_verification_cert` (Boolean) Indicates whether this is the secondary signature verification certificate. Only one certificate in the collection can have this flag set.

Read-Only:

- `cert_view` (Attributes) Certificate details. (see [below for nested schema](#nestedatt--credentials--certs--cert_view))

<a id="nestedatt--credentials--certs--x509file"></a>
### Nested Schema for `credentials.certs.x509file`

Required:

- `file_data` (String) The certificate data in PEM format. New line characters should be omitted or encoded in this value.

Optional:

- `crypto_provider` (String) Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Options are `LOCAL`, `HSM`.
- `id` (String) The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified.


<a id="nestedatt--credentials--certs--cert_view"></a>
### Nested Schema for `credentials.certs.cert_view`

Read-Only:

- `crypto_provider` (String) Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Options are `LOCAL`, `HSM`.
- `expires` (String) The end date up until which the item is valid, in ISO 8601 format (UTC).
- `id` (String) The persistent, unique ID for the certificate.
- `issuer_dn` (String) The issuer's distinguished name.
- `key_algorithm` (String) The public key algorithm.
- `key_size` (Number) The public key size.
- `serial_number` (String) The serial number assigned by the CA.
- `sha1fingerprint` (String) SHA-1 fingerprint in Hex encoding.
- `sha256fingerprint` (String) SHA-256 fingerprint in Hex encoding.
- `signature_algorithm` (String) The signature algorithm.
- `status` (String) Status of the item. Options are `VALID`, `EXPIRED`, `NOT_YET_VALID`, `REVOKED`.
- `subject_alternative_names` (Set of String) The subject alternative names (SAN).
- `subject_dn` (String) The subject's distinguished name.
- `valid_from` (String) The start date from which the item is valid, in ISO 8601 format (UTC).
- `version` (Number) The X.509 version to which the item conforms.



<a id="nestedatt--credentials--decryption_key_pair_ref"></a>
### Nested Schema for `credentials.decryption_key_pair_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--credentials--inbound_back_channel_auth"></a>
### Nested Schema for `credentials.inbound_back_channel_auth`

Optional:

- `certs` (Attributes Set) The certificates used for signature verification and XML encryption. (see [below for nested schema](#nestedatt--credentials--inbound_back_channel_auth--certs))
- `digital_signature` (Boolean) If incoming or outgoing messages must be signed.
- `http_basic_credentials` (Attributes) Username and password credentials. (see [below for nested schema](#nestedatt--credentials--inbound_back_channel_auth--http_basic_credentials))
- `require_ssl` (Boolean) Incoming HTTP transmissions must use a secure channel.
- `verification_issuer_dn` (String) If `verification_subject_dn` is provided, you can optionally restrict the issuer to a specific trusted CA by specifying its DN in this field.
- `verification_subject_dn` (String) If this property is set, the verification trust model is Anchored. The verification certificate must be signed by a trusted CA and included in the incoming message, and the subject DN of the expected certificate is specified in this property. If this property is not set, then a primary verification certificate must be specified in the certs array.

Read-Only:

- `type` (String, Deprecated) The back channel authentication type.

<a id="nestedatt--credentials--inbound_back_channel_auth--certs"></a>
### Nested Schema for `credentials.inbound_back_channel_auth.certs`

Required:

- `x509file` (Attributes) Encoded certificate data. (see [below for nested schema](#nestedatt--credentials--inbound_back_channel_auth--certs--x509file))

Optional:

- `active_verification_cert` (Boolean) Indicates whether this is an active signature verification certificate.
- `encryption_cert` (Boolean) Indicates whether to use this cert to encrypt outgoing assertions. Only one certificate in the collection can have this flag set.
- `primary_verification_cert` (Boolean) Indicates whether this is the primary signature verification certificate. Only one certificate in the collection can have this flag set.
- `secondary_verification_cert` (Boolean) Indicates whether this is the secondary signature verification certificate. Only one certificate in the collection can have this flag set.

Read-Only:

- `cert_view` (Attributes) Certificate details. (see [below for nested schema](#nestedatt--credentials--inbound_back_channel_auth--certs--cert_view))

<a id="nestedatt--credentials--inbound_back_channel_auth--certs--x509file"></a>
### Nested Schema for `credentials.inbound_back_channel_auth.certs.x509file`

Required:

- `file_data` (String) The certificate data in PEM format. New line characters should be omitted or encoded in this value.

Optional:

- `crypto_provider` (String) Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Options are `LOCAL`, `HSM`.
- `id` (String) The persistent, unique ID for the certificate. It can be any combination of `[a-z0-9._-]`. This property is system-assigned if not specified.


<a id="nestedatt--credentials--inbound_back_channel_auth--certs--cert_view"></a>
### Nested Schema for `credentials.inbound_back_channel_auth.certs.cert_view`

Read-Only:

- `crypto_provider` (String) Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Options are `LOCAL`, `HSM`.
- `expires` (String) The end date up until which the item is valid, in ISO 8601 format (UTC).
- `id` (String) The persistent, unique ID for the certificate.
- `issuer_dn` (String) The issuer's distinguished name.
- `key_algorithm` (String) The public key algorithm.
- `key_size` (Number) The public key size.
- `serial_number` (String) The serial number assigned by the CA.
- `sha1fingerprint` (String) SHA-1 fingerprint in Hex encoding.
- `sha256fingerprint` (String) SHA-256 fingerprint in Hex encoding.
- `signature_algorithm` (String) The signature algorithm.
- `status` (String) Status of the item. Options are `VALID`, `EXPIRED`, `NOT_YET_VALID`, `REVOKED`.
- `subject_alternative_names` (Set of String) The subject alternative names (SAN).
- `subject_dn` (String) The subject's distinguished name.
- `valid_from` (String) The start date from which the item is valid, in ISO 8601 format (UTC).
- `version` (Number) The X.509 version to which the item conforms.



<a id="nestedatt--credentials--inbound_back_channel_auth--http_basic_credentials"></a>
### Nested Schema for `credentials.inbound_back_channel_auth.http_basic_credentials`

Optional:

- `encrypted_password` (String, Deprecated) For GET requests, this field contains the encrypted password, if one exists.
- `password` (String, Sensitive) User password.
- `username` (String) The username.



<a id="nestedatt--credentials--outbound_back_channel_auth"></a>
### Nested Schema for `credentials.outbound_back_channel_auth`

Optional:

- `digital_signature` (Boolean) If incoming or outgoing messages must be signed.
- `http_basic_credentials` (Attributes) Username and password credentials. (see [below for nested schema](#nestedatt--credentials--outbound_back_channel_auth--http_basic_credentials))
- `ssl_auth_key_pair_ref` (Attributes) A reference to a resource. (see [below for nested schema](#nestedatt--credentials--outbound_back_channel_auth--ssl_auth_key_pair_ref))
- `validate_partner_cert` (Boolean) Validate the partner server certificate. Default is `true`.

Read-Only:

- `type` (String, Deprecated) The back channel authentication type.

<a id="nestedatt--credentials--outbound_back_channel_auth--http_basic_credentials"></a>
### Nested Schema for `credentials.outbound_back_channel_auth.http_basic_credentials`

Optional:

- `encrypted_password` (String, Deprecated) For GET requests, this field contains the encrypted password, if one exists.
- `password` (String, Sensitive) User password.
- `username` (String) The username.


<a id="nestedatt--credentials--outbound_back_channel_auth--ssl_auth_key_pair_ref"></a>
### Nested Schema for `credentials.outbound_back_channel_auth.ssl_auth_key_pair_ref`

Required:

- `id` (String) The ID of the resource.



<a id="nestedatt--credentials--secondary_decryption_key_pair_ref"></a>
### Nested Schema for `credentials.secondary_decryption_key_pair_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--credentials--signing_settings"></a>
### Nested Schema for `credentials.signing_settings`

Optional:

- `algorithm` (String) The algorithm used to sign messages sent to this partner. The default is `SHA1withDSA` for DSA certs, `SHA256withRSA` for RSA certs, and `SHA256withECDSA` for EC certs. For RSA certs, `SHA1withRSA`, `SHA384withRSA`, `SHA512withRSA`, `SHA256withRSAandMGF1`, `SHA384withRSAandMGF1` and `SHA512withRSAandMGF1` are also supported. For EC certs, `SHA384withECDSA` and `SHA512withECDSA` are also supported. If the connection is WS-Federation with JWT token type, then the possible values are RSA SHA256, RSA SHA384, RSA SHA512, RSASSA-PSS SHA256, RSASSA-PSS SHA384, RSASSA-PSS SHA512, ECDSA SHA256, ECDSA SHA384, ECDSA SHA512
- `alternative_signing_key_pair_refs` (Attributes Set) The list of IDs of alternative key pairs used to sign messages sent to this partner. The ID of the key pair is also known as the alias and can be found by viewing the corresponding certificate under 'Signing & Decryption Keys & Certificates' in the PingFederate admin console. (see [below for nested schema](#nestedatt--credentials--signing_settings--alternative_signing_key_pair_refs))
- `include_cert_in_signature` (Boolean) Determines whether the signing certificate is included in the signature <KeyInfo> element. Default is `false`.
- `include_raw_key_in_signature` (Boolean) Determines whether the <KeyValue> element with the raw public key is included in the signature <KeyInfo> element.
- `signing_key_pair_ref` (Attributes) A reference to a resource. (see [below for nested schema](#nestedatt--credentials--signing_settings--signing_key_pair_ref))

<a id="nestedatt--credentials--signing_settings--alternative_signing_key_pair_refs"></a>
### Nested Schema for `credentials.signing_settings.alternative_signing_key_pair_refs`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--credentials--signing_settings--signing_key_pair_ref"></a>
### Nested Schema for `credentials.signing_settings.signing_key_pair_ref`

Required:

- `id` (String) The ID of the resource.




<a id="nestedatt--extended_properties"></a>
### Nested Schema for `extended_properties`

Optional:

- `values` (Set of String) A List of values


<a id="nestedatt--metadata_reload_settings"></a>
### Nested Schema for `metadata_reload_settings`

Optional:

- `enable_auto_metadata_update` (Boolean) Specifies whether the metadata of the connection will be automatically reloaded. The default value is `true`.
- `metadata_url_ref` (Attributes) A reference to a resource. (see [below for nested schema](#nestedatt--metadata_reload_settings--metadata_url_ref))

<a id="nestedatt--metadata_reload_settings--metadata_url_ref"></a>
### Nested Schema for `metadata_reload_settings.metadata_url_ref`

Required:

- `id` (String) The ID of the resource.



<a id="nestedatt--outbound_provision"></a>
### Nested Schema for `outbound_provision`

Required:

- `channels` (Attributes List) Includes settings of a source data store, managing provisioning threads and mapping of attributes. (see [below for nested schema](#nestedatt--outbound_provision--channels))
- `target_settings` (Attributes Set) Configuration fields that includes credentials to target SaaS application. (see [below for nested schema](#nestedatt--outbound_provision--target_settings))
- `type` (String) The SaaS plugin type.

Optional:

- `custom_schema` (Attributes) Custom SCIM Attributes configuration. (see [below for nested schema](#nestedatt--outbound_provision--custom_schema))

Read-Only:

- `target_settings_all` (Attributes Set) Configuration fields that includes credentials to target SaaS application. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--outbound_provision--target_settings_all))

<a id="nestedatt--outbound_provision--channels"></a>
### Nested Schema for `outbound_provision.channels`

Required:

- `active` (Boolean) Indicates whether the channel is the active channel for this connection.
- `attribute_mapping` (Attributes Set) The mapping of attributes from the local data store into Fields specified by the service provider. (see [below for nested schema](#nestedatt--outbound_provision--channels--attribute_mapping))
- `channel_source` (Attributes) The source data source and LDAP settings. (see [below for nested schema](#nestedatt--outbound_provision--channels--channel_source))
- `name` (String) The name of the channel.

Optional:

- `max_threads` (Number) The number of processing threads. The default value is `1`.
- `timeout` (Number) Timeout, in seconds, for individual user and group provisioning operations on the target service provider. The default value is `60`.

Read-Only:

- `attribute_mapping_all` (Attributes Set) The mapping of attributes from the local data store into Fields specified by the service provider. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--outbound_provision--channels--attribute_mapping_all))

<a id="nestedatt--outbound_provision--channels--attribute_mapping"></a>
### Nested Schema for `outbound_provision.channels.attribute_mapping`

Required:

- `field_name` (String) The name of target field.
- `saas_field_info` (Attributes) The settings that represent how attribute values from source data store will be mapped into Fields specified by the service provider. (see [below for nested schema](#nestedatt--outbound_provision--channels--attribute_mapping--saas_field_info))

<a id="nestedatt--outbound_provision--channels--attribute_mapping--saas_field_info"></a>
### Nested Schema for `outbound_provision.channels.attribute_mapping.saas_field_info`

Optional:

- `attribute_names` (List of String) The list of source attribute names used to generate or map to a target field
- `character_case` (String) The character case of the field value.
- `create_only` (Boolean) Indicates whether this field is a create only field and cannot be updated.
- `default_value` (String) The default value for the target field
- `expression` (String) An OGNL expression to obtain a value.
- `masked` (Boolean) Indicates whether the attribute should be masked in server logs.
- `parser` (String) Indicates how the field shall be parsed. Options are `NONE`, `EXTRACT_CN_FROM_DN`, `EXTRACT_USERNAME_FROM_EMAIL`.
- `trim` (Boolean) Indicates whether field should be trimmed before provisioning.



<a id="nestedatt--outbound_provision--channels--channel_source"></a>
### Nested Schema for `outbound_provision.channels.channel_source`

Required:

- `account_management_settings` (Attributes) Account management settings. (see [below for nested schema](#nestedatt--outbound_provision--channels--channel_source--account_management_settings))
- `base_dn` (String) The base DN where the user records are located.
- `change_detection_settings` (Attributes) Setting to detect changes to a user or a group. (see [below for nested schema](#nestedatt--outbound_provision--channels--channel_source--change_detection_settings))
- `data_source` (Attributes) Reference to an LDAP datastore. (see [below for nested schema](#nestedatt--outbound_provision--channels--channel_source--data_source))
- `group_membership_detection` (Attributes) Settings to detect group memberships. (see [below for nested schema](#nestedatt--outbound_provision--channels--channel_source--group_membership_detection))
- `guid_attribute_name` (String) the GUID attribute name.
- `guid_binary` (Boolean) Indicates whether the GUID is stored in binary format.
- `user_source_location` (Attributes) The location settings that includes a DN and a LDAP filter. (see [below for nested schema](#nestedatt--outbound_provision--channels--channel_source--user_source_location))

Optional:

- `group_source_location` (Attributes) The location settings that includes a DN and a LDAP filter. (see [below for nested schema](#nestedatt--outbound_provision--channels--channel_source--group_source_location))

<a id="nestedatt--outbound_provision--channels--channel_source--account_management_settings"></a>
### Nested Schema for `outbound_provision.channels.channel_source.account_management_settings`

Required:

- `account_status_algorithm` (String) The account status algorithm name. Options are `ACCOUNT_STATUS_ALGORITHM_AD`, `ACCOUNT_STATUS_ALGORITHM_FLAG`. `ACCOUNT_STATUS_ALGORITHM_AD` -  Algorithm name for Active Directory, which uses a bitmap for each user entry. `ACCOUNT_STATUS_ALGORITHM_FLAG` - Algorithm name for Oracle Directory Server and other LDAP directories that use a separate attribute to store the user's status. When this option is selected, the Flag Comparison Value and Flag Comparison Status fields should be used.
- `account_status_attribute_name` (String) The account status attribute name.

Optional:

- `default_status` (Boolean) The default status of the account.
- `flag_comparison_status` (Boolean) The flag that represents comparison status.
- `flag_comparison_value` (String) The flag that represents comparison value.


<a id="nestedatt--outbound_provision--channels--channel_source--change_detection_settings"></a>
### Nested Schema for `outbound_provision.channels.channel_source.change_detection_settings`

Required:

- `changed_users_algorithm` (String) The changed user algorithm. Options are `ACTIVE_DIRECTORY_USN`, `TIMESTAMP`, `TIMESTAMP_NO_NEGATION`. `ACTIVE_DIRECTORY_USN` - For Active Directory only, this algorithm queries for update sequence numbers on user records that are larger than the last time records were checked. `TIMESTAMP` - Queries for timestamps on user records that are not older than the last time records were checked. This check is more efficient from the point of view of the PingFederate provisioner but can be more time consuming on the LDAP side, particularly with the Oracle Directory Server. `TIMESTAMP_NO_NEGATION` - Queries for timestamps on user records that are newer than the last time records were checked. This algorithm is recommended for the Oracle Directory Server.
- `group_object_class` (String) The group object class.
- `time_stamp_attribute_name` (String) The timestamp attribute name.
- `user_object_class` (String) The user object class.

Optional:

- `usn_attribute_name` (String) The USN attribute name.


<a id="nestedatt--outbound_provision--channels--channel_source--data_source"></a>
### Nested Schema for `outbound_provision.channels.channel_source.data_source`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--outbound_provision--channels--channel_source--group_membership_detection"></a>
### Nested Schema for `outbound_provision.channels.channel_source.group_membership_detection`

Optional:

- `group_member_attribute_name` (String) The name of the attribute that represents group members in a group, also known as group member attribute.
- `member_of_group_attribute_name` (String) The name of the attribute that indicates the entity is a member of a group, also known as member of attribute.


<a id="nestedatt--outbound_provision--channels--channel_source--user_source_location"></a>
### Nested Schema for `outbound_provision.channels.channel_source.user_source_location`

Optional:

- `filter` (String) An LDAP filter.
- `group_dn` (String) The group DN for users or groups.
- `nested_search` (Boolean) Indicates whether the search is nested. Default is `false`.


<a id="nestedatt--outbound_provision--channels--channel_source--group_source_location"></a>
### Nested Schema for `outbound_provision.channels.channel_source.group_source_location`

Optional:

- `filter` (String) An LDAP filter.
- `group_dn` (String) The group DN for users or groups.
- `nested_search` (Boolean) Indicates whether the search is nested. The default value is `false`.



<a id="nestedatt--outbound_provision--channels--attribute_mapping_all"></a>
### Nested Schema for `outbound_provision.channels.attribute_mapping_all`

Required:

- `field_name` (String) The name of target field.
- `saas_field_info` (Attributes) The settings that represent how attribute values from source data store will be mapped into Fields specified by the service provider. (see [below for nested schema](#nestedatt--outbound_provision--channels--attribute_mapping_all--saas_field_info))

<a id="nestedatt--outbound_provision--channels--attribute_mapping_all--saas_field_info"></a>
### Nested Schema for `outbound_provision.channels.attribute_mapping_all.saas_field_info`

Optional:

- `attribute_names` (List of String) The list of source attribute names used to generate or map to a target field
- `character_case` (String) The character case of the field value.
- `create_only` (Boolean) Indicates whether this field is a create only field and cannot be updated.
- `default_value` (String) The default value for the target field
- `expression` (String) An OGNL expression to obtain a value.
- `masked` (Boolean) Indicates whether the attribute should be masked in server logs.
- `parser` (String) Indicates how the field shall be parsed. Options are `NONE`, `EXTRACT_CN_FROM_DN`, `EXTRACT_USERNAME_FROM_EMAIL`.
- `trim` (Boolean) Indicates whether field should be trimmed before provisioning.




<a id="nestedatt--outbound_provision--target_settings"></a>
### Nested Schema for `outbound_provision.target_settings`

Required:

- `name` (String) The name of the configuration field.

Optional:

- `value` (String, Sensitive) The value for the configuration field.


<a id="nestedatt--outbound_provision--custom_schema"></a>
### Nested Schema for `outbound_provision.custom_schema`

Optional:

- `attributes` (Attributes Set) (see [below for nested schema](#nestedatt--outbound_provision--custom_schema--attributes))
- `namespace` (String)

<a id="nestedatt--outbound_provision--custom_schema--attributes"></a>
### Nested Schema for `outbound_provision.custom_schema.attributes`

Optional:

- `multi_valued` (Boolean) Indicates whether the attribute is multi-valued.
- `name` (String) Name of the attribute.
- `sub_attributes` (Set of String) List of sub-attributes for an attribute.
- `types` (Set of String) Represents the name of each attribute type in case of multi-valued attribute.



<a id="nestedatt--outbound_provision--target_settings_all"></a>
### Nested Schema for `outbound_provision.target_settings_all`

Required:

- `name` (String) The name of the configuration field.

Optional:

- `value` (String, Sensitive) The value for the configuration field.



<a id="nestedatt--sp_browser_sso"></a>
### Nested Schema for `sp_browser_sso`

Required:

- `adapter_mappings` (Attributes Set) A list of adapters that map to outgoing assertions. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings))
- `assertion_lifetime` (Attributes) The timeframe of validity before and after the issuance of the assertion. (see [below for nested schema](#nestedatt--sp_browser_sso--assertion_lifetime))
- `attribute_contract` (Attributes) A set of user attributes that the IdP sends in the SAML assertion. (see [below for nested schema](#nestedatt--sp_browser_sso--attribute_contract))
- `protocol` (String) The browser-based SSO protocol to use. Options are `SAML20`, `WSFED`, `SAML11`, `SAML10`, `OIDC`.
- `sso_service_endpoints` (Attributes Set) A list of possible endpoints to send assertions to. (see [below for nested schema](#nestedatt--sp_browser_sso--sso_service_endpoints))

Optional:

- `always_sign_artifact_response` (Boolean) Specify to always sign the SAML ArtifactResponse.
- `artifact` (Attributes) The settings for an Artifact binding. (see [below for nested schema](#nestedatt--sp_browser_sso--artifact))
- `authentication_policy_contract_assertion_mappings` (Attributes Set) A list of authentication policy contracts that map to outgoing assertions. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings))
- `default_target_url` (String) Default Target URL for SAML1.x connections. This default URL represents the destination on the SP where the user will be directed.
- `enabled_profiles` (Set of String) The profiles that are enabled for browser-based SSO. SAML 2.0 supports all profiles whereas SAML 1.x IdP connections support both IdP and SP (non-standard) initiated SSO. This is required for SAMLx.x Connections.
- `encryption_policy` (Attributes) Defines what to encrypt in the browser-based SSO profile. (see [below for nested schema](#nestedatt--sp_browser_sso--encryption_policy))
- `incoming_bindings` (Set of String) The SAML bindings that are enabled for browser-based SSO. This is required for SAML 2.0 connections when the enabled profiles contain the SP-initiated SSO profile or either SLO profile. For SAML 1.x based connections, it is not used for SP Connections.
- `message_customizations` (Attributes Set) The message customizations for browser-based SSO. Depending on server settings, connection type, and protocol this may or may not be supported. (see [below for nested schema](#nestedatt--sp_browser_sso--message_customizations))
- `require_signed_authn_requests` (Boolean) Require AuthN requests to be signed when received via the POST or Redirect bindings.
- `sign_assertions` (Boolean) Always sign the SAML Assertion.
- `sign_response_as_required` (Boolean) Sign SAML Response as required by the associated binding and encryption policy. Applicable to SAML2.0 only and is defaulted to `true`. It can be set to `false` only on SAML2.0 connections when `sign_assertions` is set to `true`.
- `slo_service_endpoints` (Attributes Set) A list of possible endpoints to send SLO requests and responses. (see [below for nested schema](#nestedatt--sp_browser_sso--slo_service_endpoints))
- `sp_saml_identity_mapping` (String) Process in which users authenticated by the IdP are associated with user accounts local to the SP. Options are `PSEUDONYM`, `STANDARD`, `TRANSIENT`.
- `sp_ws_fed_identity_mapping` (String) Process in which users authenticated by the IdP are associated with user accounts local to the SP for WS-Federation connection types. Options are `EMAIL_ADDRESS`, `USER_PRINCIPLE_NAME`, `COMMON_NAME`.
- `url_whitelist_entries` (Attributes Set) For WS-Federation connections, a whitelist of additional allowed domains and paths used to validate wreply for SLO, if enabled. (see [below for nested schema](#nestedatt--sp_browser_sso--url_whitelist_entries))
- `ws_fed_token_type` (String) The WS-Federation Token Type to use. Options are `SAML11`, `SAML20`, `JWT`.
- `ws_trust_version` (String) The WS-Trust version for a WS-Federation connection. The default version is `WSTRUST12`. Options are `WSTRUST12`, `WSTRUST13`.

Read-Only:

- `sso_application_endpoint` (String) Application endpoint that can be used to invoke single sign-on (SSO) for the connection. This is a read-only parameter. Supported in PF version `11.3` or later.

<a id="nestedatt--sp_browser_sso--adapter_mappings"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings`

Required:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_contract_fulfillment))

Optional:

- `abort_sso_transaction_as_fail_safe` (Boolean) If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used. By default, this value is `false`.
- `adapter_override_settings` (Attributes) (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings))
- `attribute_sources` (Attributes Set) A list of configured data stores to look up attributes from. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources))
- `idp_adapter_ref` (Attributes) Reference to the associated IdP adapter. Note: This is ignored if adapter overrides for this mapping exists. In this case, the override's parent adapter reference is used. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--idp_adapter_ref))
- `issuance_criteria` (Attributes) The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria))
- `restrict_virtual_entity_ids` (Boolean) Restricts this mapping to specific virtual entity IDs.
- `restricted_virtual_entity_ids` (Set of String) The list of virtual server IDs that this mapping is restricted to.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings`

Required:

- `configuration` (Attributes) Plugin instance configuration. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--configuration))
- `id` (String) The ID of the plugin instance. The ID cannot be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.
- `name` (String) The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.
- `plugin_descriptor_ref` (Attributes) Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--plugin_descriptor_ref))

Optional:

- `attribute_contract` (Attributes) A set of attributes exposed by an IdP adapter. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--attribute_contract))
- `attribute_mapping` (Attributes) An IdP Adapter Contract Mapping. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--attribute_mapping))
- `authn_ctx_class_ref` (String) The fixed value that indicates how the user was authenticated.
- `parent_ref` (Attributes) The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides) (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--configuration"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.configuration`

Optional:

- `fields` (Attributes Set) List of configuration fields. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--fields))
- `sensitive_fields` (Attributes Set) List of sensitive configuration fields. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--sensitive_fields))
- `tables` (Attributes List) List of configuration tables. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables))

Read-Only:

- `fields_all` (Attributes Set) List of configuration fields. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--fields_all))
- `tables_all` (Attributes List) List of configuration tables. This attribute will include any values set by default by PingFederate. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables_all))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--fields"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--sensitive_fields"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The sensitive value for the configuration field.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.tables`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables--rows))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables--rows"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.tables.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables--rows--fields))
- `sensitive_fields` (Attributes Set) The sensitive configuration fields in the row. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables--rows--sensitive_fields))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables--rows--fields"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.tables.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables--rows--sensitive_fields"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.tables.rows.sensitive_fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String, Sensitive) The sensitive value for the configuration field.




<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--fields_all"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.fields_all`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables_all"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.tables_all`

Required:

- `name` (String) The name of the table.

Optional:

- `rows` (Attributes List) List of table rows. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables_all--rows))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables_all--rows"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.tables_all.rows`

Optional:

- `default_row` (Boolean) Whether this row is the default.
- `fields` (Attributes Set) The configuration fields in the row. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables_all--rows--fields))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--tables_all--rows--fields"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.tables_all.rows.fields`

Required:

- `name` (String) The name of the configuration field.
- `value` (String) The value for the configuration field.





<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--plugin_descriptor_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.plugin_descriptor_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--attribute_contract"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.attribute_contract`

Required:

- `core_attributes` (Attributes Set) A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--core_attributes))

Optional:

- `extended_attributes` (Attributes Set) A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--extended_attributes))
- `mask_ognl_values` (Boolean) Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to `false`.
- `unique_user_key_attribute` (String) The attribute to use for uniquely identify a user's authentication sessions.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--core_attributes"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.core_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `masked` (Boolean) Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.
- `pseudonym` (Boolean) Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--extended_attributes"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.extended_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `masked` (Boolean) Specifies whether this attribute is masked in PingFederate logs. Defaults to `false`.
- `pseudonym` (Boolean) Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to `false`.



<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--attribute_mapping"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.attribute_mapping`

Required:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_contract_fulfillment))

Optional:

- `attribute_sources` (Attributes Set) A list of configured data stores to look up attributes from. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources))
- `issuance_criteria` (Attributes) The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources`

Optional:

- `custom_attribute_source` (Attributes) The configured settings used to look up attributes from a custom data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--custom_attribute_source))
- `jdbc_attribute_source` (Attributes) The configured settings used to look up attributes from a JDBC data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--jdbc_attribute_source))
- `ldap_attribute_source` (Attributes) The configured settings used to look up attributes from a LDAP data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--custom_attribute_source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.custom_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--data_store_ref))

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `filter_fields` (Attributes Set) The list of fields that can be used to filter a request to the custom data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--filter_fields))
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--type--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--type--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.type.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--filter_fields"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.filter_fields`

Required:

- `name` (String) The name of this field.

Optional:

- `value` (String) The value of this field. Whether or not the value is required will be determined by plugin validation checks.



<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--jdbc_attribute_source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.jdbc_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--data_store_ref))
- `filter` (String) The JDBC WHERE clause used to query your data store to locate a user record.
- `table` (String) The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `column_names` (List of String) A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `schema` (String) Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--type--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--type--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.type.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.




<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--data_store_ref))
- `search_filter` (String) The LDAP filter that will be used to lookup the objects from the directory.
- `search_scope` (String) Determines the node depth of the query.
- `type` (String) The data store type of this attribute source.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `binary_attribute_settings` (Attributes Map) The advanced settings for binary LDAP attributes. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--binary_attribute_settings))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `member_of_nested_group` (Boolean) Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.
- `search_attributes` (Set of String) A list of LDAP attributes returned from search and available for mapping.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--search_attributes--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--search_attributes--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.search_attributes.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--attribute_sources--ldap_attribute_source--binary_attribute_settings"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.attribute_sources.ldap_attribute_source.binary_attribute_settings`

Optional:

- `binary_encoding` (String) Get the encoding type for this attribute. If not specified, the default is BASE64.




<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.issuance_criteria`

Optional:

- `conditional_criteria` (Attributes Set) A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria--conditional_criteria))
- `expression_criteria` (Attributes Set) A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria--expression_criteria))

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria--conditional_criteria"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.issuance_criteria.conditional_criteria`

Required:

- `attribute_name` (String) The name of the attribute to use in this issuance criterion.
- `condition` (String) The condition that will be applied to the source attribute's value and the expected value. Options are `EQUALS`, `EQUALS_CASE_INSENSITIVE`, `EQUALS_DN`, `NOT_EQUAL`, `NOT_EQUAL_CASE_INSENSITIVE`, `NOT_EQUAL_DN`, `MULTIVALUE_CONTAINS`, `MULTIVALUE_CONTAINS_CASE_INSENSITIVE`, `MULTIVALUE_CONTAINS_DN`, `MULTIVALUE_DOES_NOT_CONTAIN`, `MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE`, `MULTIVALUE_DOES_NOT_CONTAIN_DN`.
- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria--expression_criteria--source))
- `value` (String) The expected value of this issuance criterion.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.

<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria--expression_criteria--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.issuance_criteria.expression_criteria.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref--issuance_criteria--expression_criteria"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref.issuance_criteria.expression_criteria`

Required:

- `expression` (String) The OGNL expression to evaluate.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.




<a id="nestedatt--sp_browser_sso--adapter_mappings--adapter_override_settings--parent_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.adapter_override_settings.parent_ref`

Required:

- `id` (String) The ID of the resource.



<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources`

Optional:

- `custom_attribute_source` (Attributes) The configured settings used to look up attributes from a custom data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--custom_attribute_source))
- `jdbc_attribute_source` (Attributes) The configured settings used to look up attributes from a JDBC data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--jdbc_attribute_source))
- `ldap_attribute_source` (Attributes) The configured settings used to look up attributes from a LDAP data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source))

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--custom_attribute_source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.custom_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--data_store_ref))

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `filter_fields` (Attributes Set) The list of fields that can be used to filter a request to the custom data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--filter_fields))
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--filter_fields"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.filter_fields`

Required:

- `name` (String) The name of this field.

Optional:

- `value` (String) The value of this field. Whether or not the value is required will be determined by plugin validation checks.



<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--jdbc_attribute_source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.jdbc_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--data_store_ref))
- `filter` (String) The JDBC WHERE clause used to query your data store to locate a user record.
- `table` (String) The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `column_names` (List of String) A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `schema` (String) Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.




<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--data_store_ref))
- `search_filter` (String) The LDAP filter that will be used to lookup the objects from the directory.
- `search_scope` (String) Determines the node depth of the query.
- `type` (String) The data store type of this attribute source.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `binary_attribute_settings` (Attributes Map) The advanced settings for binary LDAP attributes. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--binary_attribute_settings))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `member_of_nested_group` (Boolean) Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.
- `search_attributes` (Set of String) A list of LDAP attributes returned from search and available for mapping.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--attribute_sources--ldap_attribute_source--binary_attribute_settings"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.attribute_sources.ldap_attribute_source.binary_attribute_settings`

Optional:

- `binary_encoding` (String) Get the encoding type for this attribute. If not specified, the default is BASE64.




<a id="nestedatt--sp_browser_sso--adapter_mappings--idp_adapter_ref"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.idp_adapter_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.issuance_criteria`

Optional:

- `conditional_criteria` (Attributes Set) A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria--conditional_criteria))
- `expression_criteria` (Attributes Set) A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria--expression_criteria))

<a id="nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria--conditional_criteria"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.issuance_criteria.conditional_criteria`

Required:

- `attribute_name` (String) The name of the attribute to use in this issuance criterion.
- `condition` (String) The condition that will be applied to the source attribute's value and the expected value. Options are `EQUALS`, `EQUALS_CASE_INSENSITIVE`, `EQUALS_DN`, `NOT_EQUAL`, `NOT_EQUAL_CASE_INSENSITIVE`, `NOT_EQUAL_DN`, `MULTIVALUE_CONTAINS`, `MULTIVALUE_CONTAINS_CASE_INSENSITIVE`, `MULTIVALUE_CONTAINS_DN`, `MULTIVALUE_DOES_NOT_CONTAIN`, `MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE`, `MULTIVALUE_DOES_NOT_CONTAIN_DN`.
- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria--expression_criteria--source))
- `value` (String) The expected value of this issuance criterion.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.

<a id="nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria--expression_criteria--source"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.issuance_criteria.expression_criteria.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--adapter_mappings--issuance_criteria--expression_criteria"></a>
### Nested Schema for `sp_browser_sso.adapter_mappings.issuance_criteria.expression_criteria`

Required:

- `expression` (String) The OGNL expression to evaluate.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.




<a id="nestedatt--sp_browser_sso--assertion_lifetime"></a>
### Nested Schema for `sp_browser_sso.assertion_lifetime`

Required:

- `minutes_after` (Number) Assertion validity in minutes after the assertion issuance.
- `minutes_before` (Number) Assertion validity in minutes before the assertion issuance.


<a id="nestedatt--sp_browser_sso--attribute_contract"></a>
### Nested Schema for `sp_browser_sso.attribute_contract`

Optional:

- `core_attributes` (Attributes Set) A list of read-only assertion attributes (for example, SAML_SUBJECT) that are automatically populated by PingFederate. (see [below for nested schema](#nestedatt--sp_browser_sso--attribute_contract--core_attributes))
- `extended_attributes` (Attributes Set) A list of additional attributes that are added to the outgoing assertion. (see [below for nested schema](#nestedatt--sp_browser_sso--attribute_contract--extended_attributes))

<a id="nestedatt--sp_browser_sso--attribute_contract--core_attributes"></a>
### Nested Schema for `sp_browser_sso.attribute_contract.core_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `name_format` (String) The SAML Name Format for the attribute.


<a id="nestedatt--sp_browser_sso--attribute_contract--extended_attributes"></a>
### Nested Schema for `sp_browser_sso.attribute_contract.extended_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `name_format` (String) The SAML Name Format for the attribute.



<a id="nestedatt--sp_browser_sso--sso_service_endpoints"></a>
### Nested Schema for `sp_browser_sso.sso_service_endpoints`

Required:

- `url` (String) The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.

Optional:

- `binding` (String) The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`.
- `index` (Number) The priority of the endpoint.
- `is_default` (Boolean) Whether or not this endpoint is the default endpoint. Defaults to `false`.


<a id="nestedatt--sp_browser_sso--artifact"></a>
### Nested Schema for `sp_browser_sso.artifact`

Required:

- `lifetime` (Number) The lifetime of the artifact in seconds.
- `resolver_locations` (Attributes Set) Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message (see [below for nested schema](#nestedatt--sp_browser_sso--artifact--resolver_locations))

Optional:

- `source_id` (String) Source ID for SAML1.x connections

<a id="nestedatt--sp_browser_sso--artifact--resolver_locations"></a>
### Nested Schema for `sp_browser_sso.artifact.resolver_locations`

Required:

- `index` (Number) The priority of the endpoint.
- `url` (String) Remote party URLs that you will use to resolve/translate the artifact and get the actual protocol message



<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings`

Required:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_contract_fulfillment))
- `authentication_policy_contract_ref` (Attributes) Reference to the associated Authentication Policy Contract. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--authentication_policy_contract_ref))

Optional:

- `abort_sso_transaction_as_fail_safe` (Boolean) If set to true, SSO transaction will be aborted as a fail-safe when the data-store's attribute mappings fail to complete the attribute contract. Otherwise, the attribute contract with default values is used. By default, this value is `false`.
- `attribute_sources` (Attributes Set) A list of configured data stores to look up attributes from. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources))
- `issuance_criteria` (Attributes) The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria))
- `restrict_virtual_entity_ids` (Boolean) Restricts this mapping to specific virtual entity IDs.
- `restricted_virtual_entity_ids` (Set of String) The list of virtual server IDs that this mapping is restricted to.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--authentication_policy_contract_ref"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.authentication_policy_contract_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources`

Optional:

- `custom_attribute_source` (Attributes) The configured settings used to look up attributes from a custom data store. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--custom_attribute_source))
- `jdbc_attribute_source` (Attributes) The configured settings used to look up attributes from a JDBC data store. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--jdbc_attribute_source))
- `ldap_attribute_source` (Attributes) The configured settings used to look up attributes from a LDAP data store. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source))

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--custom_attribute_source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.custom_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--data_store_ref))

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `filter_fields` (Attributes Set) The list of fields that can be used to filter a request to the custom data store. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--filter_fields))
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--filter_fields"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.filter_fields`

Required:

- `name` (String) The name of this field.

Optional:

- `value` (String) The value of this field. Whether or not the value is required will be determined by plugin validation checks.



<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--jdbc_attribute_source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.jdbc_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--data_store_ref))
- `filter` (String) The JDBC WHERE clause used to query your data store to locate a user record.
- `table` (String) The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `column_names` (List of String) A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `schema` (String) Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.




<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--data_store_ref))
- `search_filter` (String) The LDAP filter that will be used to lookup the objects from the directory.
- `search_scope` (String) Determines the node depth of the query.
- `type` (String) The data store type of this attribute source.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `binary_attribute_settings` (Attributes Map) The advanced settings for binary LDAP attributes. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--binary_attribute_settings))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `member_of_nested_group` (Boolean) Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.
- `search_attributes` (Set of String) A list of LDAP attributes returned from search and available for mapping.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--attribute_sources--ldap_attribute_source--binary_attribute_settings"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.attribute_sources.ldap_attribute_source.binary_attribute_settings`

Optional:

- `binary_encoding` (String) Get the encoding type for this attribute. If not specified, the default is BASE64.




<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.issuance_criteria`

Optional:

- `conditional_criteria` (Attributes Set) A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria--conditional_criteria))
- `expression_criteria` (Attributes Set) A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria--expression_criteria))

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria--conditional_criteria"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.issuance_criteria.conditional_criteria`

Required:

- `attribute_name` (String) The name of the attribute to use in this issuance criterion.
- `condition` (String) The condition that will be applied to the source attribute's value and the expected value. Options are `EQUALS`, `EQUALS_CASE_INSENSITIVE`, `EQUALS_DN`, `NOT_EQUAL`, `NOT_EQUAL_CASE_INSENSITIVE`, `NOT_EQUAL_DN`, `MULTIVALUE_CONTAINS`, `MULTIVALUE_CONTAINS_CASE_INSENSITIVE`, `MULTIVALUE_CONTAINS_DN`, `MULTIVALUE_DOES_NOT_CONTAIN`, `MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE`, `MULTIVALUE_DOES_NOT_CONTAIN_DN`.
- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria--expression_criteria--source))
- `value` (String) The expected value of this issuance criterion.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.

<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria--expression_criteria--source"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.issuance_criteria.expression_criteria.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--sp_browser_sso--authentication_policy_contract_assertion_mappings--issuance_criteria--expression_criteria"></a>
### Nested Schema for `sp_browser_sso.authentication_policy_contract_assertion_mappings.issuance_criteria.expression_criteria`

Required:

- `expression` (String) The OGNL expression to evaluate.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.




<a id="nestedatt--sp_browser_sso--encryption_policy"></a>
### Nested Schema for `sp_browser_sso.encryption_policy`

Optional:

- `encrypt_assertion` (Boolean) Whether the outgoing SAML assertion will be encrypted.
- `encrypt_slo_subject_name_id` (Boolean) Encrypt the name-identifier attribute in outbound SLO messages. This can be set if the name id is encrypted.
- `encrypted_attributes` (Set of String) The list of outgoing SAML assertion attributes that will be encrypted. The `encrypt_assertion` property takes precedence over this.
- `slo_subject_name_id_encrypted` (Boolean) Allow the encryption of the name-identifier attribute for inbound SLO messages. This can be set if SP initiated SLO is enabled.


<a id="nestedatt--sp_browser_sso--message_customizations"></a>
### Nested Schema for `sp_browser_sso.message_customizations`

Optional:

- `context_name` (String) The context in which the customization will be applied. Depending on the connection type and protocol, this can either be `assertion`, `authn-response` or `authn-request`.
- `message_expression` (String) The OGNL expression that will be executed. Refer to the Admin Manual for a list of variables provided by PingFederate.


<a id="nestedatt--sp_browser_sso--slo_service_endpoints"></a>
### Nested Schema for `sp_browser_sso.slo_service_endpoints`

Required:

- `url` (String) The absolute or relative URL of the endpoint. A relative URL can be specified if a base URL for the connection has been defined.

Optional:

- `binding` (String) The binding of this endpoint, if applicable - usually only required for SAML 2.0 endpoints. Options are `ARTIFACT`, `POST`, `REDIRECT`, `SOAP`.
- `response_url` (String) The absolute or relative URL to which logout responses are sent. A relative URL can be specified if a base URL for the connection has been defined.


<a id="nestedatt--sp_browser_sso--url_whitelist_entries"></a>
### Nested Schema for `sp_browser_sso.url_whitelist_entries`

Optional:

- `allow_query_and_fragment` (Boolean) Allow Any Query/Fragment
- `require_https` (Boolean) Require HTTPS
- `valid_domain` (String) Valid Domain Name (leading wildcard '*.' allowed)
- `valid_path` (String) Valid Path (leave undefined to allow any path)



<a id="nestedatt--ws_trust"></a>
### Nested Schema for `ws_trust`

Required:

- `attribute_contract` (Attributes) A set of user attributes that this server will send in the token. (see [below for nested schema](#nestedatt--ws_trust--attribute_contract))
- `partner_service_ids` (Set of String) The partner service identifiers.
- `token_processor_mappings` (Attributes Set) A list of token processors to validate incoming tokens. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings))

Optional:

- `abort_if_not_fulfilled_from_request` (Boolean) If the attribute contract cannot be fulfilled using data from the Request, abort the transaction.
- `default_token_type` (String) The default token type when a web service client (WSC) does not specify in the token request which token type the STS should issue. Options are `SAML20`, `SAML11`, `SAML11_O365`. Defaults to `SAML20`.
- `encrypt_saml2_assertion` (Boolean) When selected, the STS encrypts the SAML 2.0 assertion. Applicable only to SAML 2.0 security token.  This option does not apply to OAuth assertion profiles.
- `generate_key` (Boolean) When selected, the STS generates a symmetric key to be used in conjunction with the "Holder of Key" (HoK) designation for the assertion's Subject Confirmation Method.  This option does not apply to OAuth assertion profiles.
- `message_customizations` (Attributes Set) The message customizations for WS-Trust. Depending on server settings, connection type, and protocol this may or may not be supported. (see [below for nested schema](#nestedatt--ws_trust--message_customizations))
- `minutes_after` (Number) The amount of time after the SAML token was issued during which it is to be considered valid. The default value is `30`.
- `minutes_before` (Number) The amount of time before the SAML token was issued during which it is to be considered valid. The default value is `5`.
- `oauth_assertion_profiles` (Boolean) When selected, four additional token-type requests become available.
- `request_contract_ref` (Attributes) Request Contract to be used to map attribute values into the security token. (see [below for nested schema](#nestedatt--ws_trust--request_contract_ref))

<a id="nestedatt--ws_trust--attribute_contract"></a>
### Nested Schema for `ws_trust.attribute_contract`

Optional:

- `core_attributes` (Attributes Set) A list of read-only assertion attributes that are automatically populated by PingFederate. (see [below for nested schema](#nestedatt--ws_trust--attribute_contract--core_attributes))
- `extended_attributes` (Attributes Set) A list of additional attributes that are added to the outgoing assertion. (see [below for nested schema](#nestedatt--ws_trust--attribute_contract--extended_attributes))

<a id="nestedatt--ws_trust--attribute_contract--core_attributes"></a>
### Nested Schema for `ws_trust.attribute_contract.core_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `namespace` (String) The attribute namespace.  This is required when the Default Token Type is SAML2.0 or SAML1.1 or SAML1.1 for Office 365.


<a id="nestedatt--ws_trust--attribute_contract--extended_attributes"></a>
### Nested Schema for `ws_trust.attribute_contract.extended_attributes`

Required:

- `name` (String) The name of this attribute.

Optional:

- `namespace` (String) The attribute namespace.  This is required when the Default Token Type is SAML2.0 or SAML1.1 or SAML1.1 for Office 365.



<a id="nestedatt--ws_trust--token_processor_mappings"></a>
### Nested Schema for `ws_trust.token_processor_mappings`

Required:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_contract_fulfillment))
- `idp_token_processor_ref` (Attributes) Reference to the associated token processor. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--idp_token_processor_ref))

Optional:

- `attribute_sources` (Attributes Set) A list of configured data stores to look up attributes from. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources))
- `issuance_criteria` (Attributes) The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--issuance_criteria))
- `restricted_virtual_entity_ids` (Set of String) The list of virtual server IDs that this mapping is restricted to.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_contract_fulfillment"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_contract_fulfillment--source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--ws_trust--token_processor_mappings--idp_token_processor_ref"></a>
### Nested Schema for `ws_trust.token_processor_mappings.idp_token_processor_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources`

Optional:

- `custom_attribute_source` (Attributes) The configured settings used to look up attributes from a custom data store. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--custom_attribute_source))
- `jdbc_attribute_source` (Attributes) The configured settings used to look up attributes from a JDBC data store. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--jdbc_attribute_source))
- `ldap_attribute_source` (Attributes) The configured settings used to look up attributes from a LDAP data store. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source))

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--custom_attribute_source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.custom_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--data_store_ref))

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `filter_fields` (Attributes Set) The list of fields that can be used to filter a request to the custom data store. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--filter_fields))
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--filter_fields"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.filter_fields`

Required:

- `name` (String) The name of this field.

Optional:

- `value` (String) The value of this field. Whether or not the value is required will be determined by plugin validation checks.



<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--jdbc_attribute_source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.jdbc_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--data_store_ref))
- `filter` (String) The JDBC WHERE clause used to query your data store to locate a user record.
- `table` (String) The name of the database table. The name is used to construct the SQL query to retrieve data from the data store.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `column_names` (List of String) A list of column names used to construct the SQL query to retrieve data from the specified table in the datastore.
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `schema` (String) Lists the table structure that stores information within a database. Some databases, such as Oracle, require a schema for a JDBC query. Other databases, such as MySQL, do not require a schema.

Read-Only:

- `type` (String) The data store type of this attribute source.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.




<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source`

Required:

- `data_store_ref` (Attributes) Reference to the associated data store. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--data_store_ref))
- `search_filter` (String) The LDAP filter that will be used to lookup the objects from the directory.
- `search_scope` (String) Determines the node depth of the query.
- `type` (String) The data store type of this attribute source.

Optional:

- `attribute_contract_fulfillment` (Attributes Map) Defines how an attribute in an attribute contract should be populated. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment))
- `base_dn` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `binary_attribute_settings` (Attributes Map) The advanced settings for binary LDAP attributes. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--binary_attribute_settings))
- `description` (String) The description of this attribute source. The description needs to be unique amongst the attribute sources for the mapping.<br>Note: Required for APC-to-SP Adapter Mappings
- `id` (String) The ID that defines this attribute source. Only alphanumeric characters allowed. Note: Required for OpenID Connect policy attribute sources, OAuth IdP adapter mappings, OAuth access token mappings and APC-to-SP Adapter Mappings. IdP Connections will ignore this property since it only allows one attribute source to be defined per mapping. IdP-to-SP Adapter Mappings can contain multiple attribute sources.
- `member_of_nested_group` (Boolean) Set this to true to return transitive group memberships for the 'memberOf' attribute.  This only applies for Active Directory data sources.  All other data sources will be set to false.
- `search_attributes` (Set of String) A list of LDAP attributes returned from search and available for mapping.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--data_store_ref"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.data_store_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment`

Required:

- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source))

Optional:

- `value` (String) The value for this attribute.

<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--attribute_contract_fulfillment--source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.attribute_contract_fulfillment.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--ws_trust--token_processor_mappings--attribute_sources--ldap_attribute_source--binary_attribute_settings"></a>
### Nested Schema for `ws_trust.token_processor_mappings.attribute_sources.ldap_attribute_source.binary_attribute_settings`

Optional:

- `binary_encoding` (String) Get the encoding type for this attribute. If not specified, the default is BASE64.




<a id="nestedatt--ws_trust--token_processor_mappings--issuance_criteria"></a>
### Nested Schema for `ws_trust.token_processor_mappings.issuance_criteria`

Optional:

- `conditional_criteria` (Attributes Set) A list of conditional issuance criteria where existing attributes must satisfy their conditions against expected values in order for the transaction to continue. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--issuance_criteria--conditional_criteria))
- `expression_criteria` (Attributes Set) A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue. Expressions must be enabled in PingFederate to use expression criteria. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--issuance_criteria--expression_criteria))

<a id="nestedatt--ws_trust--token_processor_mappings--issuance_criteria--conditional_criteria"></a>
### Nested Schema for `ws_trust.token_processor_mappings.issuance_criteria.conditional_criteria`

Required:

- `attribute_name` (String) The name of the attribute to use in this issuance criterion.
- `condition` (String) The condition that will be applied to the source attribute's value and the expected value. Options are `EQUALS`, `EQUALS_CASE_INSENSITIVE`, `EQUALS_DN`, `NOT_EQUAL`, `NOT_EQUAL_CASE_INSENSITIVE`, `NOT_EQUAL_DN`, `MULTIVALUE_CONTAINS`, `MULTIVALUE_CONTAINS_CASE_INSENSITIVE`, `MULTIVALUE_CONTAINS_DN`, `MULTIVALUE_DOES_NOT_CONTAIN`, `MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE`, `MULTIVALUE_DOES_NOT_CONTAIN_DN`.
- `source` (Attributes) The attribute value source. (see [below for nested schema](#nestedatt--ws_trust--token_processor_mappings--issuance_criteria--expression_criteria--source))
- `value` (String) The expected value of this issuance criterion.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.

<a id="nestedatt--ws_trust--token_processor_mappings--issuance_criteria--expression_criteria--source"></a>
### Nested Schema for `ws_trust.token_processor_mappings.issuance_criteria.expression_criteria.source`

Required:

- `type` (String) The source type of this key. Options are `TOKEN_EXCHANGE_PROCESSOR_POLICY`, `ACCOUNT_LINK`, `ADAPTER`, `ASSERTION`, `CONTEXT`, `CUSTOM_DATA_STORE`, `EXPRESSION`, `JDBC_DATA_STORE`, `LDAP_DATA_STORE`, `PING_ONE_LDAP_GATEWAY_DATA_STORE`, `MAPPED_ATTRIBUTES`, `NO_MAPPING`, `TEXT`, `TOKEN`, `REQUEST`, `OAUTH_PERSISTENT_GRANT`, `SUBJECT_TOKEN`, `ACTOR_TOKEN`, `PASSWORD_CREDENTIAL_VALIDATOR`, `IDP_CONNECTION`, `AUTHENTICATION_POLICY_CONTRACT`, `CLAIMS`, `LOCAL_IDENTITY_PROFILE`, `EXTENDED_CLIENT_METADATA`, `EXTENDED_PROPERTIES`, `TRACKED_HTTP_PARAMS`, `FRAGMENT`, `INPUTS`, `ATTRIBUTE_QUERY`, `IDENTITY_STORE_USER`, `IDENTITY_STORE_GROUP`, `SCIM_USER`, `SCIM_GROUP`.

Optional:

- `id` (String) The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.



<a id="nestedatt--ws_trust--token_processor_mappings--issuance_criteria--expression_criteria"></a>
### Nested Schema for `ws_trust.token_processor_mappings.issuance_criteria.expression_criteria`

Required:

- `expression` (String) The OGNL expression to evaluate.

Optional:

- `error_result` (String) The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.




<a id="nestedatt--ws_trust--message_customizations"></a>
### Nested Schema for `ws_trust.message_customizations`

Optional:

- `context_name` (String) The context in which the customization will be applied. Depending on the connection type and protocol, this can either be `assertion`, `authn-response` or `authn-request`.
- `message_expression` (String) The OGNL expression that will be executed. Refer to the Admin Manual for a list of variables provided by PingFederate.


<a id="nestedatt--ws_trust--request_contract_ref"></a>
### Nested Schema for `ws_trust.request_contract_ref`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "connectionId" should be the id of the SP Connection to be imported

```shell
terraform import pingfederate_idp_sp_connection.idpSpConnection connectionId
```