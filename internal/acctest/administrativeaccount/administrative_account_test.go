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

const username = "username"
const password = "2FederateM0re!"

var roles = []string{"USER_ADMINISTRATOR", "CRYPTO_ADMINISTRATOR"}

type administrativeAccountResourceModel struct {
	active       bool
	description  string
	stateId      string
	emailAddress string
}

func TestAccAdministrativeAccount(t *testing.T) {
	resourceName := "myAdministrativeAccount"
	initialResourceModel := administrativeAccountResourceModel{
		active:       false,
		description:  "example description",
		stateId:      username,
		emailAddress: "firstemail@example.com",
	}
	updatedResourceModel := administrativeAccountResourceModel{
		active:       true,
		description:  "updated description",
		stateId:      username,
		emailAddress: "secondemail@example.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAdministrativeAccount(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAdministrativeAccountAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAdministrativeAccount(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAdministrativeAccountAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccAdministrativeAccount(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_administrative_account." + resourceName,
				ImportStateId:           initialResourceModel.stateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccAdministrativeAccount(resourceName string, resourceModel administrativeAccountResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_administrative_account" "%[1]s" {
  active      = %[2]t
  description = "%[3]s"
  roles       = %[4]s
  password    = "%[5]s"
  username    = "%[6]s"
  email_address = "%[7]s"
}

data "pingfederate_administrative_account" "%[1]s" {
  id = pingfederate_administrative_account.%[1]s.username
}`, resourceName,
		resourceModel.active,
		resourceModel.description,
		acctest.StringSliceToTerraformString(roles),
		password,
		username,
		resourceModel.emailAddress,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAdministrativeAccountAttributes(config administrativeAccountResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AdministrativeAccount"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AdministrativeAccountsAPI.GetAccount(ctx, config.stateId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, &config.stateId, "active",
			config.active, *response.Active)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.stateId, "description",
			config.description, response.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringPointer(resourceType, &config.stateId, "email_address",
			config.emailAddress, response.EmailAddress)
		if err != nil {
			return err
		}

		return nil
	}
}
