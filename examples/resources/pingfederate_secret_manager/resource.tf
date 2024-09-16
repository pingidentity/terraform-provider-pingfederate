resource "pingfederate_secret_manager" "myCyberArkSecretManager" {
  manager_id = "myCyberArkSecretManager"
  name       = "My CyberArk Secret Manager"

  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.secretmanagers.cyberark.CyberArkCredentialProvider"
  }

  configuration = {
    fields = [
      {
        name  = "Connection Port",
        value = "18923"
      },
      {
        name  = "APP ID"
        value = "myappid"
      },
      {
        name  = "Connection Timeout (sec)",
        value = "30"
      },
      {
        name  = "Username Retrieval Property Name",
        value = "username"
      }
    ]
  }
}