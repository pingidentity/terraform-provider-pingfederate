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
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const idpAdaptersId = "idpAdapterId"

// Attributes to test with. Add optional properties to test here if desired.
type idpAdaptersResourceModel struct {
	name                  string
	pluginDescriptorRefId string
	configuration         client.PluginConfiguration
	attributeMapping      *client.IdpAdapterContractMapping
	attributeContract     client.IdpAdapterAttributeContract
}

func boolPointer(val bool) *bool {
	return &val
}

func stringPointer(val string) *string {
	return &val
}

func basicAttributeContract() *client.IdpAdapterAttributeContract {
	attributeContract := client.NewIdpAdapterAttributeContract([]client.IdpAdapterAttribute{})
	attributeContract.SetMaskOgnlValues(false)
	attributeContract.CoreAttributes = append(attributeContract.CoreAttributes, client.IdpAdapterAttribute{
		Name: "username",
	})
	attributeContract.CoreAttributes = append(attributeContract.CoreAttributes, client.IdpAdapterAttribute{
		Name:      "policy.action",
		Pseudonym: boolPointer(true),
	})
	return attributeContract
}

func TestAccIdpAdapters(t *testing.T) {
	resourceName := "myIdpAdapters"
	initialResourceModel := idpAdaptersResourceModel{
		name:                  "testIdpAdapter",
		pluginDescriptorRefId: "com.pingidentity.adapters.htmlform.idp.HtmlFormIdpAuthnAdapter",
		configuration: client.PluginConfiguration{
			Tables: []client.ConfigTable{
				{
					Name: "Credential Validators",
					Rows: []client.ConfigRow{
						{
							DefaultRow: boolPointer(false),
							Fields: []client.ConfigField{
								{
									Name:  "Password Credential Validator Instance",
									Value: stringPointer("pingdirectory"),
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
	attributeMapping.Inherited = boolPointer(false)
	attributeContract := basicAttributeContract()
	attributeContract.ExtendedAttributes = append(attributeContract.ExtendedAttributes, client.IdpAdapterAttribute{
		Name:      "entryUUID",
		Pseudonym: boolPointer(false),
		Masked:    boolPointer(false),
	})
	updatedResourceModel := idpAdaptersResourceModel{
		name:                  "testIdpAdapterNewName",
		pluginDescriptorRefId: "com.pingidentity.adapters.htmlform.idp.HtmlFormIdpAuthnAdapter",
		configuration: client.PluginConfiguration{
			Tables: []client.ConfigTable{
				{
					Name: "Credential Validators",
					Rows: []client.ConfigRow{
						{
							DefaultRow: boolPointer(false),
							Fields: []client.ConfigField{
								{
									Name:  "Password Credential Validator Instance",
									Value: stringPointer("pingdirectory"),
								},
							},
						},
					},
				},
			},
			Fields: []client.ConfigField{
				{
					Name:  "Challenge Retries",
					Value: stringPointer("3"),
				},
			},
		},
		attributeContract: *attributeContract,
		attributeMapping:  attributeMapping,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckIdpAdaptersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdpAdapters(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedIdpAdaptersAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccIdpAdapters(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedIdpAdaptersAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:        testAccIdpAdapters(resourceName, updatedResourceModel),
				ResourceName:  "pingfederate_idp_adapters." + resourceName,
				ImportStateId: idpAdaptersId,
				ImportState:   true,
				//TODO need to re-enable this once we have a way to store fields/tables/attrs/etc. that PF generates itself and returns
				//ImportStateVerify: true,
			},
		},
	})
}

func configurationHclBlock(configuration client.PluginConfiguration) string {
	var builder strings.Builder
	builder.WriteString("configuration = {\n")
	builder.WriteString("    tables = [\n")
	for _, table := range configuration.Tables {
		builder.WriteString("       {\n")
		builder.WriteString("           name = \"")
		builder.WriteString(table.Name)
		builder.WriteString("\"\n")
		builder.WriteString("           rows = [\n")
		for _, row := range table.Rows {
			builder.WriteString("               {\n")
			if row.DefaultRow != nil {
				builder.WriteString("                   default_row = ")
				builder.WriteString(strconv.FormatBool(*row.DefaultRow))
				builder.WriteRune('\n')
			}
			builder.WriteString("                   fields = [\n")
			for _, field := range row.Fields {
				builder.WriteString("                       {\n")
				builder.WriteString("                           name = \"")
				builder.WriteString(field.Name)
				builder.WriteString("\"\n")
				if field.Value != nil {
					builder.WriteString("                           value = \"")
					builder.WriteString(*field.Value)
					builder.WriteString("\"\n")
				}
				builder.WriteString("                       },\n")
			}
			builder.WriteString("                   ]\n")
			builder.WriteString("               }\n")
		}
		builder.WriteString("           ]\n")
		builder.WriteString("       },\n")
	}
	builder.WriteString("    ]\n")
	builder.WriteString("    fields = [\n")
	for _, field := range configuration.Fields {
		builder.WriteString("        {\n")
		builder.WriteString("            name = \"")
		builder.WriteString(field.Name)
		builder.WriteString("\"\n")
		if field.Value != nil {
			builder.WriteString("            value = \"")
			builder.WriteString(*field.Value)
			builder.WriteString("\"\n")
		}
		builder.WriteString("        },\n")
	}
	builder.WriteString("    ]\n")
	builder.WriteString("}\n")
	return builder.String()
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
	//TODO attribute sources
	/*if attributeMapping.AttributeSources != nil && len(attributeMapping.AttributeSources) > 0 {

	}*/
	if attributeMapping.AttributeContractFulfillment != nil {
		builder.WriteString("    attribute_contract_fulfillment = {\n")
		for key, val := range attributeMapping.AttributeContractFulfillment {
			builder.WriteString("        \"")
			builder.WriteString(key)
			builder.WriteString("\" = {\n")
			builder.WriteString("            source = {\n")
			builder.WriteString("                type = \"")
			builder.WriteString(val.Source.Type)
			builder.WriteString("\"\n")
			if val.Source.Id != nil {
				builder.WriteString("                id = \"")
				builder.WriteString(*val.Source.Id)
				builder.WriteString("\"\n")
			}
			builder.WriteString("            }\n")
			builder.WriteString("            value = \"")
			builder.WriteString(val.Value)
			builder.WriteString("\"\n")
			builder.WriteString("        }\n")
		}
		builder.WriteString("    }\n")
	}
	//TODO issuance_criteria
	/*if attributeMapping.IssuanceCriteria != nil {
		builder.WriteString("    issuance_criteria = ")
		builder.WriteString(strconv.FormatBool(*attributeContract.MaskOgnlValues))
		builder.WriteRune('\n')
	}*/
	if attributeMapping.Inherited != nil {
		builder.WriteString("    inherited = ")
		builder.WriteString(strconv.FormatBool(*attributeMapping.Inherited))
		builder.WriteRune('\n')
	}
	builder.WriteString("}\n")
	return builder.String()
}

func testAccIdpAdapters(resourceName string, resourceModel idpAdaptersResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_idp_adapters" "%[1]s" {
	custom_id = "%[2]s"
	name = "%[3]s"
	plugin_descriptor_ref = {
        id = "%[4]s"
    }
	%[5]s
	%[6]s
	%[7]s
}`, resourceName,
		idpAdaptersId,
		resourceModel.name,
		resourceModel.pluginDescriptorRefId,
		configurationHclBlock(resourceModel.configuration),
		attributeContractHclBlock(resourceModel.attributeContract),
		attributeMappingHclBlock(resourceModel.attributeMapping),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedIdpAdaptersAttributes(config idpAdaptersResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//resourceType := "IdpAdapters"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.IdpAdaptersAPI.GetIdpAdapters(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		//FILL THESE in!

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckIdpAdaptersDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.IdpAdaptersAPI.DeleteIdpAdapter(ctx, idpAdaptersId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("IdpAdapters", idpAdaptersId)
	}
	return nil
}
