package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// These variables cannot be modified due to resource dependent values
const customDataStoreId = "customDataStoreId"

// Attributes to test with. Add optional properties to test here if desired.
type customDataStoreResourceModel struct {
	name                            string
	maskAttributeValues             bool
	baseUrl                         string
	tag                             string
	headerName                      string
	headerValue                     string
	localAttribute                  string
	jsonResponseAttributePath       string
	oauthTokenEndpoint              string
	oauthScope                      string
	clientId                        string
	enableHttpsHostnameVerification string
	readTimeout                     string
	connectionTimeout               string
	maxPayloadSize                  string
	retryRequest                    string
	maximumRetriesLimit             string
	testConnectionUrl               string
	testConnectionBody              string
}

func TestAccCustomDataStore(t *testing.T) {
	resourceName := "myCustomDataStore"
	initialResourceModel := customDataStoreResourceModel{
		maskAttributeValues:             false,
		name:                            "initialCustomDataStore",
		baseUrl:                         "https://example.com",
		tag:                             "initialTag",
		headerName:                      "initialHeaderName",
		headerValue:                     "initialHeaderValue",
		localAttribute:                  "initialLocalAttribute",
		jsonResponseAttributePath:       "/initialJsonResponseAttributePath",
		oauthTokenEndpoint:              "https://example.com",
		oauthScope:                      "initialOauthScope",
		clientId:                        "initialClientId",
		enableHttpsHostnameVerification: "false",
		readTimeout:                     "1000",
		connectionTimeout:               "1000",
		maxPayloadSize:                  "1000",
		retryRequest:                    "false",
		maximumRetriesLimit:             "0",
		testConnectionUrl:               "https://example.com",
		testConnectionBody:              "initialTestConnectionBody",
	}
	updatedResourceModel := customDataStoreResourceModel{
		maskAttributeValues:             false,
		name:                            "updatedCustomDataStore",
		baseUrl:                         "http://example.com",
		tag:                             "updatedTag",
		headerName:                      "updatedHeaderName",
		headerValue:                     "updatedHeaderValue",
		localAttribute:                  "updatedLocalAttribute",
		jsonResponseAttributePath:       "/updatedJsonResponseAttributePath",
		oauthTokenEndpoint:              "http://example.com",
		oauthScope:                      "updatedOauthScope",
		clientId:                        "updatedClientId",
		enableHttpsHostnameVerification: "true",
		readTimeout:                     "2000",
		connectionTimeout:               "2000",
		maxPayloadSize:                  "2000",
		retryRequest:                    "true",
		maximumRetriesLimit:             "1",
		testConnectionUrl:               "http://example.com",
		testConnectionBody:              "updatedTestConnectionBody",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckCustomDataStoreDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal model
				Config: testAccCustomDataStore(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedCustomDataStoreAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccCustomDataStore(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedCustomDataStoreAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccCustomDataStore(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_data_store." + resourceName,
				ImportStateId:           customDataStoreId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_data_store.configuration.fields"},
			},
			{
				// Back to the initial minimal model
				Config: testAccCustomDataStore(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedCustomDataStoreAttributes(initialResourceModel),
			},
		},
	})
}

func testAccCustomDataStore(resourceName string, resourceModel customDataStoreResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_data_store" "%[1]s" {
  data_store_id         = "%[2]s"
  mask_attribute_values = %[3]t
  custom_data_store = {
    name = "%[4]s"
    plugin_descriptor_ref = {
      id = "com.pingidentity.pf.datastore.other.RestDataSourceDriver"
    }
    configuration = {
      tables = [
        {
          name = "Base URLs and Tags"
          rows = [
            {
              fields = [
                {
                  name  = "Base URL"
                  value = "%[5]s"
                },
                {
                  name  = "Tags"
                  value = "%[6]s"
                }
              ],
              default_row = true
            }
          ]
        },
        {
          name = "HTTP Request Headers"
          rows = [
            {
              fields = [
                {
                  name  = "Header Name"
                  value = "%[7]s"
                },
                {
                  name  = "Header Value"
                  value = "%[8]s"
                }
              ],
              default_row = false
            }
          ]
        },
        {
          name = "Attributes"
          rows = [
            {
              fields = [
                {
                  name  = "Local Attribute"
                  value = "%[9]s"
                },
                {
                  name  = "JSON Response Attribute Path"
                  value = "%[10]s"
                }
              ],
              default_row = false
            }
          ]
        }
      ],
      fields = [
        {
          name  = "Authentication Method"
          value = "Basic Authentication"
        },
        {
          name  = "HTTP Method"
          value = "GET"
        },
        {
          name  = "Username"
          value = "Administrator"
        },
        {
          name  = "Password"
          value = "2FederateM0re"
        },
        {
          name  = "Password Reference"
          value = ""
        },
        {
          name  = "OAuth Token Endpoint"
          value = "%[11]s"
        },
        {
          name  = "OAuth Scope"
          value = "%[12]s"
        },
        {
          name  = "Client ID"
          value = "%[13]s"
        },
        {
          name  = "Client Secret"
          value = "2FederateM0re"
        },
        {
          name  = "Client Secret Reference"
          value = ""
        },
        {
          name  = "Enable HTTPS Hostname Verification"
          value = "%[14]s"
        },
        {
          name  = "Read Timeout (ms)"
          value = "%[15]s"
        },
        {
          name  = "Connection Timeout (ms)"
          value = "%[16]s"
        },
        {
          name  = "Max Payload Size (KB)"
          value = "%[17]s"
        },
        {
          name  = "Retry Request"
          value = "%[18]s"
        },
        {
          name  = "Maximum Retries Limit"
          value = "%[19]s"
        },
        {
          name  = "Retry Error Codes"
          value = "429"
        },
        {
          name  = "Test Connection URL"
          value = "%[20]s"
        },
        {
          name  = "Test Connection Body"
          value = "%[21]s"
        }
      ]
    }
  }
}
data "pingfederate_data_store" "%[1]s" {
  data_store_id = pingfederate_data_store.%[1]s.id
}`, resourceName,
		customDataStoreId,
		resourceModel.maskAttributeValues,
		resourceModel.name,
		resourceModel.baseUrl,
		resourceModel.tag,
		resourceModel.headerName,
		resourceModel.headerValue,
		resourceModel.localAttribute,
		resourceModel.jsonResponseAttributePath,
		resourceModel.oauthTokenEndpoint,
		resourceModel.oauthScope,
		resourceModel.clientId,
		resourceModel.enableHttpsHostnameVerification,
		resourceModel.readTimeout,
		resourceModel.connectionTimeout,
		resourceModel.maxPayloadSize,
		resourceModel.retryRequest,
		resourceModel.maximumRetriesLimit,
		resourceModel.testConnectionUrl,
		resourceModel.testConnectionBody,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedCustomDataStoreAttributes(config customDataStoreResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "CustomDataStore"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		resp, _, err := testClient.DataStoresAPI.GetDataStore(ctx, customDataStoreId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "name", config.name, resp.CustomDataStore.Name)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, pointers.String(customDataStoreId), "mask_attribute_values", config.maskAttributeValues, *resp.CustomDataStore.MaskAttributeValues)
		if err != nil {
			return err
		}

		configuration := resp.CustomDataStore.Configuration
		configTables := configuration.Tables
		configFields := configuration.Fields

		//  configuration.tables
		for i := range configTables {
			switch configTables[i].Name {
			case "Base URLs and Tags":
				for configRow := range configTables[i].Rows {
					configFields := configTables[i].Rows[configRow].Fields
					for configField := range configFields {
						switch configFields[configField].Name {
						case "Base URL":
							err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.tables.rows.fields.name.value", config.baseUrl, *configFields[configField].Value)
							if err != nil {
								return err
							}
						case "Tags":
							err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.tables.rows.fields.name.value", config.tag, *configFields[configField].Value)
							if err != nil {
								return err
							}
						}
					}
				}
			case "HTTP Request Headers":
				for configRow := range configTables[i].Rows {
					configFields := configTables[i].Rows[configRow].Fields
					for configField := range configFields {
						switch configFields[configField].Name {
						case "Header Name":
							err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.tables.rows.fields.name.value", config.headerName, *configFields[configField].Value)
							if err != nil {
								return err
							}
						case "Header Value":
							err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.tables.rows.fields.name.value", config.headerValue, *configFields[configField].Value)
							if err != nil {
								return err
							}
						}
					}
				}
			case "Attributes":
				for configRow := range configTables[i].Rows {
					configFields := configTables[i].Rows[configRow].Fields
					for configField := range configFields {
						switch configFields[configField].Name {
						case "Local Attribute":
							err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.tables.rows.fields.name.value", config.localAttribute, *configFields[configField].Value)
							if err != nil {
								return err
							}
						case "JSON Response Attribute Path":
							err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.tables.rows.fields.name.value", config.jsonResponseAttributePath, *configFields[configField].Value)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}

		// configuration.fields
		for i := range configFields {
			switch configFields[i].Name {
			case "OAuth Token Endpoint":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.oauthTokenEndpoint, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "OAuth Scope":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.oauthScope, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Client ID":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.clientId, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Enable HTTPS Hostname Verification":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.enableHttpsHostnameVerification, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Read Timeout (ms)":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.readTimeout, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Connection Timeout (ms)":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.connectionTimeout, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Max Payload Size (KB)":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.maxPayloadSize, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Retry Request":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.retryRequest, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Maximum Retries Limit":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.maximumRetriesLimit, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Test Connection URL":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.testConnectionUrl, *configFields[i].Value)
				if err != nil {
					return err
				}
			case "Test Connection Body":
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(customDataStoreId), "configuration.fields.name.value", config.testConnectionBody, *configFields[i].Value)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckCustomDataStoreDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.DataStoresAPI.DeleteDataStore(ctx, customDataStoreId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("CustomDataStore", customDataStoreId)
	}
	return nil
}
