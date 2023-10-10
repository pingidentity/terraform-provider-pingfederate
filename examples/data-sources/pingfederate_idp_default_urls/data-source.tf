terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username               = "administrator"
  password               = "2FederateM0re"
  https_host             = "https://localhost:9999"
  insecure_trust_all_tls = true
}

resource "pingfederate_idp_default_urls" "myIdpDefaultUrl" {
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "https://example"
}

data "pingfederate_idp_default_urls" "myIdpDefaultUrl" {
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "pingfederate_idp_default_urls.myIdpDefaultUrl.idp_slo_success_url"
}
resource "pingfederate_idp_default_urls" "idpDefaultUrlExample" {
  confirm_idp_slo     = true
  idp_error_msg       = "errorDetail.idpSsoFailure"
  idp_slo_success_url = "${data.pingfederate_idp_default_urls.myIdpDefaultUrl.idp_slo_success_url}.com"
}