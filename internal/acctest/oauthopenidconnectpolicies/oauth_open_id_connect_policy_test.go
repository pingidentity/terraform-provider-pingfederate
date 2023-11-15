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

const oauthOpenIdConnectPoliciesId = "testOpenIdConnectPolicy"

// Attributes to test with. Add optional properties to test here if desired.
type oauthOpenIdConnectPoliciesResourceModel struct {
	id string
	/*name string
	accessTokenManagerRef
	idTokenLifetime int64
	includeSriInIdToken bool
	includeUserInfoInIdToken bool
	includeSHashInIdToken bool
	returnIdTokenOnRefreshGrant bool
	reissueIdTokenInHybridFlow bool
	attributeContract
	attributeMapping
	scopeAttributeMappings determine this value manually*/
}

func TestAccOauthOpenIdConnectPolicies(t *testing.T) {
	resourceName := "myOauthOpenIdConnectPolicies"
	initialResourceModel := oauthOpenIdConnectPoliciesResourceModel{
		/*name: fill in test value,
		accessTokenManagerRef: fill in test value,
		idTokenLifetime: fill in test value,
		includeSriInIdToken: fill in test value,
		includeUserInfoInIdToken: fill in test value,
		includeSHashInIdToken: fill in test value,
		returnIdTokenOnRefreshGrant: fill in test value,
		reissueIdTokenInHybridFlow: fill in test value,
		attributeContract: fill in test value,
		attributeMapping: fill in test value,
		scopeAttributeMappings: fill in test value,*/
	}
	updatedResourceModel := oauthOpenIdConnectPoliciesResourceModel{
		/*name: fill in test value,
		accessTokenManagerRef: fill in test value,
		idTokenLifetime: fill in test value,
		includeSriInIdToken: fill in test value,
		includeUserInfoInIdToken: fill in test value,
		includeSHashInIdToken: fill in test value,
		returnIdTokenOnRefreshGrant: fill in test value,
		reissueIdTokenInHybridFlow: fill in test value,
		attributeContract: fill in test value,
		attributeMapping: fill in test value,
		scopeAttributeMappings: fill in test value,*/
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOauthOpenIdConnectPoliciesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthOpenIdConnectPolicies(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthOpenIdConnectPoliciesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthOpenIdConnectPolicies(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthOpenIdConnectPoliciesAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthOpenIdConnectPolicies(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_open_id_connect_policy." + resourceName,
				ImportStateId:     oauthOpenIdConnectPoliciesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Back to minimal model
				Config: testAccOauthOpenIdConnectPolicies(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthOpenIdConnectPoliciesAttributes(initialResourceModel),
			},
		},
	})
}

func testAccOauthOpenIdConnectPolicies(resourceName string, resourceModel oauthOpenIdConnectPoliciesResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_open_id_connect_policy" "%[1]s" {
	id = "%[2]s"
	// FILL THIS IN
}`, resourceName,
		oauthOpenIdConnectPoliciesId,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthOpenIdConnectPoliciesAttributes(config oauthOpenIdConnectPoliciesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//resourceType := "OauthOpenIdConnectPolicy"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		_, _, err := testClient.OauthOpenIdConnectAPI.GetOIDCPolicy(ctx, oauthOpenIdConnectPoliciesId).Execute()

		if err != nil {
			return err
		}

		//TODO Verify that attributes have expected values

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOauthOpenIdConnectPoliciesDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthOpenIdConnectAPI.DeleteOIDCPolicy(ctx, oauthOpenIdConnectPoliciesId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthOpenIdConnectPolict", oauthOpenIdConnectPoliciesId)
	}
	return nil
}
