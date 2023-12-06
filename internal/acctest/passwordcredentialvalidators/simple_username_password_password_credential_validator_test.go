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

const simpleUsernamePasswordPasswordCredentialValidatorsId = "simpleUsernamePasswordPcv"

// Attributes to test with. Add optional properties to test here if desired.
type simpleUsernamePasswordPasswordCredentialValidatorsResourceModel struct {
	id                    string
	name                  string
	password              string
	includeOptionalFields bool
}

func TestAccSimpleUsernamePasswordCredentialValidators(t *testing.T) {
	resourceName := "mySimpleUsernamePasswordCredentialValidators"
	initialResourceModel := simpleUsernamePasswordPasswordCredentialValidatorsResourceModel{
		id:                    simpleUsernamePasswordPasswordCredentialValidatorsId,
		name:                  "example",
		password:              "2FederateM0re",
		includeOptionalFields: false,
	}
	updatedResourceModel := simpleUsernamePasswordPasswordCredentialValidatorsResourceModel{
		id:                    simpleUsernamePasswordPasswordCredentialValidatorsId,
		name:                  "updated example",
		password:              "2FederateM0re!",
		includeOptionalFields: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckPasswordCredentialValidatorsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPasswordCredentialValidators(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPasswordCredentialValidatorsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPasswordCredentialValidators(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPasswordCredentialValidatorsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccPasswordCredentialValidators(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_password_credential_validator." + resourceName,
				ImportStateId:     simpleUsernamePasswordPasswordCredentialValidatorsId,
				ImportState:       true,
				ImportStateVerify: true,
				// Tables get imported to tables_all, so can't verify them here
				ImportStateVerifyIgnore: []string{
					"configuration.tables",
				},
			},
			{
				// Back to minimal model
				Config: testAccPasswordCredentialValidators(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPasswordCredentialValidatorsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccPasswordCredentialValidators(resourceName string, resourceModel simpleUsernamePasswordPasswordCredentialValidatorsResourceModel) string {
	optionalHcl := ""
	if resourceModel.includeOptionalFields {
		optionalHcl = `
		attribute_contract = {
			extended_attributes = []
			inherited = false
		}
		`
	}
	return fmt.Sprintf(`
resource "pingfederate_password_credential_validator" "%[1]s" {
  validator_id = "%[2]s"
  name         = "%[3]s"
  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.SimpleUsernamePasswordCredentialValidator"
  }
  configuration = {
    tables = [
      {
        name = "Users"
        rows = [
          {
            fields = [
              {
                name  = "Username"
                value = "example"
              },
              {
                name = "Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "%[4]s"
              },
              {
                name = "Confirm Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "%[4]s"
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          },
          {
            fields = [
              {
                name  = "Username"
                value = "example2"
              },
              {
                name = "Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "%[4]s"
              },
              {
                name = "Confirm Password"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "%[4]s"
              },
              {
                name  = "Relax Password Requirements"
                value = "false"
              }
            ]
            default_row = false
          }
        ],
      }
    ]
  }
  %[5]s
}
data "pingfederate_password_credential_validator" "%[1]s" {
  validator_id = pingfederate_password_credential_validator.%[1]s.validator_id
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.password,
		optionalHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedPasswordCredentialValidatorsAttributes(config simpleUsernamePasswordPasswordCredentialValidatorsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "PasswordCredentialValidators"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PasswordCredentialValidatorsAPI.GetPasswordCredentialValidator(ctx, simpleUsernamePasswordPasswordCredentialValidatorsId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.name, response.Name)
		if err != nil {
			return err
		}
		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckPasswordCredentialValidatorsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.PasswordCredentialValidatorsAPI.DeletePasswordCredentialValidator(ctx, simpleUsernamePasswordPasswordCredentialValidatorsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("PasswordCredentialValidators", simpleUsernamePasswordPasswordCredentialValidatorsId)
	}
	return nil
}
