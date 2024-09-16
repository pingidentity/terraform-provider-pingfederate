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
