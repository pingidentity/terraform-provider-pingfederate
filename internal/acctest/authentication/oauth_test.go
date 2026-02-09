// Copyright © 2026 Ping Identity Corporation

package authentication_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/authentication"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccOAuthVirtualHostNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					authentication.TestEnvVarSlice([]string{"PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID", "PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET", "PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL", "PINGFEDERATE_PROVIDER_OAUTH_SCOPES"}, "oauth_test.go", t)
				},
				Config: testAccOAuthVirtualHostNames("virtualHostNames"),
				Check:  testAccOAuthGetVirtualHostNames(),
			},
		},
	})
}

func testAccOAuthVirtualHostNames(resourceName string) string {
	return fmt.Sprintf(`
resource "pingfederate_virtual_host_names" "%[1]s" {
  virtual_host_names = ["test"]
}
data "pingfederate_virtual_host_names" "%[1]s" {
  depends_on = [pingfederate_virtual_host_names.%[1]s]
}`, resourceName,
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
