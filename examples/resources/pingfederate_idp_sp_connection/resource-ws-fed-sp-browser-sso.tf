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
