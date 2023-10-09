terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username               = "administrator"
  password               = "2FederateM0re"
  https_host             = "https://localhost:9999"
  insecure_trust_all_tls = true
}

resource "pingfederate_local_identity_identity_profile" "myLocalIdentityIdentityProfile" {
  name = "yourIdentityProfileName"
  apc_id = {
    id = "apcid"
  }
}

data "pingfederate_local_identity_identity_profile" "myLocalIdentityIdentityProfile" {
  name = pingfederate_local_identity_identity_profile.myLocalIdentityIdentityProfile.name
  apc_id = {
    id = "apcid"
  }
}
resource "pingfederate_local_identity_identity_profile" "localIdentityIdentityProfileExample" {
  name = "${data.pingfederate_local_identity_identity_profile.myLocalIdentityIdentityProfile.name}2"
  apc_id = {
    id = "apcid"
  }
}