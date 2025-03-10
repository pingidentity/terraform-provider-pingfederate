// Copyright © 2025 Ping Identity Corporation

package datastore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// These variables cannot be modified due to resource dependent values
const jdbcDataStoreId = "jdbcDataStoreId"
const driverClass = "org.hsqldb.jdbcDriver"
const userName = "sa"
const password = "secretpass"
const jdbcDataStoreType = "JDBC"
const connectionUrl = "jdbc:hsqldb:$${pf.server.data.dir}$${/}hypersonic$${/}ProvisionerDefaultDB;hsqldb.lock_file=false"
const connectionUrlUnescaped = "jdbc:hsqldb:${pf.server.data.dir}${/}hypersonic${/}ProvisionerDefaultDB;hsqldb.lock_file=false"

// Attributes to test with. Add optional properties to test here if desired.
type jdbcDataStoreResourceModel struct {
	maskAttributeValues bool
	jdbcDataStore       client.JdbcDataStore
	addTags             bool
}

func initialJdbcDataStore() *client.JdbcDataStore {
	jdbcDataStore := client.NewJdbcDataStore(driverClass, jdbcDataStoreType)
	jdbcDataStore.UserName = pointers.String(userName)
	jdbcDataStore.Password = pointers.String(password)
	jdbcDataStore.ConnectionUrl = pointers.String(connectionUrl)
	jdbcDataStore.AllowMultiValueAttributes = pointers.Bool(false)
	jdbcDataStore.MinPoolSize = pointers.Int64(10)
	jdbcDataStore.MaxPoolSize = pointers.Int64(100)
	jdbcDataStore.BlockingTimeout = pointers.Int64(5000)
	jdbcDataStore.IdleTimeout = pointers.Int64(5)
	return jdbcDataStore
}

func updatedJdbcDataStore() *client.JdbcDataStore {
	jdbcDataStore := client.NewJdbcDataStore(driverClass, jdbcDataStoreType)
	jdbcDataStore.UserName = pointers.String(userName)
	jdbcDataStore.Password = pointers.String(password)
	jdbcDataStore.ConnectionUrl = pointers.String(connectionUrl)
	jdbcDataStore.Name = pointers.String("updatedJdbcDataStore")
	jdbcDataStore.AllowMultiValueAttributes = pointers.Bool(true)
	jdbcDataStore.MinPoolSize = pointers.Int64(20)
	jdbcDataStore.MaxPoolSize = pointers.Int64(200)
	jdbcDataStore.BlockingTimeout = pointers.Int64(6000)
	jdbcDataStore.IdleTimeout = pointers.Int64(10)
	return jdbcDataStore
}

func hclJdbcDataStore(jdbcDataStore *client.JdbcDataStore, addTags bool) string {
	var builder strings.Builder
	if jdbcDataStore == nil {
		return ""
	}
	var optionalHcl string
	if addTags {
		optionalHcl = fmt.Sprintf(`
		connection_url_tags = [
			{
				connection_url = "%s"
				default_source = true
			},
			{
				connection_url = "secondurl.com"
				tags = "us-east-1"
			},
			{
			    connection_url = "thirdurl.com"
				tags = "us-west-1"
				default_source = false
			}
		]
		`, *jdbcDataStore.ConnectionUrl)
	}
	if jdbcDataStore.Name != nil {
		optionalHcl += fmt.Sprintf(`
		name = "%s"
		`, *jdbcDataStore.Name)
	}
	if jdbcDataStore != nil {
		tf := `
		jdbc_data_store = {
			connection_url               = "%s"
			driver_class                 = "%s"
			user_name                    = "%s"
			password                     = "%s"
			allow_multi_value_attributes = %t
			min_pool_size    = %d
			max_pool_size    = %d
			blocking_timeout = %d
			idle_timeout     = %d
			%s
		}
	`
		builder.WriteString(fmt.Sprintf(tf,
			*jdbcDataStore.ConnectionUrl,
			jdbcDataStore.DriverClass,
			*jdbcDataStore.UserName,
			*jdbcDataStore.Password,
			*jdbcDataStore.AllowMultiValueAttributes,
			*jdbcDataStore.MinPoolSize,
			*jdbcDataStore.MaxPoolSize,
			*jdbcDataStore.BlockingTimeout,
			*jdbcDataStore.IdleTimeout,
			optionalHcl),
		)
	}
	return builder.String()
}

func TestAccJdbcDataStore(t *testing.T) {
	resourceName := "myJdbcDataStore"
	initialResourceModel := jdbcDataStoreResourceModel{
		maskAttributeValues: false,
		jdbcDataStore:       *initialJdbcDataStore(),
		addTags:             false,
	}

	updatedResourceModel := jdbcDataStoreResourceModel{
		maskAttributeValues: true,
		jdbcDataStore:       *updatedJdbcDataStore(),
		addTags:             true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckJdbcDataStoreDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal model
				Config: testAccJdbcDataStore(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedJdbcDataStoreAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.name", connectionUrlUnescaped+" (sa)"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.0.connection_url", connectionUrlUnescaped),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.0.default_source", "true"),
					resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.0.tags"),
				),
			},
			{
				// Test updating some fields
				Config: testAccJdbcDataStore(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedJdbcDataStoreAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccJdbcDataStore(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_data_store." + resourceName,
				ImportStateId:     jdbcDataStoreId,
				ImportState:       true,
				ImportStateVerify: true,
				// password can't be imported, and encrypted_password will change each time it is read
				ImportStateVerifyIgnore: []string{"jdbc_data_store.password", "jdbc_data_store.encrypted_password"},
			},
			{
				// Back to the initial minimal model
				Config: testAccJdbcDataStore(resourceName, initialResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedJdbcDataStoreAttributes(initialResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.name", connectionUrlUnescaped+" (sa)"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.#", "1"),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.0.connection_url", connectionUrlUnescaped),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.0.default_source", "true"),
					resource.TestCheckNoResourceAttr(fmt.Sprintf("pingfederate_data_store.%s", resourceName), "jdbc_data_store.connection_url_tags.0.tags"),
				),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.DataStoresAPI.DeleteDataStore(ctx, jdbcDataStoreId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Minimal model
				Config: testAccJdbcDataStore(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedJdbcDataStoreAttributes(initialResourceModel),
			},
		},
	})
}

func testAccJdbcDataStore(resourceName string, resourceModel jdbcDataStoreResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_data_store" "%[1]s" {
  data_store_id         = "%[2]s"
  mask_attribute_values = %[3]t
	%[4]s
}
data "pingfederate_data_store" "%[1]s" {
  data_store_id = pingfederate_data_store.%[1]s.id
}`, resourceName,
		jdbcDataStoreId,
		resourceModel.maskAttributeValues,
		hclJdbcDataStore(&resourceModel.jdbcDataStore, resourceModel.addTags),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedJdbcDataStoreAttributes(config jdbcDataStoreResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "JdbcDataStore"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		resp, _, err := testClient.DataStoresAPI.GetDataStore(ctx, jdbcDataStoreId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		if config.jdbcDataStore.Name != nil {
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(jdbcDataStoreId), "name", *config.jdbcDataStore.Name, *resp.JdbcDataStore.Name)
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(jdbcDataStoreId), "allow_multi_value_attributes", *config.jdbcDataStore.AllowMultiValueAttributes, *resp.JdbcDataStore.AllowMultiValueAttributes)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, pointers.String(jdbcDataStoreId), "min_pool_size", *config.jdbcDataStore.MinPoolSize, *resp.JdbcDataStore.MinPoolSize)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, pointers.String(jdbcDataStoreId), "max_pool_size", *config.jdbcDataStore.MaxPoolSize, *resp.JdbcDataStore.MaxPoolSize)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, pointers.String(jdbcDataStoreId), "blocking_timeout", *config.jdbcDataStore.BlockingTimeout, *resp.JdbcDataStore.BlockingTimeout)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchInt(resourceType, pointers.String(jdbcDataStoreId), "idle_timeout", *config.jdbcDataStore.IdleTimeout, *resp.JdbcDataStore.IdleTimeout)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckJdbcDataStoreDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.DataStoresAPI.DeleteDataStore(ctx, jdbcDataStoreId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("JdbcDataStore", jdbcDataStoreId)
	}
	return nil
}
