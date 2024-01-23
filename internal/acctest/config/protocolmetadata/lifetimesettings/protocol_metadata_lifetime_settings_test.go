package protocolmetadatalifetimesettings_test

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

// Attributes to test with. Add optional properties to test here if desired.
type protocolMetadataLifetimeSettingsResourceModel struct {
	cacheDuration int64
	reloadDelay   int64
}

func TestAccProtocolMetadataLifetimeSettings(t *testing.T) {
	resourceName := "myProtocolMetadataLifetimeSettings"
	updatedResourceModel := protocolMetadataLifetimeSettingsResourceModel{
		cacheDuration: 1440,
		reloadDelay:   1440,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProtocolMetadataLifetimeSettings(resourceName, nil),
				Check:  testAccCheckExpectedProtocolMetadataLifetimeSettingsAttributes(nil),
			},
			{
				// Test updating some fields
				Config: testAccProtocolMetadataLifetimeSettings(resourceName, &updatedResourceModel),
				Check:  testAccCheckExpectedProtocolMetadataLifetimeSettingsAttributes(&updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccProtocolMetadataLifetimeSettings(resourceName, &updatedResourceModel),
				ResourceName:      "pingfederate_protocol_metadata_lifetime_settings." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccProtocolMetadataLifetimeSettings(resourceName, nil),
				Check:  testAccCheckExpectedProtocolMetadataLifetimeSettingsAttributes(nil),
			},
		},
	})
}

func testAccProtocolMetadataLifetimeSettings(resourceName string, resourceModel *protocolMetadataLifetimeSettingsResourceModel) string {
	optionalHcl := ""
	if resourceModel != nil {
		optionalHcl = fmt.Sprintf(`
		cache_duration = %d
		reload_delay   = %d
		`,
			resourceModel.cacheDuration,
			resourceModel.reloadDelay)
	}

	return fmt.Sprintf(`
resource "pingfederate_protocol_metadata_lifetime_settings" "%s" {
  %s
}
data "pingfederate_protocol_metadata_lifetime_settings" "%[1]s" {
  depends_on = [pingfederate_protocol_metadata_lifetime_settings.%[1]s]
}`, resourceName,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedProtocolMetadataLifetimeSettingsAttributes(config *protocolMetadataLifetimeSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ProtocolMetadataLifetimeSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ProtocolMetadataAPI.GetLifetimeSettings(ctx).Execute()

		if err != nil {
			return err
		}

		if config == nil {
			return nil
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, nil, "cache_duration",
			config.cacheDuration, *response.CacheDuration)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, nil, "reload_delay",
			config.reloadDelay, *response.ReloadDelay)
		if err != nil {
			return err
		}

		return nil
	}
}
