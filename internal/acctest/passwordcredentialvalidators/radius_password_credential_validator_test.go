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

const radiusPasswordCredentialValidatorsId = "radiusPcv"

// Attributes to test with. Add optional properties to test here if desired.
type radiusPasswordCredentialValidatorsResourceModel struct {
	id                    string
	name                  string
	authPort              string
	timeout               string
	sharedSecret          string
	includeOptionalFields bool
}

func TestAccRadiusPasswordCredentialValidators(t *testing.T) {
	resourceName := "radiusPCV"
	initialResourceModel := radiusPasswordCredentialValidatorsResourceModel{
		id:                    radiusPasswordCredentialValidatorsId,
		name:                  "example",
		authPort:              "1812",
		timeout:               "3000",
		sharedSecret:          "2FederateM0re",
		includeOptionalFields: false,
	}
	updatedResourceModel := radiusPasswordCredentialValidatorsResourceModel{
		id:                    radiusPasswordCredentialValidatorsId,
		name:                  "updated example",
		authPort:              "1813",
		timeout:               "4000",
		sharedSecret:          "2FederateM0re!",
		includeOptionalFields: true,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckRadiusPasswordCredentialValidatorsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRadiusPasswordCredentialValidators(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedRadiusPasswordCredentialValidatorsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccRadiusPasswordCredentialValidators(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedRadiusPasswordCredentialValidatorsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccRadiusPasswordCredentialValidators(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_password_credential_validator." + resourceName,
				ImportStateId:     radiusPasswordCredentialValidatorsId,
				ImportState:       true,
				ImportStateVerify: true,
				// Fields get imported to fields_all, so can't check them here. Also can't check the imported shared secret
				ImportStateVerifyIgnore: []string{
					"configuration.fields",
					"configuration.tables.0.rows.0.fields.3.value",
				},
			},
			{
				// Back to minimal model
				Config: testAccRadiusPasswordCredentialValidators(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedRadiusPasswordCredentialValidatorsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccRadiusPasswordCredentialValidators(resourceName string, resourceModel radiusPasswordCredentialValidatorsResourceModel) string {
	fieldsHcl := ""
	attributeContractHcl := ""
	if resourceModel.includeOptionalFields {
		fieldsHcl = fmt.Sprintf(`
		fields = [
			{
				name  = "NAS Identifier"
				value = "PingFederate"
			},
			{
				name  = "Timeout"
				value = "%[1]s"
			},
			{
				name  = "Retry Count"
				value = "3"
			},
			{
				name  = "Allow Challenge Retries after Access-Reject"
				value = "false"
			}
		]
	`, resourceModel.timeout)
		attributeContractHcl = `
		attribute_contract = {
			extended_attributes = [
				{
					name = "contract"
				}
			]
			inherited = false
	  	}
	`
	}

	return fmt.Sprintf(`
resource "pingfederate_password_credential_validator" "%[1]s" {
  validator_id = "%[2]s"
  name         = "%[3]s"
  plugin_descriptor_ref = {
    id = "org.sourceid.saml20.domain.RadiusUsernamePasswordCredentialValidator"
  }
  configuration = {
    tables = [
      {
        name = "RADIUS Servers"
        rows = [
          {
            fields = [
              {
                name  = "Hostname"
                value = "localhost"
              },
              {
                name  = "Authentication Port"
                value = "%[4]s"
              },
              {
                name  = "Authentication Protocol"
                value = "PAP"
              },
              {
                name = "Shared Secret"
                # This value will be stored into your state file and will not detect any configuration changes made in the UI
                value = "%[5]s"
              }
            ]
            default_row = false
          }
        ]
      }
    ],
    %[6]s
  }
  %[7]s
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
		resourceModel.authPort,
		resourceModel.sharedSecret,
		fieldsHcl,
		attributeContractHcl,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedRadiusPasswordCredentialValidatorsAttributes(config radiusPasswordCredentialValidatorsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "PasswordCredentialValidators"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PasswordCredentialValidatorsAPI.GetPasswordCredentialValidator(ctx, radiusPasswordCredentialValidatorsId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.name, response.Name)
		if err != nil {
			return err
		}

		respConfig := response.Configuration
		configTables := respConfig.Tables
		for _, configTable := range configTables {
			for _, row := range configTable.Rows {
				for _, field := range row.Fields {
					if field.Name == "Authentication Port" {
						authPort := field.Value
						err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.authPort, *authPort)
						if err != nil {
							return err
						}
					}
				}
			}
		}

		configFields := respConfig.Fields
		for _, field := range configFields {
			if field.Name == "Timeout" {
				timeout := field.Value
				err = acctest.TestAttributesMatchString(resourceType, &config.id, "name", config.timeout, *timeout)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckRadiusPasswordCredentialValidatorsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.PasswordCredentialValidatorsAPI.DeletePasswordCredentialValidator(ctx, radiusPasswordCredentialValidatorsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("PasswordCredentialValidators", radiusPasswordCredentialValidatorsId)
	}
	return nil
}
