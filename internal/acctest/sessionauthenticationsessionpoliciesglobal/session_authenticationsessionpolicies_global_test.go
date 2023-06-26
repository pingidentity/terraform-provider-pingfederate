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

const sessionAuthenticationSessionPoliciesGlobalId = "id"

// Attributes to test with. Add optional properties to test here if desired.
type sessionAuthenticationSessionPoliciesGlobalResourceModel struct {
	id                         string
	enableSessions             bool
	persistentSessions         bool
	hashUniqueUserKeyAttribute bool
	idleTimeoutMins            int64
	idleTimeoutDisplayUnit     string
	maxTimeoutMins             int64
	maxTimeoutDisplayUnit      string
}

func TestAccSessionAuthenticationSessionPoliciesGlobal(t *testing.T) {
	resourceName := "mySessionAuthenticationSessionPoliciesGlobal"
	initialResourceModel := sessionAuthenticationSessionPoliciesGlobalResourceModel{
		enableSessions:             true,
		persistentSessions:         true,
		hashUniqueUserKeyAttribute: true,
		idleTimeoutMins:            60,
		idleTimeoutDisplayUnit:     "MINUTES",
		maxTimeoutMins:             60,
		maxTimeoutDisplayUnit:      "MINUTES",
	}
	updatedResourceModel := sessionAuthenticationSessionPoliciesGlobalResourceModel{
		enableSessions:             false,
		persistentSessions:         false,
		hashUniqueUserKeyAttribute: false,
		idleTimeoutMins:            120,
		idleTimeoutDisplayUnit:     "HOURS",
		maxTimeoutMins:             120,
		maxTimeoutDisplayUnit:      "HOURS",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSessionAuthenticationSessionPoliciesGlobal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSessionAuthenticationSessionPoliciesGlobalAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSessionAuthenticationSessionPoliciesGlobal(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSessionAuthenticationSessionPoliciesGlobalAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccSessionAuthenticationSessionPoliciesGlobal(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_session_authenticationsessionpolicies_global." + resourceName,
				ImportStateId:     sessionAuthenticationSessionPoliciesGlobalId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSessionAuthenticationSessionPoliciesGlobal(resourceName string, resourceModel sessionAuthenticationSessionPoliciesGlobalResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_session_authenticationsessionpolicies_global" "%[1]s" {
  enable_sessions                = %[2]t
  persistent_sessions            = %[3]t
  hash_unique_user_key_attribute = %[4]t
  idle_timeout_mins              = %[5]d
  idle_timeout_display_unit      = "%[6]s"
  max_timeout_mins               = %[7]d
  max_timeout_display_unit       = "%[8]s"
}`, resourceName,
		resourceModel.enableSessions,
		resourceModel.persistentSessions,
		resourceModel.hashUniqueUserKeyAttribute,
		resourceModel.idleTimeoutMins,
		resourceModel.idleTimeoutDisplayUnit,
		resourceModel.maxTimeoutMins,
		resourceModel.maxTimeoutDisplayUnit,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSessionAuthenticationSessionPoliciesGlobalAttributes(config sessionAuthenticationSessionPoliciesGlobalResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SessionAuthenticationSessionPoliciesGlobal"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SessionApi.GetGlobalPolicy(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "enable_sessions",
			config.enableSessions, response.EnableSessions)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "persistent_sessions",
			config.persistentSessions, *response.PersistentSessions)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "hash_unique_user_key_attribute",
			config.hashUniqueUserKeyAttribute, *response.HashUniqueUserKeyAttribute)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "idle_timeout_mins",
			config.idleTimeoutMins, *response.IdleTimeoutMins)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "idle_timeout_display_unit",
			config.idleTimeoutDisplayUnit, *response.IdleTimeoutDisplayUnit)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "max_timeout_mins",
			config.maxTimeoutMins, *response.MaxTimeoutMins)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "max_timeout_display_unit",
			config.maxTimeoutDisplayUnit, *response.MaxTimeoutDisplayUnit)
		if err != nil {
			return err
		}

		return nil
	}
}
