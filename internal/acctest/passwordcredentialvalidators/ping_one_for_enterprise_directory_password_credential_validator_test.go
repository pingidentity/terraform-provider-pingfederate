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

const pingOneForEnterpriseDirectoryPasswordCredentialValidatorsId = "pingOneForEnterpriseDirectoryPcv"

// Attributes to test with. Add optional properties to test here if desired.
type pingOneForEnterpriseDirectoryPasswordCredentialValidatorsResourceModel struct {
	id   string
	name string
}

func TestAccPingOneForEnterpriseDirectoryPasswordCredentialValidators(t *testing.T) {
	resourceName := "pingOneForEnterpriseDirectoryPCV"
	initialResourceModel := pingOneForEnterpriseDirectoryPasswordCredentialValidatorsResourceModel{
		id:   pingOneForEnterpriseDirectoryPasswordCredentialValidatorsId,
		name: "example",
	}
	updatedResourceModel := pingOneForEnterpriseDirectoryPasswordCredentialValidatorsResourceModel{
		id:   pingOneForEnterpriseDirectoryPasswordCredentialValidatorsId,
		name: "updated example",
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckPingOneForEnterpriseDirectoryPasswordCredentialValidatorsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPingOneForEnterpriseDirectoryPasswordCredentialValidators(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedPingOneForEnterpriseDirectoryPasswordCredentialValidatorsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccPingOneForEnterpriseDirectoryPasswordCredentialValidators(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedPingOneForEnterpriseDirectoryPasswordCredentialValidatorsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccPingOneForEnterpriseDirectoryPasswordCredentialValidators(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_password_credential_validators." + resourceName,
				ImportStateId:           pingOneForEnterpriseDirectoryPasswordCredentialValidatorsId,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"configuration"},
			},
		},
	})
}

func testAccPingOneForEnterpriseDirectoryPasswordCredentialValidators(resourceName string, resourceModel pingOneForEnterpriseDirectoryPasswordCredentialValidatorsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_password_credential_validators" "%[1]s" {
	id = "%[2]s"
	name = "%[3]s"
	plugin_descriptor_ref = {
		id = "com.pingconnect.alexandria.pingfed.pcv.PingOnePasswordValidator"
	}
	configuration = {
		tables = [],
		fields = [
			{
				name = "Client Id"
				value = "ping_federate_client_id"
			},
			{
				name = "Client Secret"
				value = "2FederateM0re"
			},
			{
				name = "PingOne URL"
				value = "https://directory-api.pingone.com/api"
			},
			{
				name = "Authenticate by Subject URL"
				value = "/directory/users/authenticate?by=subject"
			},
			{
				name = "Reset Password URL"
				value = "/directory/users/password-reset"
			},
			{
				name = "SCIM User URL"
				value = "/directory/user"
			},
			{
				name = "Connection Pool Size"
				value = "100"
			},
			{
				name = "Connection Pool Idle Timeout"
				value = "4000"
			}
		]
	}
}`, resourceName,
		resourceModel.id,
		resourceModel.name,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedPingOneForEnterpriseDirectoryPasswordCredentialValidatorsAttributes(config pingOneForEnterpriseDirectoryPasswordCredentialValidatorsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "PasswordCredentialValidators"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.PasswordCredentialValidatorsApi.GetPasswordCredentialValidator(ctx, pingOneForEnterpriseDirectoryPasswordCredentialValidatorsId).Execute()

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
func testAccCheckPingOneForEnterpriseDirectoryPasswordCredentialValidatorsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.PasswordCredentialValidatorsApi.DeletePasswordCredentialValidator(ctx, pingOneForEnterpriseDirectoryPasswordCredentialValidatorsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("PasswordCredentialValidators", pingOneForEnterpriseDirectoryPasswordCredentialValidatorsId)
	}
	return nil
}
