resource "pingfederate_secret_manager" "secretManager" {
  manager_id = "mySecretManager"
  configuration = {
    fields = [
      {
        name  = "APP ID"
        value = "myappid"
      }
    ]
  }
  name = "My Secret Manager"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.secretmanagers.cyberark.CyberArkCredentialProvider"
  }
}