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

const authenticationPoliciesPolicyId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type authenticationPoliciesPolicyResourceModel struct {
	id string
	id string	
	name string	
	description string	
	authenticationApiApplicationRef 	
	enabled bool	
	rootNode 	
	handleFailuresLocally bool
}

func TestAccAuthenticationPoliciesPolicy(t *testing.T) {
	resourceName := "myAuthenticationPoliciesPolicy"
	initialResourceModel := authenticationPoliciesPolicyResourceModel{
		id: fill in test value,	
		name: fill in test value,	
		description: fill in test value,	
		authenticationApiApplicationRef: fill in test value,	
		enabled: fill in test value,	
		rootNode: fill in test value,	
		handleFailuresLocally: fill in test value,
	}
	updatedResourceModel := authenticationPoliciesPolicyResourceModel{
		id: fill in test value,	
		name: fill in test value,	
		description: fill in test value,	
		authenticationApiApplicationRef: fill in test value,	
		enabled: fill in test value,	
		rootNode: fill in test value,	
		handleFailuresLocally: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckAuthenticationPoliciesPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationPoliciesPolicy(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAuthenticationPoliciesPolicy(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccAuthenticationPoliciesPolicy(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_authentication_policies_policy." + resourceName,
				ImportStateId:           authenticationPoliciesPolicyId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testAccAuthenticationPoliciesPolicy(resourceName string, resourceModel authenticationPoliciesPolicyResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_policies_policy" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationPoliciesPolicyAttributes(config authenticationPoliciesPolicyResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationPoliciesPolicy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetAuthenticationPoliciesPolicy(ctx, authenticationPoliciesPolicyId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAuthenticationPoliciesPolicyDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeleteAuthenticationPoliciesPolicy(ctx, authenticationPoliciesPolicyId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AuthenticationPoliciesPolicy", authenticationPoliciesPolicyId)
	}
	return nil
}
