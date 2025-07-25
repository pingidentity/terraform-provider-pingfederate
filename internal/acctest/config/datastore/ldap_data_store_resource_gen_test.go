// Copyright © 2025 Ping Identity Corporation
// Code generated by ping-terraform-plugin-framework-generator

package datastore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

const ldapStoreId = "ldapDataStoreId"

func TestAccLdapDataStore_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: ldapDataStore_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: ldapDataStore_MinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					ldapDataStore_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLdapDataStore_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: ldapDataStore_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: ldapDataStore_MinimalHCL(),
				Check:  ldapDataStore_CheckComputedValuesMinimal(),
			},
			{
				// Delete the minimal model
				Config:  ldapDataStore_MinimalHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: ldapDataStore_CompleteHCL(),
				Check:  ldapDataStore_CheckComputedValuesComplete(),
			},
			{
				// Back to minimal model
				Config: ldapDataStore_MinimalHCL(),
				Check:  ldapDataStore_CheckComputedValuesMinimal(),
			},
			{
				// Back to complete model
				Config: ldapDataStore_CompleteHCL(),
				Check:  ldapDataStore_CheckComputedValuesComplete(),
			},
			{
				// Test importing the resource
				Config:                               ldapDataStore_CompleteHCL(),
				ResourceName:                         "pingfederate_data_store.example",
				ImportStateId:                        ldapStoreId,
				ImportStateVerifyIdentifierAttribute: "data_store_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// password can't be imported, and encrypted_password will change each time it is read
				ImportStateVerifyIgnore: []string{"ldap_data_store.password", "ldap_data_store.encrypted_password"},
			},
		},
	})
}

// Minimal HCL with only required values set
func ldapDataStore_MinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_data_store" "example" {
  data_store_id = "%s"
  ldap_data_store = {
    ldap_type = "PING_DIRECTORY"
    user_dn   = "cn=admin"
    password  = "mypassword"
    hostnames = [
      "pingdirectory:636"
    ]
  }
}
data "pingfederate_data_store" "example" {
  data_store_id = pingfederate_data_store.example.id
}
`, ldapStoreId)
}

// Maximal HCL with all values set where possible
func ldapDataStore_CompleteHCL() string {
	var versionedHcl string
	if acctest.VersionAtLeast(version.PingFederate1210) {
		versionedHcl += `
			use_start_tls = false
			`
	}
	return fmt.Sprintf(`
resource "pingfederate_data_store" "example" {
  data_store_id         = "%s"
  mask_attribute_values = true
  ldap_data_store = {
    ldap_type             = "PING_DIRECTORY"
    user_dn               = "cn=admintwo"
    password              = "editedpassword"
    binary_attributes     = ["updatedBinaryAttribute1", "updatedBinaryAttribute2"]
    bind_anonymously      = false
    connection_timeout    = 100
    create_if_necessary   = true
    dns_ttl               = 3000
    follow_ldap_referrals = false
    hostnames = [
      "pingdirectory.example.com",
      "pingdirectory2.example.com"
    ]
    hostnames_tags = [
      {
        hostnames = [
          "pingdirectory.example.com",
          "pingdirectory2.example.com"
        ]
        default_source = true
      },
      {
        hostnames = [
          "pdeast1:1234"
        ]
        default_source = false
        tags           = "us-east-1"
      },
      {
        hostnames = [
          "pdeast2:5678"
        ]
        tags = "us-east-2"
      }
    ]
    ldap_dns_srv_prefix     = "_ldapcustom._tcp"
    ldaps_dns_srv_prefix    = "_ldapscustom._tcp"
    max_connections         = 200
    max_wait                = 500
    min_connections         = 15
    name                    = "mypddatastore"
    read_timeout            = 100
    test_on_borrow          = true
    test_on_return          = true
    time_between_evictions  = 100
    use_dns_srv_records     = false
    use_ssl                 = true
    verify_host             = false
    retry_failed_operations = true
	%s
  }
}
data "pingfederate_data_store" "example" {
  data_store_id = pingfederate_data_store.example.id
}
`, ldapStoreId, versionedHcl)
}

func checkLdapPf121ComputedAttrs() resource.TestCheckFunc {
	if acctest.VersionAtLeast(version.PingFederate1210) {
		return resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.use_start_tls", "false")
	}
	return resource.TestCheckNoResourceAttr("pingfederate_data_store.example", "ldap_data_store.use_start_tls")
}

// Validate any computed values when applying minimal HCL
func ldapDataStore_CheckComputedValuesMinimal() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_data_store.example", "ldap_data_store.binary_attributes"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.bind_anonymously", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.connection_timeout", "0"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.create_if_necessary", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.dns_ttl", "0"),
		resource.TestCheckResourceAttrSet("pingfederate_data_store.example", "ldap_data_store.encrypted_password"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.follow_ldap_referrals", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.hostnames_tags.0.default_source", "true"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.hostnames_tags.0.hostnames.0", "pingdirectory:636"),
		resource.TestCheckNoResourceAttr("pingfederate_data_store.example", "ldap_data_store.hostnames_tags.0.tags"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.ldap_dns_srv_prefix", "_ldap._tcp"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.ldaps_dns_srv_prefix", "_ldaps._tcp"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.max_connections", "100"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.max_wait", "-1"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.min_connections", "10"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.name", "pingdirectory:636 (cn=admin)"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.read_timeout", "0"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.retry_failed_operations", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.test_on_borrow", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.test_on_return", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.time_between_evictions", "0"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.use_dns_srv_records", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.use_ssl", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.verify_host", "true"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "id", ldapStoreId),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "mask_attribute_values", "false"),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.type", "LDAP"),
		checkLdapPf121ComputedAttrs(),
	)
}

// Validate any computed values when applying complete HCL
func ldapDataStore_CheckComputedValuesComplete() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("pingfederate_data_store.example", "ldap_data_store.encrypted_password"),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_data_store.example", "ldap_data_store.hostnames_tags.*",
			map[string]string{
				"default_source": "false",
				"tags":           "us-east-1",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs("pingfederate_data_store.example", "ldap_data_store.hostnames_tags.*",
			map[string]string{
				"default_source": "false",
				"tags":           "us-east-2",
			},
		),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "id", ldapStoreId),
		resource.TestCheckResourceAttr("pingfederate_data_store.example", "ldap_data_store.type", "LDAP"),
		checkLdapPf121ComputedAttrs(),
	)
}

// Delete the resource
func ldapDataStore_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.DataStoresAPI.DeleteDataStore(acctest.TestBasicAuthContext(), ldapStoreId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func ldapDataStore_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.DataStoresAPI.DeleteDataStore(acctest.TestBasicAuthContext(), ldapStoreId).Execute()
	if err == nil {
		return fmt.Errorf("data_store still exists after tests. Expected it to be destroyed")
	}
	return nil
}
