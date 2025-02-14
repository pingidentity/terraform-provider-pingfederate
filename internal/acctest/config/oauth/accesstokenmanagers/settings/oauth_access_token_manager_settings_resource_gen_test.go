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
)

func TestAccOauthAccessTokenManagerSettings_MinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: oauthAccessTokenManagerSettings_Empty(),
			},
			{
				// Create the resource with a minimal model
				Config: oauthAccessTokenManagerSettings_WithDefault(),
			},
			{
				// Test importing the resource
				Config:                               oauthAccessTokenManagerSettings_WithDefault(),
				ResourceName:                         "pingfederate_oauth_access_token_manager_settings.example",
				ImportStateVerifyIdentifierAttribute: "default_access_token_manager_ref.id",
				ImportState:                          true,
				ImportStateVerify:                    true,
			},
			{
				// Reset to the original default access token manager ref
				Config: oauthAccessTokenManagerSettings_Empty(),
			},
		},
	})
}

// Minimal HCL with only required values set
func oauthAccessTokenManagerSettings_Empty() string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_manager_settings" "example" {
}
`)
}

// Minimal HCL with only required values set
func oauthAccessTokenManagerSettings_WithDefault() string {
	atmName := "tokenManagerSettingsTestAtm"
	return fmt.Sprintf(`
%s

resource "pingfederate_oauth_access_token_manager_settings" "example" {
  default_access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.%s.id
  }
}
`, accesstokenmanager.TestAccessTokenManagerHCL(atmName), atmName)
}
