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
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// These variables cannot be modified due to resource dependent values
const pingOneLdapGatewayDataStoreId = "pingOneLdapGatewayDataStoreId"
const pingOneLdapGDSType = "PING_ONE_LDAP_GATEWAY"
const ldapTypeVal = "PING_DIRECTORY"
const pingOneConnectionRefId = "noeOvj5ltBnf4rcmtZAKdJ"
const pingOneEnvironmentId = "f5901536-2b60-4d4a-a987-3d56aadad46d"
const pingOneLdapGatewayId = "3b7b5d9d-1820-4b21-bb29-a5336af65352"

type pingOneLdapGatewayDataStoreResourceModel struct {
	dataStore *client.PingOneLdapGatewayDataStore
}

func initialPingOneLdapGatewayDataStore() *client.PingOneLdapGatewayDataStore {
	initialPingOneLdapGatewayDataStore := client.NewPingOneLdapGatewayDataStoreWithDefaults()
	initialPingOneLdapGatewayDataStore.Id = pointers.String(pingOneLdapGatewayDataStoreId)
	initialPingOneLdapGatewayDataStore.Name = pointers.String("initialPingOneLdapGatewayDataStore")
	initialPingOneLdapGatewayDataStore.LdapType = ldapTypeVal
	initialPingOneLdapGatewayDataStore.Type = pingOneLdapGDSType
	initialPingOneLdapGatewayDataStore.PingOneConnectionRef = *client.NewResourceLink(pingOneConnectionRefId)
	initialPingOneLdapGatewayDataStore.PingOneEnvironmentId = pingOneEnvironmentId
	initialPingOneLdapGatewayDataStore.PingOneLdapGatewayId = pingOneLdapGatewayId
	return initialPingOneLdapGatewayDataStore
}

func updatedPingOneLdapGatewayDataStore() *client.PingOneLdapGatewayDataStore {
	updatedPingOneLdapGatewayDataStore := client.NewPingOneLdapGatewayDataStoreWithDefaults()
	updatedPingOneLdapGatewayDataStore.Id = pointers.String(pingOneLdapGatewayDataStoreId)
	updatedPingOneLdapGatewayDataStore.Name = pointers.String("updatedPingOneLdapGatewayDataStore")
	updatedPingOneLdapGatewayDataStore.LdapType = ldapTypeVal
	updatedPingOneLdapGatewayDataStore.Type = pingOneLdapGDSType
	updatedPingOneLdapGatewayDataStore.PingOneConnectionRef = *client.NewResourceLink(pingOneConnectionRefId)
	updatedPingOneLdapGatewayDataStore.PingOneEnvironmentId = pingOneEnvironmentId
	updatedPingOneLdapGatewayDataStore.PingOneLdapGatewayId = pingOneLdapGatewayId
	updatedPingOneLdapGatewayDataStore.UseSsl = pointers.Bool(true)
	updatedPingOneLdapGatewayDataStore.Name = pointers.String("myPingOneLdapGatewayDataStore")
	updatedPingOneLdapGatewayDataStore.MaskAttributeValues = pointers.Bool(true)
	updatedPingOneLdapGatewayDataStore.BinaryAttributes = []string{"binaryAttribute1", "binaryAttribute2"}
	return updatedPingOneLdapGatewayDataStore
}

func TestAccPingOneLdapGatewayDataStore(t *testing.T) {
	resourceName := "myPingOneLdapGatewayDataStore"
	initialResourceModel := pingOneLdapGatewayDataStoreResourceModel{
		dataStore: initialPingOneLdapGatewayDataStore(),
	}

	updatedResourceModel := pingOneLdapGatewayDataStoreResourceModel{
		dataStore: updatedPingOneLdapGatewayDataStore(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckPingOneLdapGatewayDataStoreDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal model
				Config: testAccPingOneLdapGatewayDataStore(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPingOneLdapGatewayDataStoreAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPingOneLdapGatewayDataStore(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPingOneLdapGatewayDataStoreAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPingOneLdapGatewayDataStore(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_data_store." + resourceName,
				ImportStateId:     pingOneLdapGatewayDataStoreId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to the initial minimal model
				Config: testAccPingOneLdapGatewayDataStore(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPingOneLdapGatewayDataStoreAttributes(initialResourceModel),
			},
		},
	})
}

func pingOneLdapGDShcl(pingOneLdapGDS *client.PingOneLdapGatewayDataStore) string {
	var builder strings.Builder
	if pingOneLdapGDS == nil {
		return ""
	}
	if pingOneLdapGDS != nil {
		top := `
		custom_id             = "%[1]s"
		%[2]s
		`
		builder.WriteString(
			fmt.Sprintf(top,
				*pingOneLdapGDS.Id,
				acctest.TfKeyValuePairToString("mask_attribute_values", strconv.FormatBool(pingOneLdapGDS.GetMaskAttributeValues()), false),
			))
		tf := `
		ping_one_ldap_gateway_data_store = {
		  %[1]s
			%[2]s
			ping_one_connection_ref = {
				id = "%[3]s"
			}
			%[4]s
			%[5]s
			%[6]s
		}
		`
		builder.WriteString(fmt.Sprintf(tf,
			acctest.TfKeyValuePairToString("ldap_type", pingOneLdapGDS.LdapType, true),
			acctest.TfKeyValuePairToString("name", *pingOneLdapGDS.Name, true),
			pingOneLdapGDS.PingOneConnectionRef.Id,
			acctest.TfKeyValuePairToString("ping_one_environment_id", pingOneLdapGDS.PingOneEnvironmentId, true),
			acctest.TfKeyValuePairToString("ping_one_ldap_gateway_id", pingOneLdapGDS.PingOneLdapGatewayId, true),
			acctest.TfKeyValuePairToString("use_ssl", strconv.FormatBool(pingOneLdapGDS.GetUseSsl()), false),
		))
	}
	return builder.String()
}

func testAccPingOneLdapGatewayDataStore(resourceName string, pingOneLdapGatewayDataStore pingOneLdapGatewayDataStoreResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_data_store" "%[1]s" {
	%[2]s
}`, resourceName,
		pingOneLdapGDShcl(pingOneLdapGatewayDataStore.dataStore),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedPingOneLdapGatewayDataStoreAttributes(config pingOneLdapGatewayDataStoreResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "PingOneLdapGatewayDataStore"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		resp, _, err := testClient.DataStoresAPI.GetDataStore(ctx, pingOneLdapGatewayDataStoreId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		if resp.PingOneLdapGatewayDataStore.Name != nil {
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(pingOneLdapGatewayDataStoreId), "name", *config.dataStore.Name, *resp.PingOneLdapGatewayDataStore.Name)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchString(resourceType, pointers.String(pingOneLdapGatewayDataStoreId), "id", config.dataStore.PingOneConnectionRef.Id, resp.PingOneLdapGatewayDataStore.PingOneConnectionRef.Id)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(pingOneLdapGatewayDataStoreId), "mask_attribute_values", config.dataStore.GetMaskAttributeValues(), *resp.PingOneLdapGatewayDataStore.MaskAttributeValues)
		if err != nil {
			return err
		}

		if resp.PingOneLdapGatewayDataStore.UseSsl != nil {
			err = acctest.TestAttributesMatchBool(resourceType, pointers.String(pingOneLdapGatewayDataStoreId), "use_ssl", config.dataStore.GetUseSsl(), *resp.PingOneLdapGatewayDataStore.UseSsl)
			if err != nil {
				return err
			}
		}

		if resp.PingOneLdapGatewayDataStore.BinaryAttributes != nil {
			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(pingOneLdapGatewayDataStoreId), "binary_attributes", config.dataStore.GetBinaryAttributes(), resp.PingOneLdapGatewayDataStore.BinaryAttributes)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPingOneLdapGatewayDataStoreDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.DataStoresAPI.DeleteDataStore(ctx, customDataStoreId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("PingOneLdapGatewayDataStore", pingOneLdapGatewayDataStoreId)
	}
	return nil
}
