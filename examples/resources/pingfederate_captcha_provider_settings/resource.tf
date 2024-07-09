resource "pingfederate_captcha_provider_settings" "captchaProvidersSettings" {
  default_captcha_provider_ref = {
    id = "myDefaultCaptchaProvider"
  }
}