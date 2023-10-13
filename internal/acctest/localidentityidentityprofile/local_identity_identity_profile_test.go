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

const localIdentityIdentityProfilesId = "test"

// Attributes to test with. Add optional properties to test here if desired.
type localIdentityIdentityProfilesResourceModel struct {
	id                  string
	name                string
	registrationEnabled bool
	profileEnabled      bool
}

func TestAccLocalIdentityIdentityProfiles(t *testing.T) {
	resourceName := "myLocalIdentityIdentityProfiles"
	initialResourceModel := localIdentityIdentityProfilesResourceModel{
		// Test is only run on attributes that do not require a PD dataStore.
		id:                  localIdentityIdentityProfilesId,
		name:                "example",
		registrationEnabled: false,
		profileEnabled:      false,
	}
	updatedResourceModel := localIdentityIdentityProfilesResourceModel{
		id:                  localIdentityIdentityProfilesId,
		name:                "example1",
		registrationEnabled: false,
		profileEnabled:      false,
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckLocalIdentityIdentityProfilesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLocalIdentityIdentityProfiles(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedLocalIdentityIdentityProfilesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccLocalIdentityIdentityProfiles(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedLocalIdentityIdentityProfilesAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccLocalIdentityIdentityProfiles(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_local_identity_identity_profile." + resourceName,
				ImportStateId:     localIdentityIdentityProfilesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLocalIdentityIdentityProfiles(resourceName string, resourceModel localIdentityIdentityProfilesResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractsExample" {
  custom_id           = "%[2]s"
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
}

resource "pingfederate_local_identity_identity_profile" "%[1]s" {
  custom_id = "%[2]s"
  name      = "%[3]s"
  apc_id = {
    id = pingfederate_authentication_policy_contract.authenticationPolicyContractsExample.id
  }
  registration_enabled = %[4]t
  profile_enabled      = %[5]t

}

data "pingfederate_local_identity_identity_profile" "%[1]s" {
  id = pingfederate_local_identity_identity_profile.%[1]s.id
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.registrationEnabled,
		resourceModel.profileEnabled,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedLocalIdentityIdentityProfilesAttributes(config localIdentityIdentityProfilesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "LocalIdentityIdentityProfiles"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.LocalIdentityIdentityProfilesAPI.GetIdentityProfile(ctx, localIdentityIdentityProfilesId).Execute()

		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, &config.id, "id",
			config.id, *response.Id)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "registration_enabled",
			config.registrationEnabled, *response.RegistrationEnabled)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchBool(resourceType, &config.id, "profile_enabled",
			config.profileEnabled, *response.ProfileEnabled)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckLocalIdentityIdentityProfilesDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.LocalIdentityIdentityProfilesAPI.DeleteIdentityProfile(ctx, localIdentityIdentityProfilesId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("LocalIdentityIdentityProfiles", localIdentityIdentityProfilesId)
	}
	return nil
}
