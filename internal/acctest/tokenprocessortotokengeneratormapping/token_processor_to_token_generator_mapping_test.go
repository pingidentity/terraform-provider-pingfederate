package acctest_test

import (
	"fmt"
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

const tokenProcessorToTokenGeneratorMappingId = "tokenprocessor|tokengenerator"
const tokenProcSourceId = "tokenprocessor"
const tokenGenTargetId = "tokengenerator"

type tokenProcessorToTokenGeneratorMappingResourceModel struct {
	attributeSource              *client.JdbcAttributeSource
	attributeContractFulfillment client.AttributeFulfillmentValue
	issuanceCriteria             *client.ConditionalIssuanceCriteriaEntry
	sourceId                     string
	targetId                     string
}

func stringPointer(val string) *string {
	return &val
}

func initialAttributeContractFulfillment() client.AttributeFulfillmentValue {
	initialAttributecontractfulfillment := *client.NewAttributeFulfillmentValue(
		*client.NewSourceTypeIdKey("TEXT"),
		"value",
	)
	return initialAttributecontractfulfillment
}

func updatedAttributeContractFulfillment() client.AttributeFulfillmentValue {
	updatedAttributecontractfulfillment := *client.NewAttributeFulfillmentValue(
		*client.NewSourceTypeIdKey("CONTEXT"),
		"ClientIp",
	)
	return updatedAttributecontractfulfillment
}

func attributeContractFulfillmentHclBlock(aCf *client.AttributeFulfillmentValue) string {
	var builder strings.Builder
	if aCf == nil {
		return ""
	}
	if aCf != nil {
		builder.WriteString("      source = {\n")
		builder.WriteString("        type = \"")
		builder.WriteString(aCf.Source.Type)
		builder.WriteString("\"\n")
		builder.WriteString("      },\n")
		builder.WriteString("      value = \"")
		builder.WriteString(aCf.Value)
		builder.WriteString("\"")
	}
	return builder.String()
}

func attributeSource() *client.JdbcAttributeSource {
	jdbcAttributeSource := client.NewJdbcAttributeSource(
		"CHANNEL_GROUP", "$${SAML_SUBJECT}", "JDBC", *client.NewResourceLink("ProvisionerDS"),
	)
	jdbcAttributeSource.Id = stringPointer("attributesourceid")
	jdbcAttributeSource.ColumnNames = []string{"CREATED"}
	jdbcAttributeSource.Description = stringPointer("description")
	jdbcAttributeSource.Schema = stringPointer("PUBLIC")
	return jdbcAttributeSource
}

func attributeSourcesHclBlock(attrSource *client.JdbcAttributeSource) string {
	var builder strings.Builder
	if attrSource == nil {
		return ""
	}
	if attrSource != nil {
		builder.WriteString("  attribute_sources = [\n")
		builder.WriteString("    {\n")
		builder.WriteString("      jdbc_attribute_source = {\n")
		builder.WriteString("        data_store_ref = {\n")
		builder.WriteString("          id = \"")
		builder.WriteString(attrSource.DataStoreRef.Id)
		builder.WriteString("\"\n        }\n        id           = \"")
		builder.WriteString(*attrSource.Id)
		builder.WriteString("\"\n        description  = \"")
		builder.WriteString(*attrSource.Description)
		builder.WriteString("\"\n        schema       = \"")
		builder.WriteString(*attrSource.Schema)
		builder.WriteString("\"\n        table        = \"")
		builder.WriteString(attrSource.Table)
		builder.WriteString("\"\n        filter       = \"")
		builder.WriteString(attrSource.Filter)
		builder.WriteString("\"\n        column_names = ")
		builder.WriteString(acctest.StringSliceToTerraformString(attrSource.ColumnNames))
		builder.WriteString("\n      }\n    }\n  ]")
	}
	return builder.String()
}

func issuanceCriteria() *client.ConditionalIssuanceCriteriaEntry {
	conditionalIssuanceCriteriaEntry := client.NewConditionalIssuanceCriteriaEntry(
		*client.NewSourceTypeIdKey("CONTEXT"), "ClientIp", "EQUALS", "value")
	conditionalIssuanceCriteriaEntry.ErrorResult = stringPointer("error")
	return conditionalIssuanceCriteriaEntry
}

func issuanceCriteriaHclBlock(conditionalIssuanceCriteriaEntry *client.ConditionalIssuanceCriteriaEntry) string {
	var builder strings.Builder
	if conditionalIssuanceCriteriaEntry == nil {
		return ""
	}
	if conditionalIssuanceCriteriaEntry != nil {
		builder.WriteString("  issuance_criteria = {\n    conditional_criteria = [\n      {\n")
		builder.WriteString("        error_result = \"")
		builder.WriteString(*conditionalIssuanceCriteriaEntry.ErrorResult)
		builder.WriteString("\"\n        source = {\n          type = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.Source.Type)
		builder.WriteString("\"\n        }\n        attribute_name = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.AttributeName)
		builder.WriteString("\"\n        condition      = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.Condition)
		builder.WriteString("\"\n        value          = \"")
		builder.WriteString(conditionalIssuanceCriteriaEntry.Value)
		builder.WriteString("\"\n      }\n    ]\n  }")
	}
	return builder.String()
}

func TestAccTokenProcessorToTokenGeneratorMapping(t *testing.T) {
	resourceName := "myTokenProcessorToTokenGeneratorMapping"
	initialResourceModel := tokenProcessorToTokenGeneratorMappingResourceModel{
		attributeContractFulfillment: initialAttributeContractFulfillment(),
		sourceId:                     tokenProcSourceId,
		targetId:                     tokenGenTargetId,
	}
	updatedResourceModel := tokenProcessorToTokenGeneratorMappingResourceModel{
		attributeSource:              attributeSource(),
		attributeContractFulfillment: updatedAttributeContractFulfillment(),
		issuanceCriteria:             issuanceCriteria(),
		sourceId:                     tokenProcSourceId,
		targetId:                     tokenGenTargetId,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckTokenProcessorToTokenGeneratorMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTokenProcessorToTokenGeneratorMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedTokenProcessorToTokenGeneratorMappingAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccTokenProcessorToTokenGeneratorMapping(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedTokenProcessorToTokenGeneratorMappingAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccTokenProcessorToTokenGeneratorMapping(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_token_processor_to_token_generator_mapping." + resourceName,
				ImportStateId:     tokenProcessorToTokenGeneratorMappingId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTokenProcessorToTokenGeneratorMapping(resourceName string, resourceModel tokenProcessorToTokenGeneratorMappingResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_token_processor_to_token_generator_mapping" "%[1]s" {
  custom_id = "%[2]s"
  source_id = "%[3]s"
  target_id = "%[4]s"
  attribute_contract_fulfillment = {
    "SAML_SUBJECT" = {
			%[5]s
    }
  }
	%[6]s
	%[7]s
}`, resourceName,
		tokenProcessorToTokenGeneratorMappingId,
		resourceModel.sourceId,
		resourceModel.targetId,
		attributeContractFulfillmentHclBlock(&resourceModel.attributeContractFulfillment),
		attributeSourcesHclBlock(resourceModel.attributeSource),
		issuanceCriteriaHclBlock(resourceModel.issuanceCriteria),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedTokenProcessorToTokenGeneratorMappingAttributes(config tokenProcessorToTokenGeneratorMappingResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "TokenProcessorToTokenGeneratorMapping"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.TokenProcessorToTokenGeneratorMappingsAPI.GetTokenToTokenMappingById(ctx, tokenProcessorToTokenGeneratorMappingId).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "type",
			config.attributeContractFulfillment.Source.Type, response.AttributeContractFulfillment["SAML_SUBJECT"].Source.Type)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "value",
			config.attributeContractFulfillment.Value, response.AttributeContractFulfillment["SAML_SUBJECT"].Value)
		if err != nil {
			return err
		}

		attributeSources := response.AttributeSources
		for _, attributeSource := range attributeSources {
			if attributeSource.JdbcAttributeSource != nil {
				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "id",
					config.attributeSource.DataStoreRef.Id, attributeSource.JdbcAttributeSource.DataStoreRef.Id)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "description",
					*config.attributeSource.Description, *attributeSource.JdbcAttributeSource.Description)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "schema",
					*config.attributeSource.Description, *attributeSource.JdbcAttributeSource.Description)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "table",
					config.attributeSource.Table, attributeSource.JdbcAttributeSource.Table)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "filter",
					config.attributeSource.Filter, "$"+attributeSource.JdbcAttributeSource.Filter)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchStringSlice(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "column_names",
					config.attributeSource.ColumnNames, attributeSource.JdbcAttributeSource.ColumnNames)
				if err != nil {
					return err
				}
			}
		}

		if response.IssuanceCriteria != nil {
			conditionalCriteria := response.IssuanceCriteria.ConditionalCriteria
			for _, conditionalCriteriaEntry := range conditionalCriteria {
				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "type",
					config.issuanceCriteria.Source.Type, conditionalCriteriaEntry.Source.Type)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "attribute_name",
					config.issuanceCriteria.AttributeName, conditionalCriteriaEntry.AttributeName)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "condition",
					config.issuanceCriteria.Condition, conditionalCriteriaEntry.Condition)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenProcessorToTokenGeneratorMappingId), "value",
					config.issuanceCriteria.Value, conditionalCriteriaEntry.Value)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckTokenProcessorToTokenGeneratorMappingDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.TokenProcessorToTokenGeneratorMappingsAPI.DeleteTokenToTokenMappingById(ctx, tokenProcessorToTokenGeneratorMappingId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("TokenProcessorToTokenGeneratorMapping", tokenProcessorToTokenGeneratorMappingId)
	}
	return nil
}
