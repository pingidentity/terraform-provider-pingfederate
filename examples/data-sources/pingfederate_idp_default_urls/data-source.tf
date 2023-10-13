resource "pingfederate_idp_default_urls" "myIdpDefaultUrl" {
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "https://example"
}

data "pingfederate_idp_default_urls" "myIdpDefaultUrl" {
  depends_on = [
    pingfederate_idp_default_urls.myIdpDefaultUrl
  ]
}
resource "pingfederate_idp_default_urls" "idpDefaultUrlExample" {
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "${data.pingfederate_idp_default_urls.myIdpDefaultUrl.idp_slo_success_url}.com"
}