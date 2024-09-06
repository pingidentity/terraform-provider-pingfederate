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
