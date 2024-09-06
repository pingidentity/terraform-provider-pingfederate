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