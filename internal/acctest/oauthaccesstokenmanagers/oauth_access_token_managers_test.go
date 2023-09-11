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

const oauthAccessTokenManagersId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type oauthAccessTokenManagersResourceModel struct {
	id string
	attributeContract 	
	id string	
	name string	
	pluginDescriptorRef 	
	parentRef 	
	configuration 	
	selectionSettings 	
	accessControlSettings 	
	sessionValidationSettings 	
	sequenceNumber int64
}

func TestAccOauthAccessTokenManagers(t *testing.T) {
	resourceName := "myOauthAccessTokenManagers"
	initialResourceModel := oauthAccessTokenManagersResourceModel{
		attributeContract: fill in test value,	
		id: fill in test value,	
		name: fill in test value,	
		pluginDescriptorRef: fill in test value,	
		parentRef: fill in test value,	
		configuration: fill in test value,	
		selectionSettings: fill in test value,	
		accessControlSettings: fill in test value,	
		sessionValidationSettings: fill in test value,	
		sequenceNumber: fill in test value,
	}
	updatedResourceModel := oauthAccessTokenManagersResourceModel{
		attributeContract: fill in test value,	
		id: fill in test value,	
		name: fill in test value,	
		pluginDescriptorRef: fill in test value,	
		parentRef: fill in test value,	
		configuration: fill in test value,	
		selectionSettings: fill in test value,	
		accessControlSettings: fill in test value,	
		sessionValidationSettings: fill in test value,	
		sequenceNumber: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckOauthAccessTokenManagersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthAccessTokenManagers(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthAccessTokenManagersAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthAccessTokenManagers(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthAccessTokenManagersAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccOauthAccessTokenManagers(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_oauth_access_token_managers." + resourceName,
				ImportStateId:           oauthAccessTokenManagersId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testAccOauthAccessTokenManagers(resourceName string, resourceModel oauthAccessTokenManagersResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_access_token_managers" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthAccessTokenManagersAttributes(config oauthAccessTokenManagersResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthAccessTokenManagers"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetOauthAccessTokenManagers(ctx, oauthAccessTokenManagersId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOauthAccessTokenManagersDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeleteOauthAccessTokenManagers(ctx, oauthAccessTokenManagersId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthAccessTokenManagers", oauthAccessTokenManagersId)
	}
	return nil
}
