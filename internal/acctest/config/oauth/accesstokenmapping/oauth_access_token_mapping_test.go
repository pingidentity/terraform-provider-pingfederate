package accesstokenmapping_test

import (
	"fmt"
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

const oauthAccessTokenMappingId = "client_credentials|oauthAccessTokenManagerTest"

// Attributes to test with. Add optional properties to test here if desired.
type oauthAccessTokenMappingResourceModel struct {
	attributeContractFulfillment map[string]client.AttributeFulfillmentValue
	attributeSource              *client.JdbcAttributeSource
	issuanceCriteria             *client.IssuanceCriteria
}

func TestAccOauthAccessTokenMapping(t *testing.T) {
	resourceName := "myOauthAccessTokenMapping"

	// Initial Model
	// attributeContractFulfillment
	initialAttributeContractFulfillmentVal := map[string]client.AttributeFulfillmentValue{}
	initialAFV := client.NewAttributeFulfillmentValue(*client.NewSourceTypeIdKey("TEXT"), "Administrator")
	initialAttributeContractFulfillmentVal["extended_contract"] = *initialAFV

	initialResourceModel := oauthAccessTokenMappingResourceModel{
		attributeContractFulfillment: initialAttributeContractFulfillmentVal,
	}

	//  Updated Model
	// attributeContractFulfillment
	updatedAttributeContractFullfillmentVal := map[string]client.AttributeFulfillmentValue{}
	updatedAttributeContractFullfillmentVal["extended_contract"] = client.AttributeFulfillmentValue{
		Source: client.SourceTypeIdKey{
			Type: "NO_MAPPING",
		},
	}

	//  attributeSource
	jdbcAttributeSource := client.NewJdbcAttributeSource("ADMINISTRABLE_ROLE_AUTHORIZATIONS", "${client_id}", "JDBC", *client.NewResourceLink("ProvisionerDS"))
	jdbcAttributeSource.Schema = pointers.String("INFORMATION_SCHEMA")
	jdbcAttributeSource.ColumnNames = []string{"GRANTEE", "IS_GRANTABLE", "ROLE_NAME"}
	jdbcAttributeSource.Id = pointers.String("test")
	jdbcAttributeSource.Description = pointers.String("test")

	issuanceCriteria := client.NewIssuanceCriteria()
	conditionalIssuanceCriteriaEntry := client.NewConditionalIssuanceCriteriaEntry(*client.NewSourceTypeIdKey("CONTEXT"), "ClientId", "EQUALS_CASE_INSENSITIVE", "text")
	issuanceCriteria.ConditionalCriteria = append(issuanceCriteria.ConditionalCriteria, *conditionalIssuanceCriteriaEntry)

	updatedResourceModel := oauthAccessTokenMappingResourceModel{
		attributeContractFulfillment: updatedAttributeContractFullfillmentVal,
		attributeSource:              jdbcAttributeSource,
		issuanceCriteria:             issuanceCriteria,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOauthAccessTokenMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthAccessTokenMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAccessTokenMappingAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthAccessTokenMapping(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthAccessTokenMappingAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthAccessTokenMapping(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_access_token_mapping." + resourceName,
				ImportStateId:     oauthAccessTokenMappingId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOauthAccessTokenMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAccessTokenMappingAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OauthAccessTokenMappingsAPI.DeleteMapping(ctx, oauthAccessTokenMappingId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccOauthAccessTokenMapping(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAccessTokenMappingAttributes(initialResourceModel),
			},
		},
	})
}

func optionalHcl(resourceModel oauthAccessTokenMappingResourceModel) string {
	var attributeSourceHcl string
	if resourceModel.attributeSource == nil {
		attributeSourceHcl = ""
	} else {
		attributeSourceHcl = fmt.Sprintf(`
		attribute_sources = [
			{
				jdbc_attribute_source = {
					data_store_ref = {
						id = "%[1]s"
					}
					id           = "%[2]s"
					description  = "%[3]s"
					schema       = "%[4]s"
					table        = "%[5]s"
					filter       = "$%[6]s"
					column_names = %[7]s
				}
			},
			{
				custom_attribute_source = {
					data_store_ref = {
					  	id = "customDataStore"
					}
					description = "APIStubs"
					filter_fields = [
					  	{
					  	  	name = "Authorization Header"
					  	},
					  	{
					  	  	name = "Body"
					  	},
					  	{
					  	  	name  = "Resource Path"
					  	  	value = "/users/$${external_id}"
					  	},
					]
					id = "APIStubs"
				}
			},
			{
      ldap_attribute_source = {
        base_dn = "ou=Users,dc=bxretail,dc=org"
        data_store_ref = {
          id = "pingdirectory"
        }
        description            = "Directory"
        id                     = "Directory"
        member_of_nested_group = false
        search_attributes = [
          "Subject DN",
        ]
        search_filter = "(&(memberUid=example)(cn=Postman))"
        search_scope  = "SUBTREE"
        type          = "LDAP"
      }
    },
		]
		`, resourceModel.attributeSource.DataStoreRef.Id,
			*resourceModel.attributeSource.Id,
			*resourceModel.attributeSource.Description,
			*resourceModel.attributeSource.Schema,
			resourceModel.attributeSource.Table,
			resourceModel.attributeSource.Filter,
			acctest.StringSliceToTerraformString(resourceModel.attributeSource.ColumnNames))
	}

	var issuanceCriteriaHcl string
	if resourceModel.issuanceCriteria == nil {
		issuanceCriteriaHcl = ""
	} else {
		issuanceCriteriaHcl = fmt.Sprintf(`
		issuance_criteria = {
			conditional_criteria = [
				{
					attribute_name = "%[1]s"
					condition      = "%[2]s"
					error_result   = "text"
					source = {
						type = "%[3]s"
					}
					value = "%[4]s"
				}
			]
		}
		`, resourceModel.issuanceCriteria.ConditionalCriteria[0].AttributeName,
			resourceModel.issuanceCriteria.ConditionalCriteria[0].Condition,
			resourceModel.issuanceCriteria.ConditionalCriteria[0].Source.Type,
			resourceModel.issuanceCriteria.ConditionalCriteria[0].Value)
	}

	return fmt.Sprintf(`
	%s
	%s
	`, attributeSourceHcl, issuanceCriteriaHcl)
}

func testAccOauthAccessTokenMapping(resourceName string, resourceModel oauthAccessTokenMappingResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager" "oauthAccessTokenManagerTest" {
  manager_id = "oauthAccessTokenManagerTest"
  name       = "oauthAccessTokenManagerTest"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = []
    fields = [
      {
        name  = "Token Length"
        value = "56"
      },
      {
        name  = "Token Lifetime"
        value = "240"
      },
      {
        name  = "Lifetime Extension Policy"
        value = "NONE"
      },
      {
        name  = "Maximum Token Lifetime"
        value = ""
      },
      {
        name  = "Lifetime Extension Threshold Percentage"
        value = "30"
      },
      {
        name  = "Mode for Synchronous RPC"
        value = "3"
      },
      {
        name  = "RPC Timeout"
        value = "500"
      },
      {
        name  = "Expand Scope Groups"
        value = "false"
      }
    ]
  }
  attribute_contract = {
    coreAttributes = []
    extended_attributes = [
      {
        name         = "extended_contract"
        multi_valued = true
      },
      {
        name         = "subject"
        multi_valued = false
      }
    ]
  }
  selection_settings = {
    resource_uris = []
  }
  access_control_settings = {
    restrict_clients = false
    allowedClients   = []
  }
  session_validation_settings = {
    check_valid_authn_session       = false
    check_session_revocation_status = false
    update_authn_session_activity   = false
    include_session_id              = false
  }
}

resource "pingfederate_oauth_access_token_mapping" "%[1]s" {
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.oauthAccessTokenManagerTest.id
  }

  context = {
    "type" = "CLIENT_CREDENTIALS"
  }

  attribute_contract_fulfillment = {
    "extended_contract" : {
      source = {
        type = "%[2]s"
      }
      value = "%[3]s"
    }
    "subject" = {
      source = {
        type = "TEXT"
      },
      value = "subject"
    }
  }

	%[4]s
}`, resourceName,
		resourceModel.attributeContractFulfillment["extended_contract"].Source.Type,
		resourceModel.attributeContractFulfillment["extended_contract"].Value,
		optionalHcl(resourceModel),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthAccessTokenMappingAttributes(config oauthAccessTokenMappingResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAccessTokenMapping"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAccessTokenMappingsAPI.GetMapping(ctx, oauthAccessTokenMappingId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		// attribute_contract_fulfillment
		attributeContractFullfillmentResp := response.AttributeContractFulfillment["extended_contract"]
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "type", config.attributeContractFulfillment["extended_contract"].Source.Type, attributeContractFullfillmentResp.Source.Type)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "value", config.attributeContractFulfillment["extended_contract"].Value, attributeContractFullfillmentResp.Value)
		if err != nil {
			return err
		}

		if config.attributeSource != nil {
			// attribute_source
			attributeSourceResp := response.AttributeSources[0]
			if attributeSourceResp.JdbcAttributeSource == nil {
				attributeSourceResp = response.AttributeSources[1]
			}
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "id", *config.attributeSource.Id, *attributeSourceResp.JdbcAttributeSource.Id)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "description", *config.attributeSource.Description, *attributeSourceResp.JdbcAttributeSource.Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "schema", *config.attributeSource.Schema, *attributeSourceResp.JdbcAttributeSource.Schema)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "table", config.attributeSource.Table, attributeSourceResp.JdbcAttributeSource.Table)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "filter", config.attributeSource.Filter, attributeSourceResp.JdbcAttributeSource.Filter)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(oauthAccessTokenMappingId), "column_names", config.attributeSource.ColumnNames, attributeSourceResp.JdbcAttributeSource.ColumnNames)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "id", config.attributeSource.DataStoreRef.Id, attributeSourceResp.JdbcAttributeSource.DataStoreRef.Id)
			if err != nil {
				return err
			}

		}

		if config.issuanceCriteria != nil {
			// issuance_criteria
			issuanceCriteriaResp := response.IssuanceCriteria
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "attribute_name", config.issuanceCriteria.ConditionalCriteria[0].AttributeName, issuanceCriteriaResp.ConditionalCriteria[0].AttributeName)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "condition", config.issuanceCriteria.ConditionalCriteria[0].Condition, issuanceCriteriaResp.ConditionalCriteria[0].Condition)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "source_type", config.issuanceCriteria.ConditionalCriteria[0].Source.Type, issuanceCriteriaResp.ConditionalCriteria[0].Source.Type)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, pointers.String(oauthAccessTokenMappingId), "value", config.issuanceCriteria.ConditionalCriteria[0].Value, issuanceCriteriaResp.ConditionalCriteria[0].Value)
			if err != nil {
				return err
			}

		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOauthAccessTokenMappingDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthAccessTokenMappingsAPI.DeleteMapping(ctx, oauthAccessTokenMappingId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAccessTokenMapping", oauthAccessTokenMappingId)
	}
	return nil
}
