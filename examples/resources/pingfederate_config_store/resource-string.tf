resource "pingfederate_config_store" "enable_expressions" {
  bundle       = "org.sourceid.common.ExpressionManager"
  setting_id   = "evaluateExpressions"
  string_value = "true"
}