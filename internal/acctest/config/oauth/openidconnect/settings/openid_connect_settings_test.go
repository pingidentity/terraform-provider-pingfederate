// Copyright © 2025 Ping Identity Corporation

package oauthopenidconnectsettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccOpenIdConnectSettings(t *testing.T) {
	resourceName := "myOpenIdConnectSettings"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOpenIdConnectSettings(resourceName),
			},
			{
				// Test importing the resource
				Config:       testAccOpenIdConnectSettings(resourceName),
				ResourceName: "pingfederate_openid_connect_settings." + resourceName,
				ImportState:  true,
			},
		},
	})
}

func testAccOpenIdConnectSettings(resourceName string) string {
	// The dependent OIDC policy is not created in this test because prior to PF 12.1 it isn't possible to delete
	// the final OIDC policy from the server config, because it is always in use. This would also interfere
	// with other tests that make changes to OIDC policies or access token managers.
	//TODO update this test once 12.1 is the oldest supported version, and create the OIDC policy in this test.
	return fmt.Sprintf(`
resource "pingfederate_openid_connect_settings" "%s" {
}`, resourceName,
	)
}
