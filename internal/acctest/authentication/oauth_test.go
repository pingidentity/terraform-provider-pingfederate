package auth_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func getOAuthEnvVars() {
	os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID")
	os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET")
	os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL")
	os.Getenv("PINGFEDERATE_PROVIDER_OAUTH_SCOPES")
}

func TestAccOAuthVirtualHostNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: getOAuthEnvVars,
				Config:    testAccOAuthVirtualHostNames("virtualHostNames"),
				Check:     testAccOAuthGetVirtualHostNames(),
			},
		},
	})
}

func testAccOAuthVirtualHostNames(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_virtual_host_names" "%[1]s" {
  virtual_host_names = %[2]s
}
data "pingfederate_virtual_host_names" "%[1]s" {
  depends_on = [pingfederate_virtual_host_names.%[1]s]
}`, resourceName,
		acctest.StringSliceToTerraformString([]string{"test"}),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccOAuthGetVirtualHostNames() resource.TestCheckFunc {
	test := func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestOauth2Context()
		_, _, err := testClient.VirtualHostNamesAPI.GetVirtualHostNamesSettings(ctx).Execute()
		if err != nil {
			return err
		}

		return nil
	}
	return test
}
