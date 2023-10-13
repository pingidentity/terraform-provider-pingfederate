resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractsExample" {
  custom_id           = "test"
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
}

resource "pingfederate_local_identity_identity_profile" "myLocalIdentityIdentityProfile" {
  custom_id = "id"
  name      = "yourIdentityProfileName"
  apc_id = {
    id = pingfederate_authentication_policy_contract.authenticationPolicyContractsExample.id
  }
  registration_enabled = false
  profile_enabled      = false
}

data "pingfederate_local_identity_identity_profile" "myLocalIdentityIdentityProfile" {
  id = pingfederate_local_identity_identity_profile.myLocalIdentityIdentityProfile.custom_id
}

resource "pingfederate_local_identity_identity_profile" "localIdentityIdentityProfileExample" {
  name = "${data.pingfederate_local_identity_identity_profile.myLocalIdentityIdentityProfile.name}2"
  apc_id = {
    id = "apcid"
  }
}