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

const oauthIssuerId = "2"

type oauthIssuerResourceModel struct {
	description   string
	host          string
	name          string
	path          string
	oauthIssuerId string
}

func TestAccOauthIssuer(t *testing.T) {
	resourceName := "myOauthIssuer"
	initialResourceModel := oauthIssuerResourceModel{
		description:   "description",
		host:          "hostname",
		name:          "name",
		path:          "/example",
		oauthIssuerId: oauthIssuerId,
	}
	updatedResourceModel := oauthIssuerResourceModel{
		description:   "updated description",
		host:          "updatedhostname",
		name:          "updatedname",
		path:          "/updated",
		oauthIssuerId: oauthIssuerId,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckOauthIssuerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOauthIssuer(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedOauthIssuerAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccOauthIssuer(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedOauthIssuerAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccOauthIssuer(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_oauth_issuer." + resourceName,
				ImportStateId:     initialResourceModel.oauthIssuerId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// No need to go back to "minimal" here since everything is required on this resource
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.OauthIssuersAPI.DeleteOauthIssuer(ctx, oauthIssuerId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccOauthIssuer(resourceName string, resourceModel oauthIssuerResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_issuer" "%[1]s" {
  issuer_id   = "%[2]s"
  description = "%[3]s"
  host        = "%[4]s"
  name        = "%[5]s"
  path        = "%[6]s"
}
data "pingfederate_oauth_issuer" "%[1]s" {
  issuer_id = pingfederate_oauth_issuer.%[1]s.id
}`, resourceName,
		resourceModel.oauthIssuerId,
		resourceModel.description,
		resourceModel.host,
		resourceModel.name,
		resourceModel.path,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedOauthIssuerAttributes(config oauthIssuerResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "OauthIssuer"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.OauthIssuersAPI.GetOauthIssuerById(ctx, config.oauthIssuerId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.oauthIssuerId, "description",
			config.description, *response.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.oauthIssuerId, "host",
			config.host, response.Host)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.oauthIssuerId, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.oauthIssuerId, "path",
			config.path, *response.Path)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckOauthIssuerDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.OauthIssuersAPI.DeleteOauthIssuer(ctx, oauthIssuerId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("OauthIssuer", oauthIssuerId)
	}
	return nil
}
