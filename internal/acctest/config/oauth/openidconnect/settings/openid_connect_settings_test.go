// Copyright © 2026 Ping Identity Corporation

package oauthopenidconnectsettings_test

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

func TestAccOpenIdConnectSettings(t *testing.T) {
	resourceName := "myOpenIdConnectSettings"

	var steps []resource.TestStep
	if acctest.VersionAtLeast(version.PingFederate1210) {
		steps = testAccOpenIdConnectSettingsPf121(resourceName)
	} else {
		steps = testAccOpenIdConnectSettingsPrePf121(resourceName)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: steps,
	})
}

// Prior to PF 12.1 it isn't possible to delete the final OIDC policy from the server config,
// because it is always in use.
func testAccOpenIdConnectSettingsPrePf121(resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			// Set to the existing default
			Config: testAccOpenIdConnectSettingsExistingDefaultPolicy(resourceName, "acctestOidcPolicy"),
		},
		{
			// Test importing the resource
			Config:                               testAccOpenIdConnectSettingsExistingDefaultPolicy(resourceName, "acctestOidcPolicy"),
			ResourceName:                         "pingfederate_openid_connect_settings." + resourceName,
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "default_policy_ref.id",
		},
	}
}

func testAccOpenIdConnectSettingsPf121(resourceName string) []resource.TestStep {
	return []resource.TestStep{
		{
			// No policies configured and no default
			Config: testAccOpenIdConnectSettingsEmpty(resourceName),
		},
		{
			// Set a default policy
			Config: testAccOpenIdConnectSettingsBuildDefaultPolicy(resourceName),
		},
		{
			// Test importing the resource
			Config:                               testAccOpenIdConnectSettingsBuildDefaultPolicy(resourceName),
			ResourceName:                         "pingfederate_openid_connect_settings." + resourceName,
			ImportState:                          true,
			ImportStateVerify:                    true,
			ImportStateVerifyIdentifierAttribute: "default_policy_ref.id",
		},
		{
			// Reset back to no policies
			Config: testAccOpenIdConnectSettingsEmpty(resourceName),
		},
	}
}

func testAccOpenIdConnectSettingsEmpty(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_openid_connect_settings" "%s" {
}`, resourceName,
	)
}

func testAccOpenIdConnectSettingsBuildDefaultPolicy(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager" "%[1]satm" {
  manager_id = "%[1]satm"
  name       = "%[1]satm"
  plugin_descriptor_ref = {
    id = "com.pingidentity.pf.access.token.management.plugins.JwtBearerAccessTokenManagementPlugin"
  }
  configuration = {
    tables = [
      {
        name = "Symmetric Keys"
        rows = [
          {
            fields = [
              {
                name  = "Key ID"
                value = "keyidentifier"
              },
              {
                name  = "Key"
                value = "e1oDxOiC3Jboz3um8hBVmW3JRZNo9z7C0DMm/oj2V1gclQRcgi2gKM2DBj9N05G4"
              },
              {
                name  = "Encoding"
                value = "b64u"
              }
            ]
          }
        ]
      },
      {
        name = "Certificates"
        rows = []
      }
    ]
    fields = [
      {
        name  = "JWE Algorithm"
        value = "dir"
      },
      {
        name  = "JWE Content Encryption Algorithm"
        value = "A192CBC-HS384"
      },
      {
        name  = "Active Symmetric Encryption Key ID"
        value = "keyidentifier"
      },
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "contract"
      },
      {
        name         = "another"
        multi_valued = false
      }
    ]
  }
}

resource "pingfederate_openid_connect_policy" "%[1]s" {
  policy_id = "%[1]sPolicy"
  name      = "%[1]sPolicy"
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.%[1]satm.id
  }
  attribute_contract = {
    extended_attributes = []
  }
  attribute_mapping = {
    attribute_contract_fulfillment = {
      "sub" = {
        source = {
          type = "TOKEN"
        }
        value = "contract"
      }
    }
  }
}


resource "pingfederate_openid_connect_settings" "%[1]s" {
  default_policy_ref = {
    id = pingfederate_openid_connect_policy.%[1]s.policy_id
  }
}`, resourceName,
	)
}

func testAccOpenIdConnectSettingsExistingDefaultPolicy(resourceName, existingDefaultPolicyId string) string {
	return fmt.Sprintf(`
resource "pingfederate_openid_connect_settings" "%[1]s" {
  default_policy_ref = {
    id = "%[2]s"
  }
}`, resourceName, existingDefaultPolicyId)
}
