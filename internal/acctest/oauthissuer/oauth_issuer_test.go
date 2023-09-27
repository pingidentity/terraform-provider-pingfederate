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

const stateID = "2"

type oauthIssuerResourceModel struct {
	description string
	host        string
	name        string
	path        string
	stateId     string
}

func TestAccOauthIssuer(t *testing.T) {
	resourceName := "myOauthIssuer"
	initialResourceModel := oauthIssuerResourceModel{
		description: "description",
		host:        "hostname",
		name:        "name",
		path:        "/example",
		stateId:     stateID,
	}
	updatedResourceModel := oauthIssuerResourceModel{
		description: "updated description",
		host:        "updatedhostname",
		name:        "updatedname",
		path:        "/updated",
		stateId:     stateID,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
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
				ImportStateId:     initialResourceModel.stateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOauthIssuer(resourceName string, resourceModel oauthIssuerResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_issuer" "%[1]s" {
  custom_id   = "%[2]s"
  description = "%[3]s"
  host        = "%[4]s"
  name        = "%[5]s"
  path        = "%[6]s"
}`, resourceName,
		resourceModel.stateId,
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
		response, _, err := testClient.OauthIssuersApi.GetOauthIssuerById(ctx, config.stateId).Execute()
		if err != nil {
			return err
		}
		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "description",
			config.description, *response.Description)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "host",
			config.host, response.Host)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "name",
			config.name, response.Name)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchString(resourceType, &config.stateId, "path",
			config.path, *response.Path)
		if err != nil {
			return err
		}

		return nil
	}
}
