resource "pingfederate_captcha_provider" "reCAPTCHAv2ProviderExample" {
  provider_id = "myreCAPTCHAv2ProviderId"
  name        = "My reCAPTCHA v2 Provider"
  configuration = {
    fields = [
      {
        name  = "Site Key"
        value = var.recaptcha_v2_site_key
      }
    ]
    sensitive_fields = [
      {
        name  = "Secret Key"
        value = var.recaptcha_v2_secret_key
      }
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.captcha.ReCaptchaV2InvisiblePlugin"
  }
}