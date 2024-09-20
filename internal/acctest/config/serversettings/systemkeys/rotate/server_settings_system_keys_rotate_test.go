package configurationencryptionkeysrotate_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

var (
	currentSystemKeyEncrypted  string
	previousSystemKeyEncrypted string
	pendingSystemKeyEncrypted  string
)

func TestAccConfigurationEncryptionKeysRotate(t *testing.T) {
	// Get the  keys currently on the server
	configurationEncryptionKeysRotate_getCurrentKeys(t)

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
				Check:  configurationEncryptionKeysRotate_checkExpectedKeys(t, true),
			},
			{
				// Expect no additional rotation
				PreConfig: func() {
					configurationEncryptionKeysRotate_getCurrentKeys(t)
				},
				Config: configurationEncryptionKeysRotate_FirstNoRotateHCL(),
				Check:  configurationEncryptionKeysRotate_checkExpectedKeys(t, false),
			},
			{
				// Expect rotation
				PreConfig: func() {
					configurationEncryptionKeysRotate_getCurrentKeys(t)
				},
				Config: configurationEncryptionKeysRotate_SecondRotateHCL(),
				Check:  configurationEncryptionKeysRotate_checkExpectedKeys(t, true),
			},
			{
				// Expect no additional rotation
				PreConfig: func() {
					configurationEncryptionKeysRotate_getCurrentKeys(t)
				},
				Config: configurationEncryptionKeysRotate_SecondNoRotateHCL(),
				Check:  configurationEncryptionKeysRotate_checkExpectedKeys(t, false),
			},
			{
				// Test importing the resource
				Config:                               configurationEncryptionKeysRotate_SecondRotateHCL(),
				ResourceName:                         "pingfederate_server_settings_system_keys_rotate.example",
				ImportStateVerifyIdentifierAttribute: "current.creation_date",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// The rotation trigger values are terraform-only, so they can't be imported
				ImportStateVerifyIgnore: []string{"rotation_trigger_values"},
			},
		},
	})
}

func configurationEncryptionKeysRotate_getCurrentKeys(t *testing.T) {
	// Get the keys currently on the server
	testClient := acctest.TestClient()
	keys, _, err := testClient.ServerSettingsAPI.GetSystemKeys(acctest.TestBasicAuthContext()).Execute()
	if err != nil {
		t.Fatal("An error occurred while getting the current system keys on the server: ", err.Error())
	}
	currentSystemKeyEncrypted = *keys.Current.EncryptedKeyData
	if keys.Previous != nil {
		previousSystemKeyEncrypted = *keys.Previous.EncryptedKeyData
	} else {
		previousSystemKeyEncrypted = ""
	}
	pendingSystemKeyEncrypted = *keys.Pending.EncryptedKeyData
}

func configurationEncryptionKeysRotate_checkExpectedKeys(t *testing.T, expectRotation bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Expect the previous set of keys
		err := resource.TestCheckResourceAttr("pingfederate_server_settings_system_keys_rotate.example", "current.encrypted_key_data", currentSystemKeyEncrypted)(s)
		if err != nil && !expectRotation {
			return err
		} else if err == nil && expectRotation {
			return errors.New("Expected the current system key to have rotated, but the key has not changed")
		}
		err = resource.TestCheckResourceAttr("pingfederate_server_settings_system_keys_rotate.example", "previous.encrypted_key_data", previousSystemKeyEncrypted)(s)
		if err != nil && !expectRotation {
			return err
		} else if err == nil && expectRotation {
			return errors.New("Expected the previous system key to have rotated, but the key has not changed")
		}
		err = resource.TestCheckResourceAttr("pingfederate_server_settings_system_keys_rotate.example", "pending.encrypted_key_data", pendingSystemKeyEncrypted)(s)
		if err != nil && !expectRotation {
			return err
		} else if err == nil && expectRotation {
			return errors.New("Expected the pending system key to have rotated, but the key has not changed")
		}
		return nil
	}
}

func configurationEncryptionKeysRotate_FirstRotateHCL() string {
	return `
resource "pingfederate_server_settings_system_keys_rotate" "example" {
}
`
}

// Ensure that adding triggers doesn't cause a rotation
func configurationEncryptionKeysRotate_FirstNoRotateHCL() string {
	return `
resource "pingfederate_server_settings_system_keys_rotate" "example" {
  rotation_trigger_values = {
    "trigger" = "initial"
  }
}
`
}

func configurationEncryptionKeysRotate_SecondRotateHCL() string {
	return `
resource "pingfederate_server_settings_system_keys_rotate" "example" {
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
resource "pingfederate_server_settings_system_keys_rotate" "example" {
  rotation_trigger_values = {
    "trigger" = "updated"
  }
}
`
}
