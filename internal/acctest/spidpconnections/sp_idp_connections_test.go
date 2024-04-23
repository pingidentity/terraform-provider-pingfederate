package spIdpConnections_test

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

const spIdpConnectionsId = "2"

// Attributes to test with. Add optional properties to test here if desired.
type spIdpConnectionsResourceModel struct {
	id string
	oidcClientCredentials 	
	type string	
	id string	
	entityId string	
	name string	
	modificationDate string	
	creationDate string	
	active bool	
	baseUrl string	
	defaultVirtualEntityId string	
	virtualEntityIds []correct_type	
	metadataReloadSettings 	
	credentials 	
	contactInfo 	
	licenseConnectionGroup string	
	loggingMode string	
	additionalAllowedEntitiesConfiguration 	
	extendedProperties determine this value manually	
	idpBrowserSso 	
	attributeQuery 	
	idpOAuthGrantAttributeMapping 	
	wsTrust 	
	inboundProvisioning 	
	errorPageMsgId string
}

func TestAccSpIdpConnections(t *testing.T) {
	resourceName := "mySpIdpConnections"
	initialResourceModel := spIdpConnectionsResourceModel{
		oidcClientCredentials: fill in test value,	
		type: fill in test value,	
		id: fill in test value,	
		entityId: fill in test value,	
		name: fill in test value,	
		modificationDate: fill in test value,	
		creationDate: fill in test value,	
		active: fill in test value,	
		baseUrl: fill in test value,	
		defaultVirtualEntityId: fill in test value,	
		virtualEntityIds: fill in test value,	
		metadataReloadSettings: fill in test value,	
		credentials: fill in test value,	
		contactInfo: fill in test value,	
		licenseConnectionGroup: fill in test value,	
		loggingMode: fill in test value,	
		additionalAllowedEntitiesConfiguration: fill in test value,	
		extendedProperties: fill in test value,	
		idpBrowserSso: fill in test value,	
		attributeQuery: fill in test value,	
		idpOAuthGrantAttributeMapping: fill in test value,	
		wsTrust: fill in test value,	
		inboundProvisioning: fill in test value,	
		errorPageMsgId: fill in test value,
	}
	updatedResourceModel := spIdpConnectionsResourceModel{
		oidcClientCredentials: fill in test value,	
		type: fill in test value,	
		id: fill in test value,	
		entityId: fill in test value,	
		name: fill in test value,	
		modificationDate: fill in test value,	
		creationDate: fill in test value,	
		active: fill in test value,	
		baseUrl: fill in test value,	
		defaultVirtualEntityId: fill in test value,	
		virtualEntityIds: fill in test value,	
		metadataReloadSettings: fill in test value,	
		credentials: fill in test value,	
		contactInfo: fill in test value,	
		licenseConnectionGroup: fill in test value,	
		loggingMode: fill in test value,	
		additionalAllowedEntitiesConfiguration: fill in test value,	
		extendedProperties: fill in test value,	
		idpBrowserSso: fill in test value,	
		attributeQuery: fill in test value,	
		idpOAuthGrantAttributeMapping: fill in test value,	
		wsTrust: fill in test value,	
		inboundProvisioning: fill in test value,	
		errorPageMsgId: fill in test value,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckSpIdpConnectionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpIdpConnections(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSpIdpConnectionsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccSpIdpConnections(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedSpIdpConnectionsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:                  testAccSpIdpConnections(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_sp_idp_connections." + resourceName,
				ImportStateId:           spIdpConnectionsId,
				ImportState:             true,
				ImportStateVerify:       true,
			},
			{
				Config: testAccSpIdpConnections(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSpIdpConnectionsAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.ReplaceMe.DeleteResourceType(ctx, idGoesHere).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccSpIdpConnections(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedSpIdpConnectionsAttributes(initialResourceModel),
			},
		},
	})
}

func testAccSpIdpConnections(resourceName string, resourceModel spIdpConnectionsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_sp_idp_connections" "%[1]s" {
	id = "%[2]s"
	FILL THIS IN
}`, resourceName,
		resourceModel.id,
	
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedSpIdpConnectionsAttributes(config spIdpConnectionsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "SpIdpConnections"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.<RESOURCE_API>.GetSpIdpConnections(ctx, spIdpConnectionsId).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		FILL THESE in! 

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckSpIdpConnectionsDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.<RESOURCE_API>.DeleteSpIdpConnections(ctx, spIdpConnectionsId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("SpIdpConnections", spIdpConnectionsId)
	}
	return nil
}
