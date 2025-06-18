resource "pingfederate_config_store" "base64_required_plugins" {
  bundle     = "org.sourceid.oauth20.handlers.process.exchange.execution.SecurityTokenCreator"
  setting_id = "base64-required-plugins"
  list_value = ["org.sourceid.wstrust.processor.oauth.BearerAccessTokenTokenProcessor"]
}