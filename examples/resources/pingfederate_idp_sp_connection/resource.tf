terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 1.0.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username   = "administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:9999"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the server.
  insecure_trust_all_tls = true
  x_bypass_external_validation_header = true
  product_version = "12.1"
}


resource "pingfederate_idp_sp_connection" "wsFedSpBrowserSSOExample" {
#   name          = "wsfedspconn1"
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
        id = "419x9yg43rlawqwq9v6az997k"
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

# resource "pingfederate_idp_sp_connection" "outboundProvisionExample" {
#   connection_id = "connectionId"
#   name          = "PingOne Connector"
#   entity_id     = "entity"
#   active        = false
#   contact_info = {
#     company = "Example Corp"
#   }
#   base_url               = "https://api.pingone.com/v5"
#   logging_mode           = "STANDARD"
#   connection_target_type = "STANDARD"
#   credentials = {
#     certs = []
#     signing_settings = {
#       signing_key_pair_ref = {
#         id = ""
#       }
#       include_raw_key_in_signature = false
#       include_cert_in_signature    = false
#       algorithm                    = "SHA256withRSA"
#     }
#   }
#   outbound_provision = {
#     type = "PingOne"
#     target_settings = [
#       {
#         name  = "PINGONE_ENVIRONMENT"
#         value = "example"
#       }
#     ]
#     channels = [
#       {
#         name        = "Channel1"
#         max_threads = 1
#         timeout     = 120
#         active      = false
#         channel_source = {
#           base_dn = "dc=example,dc=com"
#           data_source = {
#             id = "pingdirectory"
#           }
#           guid_attribute_name = "entry_uuid"
#           change_detection_settings = {
#             user_object_class         = "inetOrgPerson"
#             changed_users_algorithm   = "TIMESTAMP_NO_NEGATION"
#             group_object_class        = "groupOfUniqueNames"
#             time_stamp_attribute_name = "modifyTimestamp"
#           }
#           account_management_settings = {
#             account_status_algorithm      = "ACCOUNT_STATUS_ALGORITHM_FLAG"
#             account_status_attribute_name = "nsaccountlock"
#             flag_comparison_value         = "true"
#             flag_comparison_status        = true
#             default_status                = true
#           }
#           group_membership_detection = {
#             group_member_attribute_name = "uniqueMember"
#           }
#           guid_binary = false
#           user_source_location = {
#             filter = "cn=John"

#           }
#         }
#         attribute_mapping = [
#           {
#             field_name = "username"
#             saas_field_info = {
#               attribute_names = [
#                 "uid"
#               ]
#             }
#           },
#           {
#             field_name = "email"
#             saas_field_info = {
#               attribute_names = [
#                 "mail"
#               ]
#             }
#           },
#           {
#             field_name = "populationID"
#             saas_field_info = {
#               default_value = "example"
#             }
#           }
#         ]
#       }
#     ]
#   }
# }

# resource "pingfederate_idp_sp_connection" "wsTrustExample" {
#   connection_id      = "connection"
#   name               = "connection"
#   entity_id          = "entity"
#   active             = true
#   contact_info       = {}
#   base_url           = "https://localhost:9031"
#   logging_mode       = "STANDARD"
#   virtual_entity_ids = []
#   credentials = {
#     certs = []
#     signing_settings = {
#       signing_key_pair_ref = {
#         id = "signingKey"
#       }
#       include_raw_key_in_signature = false
#       include_cert_in_signature    = false
#       algorithm                    = "SHA256withRSA"
#     }
#   }
#   ws_trust = {
#     partner_service_ids = [
#       "id"
#     ]
#     oauth_assertion_profiles = true
#     default_token_type       = "SAML20"
#     generate_key             = false
#     encrypt_saml2_assertion  = false
#     minutes_before           = 5
#     minutes_after            = 30
#     attribute_contract = {
#       core_attributes = [
#         {
#           name = "TOKEN_SUBJECT"
#         }
#       ]
#       extended_attributes = []
#     }
#     token_processor_mappings = [
#       {
#         attribute_sources = []
#         attribute_contract_fulfillment = {
#           "TOKEN_SUBJECT" : {
#             source = {
#               type = "TOKEN"
#             }
#             value = "username"
#           }
#         }
#         issuance_criteria = {
#           conditional_criteria = []
#         }
#         idp_token_processor_ref = {
#           id = "UsernameTokenProcessor"
#         }
#         restricted_virtual_entity_ids = []
#       }
#     ]
#   }
#   connection_target_type = "STANDARD"
# }

# resource "pingfederate_idp_sp_connection" "samlSpBrowserSSOExample" {
#   connection_id      = "connection"
#   name               = "connection"
#   entity_id          = "entity"
#   active             = true
#   contact_info       = {}
#   base_url           = "https://localhost:9032"
#   logging_mode       = "STANDARD"
#   virtual_entity_ids = []
#   credentials = {
#     certs = []
#     signing_settings = {
#       signing_key_pair_ref = {
#         id = "signingKey"
#       }
#       include_raw_key_in_signature = false
#       include_cert_in_signature    = false
#       algorithm                    = "SHA256withRSA"
#     }
#   }
#   sp_browser_sso = {
#     protocol                      = "SAML20"
#     require_signed_authn_requests = false
#     sp_saml_identity_mapping      = "STANDARD"
#     sign_assertions               = false
#     authentication_policy_contract_assertion_mappings = [
#       {
#         abort_sso_transaction_as_fail_safe = false
#         authentication_policy_contract_ref = {
#           id = "contractId"
#         }
#         restricted_virtual_entity_ids = []
#         attribute_contract_fulfillment = {
#           "SAML_SUBJECT" = {
#             source = {
#               type = "AUTHENTICATION_POLICY_CONTRACT"
#             }
#             value = "subject"
#           }
#         }
#         restrict_virtual_entity_ids = false
#         attribute_sources           = []
#         issuance_criteria = {
#           conditional_criteria = []
#         }
#       }
#     ]
#     encryption_policy = {
#       encrypt_slo_subject_name_id   = false
#       encrypt_assertion             = false
#       encrypted_attributes          = []
#       slo_subject_name_id_encrypted = false
#     }
#     enabled_profiles = [
#       "IDP_INITIATED_SSO"
#     ]
#     sign_response_as_required = true
#     sso_service_endpoints = [
#       {
#         is_default = true
#         binding    = "POST"
#         index      = 0
#         url        = "https://httpbin.org/anything"
#       }
#     ]
#     adapter_mappings = []
#     assertion_lifetime = {
#       minutes_after  = 5
#       minutes_before = 5
#     }
#     attribute_contract = {
#       core_attributes = [
#         {
#           name_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
#           name        = "SAML_SUBJECT"
#         }
#       ]
#       extended_attributes = []
#     }
#   }
# }
