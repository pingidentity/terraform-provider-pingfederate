resource "pingfederate_default_urls" "defaultUrlsExample" {
  confirm_sp_slo      = true
  sp_slo_success_url  = "https://example.com/slo_success_url"
  sp_sso_success_url  = "https://example.com/sso_success_url"
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "https://example.com/idp_slo_success_url"
}
