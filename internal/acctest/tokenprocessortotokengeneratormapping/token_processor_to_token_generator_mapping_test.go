package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const tokenProcessorToTokenGeneratorMappingId = "tokenprocessor|tokengenerator"
const tokenProcSourceId = "tokenprocessor"
const tokenGenTargetId = "tokengenerator"

type tokenProcessorToTokenGeneratorMappingResourceModel struct {
	attributeSource              *client.LdapAttributeSource
	attributeContractFulfillment client.AttributeFulfillmentValue
	issuanceCriteria             *client.ConditionalIssuanceCriteriaEntry
	sourceId                     string
	targetId                     string
	defaultTargetResource        *string
}

func TestAccTokenProcessorToTokenGeneratorMapping(t *testing.T) {
	resourceName := "myTokenProcessorToTokenGeneratorMapping"
	initialResourceModel := tokenProcessorToTokenGeneratorMappingResourceModel{
		attributeContractFulfillment: attributecontractfulfillment.InitialAttributeContractFulfillment(),
		sourceId:                     tokenProcSourceId,
		targetId:                     tokenGenTargetId,
	}
	updatedResourceModel := tokenProcessorToTokenGeneratorMappingResourceModel{
		attributeSource:              attributesources.LdapClientStruct("(cn=example)", "SUBTREE", *client.NewResourceLink("pingdirectory")),
		attributeContractFulfillment: attributecontractfulfillment.UpdatedAttributeContractFulfillment(),
		issuanceCriteria:             issuancecriteria.ConditionalCriteria(),
		sourceId:                     tokenProcSourceId,
		targetId:                     tokenGenTargetId,
		defaultTargetResource:        pointers.String("https://example.com"),
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
			{
				// Back to minimal model
				Config: testAccTokenProcessorToTokenGeneratorMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedTokenProcessorToTokenGeneratorMappingAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.TokenProcessorToTokenGeneratorMappingsAPI.DeleteTokenToTokenMappingById(ctx, tokenProcessorToTokenGeneratorMappingId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccTokenProcessorToTokenGeneratorMapping(resourceName string, resourceModel tokenProcessorToTokenGeneratorMappingResourceModel) string {
	defaultTargetResourceHcl := ""
	if resourceModel.defaultTargetResource != nil {
		defaultTargetResourceHcl = fmt.Sprintf("default_target_resource = \"%[1]s\"", *resourceModel.defaultTargetResource)
	}

	// license_connection_group can't be tested without some changes to the license
	return fmt.Sprintf(`
resource "pingfederate_token_processor_to_token_generator_mapping" "%[1]s" {
  source_id = "%[2]s"
  target_id = "%[3]s"
  attribute_contract_fulfillment = {
    "SAML_SUBJECT" = {
			%[4]s
    }
  }
	%[5]s
	%[6]s
	%[7]s
}
data "pingfederate_token_processor_to_token_generator_mapping" "%[1]s" {
  mapping_id = pingfederate_token_processor_to_token_generator_mapping.%[1]s.id
}`, resourceName,
		resourceModel.sourceId,
		resourceModel.targetId,
		attributecontractfulfillment.Hcl(&resourceModel.attributeContractFulfillment),
		attributesources.Hcl(nil, resourceModel.attributeSource),
		issuancecriteria.Hcl(resourceModel.issuanceCriteria),
		defaultTargetResourceHcl,
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
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(tokenProcessorToTokenGeneratorMappingId), "type",
			config.attributeContractFulfillment.Source.Type, response.AttributeContractFulfillment["SAML_SUBJECT"].Source.Type)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, pointers.String(tokenProcessorToTokenGeneratorMappingId), "value",
			config.attributeContractFulfillment.Value, response.AttributeContractFulfillment["SAML_SUBJECT"].Value)
		if err != nil {
			return err
		}

		err = attributesources.ValidateResponseAttributes(resourceType, pointers.String(tokenProcessorToTokenGeneratorMappingId), nil,
			config.attributeSource, response.AttributeSources)
		if err != nil {
			return err
		}

		if response.IssuanceCriteria != nil {
			conditionalCriteria := response.IssuanceCriteria.ConditionalCriteria
			for _, conditionalCriteriaEntry := range conditionalCriteria {
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(tokenProcessorToTokenGeneratorMappingId), "type",
					config.issuanceCriteria.Source.Type, conditionalCriteriaEntry.Source.Type)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, pointers.String(tokenProcessorToTokenGeneratorMappingId), "attribute_name",
					config.issuanceCriteria.AttributeName, conditionalCriteriaEntry.AttributeName)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, pointers.String(tokenProcessorToTokenGeneratorMappingId), "condition",
					config.issuanceCriteria.Condition, conditionalCriteriaEntry.Condition)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, pointers.String(tokenProcessorToTokenGeneratorMappingId), "value",
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
