resource "pingfederate_certificates_revocation_settings" "revocationSettings" {
  ocsp_settings = {
    requester_add_nonce             = false
    action_on_responder_unavailable = "CONTINUE"
    action_on_status_unknown        = "FAIL"
    action_on_unsuccessful_response = "FAIL"
    current_update_grace_period     = 5
    next_update_grace_period        = 5
    response_cache_period           = 48
    responder_timeout               = 5
  }
}
