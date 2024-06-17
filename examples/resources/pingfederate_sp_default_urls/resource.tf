resource "pingfederate_sp_default_urls" "spDefaultUrlsExample" {
  confirm_slo     = true
  slo_success_url = "https://example.com/slo_success_url"
  sso_success_url = "https://example.com/sso_success_url"
}
