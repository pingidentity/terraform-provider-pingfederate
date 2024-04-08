package oauthaccesstokenmanager_test

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

const internallyManagedReferenceOauthAccessTokenManagerId = "internallyManagedReferenceOatm"
const internallyManagedReferenceOauthAccessTokenManagerName = "internallyManagedReferenceExample"

// Attributes to test with. Add optional properties to test here if desired.
type internallyManagedReferenceOauthAccessTokenManagerResourceModel struct {
	id                           string
	name                         string
	tokenLength                  *string
	tokenLifetime                *string
	checkSessionRevocationStatus *bool
}

func TestAccInternallyManagedReferenceOauthAccessTokenManager(t *testing.T) {
	resourceName := "myInternallyManagedReferenceOauthAccessTokenManager"
	initialResourceModel := internallyManagedReferenceOauthAccessTokenManagerResourceModel{
		id:   internallyManagedReferenceOauthAccessTokenManagerId,
		name: internallyManagedReferenceOauthAccessTokenManagerName,
	}
	updatedResourceModel := internallyManagedReferenceOauthAccessTokenManagerResourceModel{
		id:                           internallyManagedReferenceOauthAccessTokenManagerId,
		name:                         internallyManagedReferenceOauthAccessTokenManagerName,
		tokenLength:                  pointers.String("56"),
		tokenLifetime:                pointers.String("240"),
		checkSessionRevocationStatus: pointers.Bool(true),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckInternallyManagedReferenceOauthAccessTokenManagerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternallyManagedReferenceOauthAccessTokenManagerMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagerAttributes(initialResourceModel, true),
			},
			{
				// Test updating some fields
				Config: testAccInternallyManagedReferenceOauthAccessTokenManager(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagerAttributes(updatedResourceModel, false),
			},
			{
				// Test importing the resource
				Config:                  testAccInternallyManagedReferenceOauthAccessTokenManager(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_oauth_access_token_manager." + resourceName,
				ImportStateId:           internallyManagedReferenceOauthAccessTokenManagerId,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"configuration.fields.value"},
			},
			{
				// Back to minimal model
				Config: testAccInternallyManagedReferenceOauthAccessTokenManagerMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagerAttributes(initialResourceModel, true),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OauthAccessTokenManagersAPI.DeleteTokenManager(ctx, updatedResourceModel.id).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccInternallyManagedReferenceOauthAccessTokenManagerMinimal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagerAttributes(initialResourceModel, true),
			},
		},
	})
}

func testAccInternallyManagedReferenceOauthAccessTokenManagerMinimal(resourceName string, resourceModel internallyManagedReferenceOauthAccessTokenManagerResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = "%[2]s"
  name       = "%[3]s"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
  }
  attribute_contract = {
    coreAttributes = []
    extended_attributes = [
      {
        name         = "extended_contract"
        multi_valued = true
      }
    ]
  }
}

data "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = pingfederate_oauth_access_token_manager.%[1]s.id
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
	)
}

func testAccInternallyManagedReferenceOauthAccessTokenManager(resourceName string, resourceModel internallyManagedReferenceOauthAccessTokenManagerResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = "%[2]s"
  name       = "%[3]s"
  plugin_descriptor_ref = {
    id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = []
    fields = [
      {
        name  = "Token Length"
        value = "%[4]s"
      },
      {
        name  = "Token Lifetime"
        value = "%[5]s"
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
    check_session_revocation_status = %[6]t
    update_authn_session_activity   = false
    include_session_id              = false
  }
}

data "pingfederate_oauth_access_token_manager" "%[1]s" {
  manager_id = pingfederate_oauth_access_token_manager.%[1]s.id
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		*resourceModel.tokenLength,
		*resourceModel.tokenLifetime,
		*resourceModel.checkSessionRevocationStatus,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagerAttributes(config internallyManagedReferenceOauthAccessTokenManagerResourceModel, minimal bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAccessTokenManager"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAccessTokenManagersAPI.GetTokenManager(ctx, internallyManagedReferenceOauthAccessTokenManagerId).Execute()

		if err != nil {
			return err
		}

		// Check for the always-defined extended attribute
		for _, extendedAttr := range response.AttributeContract.ExtendedAttributes {
			err = acctest.TestAttributesMatchString(resourceType, &config.id, "extended_attribute.name", "extended_contract", extendedAttr.Name)
			if err != nil {
				return err
			}
		}

		// When checking the minimal model, there's nothing else to verify
		if minimal {
			return nil
		}

		// Verify that attributes have expected values
		getFields := response.Configuration.Fields
		for _, field := range getFields {
			if field.Name == "Token Length" {
				err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", *config.tokenLength, *field.Value)
				if err != nil {
					return err
				}
			}
			if field.Name == "Token Lifetime" {
				err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", *config.tokenLifetime, *field.Value)
				if err != nil {
					return err
				}
			}
		}

		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "check_session_revocation_status", *config.checkSessionRevocationStatus, *response.SessionValidationSettings.CheckSessionRevocationStatus)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckInternallyManagedReferenceOauthAccessTokenManagerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthAccessTokenManagersAPI.DeleteTokenManager(ctx, internallyManagedReferenceOauthAccessTokenManagerId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAccessTokenManager", internallyManagedReferenceOauthAccessTokenManagerId)
	}
	return nil
}
