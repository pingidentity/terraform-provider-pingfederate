resource "pingfederate_identity_store_provisioner" "identityStoreProvisioner" {
  provisioner_id = "provisioner-id"
  attribute_contract = {
    core_attributes = [
      {
        name = "username"
      }
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Delete user behavior"
        value = "Disable User"
      }
    ]
  }
  group_attribute_contract = {
    core_attributes = [
      {
        name = "groupname"
    }]
  }
  name = "My Identity Store Provisioner"
  plugin_descriptor_ref = {
    id = "com.pingidentity.identitystoreprovisioners.sample.SampleIdentityStoreProvisioner"
  }
}