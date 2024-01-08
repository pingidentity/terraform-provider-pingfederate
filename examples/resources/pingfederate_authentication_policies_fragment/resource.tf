terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 1.0.0"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username   = "administrator"
  password   = "2FederateM0re"
  https_host = "https://localhost:9999"
  # Warning: The insecure_trust_all_tls attribute configures the provider to trust any certificate presented by the PingDirectory server.
  insecure_trust_all_tls              = true
  x_bypass_external_validation_header = true
  product_version                     = "12.0.0"
}


resource "pingfederate_authentication_policies_fragment" "myAuthenticationPolicyFragment" {
  fragment_id = "myPolicyFragment"
  name        = "myFragment"
  root_node = {
    policy_action = {
      authn_source_policy_action = {
        authentication_source = {
          type = "IDP_ADAPTER"
          source_ref = {
            id = "OTIdPJava"
          }
        }
      }
    },
    children = [
      {
        policy_action = {
          done_policy_action = {
            context = "Fail"
          }
        }
      },
      {
        policy_action = {
          done_policy_action = {
            context = "Success"
          }
        }
      }
    ]
  }
}
