// Copyright © 2026 Ping Identity Corporation

package oauthcibaserverpolicysettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccOauthCibaServerPolicySettings(t *testing.T) {
	resourceName := acctest.ResourceIdGen()

	var steps []resource.TestStep
	if acctest.VersionAtLeast(version.PingFederate1210) {
		steps = testAccOauthCibaServerPolicySettingsPf121(resourceName)
	} else {
		steps = testAccOauthCibaServerPolicySettingsPrePf121(resourceName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: steps,
	})
}

// Prior to PF 12.1 it isn't possible to delete the final ciba policy from the server config,
// because it is always in use.
func testAccOauthCibaServerPolicySettingsPrePf121(resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			// Set to the existing default
			Config: testAccOauthCibaServerPolicySettingsExistingDefaultPolicy(resourceName, "acctestCibaPolicy"),
		},
		{
			// Test importing the resource
			Config:                               testAccOauthCibaServerPolicySettingsExistingDefaultPolicy(resourceName, "acctestCibaPolicy"),
			ResourceName:                         "pingfederate_oauth_ciba_server_policy_settings." + resourceName,
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "default_request_policy_ref.id",
		},
	}
}

func testAccOauthCibaServerPolicySettingsPf121(resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			// No policies configured and no default
			Config: testAccOauthCibaServerPolicySettingsEmpty(resourceName),
		},
		{
			// Set a default policy
			Config: testAccOauthCibaServerPolicySettingsBuildDefaultPolicy(resourceName),
		},
		{
			// Test importing the resource
			Config:                               testAccOauthCibaServerPolicySettingsBuildDefaultPolicy(resourceName),
			ResourceName:                         "pingfederate_oauth_ciba_server_policy_settings." + resourceName,
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "default_request_policy_ref.id",
		},
		{
			// Reset back to no policies
			Config: testAccOauthCibaServerPolicySettingsEmpty(resourceName),
		},
	}
}

func testAccOauthCibaServerPolicySettingsEmpty(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_ciba_server_policy_settings" "%s" {
}`, resourceName)
}

func testAccOauthCibaServerPolicySettingsBuildDefaultPolicy(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_ciba_server_policy_request_policy" "%[1]s-policy" {
  allow_unsigned_login_hint_token = false
  alternative_login_hint_token_issuers = [
  ]
  authenticator_ref = {
    id = "exampleCibaAuthenticator"
  }
  identity_hint_contract = {
    extended_attributes = [
    ]
  }
  identity_hint_contract_fulfillment = {
    attribute_contract_fulfillment = {
      IDENTITY_HINT_SUBJECT = {
        source = {
          id   = null
          type = "REQUEST"
        }
        value = "IDENTITY_HINT_SUBJECT"
      }
    }
  }
  identity_hint_mapping = {
    attribute_contract_fulfillment = {
      USER_KEY = {
        source = {
          id   = null
          type = "NO_MAPPING"
        }
        value = null
      }
      subject = {
        source = {
          id   = null
          type = "NO_MAPPING"
        }
        value = null
      }
    }
  }
  name                            = "%[1]s-policy"
  policy_id                       = "%[1]s-policy"
  require_token_for_identity_hint = false
  transaction_lifetime            = 120
}

resource "pingfederate_oauth_ciba_server_policy_settings" "%[1]s" {
  default_request_policy_ref = {
    id = pingfederate_oauth_ciba_server_policy_request_policy.%[1]s-policy.id
  }
}`, resourceName)
}

func testAccOauthCibaServerPolicySettingsExistingDefaultPolicy(resourceName, existingPolicyName string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_ciba_server_policy_settings" "%[1]s" {
  default_request_policy_ref = {
    id = "%[2]s"
  }
}`, resourceName, existingPolicyName)
}
