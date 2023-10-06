resource "pingfederate_idp_default_urls" "myIdpDefaultUrl" {
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "https://example"
}
