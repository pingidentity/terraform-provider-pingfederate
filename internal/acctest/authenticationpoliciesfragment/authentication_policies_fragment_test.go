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

func TestAccAuthenticationPoliciesFragment(t *testing.T) {
	resourceName := "myAuthenticationPoliciesFragment"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			// Test simple policy fragment
			{
				Config: testAccAuthenticationPoliciesFragmentSimple(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(resourceName, false),
			},
			// Test a more complex fragment
			{
				Config: testAccAuthenticationPoliciesFragmentComplex(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(resourceName, true),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationPoliciesFragmentComplex(resourceName),
				ResourceName:      "pingfederate_authentication_policies_fragment." + resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccAuthenticationPoliciesFragmentSimple(resourceName),
				Check:  testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(resourceName, false),
			},
		},
	})
}

func testAccAuthenticationPoliciesFragmentSimple(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_fragment" "%[1]s" {
		fragment_id = "%[1]s"
		name = "%[1]s"
		root_node = {
			action = {
			  authn_source_policy_action = {
				authentication_source = {
				  type = "IDP_ADAPTER"
				  source_ref = {
					id = "OTIdPJava"
				  }
				}
			  }
			},
			children = [
			  {
				action = {
				  done_policy_action = {
					context = "Fail"
				  }
				}
			  },
			  {
				action = {
				  done_policy_action = {
					context = "Success"
				  }
				}
			  }
			]
		  }
}`, resourceName,
	)
}

// TODO
func testAccAuthenticationPoliciesFragmentComplex(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_fragment" "%[1]s" {
		fragment_id = "%[1]s"
		name = "%[1]s"
		root_node = {
			action = {
			  authn_source_policy_action = {
				authentication_source = {
				  type = "IDP_ADAPTER"
				  source_ref = {
					id = "OTIdPJava"
				  }
				}
			  }
			},
			children = [
			  {
				action = {
				  done_policy_action = {
					context = "Fail"
				  }
				}
			  },
			  {
				action = {
				  done_policy_action = {
					context = "Success"
				  }
				}
			  }
			]
		  }
}`, resourceName,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPoliciesFragmentAttributes(id string, isComplexFragment bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPoliciesFragment"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationPoliciesAPI.GetFragment(ctx, id).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchStringPointer(resourceType, nil, "name",
			id, response.Name)
		if err != nil {
			return err
		}

		return nil
	}
}
