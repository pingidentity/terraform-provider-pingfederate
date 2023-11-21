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

// Attributes to test with. Add optional properties to test here if desired.
type serverSettingsSystemKeysResourceModel struct {
	currentEncryptedKeyData string
	pendingEncryptedKeyData string
}

func TestAccServerSettingsSystemKeys(t *testing.T) {
	resourceName := "myServerSettingsSystemKeys"
	initialResourceModel := serverSettingsSystemKeysResourceModel{
		currentEncryptedKeyData: "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMS4yLjUuMCIsInppcCI6IkRFRiJ9..4Q-LeikGMQ-5dVVRMMDyfw.JLR4Yg1FfmaTdOpVHZ1V1BypiguCuKawnJsUD33weL3nYRvyEPFgMCuBV72GC-HG.2b2T22iR040xI4ro-Iemeg",
		pendingEncryptedKeyData: "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMS4yLjUuMCIsInppcCI6IkRFRiJ9..J1yaOm2OdYCUDN402iIKPQ.LlpjecXwfHDiFJl_K6O57Mzp1RZxHN-TAbpKnypkRfeL1XgTHZrUkPgxO3ZcU7fb.q-X1zzd-de5svqDRbAE0lw",
	}
	updatedResourceModel := serverSettingsSystemKeysResourceModel{
		currentEncryptedKeyData: "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMS4yLjUuMCIsInppcCI6IkRFRiJ9..J1yaOm2OdYCUDN402iIKPQ.LlpjecXwfHDiFJl_K6O57Mzp1RZxHN-TAbpKnypkRfeL1XgTHZrUkPgxO3ZcU7fb.q-X1zzd-de5svqDRbAE0lw",
		pendingEncryptedKeyData: "eyJhbGciOiJkaXIiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2Iiwia2lkIjoiUWVzOVR5eTV5WiIsInZlcnNpb24iOiIxMS4yLjUuMCIsInppcCI6IkRFRiJ9..4Q-LeikGMQ-5dVVRMMDyfw.JLR4Yg1FfmaTdOpVHZ1V1BypiguCuKawnJsUD33weL3nYRvyEPFgMCuBV72GC-HG.2b2T22iR040xI4ro-Iemeg",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServerSettingsSystemKeys(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedServerSettingsSystemKeysAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccServerSettingsSystemKeys(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedServerSettingsSystemKeysAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccServerSettingsSystemKeys(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_server_settings_system_keys." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServerSettingsSystemKeys(resourceName string, resourceModel serverSettingsSystemKeysResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_server_settings_system_keys" "%[1]s" {
  current = {
    encrypted_key_data = "%[2]s"
  }
  pending = {
    encrypted_key_data = "%[3]s"
  }
}

data "pingfederate_server_settings_system_keys" "%[1]s" {
}`, resourceName,
		resourceModel.currentEncryptedKeyData,
		resourceModel.pendingEncryptedKeyData,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedServerSettingsSystemKeysAttributes(config serverSettingsSystemKeysResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ServerSettingsSystemKeys"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ServerSettingsAPI.GetSystemKeys(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		currentEncryptedKeyData := response.Current.EncryptedKeyData
		err = acctest.TestAttributesMatchString(resourceType, nil, "encrypted_key_data", config.currentEncryptedKeyData, *currentEncryptedKeyData)
		if err != nil {
			return err
		}
		pendingEncryptedKeyData := response.Pending.EncryptedKeyData
		err = acctest.TestAttributesMatchString(resourceType, nil, "encrypted_key_data", config.pendingEncryptedKeyData, *pendingEncryptedKeyData)
		if err != nil {
			return err
		}

		return nil
	}
}
