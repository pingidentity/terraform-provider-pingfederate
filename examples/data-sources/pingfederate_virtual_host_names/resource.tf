terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username = "administrator"
  password = "2FederateM0re"
  https_host = "https://localhost:9999"
  insecure_trust_all_tls = true
}


data "pingfederate_virtual_host_names" "myVirtualHostNamesExample" {
}

resource "pingfederate_virtual_host_names" "myVirtualHostNamesExample" {
  virtual_host_names = ["example1", "example2"]
}
