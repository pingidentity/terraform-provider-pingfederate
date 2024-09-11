resource "pingfederate_oauth_client_registration_policy" "registrationPolicy" {
  policy_id = "myRegistrationPolicy"
  name      = "My client registration policy"

  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.client.registration.ResponseTypesConstraintsPlugin"
  }

  configuration = {
    fields = [
      {
        name  = "code"
        value = "true"
      },
      {
        name  = "code id_token"
        value = "true"
      },
      {
        name  = "code id_token token"
        value = "true"
      },
      {
        name  = "code token"
        value = "true"
      },
      {
        name  = "id_token"
        value = "true"
      },
      {
        name  = "id_token token"
        value = "true"
      },
      {
        name  = "token"
        value = "true"
      }
    ]
  }
}
