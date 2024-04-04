package sessionauthenticationsessionpoliciesglobal_test

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

// Attributes to test with. Add optional properties to test here if desired.
type sessionAuthenticationSessionPoliciesGlobalResourceModel struct {
	enableSessions             bool
	persistentSessions         *bool
	hashUniqueUserKeyAttribute *bool
	idleTimeoutMins            *int64
	idleTimeoutDisplayUnit     *string
	maxTimeoutMins             *int64
	maxTimeoutDisplayUnit      *string
}

func TestAccSessionAuthenticationSessionPoliciesGlobal(t *testing.T) {
	resourceName := "mySessionAuthenticationSessionPoliciesGlobal"
	initialResourceModel := sessionAuthenticationSessionPoliciesGlobalResourceModel{
		enableSessions: true,
	}
	updatedResourceModel := sessionAuthenticationSessionPoliciesGlobalResourceModel{
		enableSessions:             false,
		persistentSessions:         pointers.Bool(false),
		hashUniqueUserKeyAttribute: pointers.Bool(false),
		idleTimeoutMins:            pointers.Int64(120),
		idleTimeoutDisplayUnit:     pointers.String("HOURS"),
		maxTimeoutMins:             pointers.Int64(120),
		maxTimeoutDisplayUnit:      pointers.String("HOURS"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSessionAuthenticationSessionPoliciesGlobal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSessionAuthenticationSessionPoliciesGlobalAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSessionAuthenticationSessionPoliciesGlobal(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedSessionAuthenticationSessionPoliciesGlobalAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_authentication_session_policies_global.%s", resourceName), "enable_sessions", fmt.Sprintf("%t", updatedResourceModel.enableSessions)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_authentication_session_policies_global.%s", resourceName), "persistent_sessions", fmt.Sprintf("%t", *updatedResourceModel.persistentSessions)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_authentication_session_policies_global.%s", resourceName), "hash_unique_user_key_attribute", fmt.Sprintf("%t", *updatedResourceModel.hashUniqueUserKeyAttribute)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_authentication_session_policies_global.%s", resourceName), "idle_timeout_mins", fmt.Sprintf("%d", *updatedResourceModel.idleTimeoutMins)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_authentication_session_policies_global.%s", resourceName), "idle_timeout_display_unit", *updatedResourceModel.idleTimeoutDisplayUnit),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_authentication_session_policies_global.%s", resourceName), "max_timeout_mins", fmt.Sprintf("%d", *updatedResourceModel.maxTimeoutMins)),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_session_authentication_session_policies_global.%s", resourceName), "max_timeout_display_unit", *updatedResourceModel.maxTimeoutDisplayUnit),
				),
			},
			{
				// Test importing the resource
				Config:            testAccSessionAuthenticationSessionPoliciesGlobal(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_session_authentication_session_policies_global." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccSessionAuthenticationSessionPoliciesGlobal(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSessionAuthenticationSessionPoliciesGlobalAttributes(initialResourceModel),
			},
		},
	})
}

func testAccSessionAuthenticationSessionPoliciesGlobal(resourceName string, resourceModel sessionAuthenticationSessionPoliciesGlobalResourceModel) string {
	optionalHcl := ""
	// Just assuming that if the first one is set, the rest will be for this test
	if resourceModel.persistentSessions != nil {
		optionalHcl = fmt.Sprintf(`
		persistent_sessions            = %t
		hash_unique_user_key_attribute = %t
		idle_timeout_mins              = %d
		idle_timeout_display_unit      = "%s"
		max_timeout_mins               = %d
		max_timeout_display_unit       = "%s"
		`, *resourceModel.persistentSessions,
			*resourceModel.hashUniqueUserKeyAttribute,
			*resourceModel.idleTimeoutMins,
			*resourceModel.idleTimeoutDisplayUnit,
			*resourceModel.maxTimeoutMins,
			*resourceModel.maxTimeoutDisplayUnit)
	}

	return fmt.Sprintf(`
resource "pingfederate_session_authentication_session_policies_global" "%s" {
  enable_sessions = %t
  %s
}

data "pingfederate_session_authentication_session_policies_global" "%[1]s" {
  depends_on = [pingfederate_session_authentication_session_policies_global.%[1]s]
}`, resourceName,
		resourceModel.enableSessions,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSessionAuthenticationSessionPoliciesGlobalAttributes(config sessionAuthenticationSessionPoliciesGlobalResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SessionAuthenticationSessionPoliciesGlobal"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.SessionAPI.GetGlobalPolicy(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_sessions",
			config.enableSessions, response.EnableSessions)
		if err != nil {
			return err
		}

		if config.persistentSessions != nil {
			err = acctest.TestAttributesMatchBool(resourceType, nil, "persistent_sessions",
				*config.persistentSessions, *response.PersistentSessions)
			if err != nil {
				return err
			}
		}
		if config.hashUniqueUserKeyAttribute != nil {
			err = acctest.TestAttributesMatchBool(resourceType, nil, "hash_unique_user_key_attribute",
				*config.hashUniqueUserKeyAttribute, *response.HashUniqueUserKeyAttribute)
			if err != nil {
				return err
			}
		}
		if config.idleTimeoutMins != nil {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "idle_timeout_mins",
				*config.idleTimeoutMins, *response.IdleTimeoutMins)
			if err != nil {
				return err
			}
		}
		if config.idleTimeoutDisplayUnit != nil {
			err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "idle_timeout_display_unit",
				*config.idleTimeoutDisplayUnit, response.IdleTimeoutDisplayUnit)
			if err != nil {
				return err
			}
		}
		if config.maxTimeoutMins != nil {
			err = acctest.TestAttributesMatchInt(resourceType, nil, "max_timeout_mins",
				*config.maxTimeoutMins, *response.MaxTimeoutMins)
			if err != nil {
				return err
			}
		}
		if config.maxTimeoutDisplayUnit != nil {
			err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "max_timeout_display_unit",
				*config.maxTimeoutDisplayUnit, response.MaxTimeoutDisplayUnit)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
