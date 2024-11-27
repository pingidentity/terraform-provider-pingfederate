package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccExtendedProperties(t *testing.T) {
	resourceName := "myExtendedProperties"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test minimal object sent
			{
				Config: testAccExtendedProperties_MinimalHCL(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pingfederate_extended_properties."+resourceName, "items.0.multi_valued", "false"),
					resource.TestCheckResourceAttr("pingfederate_extended_properties."+resourceName, "items.1.multi_valued", "false"),
				),
			},
			{
				// Test updating some fields
				Config: testAccExtendedProperties_CompleteHCL(resourceName),
			},
			{
				// Test importing the resource
				Config:                               testAccExtendedProperties_CompleteHCL(resourceName),
				ResourceName:                         "pingfederate_extended_properties." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "items.0.name",
			},
			// Test minimal object sent
			{
				Config: testAccExtendedProperties_MinimalHCL(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pingfederate_extended_properties."+resourceName, "items.0.multi_valued", "false"),
					resource.TestCheckResourceAttr("pingfederate_extended_properties."+resourceName, "items.1.multi_valued", "false"),
				),
			},
		},
	})
}

func testAccExtendedProperties_MinimalHCL(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_extended_properties" "%s" {
  items = [
  {
	  name = "authNexp",
	  description = "Authentication Experience [Single_Factor | Internal | ID-First | Multi_Factor]",
	},
	{
	  name = "useAuthnApi",
	  description = "Use the AuthN API",
	}
  ]
}`, resourceName)
}

func testAccExtendedProperties_CompleteHCL(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_extended_properties" "%s" {
  items = [
		{
	  name = "authNexp",
	  description = "Authentication Experience [Single_Factor | Internal | ID-First | Multi_Factor]",
	  multi_valued = false
	},
	{
	  name = "useAuthnApi",
	  description = "Use the AuthN API",
	  multi_valued = false
	}
  ]
}`, resourceName)
}
