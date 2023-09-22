package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const internallyManagedReferenceOauthAccessTokenManagersId = "internallyManagedReferenceOatm"
const internallyManagedReferenceOauthAccessTokenManagersName = "internallyManagedReferenceExample"

// Attributes to test with. Add optional properties to test here if desired.
type internallyManagedReferenceOauthAccessTokenManagersResourceModel struct {
	id            string
	name          string
	tokenLength   string
	tokenLifetime string
}

func TestAccInternallyManagedReferenceOauthAccessTokenManagers(t *testing.T) {
	resourceName := "myInternallyManagedReferenceOauthAccessTokenManagers"
	initialResourceModel := internallyManagedReferenceOauthAccessTokenManagersResourceModel{
		id:            internallyManagedReferenceOauthAccessTokenManagersId,
		name:          internallyManagedReferenceOauthAccessTokenManagersName,
		tokenLength:   "28",
		tokenLifetime: "120",
	}
	updatedResourceModel := internallyManagedReferenceOauthAccessTokenManagersResourceModel{
		id:            internallyManagedReferenceOauthAccessTokenManagersId,
		name:          internallyManagedReferenceOauthAccessTokenManagersName,
		tokenLength:   "56",
		tokenLifetime: "240",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckInternallyManagedReferenceOauthAccessTokenManagersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternallyManagedReferenceOauthAccessTokenManagers(resourceName, initialResourceModel),
				// Check:  testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagersAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccInternallyManagedReferenceOauthAccessTokenManagers(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagersAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccInternallyManagedReferenceOauthAccessTokenManagers(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_access_token_managers." + resourceName,
				ImportStateId:     internallyManagedReferenceOauthAccessTokenManagersId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInternallyManagedReferenceOauthAccessTokenManagers(resourceName string, resourceModel internallyManagedReferenceOauthAccessTokenManagersResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_managers" "%[1]s" {
  id = "%[2]s"
  name = "%[3]s"
  plugin_descriptor_ref = {
	id = "org.sourceid.oauth20.token.plugin.impl.ReferenceBearerAccessTokenManagementPlugin"
  }
  configuration = {
	tables = []
	fields = [
	  {
	    name = "Token Length"
	    value = "%[4]s"
	  },
	  {
	    name = "Token Lifetime"
	    value = "%[5]s"
	  },
	  {
	    name = "Lifetime Extension Policy"
	    value = "NONE"
	  },
	  {
	    name = "Maximum Token Lifetime"
	    value = ""
	  },
	  {
	    name = "Lifetime Extension Threshold Percentage"
	    value = "30"
	  },
	  {
	    name = "Mode for Synchronous RPC"
	    value = "3"
	  },
	  {
	    name = "RPC Timeout"
	    value = "500"
	  },
	  {
	    name = "Expand Scope Groups"
	    value = "false"
	  }
	]
  }
  attribute_contract = {
	coreAttributes = []
	extended_attributes = [
	  {
		name = "extended_contract"
		multi_valued = true
	  }
	]
  }
  selection_settings = {
	resource_uris = []
  }
  access_control_settings = {
	restrict_clients = false
	allowedClients = []
  }
  session_validation_settings = {
	check_valid_authn_session = false
	check_session_revocation_status = false
	update_authn_session_activity = false
	include_session_id = false
  }
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.tokenLength,
		resourceModel.tokenLifetime,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedInternallyManagedReferenceOauthAccessTokenManagersAttributes(config internallyManagedReferenceOauthAccessTokenManagersResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAccessTokenManagers"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthAccessTokenManagersApi.GetTokenManager(ctx, internallyManagedReferenceOauthAccessTokenManagersId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.name, response.Name)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckInternallyManagedReferenceOauthAccessTokenManagersDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthAccessTokenManagersApi.DeleteTokenManager(ctx, internallyManagedReferenceOauthAccessTokenManagersId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAccessTokenManagers", internallyManagedReferenceOauthAccessTokenManagersId)
	}
	return nil
}
