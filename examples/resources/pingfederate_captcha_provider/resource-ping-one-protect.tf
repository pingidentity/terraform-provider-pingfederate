resource "pingfederate_pingone_connection" "example" {
  name        = "My PingOne Environment"
  description = "My environment"
  credential  = var.pingone_gateway_credential
}

resource "pingfederate_captcha_provider" "riskProviderExample" {
  provider_id = "pingoneProtectProviderId"
  name        = "PingOne Protect Provider"
  configuration = {
    fields = [
      {
        "name" : "PingOne Environment",
        "value" : format("%s|%s", pingfederate_pingone_connection.example.id, var.pingone_environment_id)
      },
      {
        "name" : "PingOne Risk Policy",
        "value" : var.pingone_risk_policy_id
      },
      {
        "name" : "Enable Risk Evaluation",
        "value" : "true"
      },
      {
        "name" : "Password Encryption",
        "value" : "SHA-256"
      },
      {
        "name" : "Follow Recommended Action",
        "value" : "true"
      },
      {
        "name" : "Failure Mode",
        "value" : "Continue with fallback policy decision"
      },
      {
        "name" : "Fallback Policy Decision Value",
        "value" : "MEDIUM"
      },
      {
        "name" : "API Request Timeout",
        "value" : "2000"
      },
      {
        "name" : "Proxy Settings",
        "value" : "System Defaults"
      },
      {
        "name" : "Custom Proxy Host",
        "value" : ""
      },
      {
        "name" : "Custom Proxy Port",
        "value" : ""
      },
      {
        "name" : "Custom connection pool",
        "value" : "50"
      }
    ]
  }
  plugin_descriptor_ref = {
    id = "com.pingidentity.adapters.pingone.protect.PingOneProtectProvider"
  }
}
