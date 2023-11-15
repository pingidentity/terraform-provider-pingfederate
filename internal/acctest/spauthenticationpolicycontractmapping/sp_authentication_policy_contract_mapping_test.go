package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const spAuthenticationPolicyContractMappingId = "authenticationpolicycontractid|spadapter"
const apcSourceId = "authenticationpolicycontractid"
const spTargetId = "spadapter"

type spAuthenticationPolicyContractMappingResourceModel struct {
	jdbcAttributeSource          *client.JdbcAttributeSource
	ldapAttributeSource          *client.LdapAttributeSource
	attributeContractFulfillment client.AttributeFulfillmentValue
	issuanceCriteria             *client.ConditionalIssuanceCriteriaEntry
	sourceId                     string
	targetId                     string
}

func TestAccSpAuthenticationPolicyContractMapping(t *testing.T) {
	resourceName := "spAuthenticationPolicyContractMapping"
	initialResourceModel := spAuthenticationPolicyContractMappingResourceModel{
		attributeContractFulfillment: attributecontractfulfillment.InitialAttributeContractFulfillment(),
		sourceId:                     apcSourceId,
		targetId:                     spTargetId,
	}
	updatedResourceModel := spAuthenticationPolicyContractMappingResourceModel{
		jdbcAttributeSource:          attributesources.JdbcClientStruct("CHANNEL_GROUP", "$${subject}", "JDBC", *client.NewResourceLink("ProvisionerDS")),
		ldapAttributeSource:          attributesources.LdapClientStruct("(cn=John)", "SUBTREE", *client.NewResourceLink("pingdirectory")),
		attributeContractFulfillment: attributecontractfulfillment.UpdatedAttributeContractFulfillment(),
		issuanceCriteria:             issuancecriteria.ConditionalCriteria(),
		sourceId:                     apcSourceId,
		targetId:                     spTargetId,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckSpAuthenticationPolicyContractMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpAuthenticationPolicyContractMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSpAuthenticationPolicyContractMappingAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSpAuthenticationPolicyContractMapping(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSpAuthenticationPolicyContractMappingAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSpAuthenticationPolicyContractMapping(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_sp_authentication_policy_contract_mapping." + resourceName,
				ImportStateId:     spAuthenticationPolicyContractMappingId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpAuthenticationPolicyContractMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSpAuthenticationPolicyContractMappingAttributes(initialResourceModel),
			},
		},
	})
}

func testAccSpAuthenticationPolicyContractMapping(resourceName string, resourceModel spAuthenticationPolicyContractMappingResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractExample" {
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
  contract_id         = "authenticationpolicycontractid"
}
resource "pingfederate_sp_authentication_policy_contract_mapping" "%[1]s" {
  source_id = pingfederate_authentication_policy_contract.authenticationPolicyContractExample.id
  target_id = "%[3]s"
  attribute_contract_fulfillment = {
    "subject" = {
			%[4]s
    }
  }
	%[5]s
	%[6]s
}`, resourceName,
		resourceModel.sourceId,
		resourceModel.targetId,
		attributecontractfulfillment.Hcl(&resourceModel.attributeContractFulfillment),
		attributesources.Hcl(resourceModel.jdbcAttributeSource, resourceModel.ldapAttributeSource),
		issuancecriteria.Hcl(resourceModel.issuanceCriteria),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSpAuthenticationPolicyContractMappingAttributes(config spAuthenticationPolicyContractMappingResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SpAuthenticationPolicyContractMapping"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SpAuthenticationPolicyContractMappingsAPI.GetApcToSpAdapterMappingById(ctx, spAuthenticationPolicyContractMappingId).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(spAuthenticationPolicyContractMappingId), "type",
			config.attributeContractFulfillment.Source.Type, response.AttributeContractFulfillment["subject"].Source.Type)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, pointers.String(spAuthenticationPolicyContractMappingId), "value",
			config.attributeContractFulfillment.Value, response.AttributeContractFulfillment["subject"].Value)
		if err != nil {
			return err
		}

		err = attributesources.ValidateResponseAttributes(resourceType, pointers.String(spAuthenticationPolicyContractMappingId),
			config.jdbcAttributeSource, config.ldapAttributeSource, response.AttributeSources)
		if err != nil {
			return err
		}

		if response.IssuanceCriteria != nil {
			conditionalCriteria := response.IssuanceCriteria.ConditionalCriteria
			for _, conditionalCriteriaEntry := range conditionalCriteria {
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(spAuthenticationPolicyContractMappingId), "type",
					config.issuanceCriteria.Source.Type, conditionalCriteriaEntry.Source.Type)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, pointers.String(spAuthenticationPolicyContractMappingId), "attribute_name",
					config.issuanceCriteria.AttributeName, conditionalCriteriaEntry.AttributeName)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, pointers.String(spAuthenticationPolicyContractMappingId), "condition",
					config.issuanceCriteria.Condition, conditionalCriteriaEntry.Condition)
				if err != nil {
					return err
				}

				err = acctest.TestAttributesMatchString(resourceType, pointers.String(spAuthenticationPolicyContractMappingId), "value",
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
func testAccCheckSpAuthenticationPolicyContractMappingDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.SpAuthenticationPolicyContractMappingsAPI.DeleteApcToSpAdapterMappingById(ctx, spAuthenticationPolicyContractMappingId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("spAuthenticationPolicyContractMapping", spAuthenticationPolicyContractMappingId)
	}
	return nil
}
