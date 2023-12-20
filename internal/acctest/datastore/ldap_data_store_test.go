package acctest_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

// These variables cannot be modified due to resource dependent values
const ldapDataStoreId = "ldapDataStoreId"
const dataStoreType = "LDAP"
const ldapType = "PING_DIRECTORY"
const verifyHost = false
const hostnames = "pingdirectory:1389"
const userDn = "cn=userDN"
const passwordVal = "password"

type ldapDataStoreResourceModel struct {
	dataStore *client.LdapDataStore
}

func initialLdapDataStore() *client.LdapDataStore {
	initialLdapDataStore := client.NewLdapDataStoreWithDefaults()
	initialLdapDataStore.Id = pointers.String(ldapDataStoreId)
	initialLdapDataStore.Type = dataStoreType
	initialLdapDataStore.LdapType = ldapType
	initialLdapDataStore.BindAnonymously = pointers.Bool(false)
	initialLdapDataStore.UserDN = pointers.String(userDn)
	initialLdapDataStore.Password = pointers.String(passwordVal)
	initialLdapDataStore.Hostnames = []string{hostnames}
	initialLdapDataStore.VerifyHost = pointers.Bool(verifyHost)
	return initialLdapDataStore
}

func updatedLdapDataStore() *client.LdapDataStore {
	updatedLdapDataStore := client.NewLdapDataStoreWithDefaults()
	updatedLdapDataStore.Id = pointers.String(ldapDataStoreId)
	updatedLdapDataStore.Name = pointers.String("updatedLdapDataStoreName")
	updatedLdapDataStore.Type = dataStoreType
	updatedLdapDataStore.LdapType = ldapType
	updatedLdapDataStore.BindAnonymously = pointers.Bool(true)
	updatedLdapDataStore.UserDN = pointers.String(userDn)
	updatedLdapDataStore.Password = pointers.String(passwordVal)
	updatedLdapDataStore.Hostnames = []string{hostnames}
	updatedLdapDataStore.VerifyHost = pointers.Bool(verifyHost)
	updatedLdapDataStore.MaskAttributeValues = pointers.Bool(true)
	updatedLdapDataStore.UseSsl = pointers.Bool(true)
	updatedLdapDataStore.UseDnsSrvRecords = pointers.Bool(true)
	updatedLdapDataStore.TestOnBorrow = pointers.Bool(true)
	updatedLdapDataStore.TestOnReturn = pointers.Bool(true)
	updatedLdapDataStore.CreateIfNecessary = pointers.Bool(true)
	updatedLdapDataStore.MinConnections = pointers.Int64(1)
	updatedLdapDataStore.MaxConnections = pointers.Int64(200)
	updatedLdapDataStore.MaxWait = pointers.Int64(1000)
	updatedLdapDataStore.TimeBetweenEvictions = pointers.Int64(3000)
	updatedLdapDataStore.ReadTimeout = pointers.Int64(600)
	updatedLdapDataStore.ConnectionTimeout = pointers.Int64(600)
	updatedLdapDataStore.BinaryAttributes = []string{"updatedBinaryAttribute1", "updatedBinaryAttribute2"}
	updatedLdapDataStore.DnsTtl = pointers.Int64(3000)
	updatedLdapDataStore.LdapDnsSrvPrefix = pointers.String("_ldap._tcp")
	updatedLdapDataStore.LdapsDnsSrvPrefix = pointers.String("_ldaps._tcp")
	return updatedLdapDataStore
}

func TestAccLdapDataStore(t *testing.T) {
	resourceName := "myLdapDataStore"
	initialResourceModel := ldapDataStoreResourceModel{
		dataStore: initialLdapDataStore(),
	}

	updatedResourceModel := ldapDataStoreResourceModel{
		dataStore: updatedLdapDataStore(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckLdapDataStoreDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal model
				Config: testAccLdapDataStore(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLdapDataStoreAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccLdapDataStore(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLdapDataStoreAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccLdapDataStore(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_data_store." + resourceName,
				ImportStateId:           ldapDataStoreId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ldap_data_store.user_dn", "ldap_data_store.password"},
			},
			{
				// Back to the initial minimal model
				Config: testAccLdapDataStore(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLdapDataStoreAttributes(initialResourceModel),
			},
		},
	})
}

func hcl(lds *client.LdapDataStore) string {
	var builder strings.Builder
	if lds == nil {
		return ""
	}
	if lds != nil {
		versionedHcl := ""
		if acctest.VersionAtLeast(version.PingFederate1130) {
			versionedHcl += `
			retry_failed_operations = true
			`
		}

		top := `
		data_store_id             = "%[1]s"
		%[2]s
		`
		builder.WriteString(
			fmt.Sprintf(top,
				*lds.Id,
				acctest.TfKeyValuePairToString("mask_attribute_values", strconv.FormatBool(lds.GetMaskAttributeValues()), false),
			))
		tf := `
		ldap_data_store = {
		  %[1]s
		  %[2]s
		  %[3]s
		  %[4]s
		  %[5]s
		  %[6]s
		  %[7]s
		  %[8]s
		  %[9]s
		  %[10]s
		  %[11]s
		  %[12]s
		  %[13]s
		  %[14]s
		  %[15]s
		  %[16]s
		  %[17]s
		  %[18]s
		  %[19]s
		  %[20]s
			%[21]s
			%[22]s
			%[23]s
		}
		`
		hostnames := func() string {
			if len(lds.GetHostnames()) > 0 {
				return acctest.StringSliceToTerraformString(lds.GetHostnames())
			} else {
				return ""
			}
		}
		binaryAttributes := func() string {
			if len(lds.GetBinaryAttributes()) > 0 {
				return acctest.StringSliceToTerraformString(lds.GetBinaryAttributes())
			} else {
				return ""
			}
		}
		builder.WriteString(fmt.Sprintf(tf,
			acctest.TfKeyValuePairToString("ldap_type", lds.LdapType, true),
			acctest.TfKeyValuePairToString("bind_anonymously", strconv.FormatBool(lds.GetBindAnonymously()), false),
			acctest.TfKeyValuePairToString("user_dn", *lds.UserDN, true),
			acctest.TfKeyValuePairToString("password", *lds.Password, true),
			acctest.TfKeyValuePairToString("use_ssl", strconv.FormatBool(lds.GetUseSsl()), false),
			acctest.TfKeyValuePairToString("use_dns_srv_records", strconv.FormatBool(lds.GetUseDnsSrvRecords()), false),
			acctest.TfKeyValuePairToString("name", lds.GetName(), true),
			acctest.TfKeyValuePairToString("hostnames", hostnames(), false),
			acctest.TfKeyValuePairToString("test_on_borrow", strconv.FormatBool(lds.GetTestOnBorrow()), false),
			acctest.TfKeyValuePairToString("test_on_return", strconv.FormatBool(lds.GetTestOnReturn()), false),
			acctest.TfKeyValuePairToString("create_if_necessary", strconv.FormatBool(lds.GetCreateIfNecessary()), false),
			acctest.TfKeyValuePairToString("verify_host", strconv.FormatBool(lds.GetVerifyHost()), false),
			acctest.TfKeyValuePairToString("min_connections", strconv.FormatInt(lds.GetMinConnections(), 10), false),
			acctest.TfKeyValuePairToString("max_connections", strconv.FormatInt(lds.GetMaxConnections(), 10), false),
			acctest.TfKeyValuePairToString("max_wait", strconv.FormatInt(lds.GetMaxWait(), 10), false),
			acctest.TfKeyValuePairToString("time_between_evictions", strconv.FormatInt(lds.GetTimeBetweenEvictions(), 10), false),
			acctest.TfKeyValuePairToString("read_timeout", strconv.FormatInt(lds.GetReadTimeout(), 10), false),
			acctest.TfKeyValuePairToString("connection_timeout", strconv.FormatInt(lds.GetConnectionTimeout(), 10), false),
			acctest.TfKeyValuePairToString("binary_attributes", binaryAttributes(), false),
			acctest.TfKeyValuePairToString("dns_ttl", strconv.FormatInt(lds.GetDnsTtl(), 10), false),
			acctest.TfKeyValuePairToString("ldap_dns_srv_prefix", lds.GetLdapDnsSrvPrefix(), true),
			acctest.TfKeyValuePairToString("ldaps_dns_srv_prefix", lds.GetLdapsDnsSrvPrefix(), true),
			versionedHcl),
		)
	}
	return builder.String()
}

func testAccLdapDataStore(resourceName string, ldapDataStore ldapDataStoreResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_data_store" "%[1]s" {
	%[2]s
}
data "pingfederate_data_store" "%[1]s" {
  data_store_id = pingfederate_data_store.%[1]s.id
}`, resourceName,
		hcl(ldapDataStore.dataStore),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedLdapDataStoreAttributes(config ldapDataStoreResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "LdapDataStore"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		resp, _, err := testClient.DataStoresAPI.GetDataStore(ctx, ldapDataStoreId).Execute()

		if err != nil {
			return err
		}

		nameValue := func() string {
			if config.dataStore.Name != nil {
				return *config.dataStore.Name
			} else {
				return config.dataStore.Hostnames[0] + " (" + *config.dataStore.UserDN + ")"
			}
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(ldapDataStoreId), "name", nameValue(), *resp.LdapDataStore.Name)
		if err != nil {
			return err
		}

		if config.dataStore.MaskAttributeValues != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "mask_attribute_values", *config.dataStore.MaskAttributeValues, *resp.LdapDataStore.MaskAttributeValues)
			if err != nil {
				return err
			}
		}

		if config.dataStore.BindAnonymously != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "bind_anonymously", *config.dataStore.BindAnonymously, *resp.LdapDataStore.BindAnonymously)
			if err != nil {
				return err
			}
		}

		if config.dataStore.UseSsl != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "use_ssl", *config.dataStore.UseSsl, *resp.LdapDataStore.UseSsl)
			if err != nil {
				return err
			}
		}

		if config.dataStore.UseDnsSrvRecords != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "use_dns_srv_records", *config.dataStore.UseDnsSrvRecords, *resp.LdapDataStore.UseDnsSrvRecords)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(ldapDataStoreId), "hostnames", config.dataStore.Hostnames, resp.LdapDataStore.Hostnames)
		if err != nil {
			return err
		}

		if config.dataStore.TestOnBorrow != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "test_on_borrow", *config.dataStore.TestOnBorrow, *resp.LdapDataStore.TestOnBorrow)
			if err != nil {
				return err
			}
		}

		if config.dataStore.TestOnReturn != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "test_on_return", *config.dataStore.TestOnReturn, *resp.LdapDataStore.TestOnReturn)
			if err != nil {
				return err
			}
		}

		if config.dataStore.CreateIfNecessary != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "create_if_necessary", *config.dataStore.CreateIfNecessary, *resp.LdapDataStore.CreateIfNecessary)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(ldapDataStoreId), "verify_host", *config.dataStore.VerifyHost, *resp.LdapDataStore.VerifyHost)
		if err != nil {
			return err
		}

		if config.dataStore.MinConnections != nil {
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(ldapDataStoreId), "min_connections", *config.dataStore.MinConnections, *resp.LdapDataStore.MinConnections)
			if err != nil {
				return err
			}
		}

		if config.dataStore.MaxConnections != nil {
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(ldapDataStoreId), "max_connections", *config.dataStore.MaxConnections, *resp.LdapDataStore.MaxConnections)
			if err != nil {
				return err
			}
		}

		if config.dataStore.MaxWait != nil {
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(ldapDataStoreId), "max_wait", *config.dataStore.MaxWait, *resp.LdapDataStore.MaxWait)
			if err != nil {
				return err
			}
		}

		if config.dataStore.TimeBetweenEvictions != nil {
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(ldapDataStoreId), "time_between_evictions", *config.dataStore.TimeBetweenEvictions, *resp.LdapDataStore.TimeBetweenEvictions)
			if err != nil {
				return err
			}
		}

		if config.dataStore.ReadTimeout != nil {
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(ldapDataStoreId), "read_timeout", *config.dataStore.ReadTimeout, *resp.LdapDataStore.ReadTimeout)
			if err != nil {
				return err
			}
		}

		if config.dataStore.ConnectionTimeout != nil {
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(ldapDataStoreId), "connection_timeout", *config.dataStore.ConnectionTimeout, *resp.LdapDataStore.ConnectionTimeout)
			if err != nil {
				return err
			}
		}

		if config.dataStore.BinaryAttributes != nil {
			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(ldapDataStoreId), "binary_attributes", config.dataStore.BinaryAttributes, resp.LdapDataStore.BinaryAttributes)
			if err != nil {
				return err
			}
		}

		if config.dataStore.DnsTtl != nil {
			err = acctest.TestAttributesMatchInt(resourceType, pointers.String(ldapDataStoreId), "dns_ttl", *config.dataStore.DnsTtl, *resp.LdapDataStore.DnsTtl)
			if err != nil {
				return err
			}
		}

		if config.dataStore.LdapDnsSrvPrefix != nil {
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(ldapDataStoreId), "ldap_dns_srv_prefix", *config.dataStore.LdapDnsSrvPrefix, *resp.LdapDataStore.LdapDnsSrvPrefix)
			if err != nil {
				return err
			}
		}

		if config.dataStore.LdapsDnsSrvPrefix != nil {
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(ldapDataStoreId), "ldaps_dns_srv_prefix", *config.dataStore.LdapsDnsSrvPrefix, *resp.LdapDataStore.LdapsDnsSrvPrefix)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLdapDataStoreDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.DataStoresAPI.DeleteDataStore(ctx, customDataStoreId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("LdapDataStore", ldapDataStoreId)
	}
	return nil
}
