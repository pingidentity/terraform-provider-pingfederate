resource "pingfederate_sp_idp_connection" "spIdpConnection" {
  name      = "connection name"
  entity_id = "entity_id"
  active    = true
  contact_info = {
    first_name = "FirstName"
  }
  logging_mode = "STANDARD"
  credentials = {
    certs = [
      {
        x509_file = {
          id        = "fileId"
          file_data = file("./assets/x509.crt")
        }
        encryption_cert          = false
        active_verification_cert = true
      }
    ]
    inbound_back_channel_auth = {
      type = "INBOUND"
      http_basic_credentials = {
        username = "admin"
        password = var.sp_idp_connection_inbound_back_channel_auth_password
      }
      digital_signature = true
      require_ssl       = false
    }
  }

  error_page_msg_id = "errorDetail.spSsoFailure"
  idp_browser_sso = {
    protocol              = "SAML20"
    enabled_profiles      = ["IDP_INITIATED_SSO"]
    incoming_bindings     = ["POST"]
    default_target_url    = ""
    sso_service_endpoints = []
    idp_identity_mapping  = "ACCOUNT_MAPPING"
    adapter_mappings = [
      {
        attribute_sources = []
        attribute_contract_fulfillment = {
          subject = {
            source = {
              type = "NO_MAPPING"
            }
          }
        }
        issuance_criteria = {
          conditional_criteria = []
        }
        restrict_virtual_entity_ids   = false
        restricted_virtual_entity_ids = []
        sp_adapter_ref = {
          id = "spadapterid"
        }
      }
    ]
    authentication_policy_contract_mappings = []
    assertions_signed                       = false
    sign_authn_requests                     = false
  }
  idp_oauth_grant_attribute_mapping = {
    access_token_manager_mappings = [
      {
        attribute_contract_fulfillment = {
          "username" = {
            source = {
              type = "NO_MAPPING"
            }
          }
          "org_name" = {
            source = {
              type = "NO_MAPPING"
            }
          }
        }
        access_token_manager_ref = {
          id = "jwt"
        }
      }
    ]
  }
  ws_trust = {
    attribute_contract = {
      core_attributes = [
        {
          name   = "TOKEN_SUBJECT"
          masked = false
        }
      ]
      extended_attributes = [
        {
          name   = "test"
          masked = false
        }
      ]
    }
    token_generator_mappings = [
      {
        attribute_sources = []
        attribute_contract_fulfillment = {
          saml_subject = {
            source = {
              type = "NO_MAPPING"
            }
          }
        }
        sp_token_generator_ref = {
          id = "tokengeneratorid"
        }
        restricted_virtual_entity_ids = []
        default_mapping               = true
      }
    ]
    generate_local_token = true
  }
  inbound_provisioning = {
    group_support = false
    user_repository = {
      identity_store = {
        identity_store_provisioner_ref = {
          id = "identityStoreProvisioner"
        }
      }
    }
    custom_schema = {
      namespace  = "urn:scim:schemas:extension:custom:1.0"
      attributes = []
    }
    users = {
      write_users = {
        attribute_fulfillment = {
          "username" = {
            source = {
              type = "TEXT"
            }
            value = "username"
          }
        }
      }
      read_users = {
        attribute_contract = {
          core_attributes = []
          extended_attributes = [
            {
              name   = "userName"
              masked = false
            }
          ]
        }
        attributes = []
        attribute_fulfillment = {
          user_name = {
            source = {
              type = "TEXT"
            }
            value = "username"
          }
        }
      }
    }
  }
}
