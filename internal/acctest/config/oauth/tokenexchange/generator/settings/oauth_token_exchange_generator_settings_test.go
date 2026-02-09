// Copyright © 2026 Ping Identity Corporation

package oauthtokenexchangegeneratorsettings_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

func TestAccOauthTokenExchangeGeneratorSettings(t *testing.T) {
	resourceName := "myOauthTokenExchangeGeneratorSettings"
	var testClient *configurationapi.APIClient
	ctx := acctest.TestBasicAuthContext()
	defaultGeneratorGroupId := acctest.ResourceIdGen()

	//TODO currently token exchange generator groups are not supported by the provider.
	// When they are, this should be created with terraform rather than direct API requests.
	defaultGroup := configurationapi.TokenExchangeGeneratorGroup{
		Name: defaultGeneratorGroupId,
		Id:   defaultGeneratorGroupId,
		GeneratorMappings: []configurationapi.TokenExchangeGeneratorMapping{
			{
				RequestedTokenType: "urn:ietf:params:oauth:token-type:saml2",
				DefaultMapping:     utils.Pointer(true),
				TokenGenerator: configurationapi.ResourceLink{
					Id: "tokengenerator",
				},
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			testClient = acctest.TestClient()
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccOauthTokenExchangeGeneratorSettingsEmpty(resourceName),
			},
			{
				PreConfig: func() {
					_, _, err := testClient.OauthTokenExchangeGeneratorAPI.CreateGroup(ctx).Body(defaultGroup).Execute()
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccOauthTokenExchangeGeneratorSettingsDefaultRef(resourceName, defaultGeneratorGroupId),
			},
			{
				// Test importing the resource
				Config:                               testAccOauthTokenExchangeGeneratorSettingsDefaultRef(resourceName, defaultGeneratorGroupId),
				ResourceName:                         "pingfederate_oauth_token_exchange_generator_settings." + resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "default_generator_group_ref.id",
			},
			{
				PreConfig: func() {
					_, err := testClient.OauthTokenExchangeGeneratorAPI.DeleteOauthTokenExchangeGroup(ctx, defaultGeneratorGroupId).Execute()
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testAccOauthTokenExchangeGeneratorSettingsEmpty(resourceName),
			},
		},
	})
}

func testAccOauthTokenExchangeGeneratorSettingsEmpty(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_token_exchange_generator_settings" "%s" {
}`, resourceName)
}

func testAccOauthTokenExchangeGeneratorSettingsDefaultRef(resourceName string, defaultRef string) string {
	return fmt.Sprintf(`
resource "pingfederate_oauth_token_exchange_generator_settings" "%[1]s" {
  default_generator_group_ref = {
    id = "%[2]s"
  }
}`, resourceName, defaultRef)
}
