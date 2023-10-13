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
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/configuration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const idpAdapterId = "idpAdapterId"

// Attributes to test with. Add optional properties to test here if desired.
type idpAdapterResourceModel struct {
	name                  string
	pluginDescriptorRefId string
	configuration         client.PluginConfiguration
	attributeMapping      *client.IdpAdapterContractMapping
	attributeContract     client.IdpAdapterAttributeContract
}

func basicAttributeContract() *client.IdpAdapterAttributeContract {
	attributeContract := client.NewIdpAdapterAttributeContract([]client.IdpAdapterAttribute{})
	attributeContract.SetMaskOgnlValues(false)
	attributeContract.CoreAttributes = append(attributeContract.CoreAttributes, client.IdpAdapterAttribute{
		Name: "username",
	})
	attributeContract.CoreAttributes = append(attributeContract.CoreAttributes, client.IdpAdapterAttribute{
		Name:      "policy.action",
		Pseudonym: pointers.Bool(true),
	})
	return attributeContract
}

func updatedAttributeContract() *client.IdpAdapterAttributeContract {
	contract := basicAttributeContract()
	contract.ExtendedAttributes = append(contract.ExtendedAttributes, client.IdpAdapterAttribute{
		Name:      "entryUUID",
		Pseudonym: pointers.Bool(false),
		Masked:    pointers.Bool(false),
	})
	return contract
}

func updatedAttributeMapping() *client.IdpAdapterContractMapping {
	attributeMapping := client.NewIdpAdapterContractMapping(map[string]client.AttributeFulfillmentValue{
		"entryUUID": {
			Source: client.SourceTypeIdKey{
				Type: "ADAPTER",
			},
			Value: "entryUUID",
		},
		"policy.action": {
			Source: client.SourceTypeIdKey{
				Type: "ADAPTER",
			},
			Value: "policy.action",
		},
		"username": {
			Source: client.SourceTypeIdKey{
				Type: "ADAPTER",
			},
			Value: "username",
		},
	})
	attributeMapping.Inherited = pointers.Bool(false)

	attributeMapping.AttributeSources = []client.AttributeSourceAggregation{
		{
			JdbcAttributeSource: attributesources.JdbcClientStruct(),
		},
	}
	attributeMapping.IssuanceCriteria = client.NewIssuanceCriteria()
	criteria := issuancecriteria.ConditionalCriteria()
	attributeMapping.IssuanceCriteria.ConditionalCriteria = []client.ConditionalIssuanceCriteriaEntry{
		*criteria,
	}
	return attributeMapping
}

func TestAccIdpAdapter(t *testing.T) {
	resourceName := "myIdpAdapter"
	initialResourceModel := idpAdapterResourceModel{
		name:                  "testIdpAdapter",
		pluginDescriptorRefId: "com.pingidentity.adapters.htmlform.idp.HtmlFormIdpAuthnAdapter",
		configuration: client.PluginConfiguration{
			Tables: []client.ConfigTable{
				{
					Name: "Credential Validators",
					Rows: []client.ConfigRow{
						{
							DefaultRow: pointers.Bool(false),
							Fields: []client.ConfigField{
								{
									Name:  "Password Credential Validator Instance",
									Value: pointers.String("pingdirectory"),
								},
							},
						},
					},
				},
			},
			Fields: []client.ConfigField{},
		},
		attributeContract: *basicAttributeContract(),
	}

	updatedResourceModel := idpAdapterResourceModel{
		name:                  "testIdpAdapterNewName",
		pluginDescriptorRefId: "com.pingidentity.adapters.htmlform.idp.HtmlFormIdpAuthnAdapter",
		configuration: client.PluginConfiguration{
			Tables: []client.ConfigTable{
				{
					Name: "Credential Validators",
					Rows: []client.ConfigRow{
						{
							DefaultRow: pointers.Bool(false),
							Fields: []client.ConfigField{
								{
									Name:  "Password Credential Validator Instance",
									Value: pointers.String("pingdirectory"),
								},
							},
						},
					},
				},
			},
			Fields: []client.ConfigField{
				{
					Name:  "Challenge Retries",
					Value: pointers.String("3"),
				},
			},
		},
		attributeContract: *updatedAttributeContract(),
		attributeMapping:  updatedAttributeMapping(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckIdpAdapterDestroy,
		Steps: []resource.TestStep{
			{
				// Minimal model
				Config: testAccIdpAdapter(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpAdapterAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIdpAdapter(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIdpAdapterAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccIdpAdapter(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_idp_adapter." + resourceName,
				ImportStateId:     idpAdapterId,
				ImportState:       true,
				ImportStateVerify: true,
				// Can't verify fields and core_attributes because the computed ones from the server will go into the
				// corresponding fields_all and core_attributes_all fields
				ImportStateVerifyIgnore: []string{"configuration.fields", "attribute_contract.core_attributes"},
			},
			{
				// Back to the initial minimal model
				Config: testAccIdpAdapter(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpAdapterAttributes(initialResourceModel),
			},
		},
	})
}

func attributeContractHclBlock(attributeContract client.IdpAdapterAttributeContract) string {
	var builder strings.Builder
	builder.WriteString("attribute_contract = {\n")
	if attributeContract.MaskOgnlValues != nil {
		builder.WriteString("    mask_ognl_values = ")
		builder.WriteString(strconv.FormatBool(*attributeContract.MaskOgnlValues))
		builder.WriteRune('\n')
	}
	if attributeContract.Inherited != nil {
		builder.WriteString("    inherited = ")
		builder.WriteString(strconv.FormatBool(*attributeContract.Inherited))
		builder.WriteRune('\n')
	}
	if attributeContract.UniqueUserKeyAttribute != nil {
		builder.WriteString("    unique_user_key_attribute = \"")
		builder.WriteString(*attributeContract.UniqueUserKeyAttribute)
		builder.WriteString("\"\n")
	}
	builder.WriteString("    core_attributes = [\n")
	for _, attr := range attributeContract.CoreAttributes {
		builder.WriteString("        {\n")
		builder.WriteString("            name = \"")
		builder.WriteString(attr.Name)
		builder.WriteString("\"\n")
		if attr.Masked != nil {
			builder.WriteString("            masked = ")
			builder.WriteString(strconv.FormatBool(*attr.Masked))
			builder.WriteRune('\n')
		}
		if attr.Pseudonym != nil {
			builder.WriteString("            pseudonym = ")
			builder.WriteString(strconv.FormatBool(*attr.Pseudonym))
			builder.WriteRune('\n')
		}
		builder.WriteString("        },\n")
	}
	builder.WriteString("    ]\n")
	builder.WriteString("    extended_attributes = [\n")
	for _, attr := range attributeContract.ExtendedAttributes {
		builder.WriteString("        {\n")
		builder.WriteString("            name = \"")
		builder.WriteString(attr.Name)
		builder.WriteString("\"\n")
		if attr.Masked != nil {
			builder.WriteString("            masked = ")
			builder.WriteString(strconv.FormatBool(*attr.Masked))
			builder.WriteRune('\n')
		}
		if attr.Pseudonym != nil {
			builder.WriteString("            pseudonym = ")
			builder.WriteString(strconv.FormatBool(*attr.Pseudonym))
			builder.WriteRune('\n')
		}
		builder.WriteString("        },\n")
	}
	builder.WriteString("    ]\n")
	builder.WriteString("}\n")
	return builder.String()
}

func attributeMappingHclBlock(attributeMapping *client.IdpAdapterContractMapping) string {
	if attributeMapping == nil {
		return ""
	}
	var builder strings.Builder
	builder.WriteString("attribute_mapping = {\n")
	if len(attributeMapping.AttributeSources) > 0 {
		// Only have logic for JDBC attribute sources right now, assume it is the right type
		attributesources.JdbcHcl(attributeMapping.AttributeSources[0].JdbcAttributeSource)
	}
	if attributeMapping.AttributeContractFulfillment != nil {
		builder.WriteString("    attribute_contract_fulfillment = {\n")
		for key, val := range attributeMapping.AttributeContractFulfillment {
			builder.WriteString("        \"")
			builder.WriteString(key)
			builder.WriteString("\" = {\n")
			// Avoid taking address of for loop variable
			innerVal := val
			builder.WriteString(attributecontractfulfillment.Hcl(&innerVal))
			builder.WriteString("        }\n")
		}
		builder.WriteString("    }\n")
	}
	if attributeMapping.IssuanceCriteria != nil && len(attributeMapping.IssuanceCriteria.ConditionalCriteria) > 0 {
		// Onlye have logic for one conditional criteria right now
		builder.WriteString(issuancecriteria.Hcl(&attributeMapping.IssuanceCriteria.ConditionalCriteria[0]))
	}
	if attributeMapping.Inherited != nil {
		builder.WriteString("    inherited = ")
		builder.WriteString(strconv.FormatBool(*attributeMapping.Inherited))
		builder.WriteRune('\n')
	}
	builder.WriteString("}\n")
	return builder.String()
}

func testAccIdpAdapter(resourceName string, resourceModel idpAdapterResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_adapter" "%[1]s" {
  custom_id = "%[2]s"
  name      = "%[3]s"
  plugin_descriptor_ref = {
    id = "%[4]s"
  }
	%[5]s
	%[6]s
	%[7]s
}`, resourceName,
		idpAdapterId,
		resourceModel.name,
		resourceModel.pluginDescriptorRefId,
		configuration.Hcl(resourceModel.configuration),
		attributeContractHclBlock(resourceModel.attributeContract),
		attributeMappingHclBlock(resourceModel.attributeMapping),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedIdpAdapterAttributes(config idpAdapterResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "IdpAdapter"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		resp, _, err := testClient.IdpAdaptersAPI.GetIdpAdapter(ctx, idpAdapterId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(idpAdapterId), "name", config.name, resp.Name)
		if err != nil {
			return err
		}

		// Plugin descriptor
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(idpAdapterId), "plugin_descriptor_ref.id", config.pluginDescriptorRefId, resp.PluginDescriptorRef.Id)
		if err != nil {
			return err
		}

		// JDBC attribute sources
		if config.attributeMapping != nil && len(config.attributeMapping.AttributeSources) > 0 {
			configAttrSource := config.attributeMapping.AttributeSources[0].JdbcAttributeSource
			attributeSources := resp.AttributeMapping.AttributeSources
			for _, attributeSource := range attributeSources {
				if attributeSource.JdbcAttributeSource != nil {
					err = acctest.TestAttributesMatchString(resourceType, pointers.String(idpAdapterId), "id",
						configAttrSource.DataStoreRef.Id, attributeSource.JdbcAttributeSource.DataStoreRef.Id)
					if err != nil {
						return err
					}

					err = acctest.TestAttributesMatchStringPointer(resourceType, pointers.String(idpAdapterId), "description",
						*configAttrSource.Description, attributeSource.JdbcAttributeSource.Description)
					if err != nil {
						return err
					}

					err = acctest.TestAttributesMatchStringPointer(resourceType, pointers.String(idpAdapterId), "schema",
						*configAttrSource.Description, attributeSource.JdbcAttributeSource.Description)
					if err != nil {
						return err
					}

					err = acctest.TestAttributesMatchString(resourceType, pointers.String(idpAdapterId), "table",
						configAttrSource.Table, attributeSource.JdbcAttributeSource.Table)
					if err != nil {
						return err
					}

					err = acctest.TestAttributesMatchString(resourceType, pointers.String(idpAdapterId), "filter",
						configAttrSource.Filter, attributeSource.JdbcAttributeSource.Filter)
					if err != nil {
						return err
					}

					err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(idpAdapterId), "column_names",
						configAttrSource.ColumnNames, attributeSource.JdbcAttributeSource.ColumnNames)
					if err != nil {
						return err
					}
				}
			}
		}

		// Attribute mapping attribute contract fullfilment - verify the keys are present in the response
		if config.attributeMapping != nil {
			for configKey := range config.attributeMapping.AttributeContractFulfillment {
				_, ok := resp.AttributeMapping.AttributeContractFulfillment[configKey]
				if !ok {
					return fmt.Errorf("Attribute contract fullfilment %s not found in server response", configKey)
				}
			}
		}

		// Attribute contract - verify the attribute names are in the response
		for _, configAttr := range config.attributeContract.CoreAttributes {
			found := false
			for _, respAttr := range resp.AttributeContract.CoreAttributes {
				if respAttr.Name == configAttr.Name {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("Core attribute %s not found in server response", configAttr.Name)
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckIdpAdapterDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.IdpAdaptersAPI.DeleteIdpAdapter(ctx, idpAdapterId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("IdpAdapters", idpAdapterId)
	}
	return nil
}
