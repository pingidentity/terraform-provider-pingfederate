package acctest_test

import (
	"fmt"
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

func extendedPropertiesHCLObj(extProps *client.ExtendedProperty) string {
	name := func() string {
		if extProps.GetName() == "" {
			return ""
		}
		return fmt.Sprintf("{\n\t\t\tname = \"%s\"", *extProps.Name)
	}

	description := func() string {
		if extProps.GetDescription() == "" {
			return ""
		}
		return fmt.Sprintf("\tdescription = \"%s\"", *extProps.Description)
	}

	multi_valued := func() string {
		if extProps.MultiValued == nil || !extProps.GetMultiValued() {
			return ""
		}
		return fmt.Sprintf("\tmulti_valued = %t\n\t\t}", extProps.GetMultiValued())
	}

	return fmt.Sprintf("%s\n\t\t%s\n\t\t%s", name(), description(), multi_valued())
}

func TestAccExtendedProperties(t *testing.T) {
	resourceName := "myExtendedProperties"
	initialResourceModel := client.ExtendedProperty{}
	updatedExtendedProperty := client.NewExtendedProperty()
	updatedExtendedProperty.Name = pointers.String("example")
	updatedExtendedProperty.Description = pointers.String("example description")
	updatedExtendedProperty.MultiValued = pointers.Bool(true)
	updatedResourceModel := updatedExtendedProperty

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test empty object sent
			{
				Config: testAccExtendedProperties(resourceName, &initialResourceModel),
				Check:  testAccCheckExpectedExtendedPropertiesAttributes(&initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccExtendedProperties(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedExtendedPropertiesAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccExtendedProperties(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_extended_properties." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test empty object sent
			{
				Config: testAccExtendedProperties(resourceName, &initialResourceModel),
				Check:  testAccCheckExpectedExtendedPropertiesAttributes(&initialResourceModel),
			},
		},
	})
}

func testAccExtendedProperties(resourceName string, extendedProperty *client.ExtendedProperty) string {
	return fmt.Sprintf(`
resource "pingfederate_extended_properties" "%[1]s" {
  items = [
		%s
  ]
}`, resourceName,
		extendedPropertiesHCLObj(extendedProperty),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedExtendedPropertiesAttributes(extendedProperty *client.ExtendedProperty) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ExtendedProperties"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ExtendedPropertiesAPI.GetExtendedProperties(ctx).Execute()

		if err != nil {
			return err
		}

		// Check for no items if empty object sent
		if extendedProperty.Name == nil {
			if len(response.GetItems()) > 0 {
				return fmt.Errorf("Empty items object sent, expected no items to be returned, but got %d", len(response.GetItems()))
			}
			return nil
		} else if len(response.GetItems()) == 0 {
			return fmt.Errorf("Expected items to be returned, but got none")
		}

		// Verify that attributes have expected values
		if len(response.GetItems()) > 0 {
			err = acctest.TestAttributesMatchString(resourceType, nil, "name", *extendedProperty.Name, *response.Items[0].Name)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchString(resourceType, nil, "description", *extendedProperty.Description, *response.Items[0].Description)
			if err != nil {
				return err
			}

			err = acctest.TestAttributesMatchBool(resourceType, nil, "multi_valued", *extendedProperty.MultiValued, *response.Items[0].MultiValued)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
