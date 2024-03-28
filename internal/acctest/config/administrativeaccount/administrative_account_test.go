package administrativeaccount_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
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
	initialAdministrativeAccount.Active = pointers.Bool(true)
	initialAdministrativeAccount.Roles = []string{"USER_ADMINISTRATOR"}
	return initialAdministrativeAccount
}

func updateAdministrativeAccount() *client.AdministrativeAccount {
	updateAdministrativeAccount := client.NewAdministrativeAccountWithDefaults()
	updateAdministrativeAccount.Username = username
	updateAdministrativeAccount.Active = pointers.Bool(false)
	updateAdministrativeAccount.Description = pointers.String("updated description")
	updateAdministrativeAccount.Department = pointers.String("updated department")
	updateAdministrativeAccount.EmailAddress = pointers.String("test@example.com")
	updateAdministrativeAccount.PhoneNumber = pointers.String("555-555-5555")
	updateAdministrativeAccount.Roles = []string{}
	updateAdministrativeAccount.Auditor = pointers.Bool(true)
	updateAdministrativeAccount.Department = pointers.String("department")
	return updateAdministrativeAccount
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
		%[10]s
		`
		passwords := func() (string, string) {
			if aa.EncryptedPassword != nil {
				encryptedPasswordVal := aa.GetEncryptedPassword()
				passwordVal := ""
				return encryptedPasswordVal, passwordVal
			} else {
				encryptedPasswordVal := ""
				passwordVal := password
				return encryptedPasswordVal, passwordVal
			}
		}
		encryptedPasswordTfVal, passwordTfVal := passwords()
		auditor := ""
		if aa.Auditor != nil {
			auditor = strconv.FormatBool(*aa.Auditor)
		}
		builder.WriteString(
			fmt.Sprintf(tf,
				acctest.TfKeyValuePairToString("active", strconv.FormatBool(aa.GetActive()), true),
				acctest.TfKeyValuePairToString("description", aa.GetDescription(), true),
				acctest.TfKeyValuePairToString("department", aa.GetDepartment(), true),
				acctest.TfKeyValuePairToString("email_address", aa.GetEmailAddress(), true),
				acctest.TfKeyValuePairToString("encrypted_password", encryptedPasswordTfVal, true),
				acctest.TfKeyValuePairToString("roles", acctest.StringSliceToTerraformString(aa.Roles), false),
				acctest.TfKeyValuePairToString("password", passwordTfVal, true),
				acctest.TfKeyValuePairToString("phone_number", aa.GetPhoneNumber(), true),
				acctest.TfKeyValuePairToString("username", aa.Username, true),
				acctest.TfKeyValuePairToString("auditor", auditor, false),
			),
		)
	}
	return builder.String()
}

func testAccAdministrativeAccount(resourceName string, resourceModel administrativeAccountResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_administrative_account" "%[1]s" {
	%[2]s
}

data "pingfederate_administrative_account" "%[1]s" {
  username = pingfederate_administrative_account.%[1]s.username
}`,
		resourceName,
		hcl(resourceModel.administrativeAccount),
	)
}

func testAccCheckExpectedAdministrativeAccountAttributes(config administrativeAccountResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AdministrativeAccountsAPI.GetAccount(ctx, config.administrativeAccount.Username).Execute()
		if err != nil {
			return err
		}
		resourceType := "AdministrativeAccount"
		// Verify that attributes have expected values
		stateAttributeValues := s.Modules[0].Resources["pingfederate_administrative_account.myAdministrativeAccount"].Primary.Attributes
		err = acctest.TestAttributesMatchBool(resourceType, &config.administrativeAccount.Username, "active",
			*config.administrativeAccount.Active, response.GetActive())
		if err != nil {
			return err
		}

		err = acctest.VerifyStateAttributeValue(stateAttributeValues, "active", strconv.FormatBool(*config.administrativeAccount.Active))
		if err != nil {
			return err
		}

		if config.administrativeAccount.Description != nil {
			err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "description",
				*config.administrativeAccount.Description, response.GetDescription())
			if err != nil {
				return err
			}
		}

		if config.administrativeAccount.Department != nil {
			err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "department",
				*config.administrativeAccount.Department, response.GetDepartment())
			if err != nil {
				return err
			}
		}

		if config.administrativeAccount.EmailAddress != nil {
			err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "email_address",
				*config.administrativeAccount.EmailAddress, response.GetEmailAddress())
			if err != nil {
				return err
			}
		}

		if config.administrativeAccount.PhoneNumber != nil {
			err = acctest.TestAttributesMatchString(resourceType, &config.administrativeAccount.Username, "phone_number",
				*config.administrativeAccount.PhoneNumber, response.GetPhoneNumber())
			if err != nil {
				return err
			}
		}

		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.administrativeAccount.Username, "roles",
			config.administrativeAccount.Roles, response.Roles)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAdministrativeAccountDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.AdministrativeAccountsAPI.DeleteAccount(ctx, username).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AdministrativeAccount", username)
	}
	return nil
}

func TestAccAdministrativeAccount(t *testing.T) {
	resourceName := "myAdministrativeAccount"
	initialResourceModel := administrativeAccountResourceModel{
		administrativeAccount: initialAdministrativeAccount(),
	}

	updatedResourceModel := administrativeAccountResourceModel{
		administrativeAccount: updateAdministrativeAccount(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckAdministrativeAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAdministrativeAccount(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAdministrativeAccountAttributes(initialResourceModel),
			},
			{
				Config: testAccAdministrativeAccount(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAdministrativeAccountAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccAdministrativeAccount(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_administrative_account." + resourceName,
				ImportStateId:           initialResourceModel.administrativeAccount.Username,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encrypted_password", "password"},
			},
			{
				Config: testAccAdministrativeAccount(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAdministrativeAccountAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.AdministrativeAccountsAPI.DeleteAccount(ctx, updatedResourceModel.administrativeAccount.Username).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAdministrativeAccount(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAdministrativeAccountAttributes(initialResourceModel),
			},
		},
	})
}
