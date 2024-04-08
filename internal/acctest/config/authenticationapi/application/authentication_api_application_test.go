package authenticationapiapplication_test

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

const authenticationApiApplicationId = "myAuthenticationApiApplication"
const authenticationApiApplicationName = "myAuthenticationApiApplicationName"
const authenticationApiApplicationUrl = "https://example.com"

// Attributes to test with. Add optional properties to test here if desired.
type authenticationApiApplicationResourceModel struct {
	applicationId                string
	name                         string
	url                          string
	description                  string
	additionalAllowedOrigins     []string
	clientForRedirectlessModeRef *client.ResourceLink
}

func TestAccAuthenticationApiApplication(t *testing.T) {
	resourceName := "myAuthenticationApiApplication"
	initialResourceModel := authenticationApiApplicationResourceModel{
		applicationId:                authenticationApiApplicationId,
		name:                         authenticationApiApplicationName,
		url:                          authenticationApiApplicationUrl,
		additionalAllowedOrigins:     []string{},
		clientForRedirectlessModeRef: nil,
	}

	clientForRedirectlessModeRefResourceLink := client.NewResourceLink("myOauthClientExample")

	updatedResourceModel := authenticationApiApplicationResourceModel{
		applicationId:                authenticationApiApplicationId,
		name:                         authenticationApiApplicationName,
		url:                          authenticationApiApplicationUrl,
		description:                  "myAuthenticationApiApplicationDescription",
		additionalAllowedOrigins:     []string{"https://example.com", "https://example2.com"},
		clientForRedirectlessModeRef: clientForRedirectlessModeRefResourceLink,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: testAccCheckAuthenticationApiApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationApiApplication(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationApiApplicationAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccAuthenticationApiApplication(resourceName, updatedResourceModel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExpectedAuthenticationApiApplicationAttributes(updatedResourceModel),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_api_application.%s", resourceName), "additional_allowed_origins.0", updatedResourceModel.additionalAllowedOrigins[0]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_api_application.%s", resourceName), "additional_allowed_origins.1", updatedResourceModel.additionalAllowedOrigins[1]),
					resource.TestCheckResourceAttr(fmt.Sprintf("pingfederate_authentication_api_application.%s", resourceName), "client_for_redirectless_mode_ref.id", updatedResourceModel.clientForRedirectlessModeRef.Id),
				),
			},
			{
				// Test importing the resource
				Config:                  testAccAuthenticationApiApplication(resourceName, updatedResourceModel),
				ResourceName:            "pingfederate_authentication_api_application." + resourceName,
				ImportStateId:           authenticationApiApplicationId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_for_redirectless_mode_ref"},
			},
			{
				Config: testAccAuthenticationApiApplication(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationApiApplicationAttributes(initialResourceModel),
			},
			{
				PreConfig: func() {
					testClient := acctest.TestClient()
					ctx := acctest.TestBasicAuthContext()
					_, err := testClient.AuthenticationApiAPI.DeleteApplication(ctx, updatedResourceModel.applicationId).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config: %v", err)
					}
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAuthenticationApiApplication(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationApiApplicationAttributes(initialResourceModel),
			},
		},
	})
}

func optionalHcl(model authenticationApiApplicationResourceModel) string {
	var hclDescription string
	var hclAdditionalAllowedOrigins string
	var oauthClientDependency string
	var clientForRedirectlessModeRef string

	if model.description != "" {
		hclDescription = fmt.Sprintf("description = \"%s\"", model.description)
	}

	if len(model.additionalAllowedOrigins) > 0 {
		hclAdditionalAllowedOrigins = fmt.Sprintf("additional_allowed_origins = %s", acctest.StringSliceToTerraformString(model.additionalAllowedOrigins))
	}

	if model.clientForRedirectlessModeRef != nil {
		clientForRedirectlessModeRef = `
	client_for_redirectless_mode_ref = {
	  id = pingfederate_oauth_client.myOauthClientExample.id
	}`
	}

	return fmt.Sprintf(`
	%s
	%s
	%s
	%s
	`, hclDescription, hclAdditionalAllowedOrigins, oauthClientDependency, clientForRedirectlessModeRef)
}

func testAccAuthenticationApiApplication(resourceName string, resourceModel authenticationApiApplicationResourceModel) string {
	optionalFields := optionalHcl(resourceModel)
	return fmt.Sprintf(`
resource "pingfederate_oauth_client" "myOauthClientExample" {
  client_id                     = "myOauthClientExample"
  name                          = "myOauthClientExample"
  grant_types                   = ["EXTENSION"]
  allow_authentication_api_init = true
}

resource "pingfederate_authentication_api_settings" "myAuthenticationApiSettingsExample" {
  restrict_access_to_redirectless_mode = true
}

resource "pingfederate_authentication_api_application" "%[1]s" {
  application_id = "%[2]s"
  name           = "%[3]s"
  url            = "%[4]s"
	%[5]s
}`, resourceName,
		resourceModel.applicationId,
		resourceModel.name,
		resourceModel.url,
		optionalFields,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationApiApplicationAttributes(config authenticationApiApplicationResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationApiApplication"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationApiAPI.GetApplication(ctx, authenticationApiApplicationId).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationApiApplicationName), "application_id", config.applicationId, response.Id)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationApiApplicationName), "name", config.name, response.Name)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationApiApplicationName), "url", config.url, response.Url)
		if err != nil {
			return err
		}

		if config.description != "" {
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationApiApplicationName), "description", config.description, *response.Description)
			if err != nil {
				return err
			}
		}

		if config.additionalAllowedOrigins != nil {
			err = acctest.TestAttributesMatchStringSlice(resourceType, pointers.String(authenticationApiApplicationName), "additional_allowed_origins", config.additionalAllowedOrigins, response.AdditionalAllowedOrigins)
			if err != nil {
				return err
			}
		}

		if config.clientForRedirectlessModeRef != nil {
			err = acctest.TestAttributesMatchString(resourceType, pointers.String(authenticationApiApplicationName), "client_for_redirectless_mode_ref", config.clientForRedirectlessModeRef.Id, response.ClientForRedirectlessModeRef.Id)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckAuthenticationApiApplicationDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.AuthenticationApiAPI.DeleteApplication(ctx, authenticationApiApplicationId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("AuthenticationApiApplication", authenticationApiApplicationId)
	}
	return nil
}
