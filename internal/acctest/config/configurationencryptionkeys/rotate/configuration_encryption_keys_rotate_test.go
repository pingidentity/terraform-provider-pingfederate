// Copyright © 2026 Ping Identity Corporation

package configurationencryptionkeysrotate_test

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccConfigurationEncryptionKeysRotate(t *testing.T) {
	// Get the number of keys currently on the server
	testClient := acctest.TestClient()
	keys, _, err := testClient.ConfigurationEncryptionKeysAPI.GetConfigurationEncryptionKeys(acctest.TestBasicAuthContext()).Execute()
	if err != nil {
		t.Fatal("An error occurred while checking the number of encryption keys on the server: ", err.Error())
	}
	numStartingKeys := len(keys.Items)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Initial rotation on create
				Config: configurationEncryptionKeysRotate_FirstRotateHCL(),
				Check:  resource.TestCheckResourceAttr("pingfederate_configuration_encryption_keys_rotate.example", "keys.#", strconv.FormatInt(int64(numStartingKeys+1), 10)),
			},
			{
				// Expect no additional rotation
				Config: configurationEncryptionKeysRotate_FirstNoRotateHCL(),
				Check:  resource.TestCheckResourceAttr("pingfederate_configuration_encryption_keys_rotate.example", "keys.#", strconv.FormatInt(int64(numStartingKeys+1), 10)),
			},
			{
				// Expect rotation
				Config: configurationEncryptionKeysRotate_SecondRotateHCL(),
				Check:  resource.TestCheckResourceAttr("pingfederate_configuration_encryption_keys_rotate.example", "keys.#", strconv.FormatInt(int64(numStartingKeys+2), 10)),
			},
			{
				// Expect no additional rotation
				Config: configurationEncryptionKeysRotate_SecondNoRotateHCL(),
				Check:  resource.TestCheckResourceAttr("pingfederate_configuration_encryption_keys_rotate.example", "keys.#", strconv.FormatInt(int64(numStartingKeys+2), 10)),
			},
			{
				// Test importing the resource
				Config:                               configurationEncryptionKeysRotate_SecondRotateHCL(),
				ResourceName:                         "pingfederate_configuration_encryption_keys_rotate.example",
				ImportStateVerifyIdentifierAttribute: "keys.#",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// The rotation trigger values are terraform-only, so they can't be imported
				ImportStateVerifyIgnore: []string{"rotation_trigger_values"},
			},
		},
	})
}

func configurationEncryptionKeysRotate_FirstRotateHCL() string {
	return `
resource "pingfederate_configuration_encryption_keys_rotate" "example" {
}
`
}

// Ensure that adding triggers doesn't cause a rotation
func configurationEncryptionKeysRotate_FirstNoRotateHCL() string {
	return `
resource "pingfederate_configuration_encryption_keys_rotate" "example" {
  rotation_trigger_values = {
    "trigger" = "initial"
  }
}
`
}

func configurationEncryptionKeysRotate_SecondRotateHCL() string {
	return `
resource "pingfederate_configuration_encryption_keys_rotate" "example" {
  rotation_trigger_values = {
    "trigger"    = "updated"
    "newtrigger" = "new"
  }
}
`
}

// Ensure that removing triggers doesn't cause a rotation
func configurationEncryptionKeysRotate_SecondNoRotateHCL() string {
	return `
resource "pingfederate_configuration_encryption_keys_rotate" "example" {
  rotation_trigger_values = {
    "trigger" = "updated"
  }
}
`
}
