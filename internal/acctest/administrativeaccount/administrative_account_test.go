package acctest_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const username = "username"

var password = "2FederateM0re!"

type administrativeAccountResourceModel struct {
	administrativeAccount *client.AdministrativeAccount
}

func initialAdministrativeAccount() *client.AdministrativeAccount {
	initialAdministrativeAccount := client.NewAdministrativeAccountWithDefaults()
	initialAdministrativeAccount.Username = username
	initialAdministrativeAccount.Password = &password
	return initialAdministrativeAccount
}

func updateAdministrativeAccount(encryptedPassword string) *client.AdministrativeAccount {
	updateAdministrativeAccount := client.NewAdministrativeAccountWithDefaults()
	updateAdministrativeAccount.Username = username
	updateAdministrativeAccount.EncryptedPassword = &encryptedPassword
	updateAdministrativeAccount.Active = pointers.Bool(false)
	updateAdministrativeAccount.Description = pointers.String("updated description")
	updateAdministrativeAccount.Department = pointers.String("department")
	updateAdministrativeAccount.EmailAddress = pointers.String("test@example.com")
	updateAdministrativeAccount.PhoneNumber = pointers.String("555-555-5555")
	updateAdministrativeAccount.Roles = []string{"USER_ADMINISTRATOR", "CRYPTO_ADMINISTRATOR"}
	return updateAdministrativeAccount
}

func TestAccAdministrativeAccount(t *testing.T) {
	resourceName := "myAdministrativeAccount"
	initialResourceModel := administrativeAccountResourceModel{
		administrativeAccount: initialAdministrativeAccount(),
	}

	initialTest, encryptedPassword := testAccCheckExpectedAdministrativeAccountAttributes(initialResourceModel)
	updatedResourceModel := administrativeAccountResourceModel{
		administrativeAccount: updateAdministrativeAccount(encryptedPassword),
	}
	updatedTest, _ := testAccCheckExpectedAdministrativeAccountAttributes(updatedResourceModel)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAdministrativeAccount(resourceName, initialResourceModel),
				Check:  initialTest,
			},
			{
				// Test updating some fields
				Config: testAccAdministrativeAccount(resourceName, updatedResourceModel),
				Check:  updatedTest,
			},
			{
				// Test importing the resource
				Config:            testAccAdministrativeAccount(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_administrative_account." + resourceName,
				ImportStateId:     initialResourceModel.administrativeAccount.Username,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAdministrativeAccount(resourceName, initialResourceModel),
				Check:  initialTest,
			},
		},
	})
}

func hcl(aa *client.AdministrativeAccount) string {
	var builder strings.Builder
	if aa == nil {
		return ""
	}
	if aa != nil {
		tf := `
		%[1]s
		%[2]s
		%[3]s
		%[4]s
		%[5]s
		%[6]s
		%[7]s
		%[8]s
		%[9]s
		`
		builder.WriteString(
			fmt.Sprintf(tf,
				acctest.TfKeyValuePairToString("active", strconv.FormatBool(*aa.Active), true),
				acctest.TfKeyValuePairToString("description", *aa.Description, true),
				acctest.TfKeyValuePairToString("department", *aa.Department, true),
				acctest.TfKeyValuePairToString("email_address", *aa.EmailAddress, true),
				acctest.TfKeyValuePairToString("encrypted_password", *aa.EncryptedPassword, true),
				acctest.TfKeyValuePairToString("roles", acctest.StringSliceToTerraformString(aa.Roles), false),
				acctest.TfKeyValuePairToString("password", *aa.Password, true),
				acctest.TfKeyValuePairToString("phone_number", *aa.PhoneNumber, true),
				acctest.TfKeyValuePairToString("username", aa.Username, true),
			),
		)
	}
	fmt.Print(builder.String())
	return builder.String()
}

func testAccAdministrativeAccount(resourceName string, resourceModel administrativeAccountResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_administrative_account" "%[1]s" {
	%[2]s
}

data "pingfederate_administrative_account" "%[1]s" {
  id = pingfederate_administrative_account.%[1]s.username
}`,
		resourceName,
		hcl(resourceModel.administrativeAccount),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAdministrativeAccountAttributes(config administrativeAccountResourceModel) (resource.TestCheckFunc, string) {
	var err error
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	response, _, responseErr := testClient.AdministrativeAccountsAPI.GetAccount(ctx, config.administrativeAccount.Username).Execute()
	// encryptedPassword := *response.EncryptedPassword
	return func(s *terraform.State) error {
		resourceType := "AdministrativeAccount"
		if responseErr != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, &config.administrativeAccount.Username, "active",
			*config.administrativeAccount.Active, *response.Active)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "description",
			*config.administrativeAccount.Description, *response.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "department",
			*config.administrativeAccount.Department, *response.Department)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "email_address",
			*config.administrativeAccount.EmailAddress, *response.EmailAddress)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "phone_number",
			*config.administrativeAccount.PhoneNumber, *response.PhoneNumber)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "username",
			config.administrativeAccount.Username, response.Username)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.administrativeAccount.Username, "roles",
			config.administrativeAccount.Roles, response.Roles)
		if err != nil {
			return err
		}

		// fmt.Println(*response.EncryptedPassword)

		return nil
	}, ""
}
