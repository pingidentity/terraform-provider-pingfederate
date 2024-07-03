resource "pingfederate_captcha_providers" "captchaProvidersExample" {
  providers_id = "myCaptchaProvider"
  name         = "My Captcha Provider"
  configuration = {
    tables = [],
    fields = [
      {
        name  = "Site Key"
        value = "exampleSiteKey"
      },
      {
        name  = "Secret Key"
        value = "exampleSecretKey"
      }
    ]
  }
  // class name of the plugin
  // Captcha V2 used here
  plugin_descriptor_ref = {
    id = "com.pingidentity.captcha.ReCaptchaV2InvisiblePlugin"
  }
}