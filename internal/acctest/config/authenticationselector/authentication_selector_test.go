package authenticationselector_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/common/pointers"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

const authenticationSelectorsId = "selectorTest"

// Attributes to test with. Add optional properties to test here if desired.
type authenticationSelectorsResourceModel struct {
	attributeContract                *client.AuthenticationSelectorAttributeContract
	pluginDescriptorRef              client.ResourceLink
	addOrUpdateAuthNContextAttribute string
	enableNoMatchResultValue         string
	enableNotInRequestResultValue    string
	overrideAuthnContextForFlow      string
}

func TestAccAuthenticationSelector(t *testing.T) {
	resourceName := "myAuthenticationSelector"

	initialResourceModel := authenticationSelectorsResourceModel{
		pluginDescriptorRef: client.ResourceLink{
			Id: "com.pingidentity.pf.selectors.saml.SamlAuthnContextAdapterSelector",
		},
		addOrUpdateAuthNContextAttribute: "true",
		enableNoMatchResultValue:         "false",
		enableNotInRequestResultValue:    "false",
		overrideAuthnContextForFlow:      "true",
		attributeContract: &client.AuthenticationSelectorAttributeContract{
			ExtendedAttributes: []client.AuthenticationSelectorAttribute{
				{
					Name: "result_value",
				},
			},
		},
	}

	updatedResourceModel := authenticationSelectorsResourceModel{
		pluginDescriptorRef: client.ResourceLink{
			Id: "com.pingidentity.pf.selectors.saml.SamlAuthnContextAdapterSelector",
		},
		addOrUpdateAuthNContextAttribute: "false",
		enableNoMatchResultValue:         "true",
		enableNotInRequestResultValue:    "true",
		overrideAuthnContextForFlow:      "false",
		attributeContract: &client.AuthenticationSelectorAttributeContract{
			ExtendedAttributes: []client.AuthenticationSelectorAttribute{
				{
					Name: "result_value2",
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckAuthenticationSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationSelector(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationSelectorAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAuthenticationSelector(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAuthenticationSelectorAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationSelector(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_authentication_selector." + resourceName,
				ImportStateId:     authenticationSelectorsId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAuthenticationSelector(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationSelectorAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.AuthenticationSelectorsAPI.DeleteAuthenticationSelector(ctx, authenticationSelectorsId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAuthenticationSelector(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationSelectorAttributes(initialResourceModel),
			},
		},
	})
}

func overrideAuthnContextForFlow(overrideAuthnContextForFlowVal string) string {
	if acctest.VersionAtLeast("11.3.0") {
		return fmt.Sprintf(`
		{
			name  = "Override AuthN Context for Flow"
			value = "%s"
		},`, overrideAuthnContextForFlowVal)
	}
	return ""
}

func testAccAuthenticationSelector(resourceName string, resourceModel authenticationSelectorsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_selector" "%[1]s" {
  selector_id = "%[2]s"
  name        = "%[2]s"
  plugin_descriptor_ref = {
    id = "%[3]s"
  }
  configuration = {
    tables = []
    fields = [
			%[4]s
      {
        name  = "Add or Update AuthN Context Attribute"
        value = "%[5]s"
      },
      {
        name  = "Enable 'No Match' Result Value"
        value = "%[6]s"
      },
      {
        name  = "Enable 'Not in Request' Result Value"
        value = "%[7]s"
      }
    ]
  }
  attribute_contract = {
    extended_attributes = [
      {
        name = "%[8]s"
      }
    ]
  }
}`, resourceName,
		authenticationSelectorsId,
		resourceModel.pluginDescriptorRef.Id,
		overrideAuthnContextForFlow(resourceModel.overrideAuthnContextForFlow),
		resourceModel.addOrUpdateAuthNContextAttribute,
		resourceModel.enableNoMatchResultValue,
		resourceModel.enableNotInRequestResultValue,
		resourceModel.attributeContract.ExtendedAttributes[0].Name,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationSelectorAttributes(config authenticationSelectorsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationSelector"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationSelectorsAPI.GetAuthenticationSelector(ctx, authenticationSelectorsId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationSelectorsId), "selector_id", authenticationSelectorsId, response.Id)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationSelectorsId), "name", authenticationSelectorsId, response.Name)
		if err != nil {
			return err
		}

		configFields := response.Configuration.Fields
		for _, field := range configFields {

			if field.Name == "Add or Update AuthN Context Attribute" {
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationSelectorsId), "value", config.addOrUpdateAuthNContextAttribute, *field.Value)
				if err != nil {
					return err
				}
			}

			if field.Name == "Enable No Match Result Value" {
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationSelectorsId), "value", config.enableNoMatchResultValue, *field.Value)
				if err != nil {
					return err
				}
			}

			if field.Name == "Enable Not In Request Result Value" {
				err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationSelectorsId), "value", config.enableNotInRequestResultValue, *field.Value)
				if err != nil {
					return err
				}
			}

		}

		if config.attributeContract != nil {
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationSelectorsId), "name", config.attributeContract.ExtendedAttributes[0].Name, response.AttributeContract.ExtendedAttributes[0].Name)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAuthenticationSelectorDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.AuthenticationSelectorsAPI.DeleteAuthenticationSelector(ctx, authenticationSelectorsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AuthenticationSelector", authenticationSelectorsId)
	}
	return nil
}
