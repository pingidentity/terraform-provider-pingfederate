package acctest_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1130/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const tokenExchangeProcessorPolicyToTokenGeneratorMappingId = "tokenexchangeprocessorpolicy|tokengenerator"
const tokenExchangeProccessorSourceId = "tokenexchangeprocessorpolicy"
const tokenGenTargetId = "tokengenerator"

type tokenExchangeProcessorPolicyToTokenGeneratorMappingResourceModel struct {
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
		"CHANNEL_GROUP", "CONDITION", "JDBC", *client.NewResourceLink("ProvisionerDS"),
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

func TestAccOauthTokenExchangeProcessorTokenGeneratorMapping(t *testing.T) {
	resourceName := "myOauthTokenExchangeProcessorTokenGeneratorMapping"
	initialResourceModel := tokenExchangeProcessorPolicyToTokenGeneratorMappingResourceModel{
		attributeContractFulfillment: initialAttributeContractFulfillment(),
		sourceId:                     tokenExchangeProccessorSourceId,
		targetId:                     tokenGenTargetId,
	}
	updatedResourceModel := tokenExchangeProcessorPolicyToTokenGeneratorMappingResourceModel{
		attributeSource:              attributeSource(),
		attributeContractFulfillment: updatedAttributeContractFulfillment(),
		issuanceCriteria:             issuanceCriteria(),
		sourceId:                     tokenExchangeProccessorSourceId,
		targetId:                     tokenGenTargetId,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOauthTokenExchangeProcessorTokenGeneratorMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthTokenExchangeProcessorTokenGeneratorMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeProcessorTokenGeneratorMappingAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthTokenExchangeProcessorTokenGeneratorMapping(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeProcessorTokenGeneratorMappingAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthTokenExchangeProcessorTokenGeneratorMapping(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_token_exchange_token_generator_mapping." + resourceName,
				ImportStateId:     tokenExchangeProcessorPolicyToTokenGeneratorMappingId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccOauthTokenExchangeProcessorTokenGeneratorMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthTokenExchangeProcessorTokenGeneratorMappingAttributes(initialResourceModel),
			},
		},
	})
}

func testAccOauthTokenExchangeProcessorTokenGeneratorMapping(resourceName string, resourceModel tokenExchangeProcessorPolicyToTokenGeneratorMappingResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_token_exchange_token_generator_mapping" "%[1]s" {
  source_id = "%[2]s"
  target_id = "%[3]s"
  attribute_contract_fulfillment = {
    "SAML_SUBJECT" = {
			%[4]s
    }
  }
	%[5]s
	%[6]s
}
data "pingfederate_oauth_token_exchange_token_generator_mapping" "%[1]s" {
  mapping_id = "${pingfederate_oauth_token_exchange_token_generator_mapping.%[1]s.source_id}|${pingfederate_oauth_token_exchange_token_generator_mapping.%[1]s.target_id}"
}`, resourceName,
		resourceModel.sourceId,
		resourceModel.targetId,
		attributeContractFulfillmentHclBlock(&resourceModel.attributeContractFulfillment),
		attributeSourcesHclBlock(resourceModel.attributeSource),
		issuanceCriteriaHclBlock(resourceModel.issuanceCriteria),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthTokenExchangeProcessorTokenGeneratorMappingAttributes(config tokenExchangeProcessorPolicyToTokenGeneratorMappingResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthTokenExchangeProcessorTokenGeneratorMapping"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthTokenExchangeTokenGeneratorMappingsAPI.GetTokenGeneratorMappingById(ctx, tokenExchangeProcessorPolicyToTokenGeneratorMappingId).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenExchangeProcessorPolicyToTokenGeneratorMappingId), "type",
			config.attributeContractFulfillment.Source.Type, response.AttributeContractFulfillment["SAML_SUBJECT"].Source.Type)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenExchangeProcessorPolicyToTokenGeneratorMappingId), "value",
			config.attributeContractFulfillment.Value, response.AttributeContractFulfillment["SAML_SUBJECT"].Value)
		if err != nil {
			return err
		}

		err = attributesources.ValidateResponseAttributes(resourceType, pointers.String(tokenExchangeProcessorPolicyToTokenGeneratorMappingId),
			config.attributeSource, nil, response.AttributeSources)
		if err != nil {
			return err
		}

		if response.IssuanceCriteria != nil {
			conditionalCriteria := response.IssuanceCriteria.ConditionalCriteria
			for _, conditionalCriteriaEntry := range conditionalCriteria {
				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenExchangeProcessorPolicyToTokenGeneratorMappingId), "type",
					config.issuanceCriteria.Source.Type, conditionalCriteriaEntry.Source.Type)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenExchangeProcessorPolicyToTokenGeneratorMappingId), "attribute_name",
					config.issuanceCriteria.AttributeName, conditionalCriteriaEntry.AttributeName)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenExchangeProcessorPolicyToTokenGeneratorMappingId), "condition",
					config.issuanceCriteria.Condition, conditionalCriteriaEntry.Condition)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, stringPointer(tokenExchangeProcessorPolicyToTokenGeneratorMappingId), "value",
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
func testAccCheckOauthTokenExchangeProcessorTokenGeneratorMappingDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthTokenExchangeTokenGeneratorMappingsAPI.DeleteTokenGeneratorMappingById(ctx, tokenExchangeProcessorPolicyToTokenGeneratorMappingId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthTokenExchangeProcessorTokenGeneratorMapping", tokenExchangeProcessorPolicyToTokenGeneratorMappingId)
	}
	return nil
}
