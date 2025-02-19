// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package oauthaccesstokenmanagerssettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/accesstokenmanager"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccOauthAccessTokenManagerSettings_MinimalMaximal(t *testing.T) {
	var steps []resource.TestStep
	if acctest.VersionAtLeast(version.PingFederate1210) {
		steps = testAccOauthAccessTokenManagerSettingsPf121()
	} else {
		steps = testAccOauthAccessTokenManagerSettingsPrePf121()
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: steps,
	})
}

// Prior to PF 12.1 we have to create an access token manager for use in other tests in the data.json,
// so we can't set the default ref to null in this test.
func testAccOauthAccessTokenManagerSettingsPrePf121() []resource.TestStep {
	return []resource.TestStep{
		{
			// Set to the existing default
			Config: oauthAccessTokenManagerSettings_ExistingDefaultAtm("acctestAtm"),
		},
		{
			// Test importing the resource
			Config:                               oauthAccessTokenManagerSettings_ExistingDefaultAtm("acctestAtm"),
			ResourceName:                         "pingfederate_oauth_access_token_manager_settings.example",
			ImportStateVerifyIdentifierAttribute: "default_access_token_manager_ref.id",
			ImportState:                          true,
			ImportStateVerify:                    true,
		},
	}
}

func testAccOauthAccessTokenManagerSettingsPf121() []resource.TestStep {
	return []resource.TestStep{
		{
			// No atms configured and no default
			Config: oauthAccessTokenManagerSettings_Empty(),
		},
		{
			// Set a default atm
			Config: oauthAccessTokenManagerSettings_BuildDefaultAtm(),
		},
		{
			// Test importing the resource
			Config:                               oauthAccessTokenManagerSettings_BuildDefaultAtm(),
			ResourceName:                         "pingfederate_oauth_access_token_manager_settings.example",
			ImportStateVerifyIdentifierAttribute: "default_access_token_manager_ref.id",
			ImportState:                          true,
			ImportStateVerify:                    true,
		},
		{
			// Reset back to no atms
			Config: oauthAccessTokenManagerSettings_Empty(),
		},
	}
}

// Minimal HCL with only required values set
func oauthAccessTokenManagerSettings_Empty() string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager_settings" "example" {
}
`)
}

// Minimal HCL with only required values set
func oauthAccessTokenManagerSettings_BuildDefaultAtm() string {
	atmName := "tokenManagerSettingsTestAtm"
	return fmt.Sprintf(`
%s

resource "pingfederate_oauth_access_token_manager_settings" "example" {
  default_access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.%s.id
  }
}
`, accesstokenmanager.AccessTokenManagerTestHCL(atmName), atmName)
}

func oauthAccessTokenManagerSettings_ExistingDefaultAtm(existingAtmName string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager_settings" "example" {
  default_access_token_manager_ref = {
    id = "%s"
  }
}`, existingAtmName)
}
