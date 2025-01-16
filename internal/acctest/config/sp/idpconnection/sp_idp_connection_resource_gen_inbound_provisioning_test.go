// Code generated by ping-terraform-plugin-framework-generator

package resource_sp_idp_connection_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const idpConnInboundProvisioningId = "inboundproconn"

func TestAccSpIdpConnection_InboundProvisioningMinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: spIdpConnection_InboundProvisioningCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: spIdpConnection_InboundProvisioningMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesInboundProvisioningMinimal(),
			},
			{
				// Delete the minimal model
				Config:  spIdpConnection_InboundProvisioningMinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: spIdpConnection_InboundProvisioningCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesInboundProvisioningComplete(),
			},
			{
				// Back to minimal model
				Config: spIdpConnection_InboundProvisioningMinimalHCL(),
				Check:  spIdpConnection_CheckComputedValuesInboundProvisioningMinimal(),
			},
			{
				// Back to complete model
				Config: spIdpConnection_InboundProvisioningCompleteHCL(),
				Check:  spIdpConnection_CheckComputedValuesInboundProvisioningComplete(),
			},
			{
				// Test importing the resource
				Config:            spIdpConnection_InboundProvisioningCompleteHCL(),
				ResourceName:      "pingfederate_sp_idp_connection.example",
				ImportStateId:     idpConnInboundProvisioningId,
				ImportState:       true,
				ImportStateVerify: true,
				// passwords won't be returned by the API
				// encrypted_passwords change on each get
				ImportStateVerifyIgnore: []string{
					"credentials.inbound_back_channel_auth.http_basic_credentials.password",
					"credentials.inbound_back_channel_auth.http_basic_credentials.encrypted_password",
				},
			},
		},
	})
}

func spIdpConnection_InboundProvisioningDependencyHCL() string {
	return `
resource "pingfederate_data_store" "example" {
  data_store_id = "addatastore"
  ldap_data_store = {
    binary_attributes     = ["objectGUID"]
    bind_anonymously      = false
    connection_timeout    = -1
    create_if_necessary   = true
    dns_ttl               = 60000
    password              = "2FederateM0re"
    follow_ldap_referrals = false
    hostnames             = ["ldaps2.pf.ping-eng.com"]
    hostnames_tags = [
      {
        default_source = true
        hostnames      = ["ldaps2.pf.ping-eng.com"]
      },
    ]
    ldap_dns_srv_prefix     = "_ldap._tcp"
    ldap_type               = "ACTIVE_DIRECTORY"
    max_connections         = 75
    max_wait                = -1
    min_connections         = 15
    name                    = "ldaps3.pf.ping-eng.com (cn=localadmin,cn=users,dc=example,dc=com)"
    read_timeout            = -1
    retry_failed_operations = false
    test_on_borrow          = true
    test_on_return          = false
    time_between_evictions  = -1
    use_dns_srv_records     = false
    use_ssl                 = true
    use_start_tls           = false
    user_dn                 = "cn=localadmin,cn=users,dc=ldaps2,dc=com"
    verify_host             = false
  }
  mask_attribute_values = false
}
  `
}

// Minimal HCL with only required values set
func spIdpConnection_InboundProvisioningMinimalHCL() string {
	return fmt.Sprintf(`
%s

resource "pingfederate_sp_idp_connection" "example" {
  connection_id = "%s"
  credentials = {
    inbound_back_channel_auth = {
      http_basic_credentials = {
        password = "2FederateM0re"
        username = "uname2012"
      }
    }
  }
  entity_id = "inbound_AD2012"
  inbound_provisioning = {
    action_on_delete = "PERMANENTLY_DELETE_USER"
    custom_schema = {
      attributes = [
      ]
    }
    group_support = false
    user_repository = {
      ldap = {
        base_dn = "OU=SQE Testing,OU=Resources,DC=ldaps2,DC=com"
        data_store_ref = {
          id = pingfederate_data_store.example.id
        }
        unique_user_id_filter = "CN=$${userName}"
      }
    }
    users = {
      read_users = {
        attribute_contract = {
        }
        attribute_fulfillment = {
          "addresses.work.country" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "c"
          }
        }
        attributes = [
          {
            name = "c"
          },
        ]
      }
      write_users = {
        attribute_fulfillment = {
          c = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.work.country"
          }
        }
      }
    }
  }
  name = "inbound_AD2012"
}
`, spIdpConnection_InboundProvisioningDependencyHCL(), idpConnInboundProvisioningId)
}

// Maximal HCL with all values set where possible
func spIdpConnection_InboundProvisioningCompleteHCL() string {
	return fmt.Sprintf(`
%s

resource "pingfederate_sp_idp_connection" "example" {
  active        = true
  base_url      = "https://example.com"
  connection_id = "%s"
  contact_info = {
    company    = "Ping Identity"
    email      = "test@test.com"
    first_name = "test"
    last_name  = "test"
    phone      = "555-5555"
  }
  credentials = {
    block_encryption_algorithm = null
    certs = [
      {
        active_verification_cert    = true
        encryption_cert             = false
        primary_verification_cert   = true
        secondary_verification_cert = false
        x509_file = {
          file_data = "-----BEGIN CERTIFICATE-----\nMIIDOjCCAiICCQCjbB7XBVkxCzANBgkqhkiG9w0BAQsFADBfMRIwEAYDVQQDDAlsb2NhbGhvc3Qx\nDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoM\nA0NEUjELMAkGA1UEBhMCVVMwHhcNMjMwNzE0MDI1NDUzWhcNMjQwNzEzMDI1NDUzWjBfMRIwEAYD\nVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsM\nBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\nggEKAoIBAQC5yFrh9VR2wk9IjzMz+Ei80K453g1j1/Gv3EQ/SC9h7HZBI6aV9FaEYhGnaquRT5q8\n7p8lzCphKNXVyeL6T/pDJOW70zXItkl8Ryoc0tIaknRQmj8+YA0Hr9GDdmYev2yrxSoVS7s5Bl8p\noasn3DljgnWT07vsQz+hw3NY4SPp7IFGP2PpGUBBIIvrOaDWpPGsXeznBxSFtis6Qo+JiEoaVql9\nb9/XyKZj65wOsVyZhFWeM1nCQITSP9OqOc9FSoDFYQ1AVogm4A2AzUrkMnT1SrN2dCuTmNbeVw7g\nOMqMrVf0CiTv9hI0cATbO5we1sPAlJxscSkJjsaI+sQfjiAnAgMBAAEwDQYJKoZIhvcNAQELBQAD\nggEBACgwoH1qklPF1nI9+WbIJ4K12Dl9+U3ZMZa2lP4hAk1rMBHk9SHboOU1CHDQKT1Z6uxi0NI4\nJZHmP1qP8KPNEWTI8Q76ue4Q3aiA53EQguzGb3SEtyp36JGBq05Jor9erEebFftVl83NFvio72Fn\n0N2xvu8zCnlylf2hpz9x1i01Xnz5UNtZ2ppsf2zzT+4U6w3frH+pkp0RDPuoe9mnBF001AguP31h\nSBZyZzWcwQltuNELnSRCcgJl4kC2h3mAgaVtYalrFxLRa3tA2XF2BHRHmKgocedVhTq+81xrqj+W\nQuDmUe06DnrS3Ohmyj3jhsCCluznAolmrBhT/SaDuGg=\n-----END CERTIFICATE-----\n"
          id        = "4qrossmq1vxa4p836kyqzp48h"
        }
      },
    ]
    inbound_back_channel_auth = {
      certs = [
        {
          active_verification_cert    = true
          encryption_cert             = false
          primary_verification_cert   = true
          secondary_verification_cert = false
          x509_file = {
            file_data = "-----BEGIN CERTIFICATE-----\nMIIDOjCCAiICCQCjbB7XBVkxCzANBgkqhkiG9w0BAQsFADBfMRIwEAYDVQQDDAlsb2NhbGhvc3Qx\nDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsMBFBJTkcxDDAKBgNVBAoM\nA0NEUjELMAkGA1UEBhMCVVMwHhcNMjMwNzE0MDI1NDUzWhcNMjQwNzEzMDI1NDUzWjBfMRIwEAYD\nVQQDDAlsb2NhbGhvc3QxDjAMBgNVBAgMBVRFWEFTMQ8wDQYDVQQHDAZBVVNUSU4xDTALBgNVBAsM\nBFBJTkcxDDAKBgNVBAoMA0NEUjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw\nggEKAoIBAQC5yFrh9VR2wk9IjzMz+Ei80K453g1j1/Gv3EQ/SC9h7HZBI6aV9FaEYhGnaquRT5q8\n7p8lzCphKNXVyeL6T/pDJOW70zXItkl8Ryoc0tIaknRQmj8+YA0Hr9GDdmYev2yrxSoVS7s5Bl8p\noasn3DljgnWT07vsQz+hw3NY4SPp7IFGP2PpGUBBIIvrOaDWpPGsXeznBxSFtis6Qo+JiEoaVql9\nb9/XyKZj65wOsVyZhFWeM1nCQITSP9OqOc9FSoDFYQ1AVogm4A2AzUrkMnT1SrN2dCuTmNbeVw7g\nOMqMrVf0CiTv9hI0cATbO5we1sPAlJxscSkJjsaI+sQfjiAnAgMBAAEwDQYJKoZIhvcNAQELBQAD\nggEBACgwoH1qklPF1nI9+WbIJ4K12Dl9+U3ZMZa2lP4hAk1rMBHk9SHboOU1CHDQKT1Z6uxi0NI4\nJZHmP1qP8KPNEWTI8Q76ue4Q3aiA53EQguzGb3SEtyp36JGBq05Jor9erEebFftVl83NFvio72Fn\n0N2xvu8zCnlylf2hpz9x1i01Xnz5UNtZ2ppsf2zzT+4U6w3frH+pkp0RDPuoe9mnBF001AguP31h\nSBZyZzWcwQltuNELnSRCcgJl4kC2h3mAgaVtYalrFxLRa3tA2XF2BHRHmKgocedVhTq+81xrqj+W\nQuDmUe06DnrS3Ohmyj3jhsCCluznAolmrBhT/SaDuGg=\n-----END CERTIFICATE-----\n"
            id        = "4qrossmq1vxa4p836kyqzp48h"
          }
        },
      ]
      digital_signature = false
      http_basic_credentials = {
        password = "2FederateM0re"
        username = "uname2012"
      }
      require_ssl = true
    }
    signing_settings = {
      signing_key_pair_ref = {
        id = "419x9yg43rlawqwq9v6az997k"
      }
      include_raw_key_in_signature = false
      include_cert_in_signature    = false
      algorithm                    = "SHA256withRSA"
    }
  }
  default_virtual_entity_id = "virtual_server_id_1"
  entity_id                 = "inbound_AD2012"
  extended_properties = {
    authNexp = {
      values = ["val1"]
    }
    useAuthnApi = {
      values = ["val2"]
    }
  }
  inbound_provisioning = {
    action_on_delete = "PERMANENTLY_DELETE_USER"
    custom_schema = {
      attributes = [
        {
          multi_valued = false
          name         = "customAttribute"
          sub_attributes = [
            "subAttribute1",
            "subAttribute2"
          ]
          types = []
        }
      ]
      namespace = "urn:scim:schemas:extension:custom:1.0"
    }
    group_support = true
    groups = {
      read_groups = {
        attribute_contract = {
          extended_attributes = []
        }
        attribute_fulfillment = {
          displayName = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "displayName"
          }
        }
        attributes = [
          {
            name = "displayName"
          },
        ]
      }
      write_groups = {
        attribute_fulfillment = {
          displayName = {
            source = {
              id   = null
              type = "SCIM_GROUP"
            }
            value = "displayName"
          }
        }
      }
    }
    user_repository = {
      ldap = {
        base_dn = "OU=SQE Testing,OU=Resources,DC=ldaps2,DC=com"
        data_store_ref = {
          id = pingfederate_data_store.example.id
        }
        unique_group_id_filter = "ou=groupid"
        unique_user_id_filter  = "CN=$${userName}"
      }
    }
    users = {
      read_users = {
        attribute_contract = {
          extended_attributes = [
          ]
        }
        attribute_fulfillment = {
          "addresses.home.formatted" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "homePostalAddress"
          }
          "addresses.other.formatted" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "postalAddress"
          }
          "addresses.work.country" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "c"
          }
          "addresses.work.locality" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "l"
          }
          "addresses.work.postalCode" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "postalCode"
          }
          "addresses.work.region" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "st"
          }
          "addresses.work.streetAddress" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "streetAddress"
          }
          displayName = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "displayName"
          }
          "emails.home.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "mail"
          }
          "emails.work.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "mail"
          }
          externalId = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "uid"
          }
          "name.familyName" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "sn"
          }
          "name.givenName" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "givenName"
          }
          "name.middleName" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "initials"
          }
          "phoneNumbers.fax.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "facsimileTelephoneNumber"
          }
          "phoneNumbers.home.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "homePhone"
          }
          "phoneNumbers.mobile.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "mobile"
          }
          "phoneNumbers.other.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "otherTelephone"
          }
          "phoneNumbers.pager.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "pager"
          }
          "phoneNumbers.work.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "telephoneNumber"
          }
          "photos.photo.value" = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "url"
          }
          preferredLanguage = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "preferredLanguage"
          }
          profileUrl = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "wWWHomePage"
          }
          title = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "title"
          }
          userName = {
            source = {
              id   = null
              type = "LDAP_DATA_STORE"
            }
            value = "sAMAccountName"
          }
        }
        attributes = [
          {
            name = "c"
          },
          {
            name = "displayName"
          },
          {
            name = "facsimileTelephoneNumber"
          },
          {
            name = "givenName"
          },
          {
            name = "homePhone"
          },
          {
            name = "homePostalAddress"
          },
          {
            name = "initials"
          },
          {
            name = "l"
          },
          {
            name = "mail"
          },
          {
            name = "mobile"
          },
          {
            name = "otherTelephone"
          },
          {
            name = "pager"
          },
          {
            name = "postalAddress"
          },
          {
            name = "postalCode"
          },
          {
            name = "preferredLanguage"
          },
          {
            name = "sAMAccountName"
          },
          {
            name = "sn"
          },
          {
            name = "st"
          },
          {
            name = "streetAddress"
          },
          {
            name = "telephoneNumber"
          },
          {
            name = "title"
          },
          {
            name = "uid"
          },
          {
            name = "url"
          },
          {
            name = "wWWHomePage"
          },
        ]
      }
      write_users = {
        attribute_fulfillment = {
          c = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.work.country"
          }
          displayName = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "displayName"
          }
          facsimileTelephoneNumber = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "phoneNumbers.fax.value"
          }
          givenName = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "name.givenName"
          }
          homePhone = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "phoneNumbers.home.value"
          }
          homePostalAddress = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.home.formatted"
          }
          initials = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "name.middleName"
          }
          l = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.work.locality"
          }
          mail = {
            source = {
              id   = null
              type = "EXPRESSION"
            }
            value = "emails.work.value"
          }
          mobile = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "phoneNumbers.mobile.value"
          }
          otherTelephone = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "phoneNumbers.other.value"
          }
          pager = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "phoneNumbers.pager.value"
          }
          postalAddress = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.other.formatted"
          }
          postalCode = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.work.postalCode"
          }
          preferredLanguage = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "preferredLanguage"
          }
          sAMAccountName = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "userName"
          }
          sn = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "name.familyName"
          }
          st = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.work.region"
          }
          streetAddress = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "addresses.work.streetAddress"
          }
          telephoneNumber = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "phoneNumbers.work.value"
          }
          title = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "title"
          }
          uid = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "externalId"
          }
          url = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "photos.photo.value"
          }
          wWWHomePage = {
            source = {
              id   = null
              type = "SCIM_USER"
            }
            value = "profileUrl"
          }
        }
      }
    }
  }
  logging_mode       = "STANDARD"
  name               = "inbound_AD2012"
  virtual_entity_ids = ["virtual_server_id_1", "virtual_server_id_2", "virtual_server_id_3"]
}
`, spIdpConnection_InboundProvisioningDependencyHCL(), idpConnInboundProvisioningId)
}

// Validate any computed values when applying minimal HCL
func spIdpConnection_CheckComputedValuesInboundProvisioningMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "active", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.digital_signature", "false"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.http_basic_credentials.encrypted_password"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.require_ssl", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "id", idpConnInboundProvisioningId),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.custom_schema.namespace", "urn:scim:schemas:extension:custom:1.0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.core_attributes.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.core_attributes.0.name", "addresses.work.country"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.extended_attributes.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "logging_mode", "STANDARD"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "virtual_entity_ids.#", "0"),
	)
}

// Validate any computed values when applying complete HCL
func spIdpConnection_CheckComputedValuesInboundProvisioningComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.expires", "2024-07-13T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.id", "4qrossmq1vxa4p836kyqzp48h"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.issuer_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.serial_number", "11775821034523537675"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.sha1_fingerprint", "3CFE421ED628F7CEFE08B02DEB3EB4FB5DE9B92D"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.sha256_fingerprint", "633FF42A14E808AEEE5810D78F2C68358AD27787CDDADA302A7E201BA7F2A046"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.status", "EXPIRED"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.subject_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.valid_from", "2023-07-14T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.certs.0.cert_view.version", "1"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "credentials.certs.0.x509_file.formatted_file_data"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.expires", "2024-07-13T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.id", "4qrossmq1vxa4p836kyqzp48h"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.issuer_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.serial_number", "11775821034523537675"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.sha1_fingerprint", "3CFE421ED628F7CEFE08B02DEB3EB4FB5DE9B92D"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.sha256_fingerprint", "633FF42A14E808AEEE5810D78F2C68358AD27787CDDADA302A7E201BA7F2A046"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.status", "EXPIRED"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.subject_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS, CN=localhost"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.valid_from", "2023-07-14T02:54:53Z"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.cert_view.version", "1"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.certs.0.x509_file.formatted_file_data"),
		resource.TestCheckResourceAttrSet("pingfederate_sp_idp_connection.example", "credentials.inbound_back_channel_auth.http_basic_credentials.encrypted_password"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.groups.read_groups.attribute_contract.core_attributes.#", "1"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.groups.read_groups.attribute_contract.core_attributes.0.name", "displayName"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.groups.read_groups.attribute_contract.core_attributes.0.masked", "false"),
		resource.TestCheckResourceAttr("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.core_attributes.#", "23"),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.core_attributes.*",
			map[string]string{
				"name":   "userName",
				"masked": "false",
			},
		), resource.TestCheckTypeSetElemNestedAttrs("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.core_attributes.*",
			map[string]string{
				"name":   "phoneNumbers.fax.value",
				"masked": "false",
			},
		), resource.TestCheckTypeSetElemNestedAttrs("pingfederate_sp_idp_connection.example", "inbound_provisioning.users.read_users.attribute_contract.core_attributes.*",
			map[string]string{
				"name":   "addresses.home.formatted",
				"masked": "false",
			},
		),
	)
}

// Test that any objects created by the test are destroyed
func spIdpConnection_InboundProvisioningCheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.SpIdpConnectionsAPI.DeleteConnection(acctest.TestBasicAuthContext(), idpConnInboundProvisioningId).Execute()
	if err == nil {
		return fmt.Errorf("sp_idp_connection still exists after tests. Expected it to be destroyed")
	}
	return nil
}
