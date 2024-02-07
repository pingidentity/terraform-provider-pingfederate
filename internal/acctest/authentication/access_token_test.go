package authentication_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest/authentication"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

type Response struct {
	AccessToken string `json:"access_token"`
}

func getAccessToken(t *testing.T) {
	//#nosec G101
	tokenRequestUrl := "https://localhost:9031/as/token.oauth2"

	//#nosec G402
	client := &http.Client{Transport: acctest.GetTransport()}
	envVars, errors := authentication.TestEnvVarSlice([]string{"PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID", "PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET"}, "access_token_test.go")
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		t.FailNow()
	}

	clientInfo := fmt.Sprintf("client_id=%s&grant_type=client_credentials&client_secret=%s&scope=email", envVars["PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID"], envVars["PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET"])
	jsonBodyReader := strings.NewReader(clientInfo)
	resp, err := client.Post(tokenRequestUrl, "application/x-www-form-urlencoded", jsonBodyReader)
	if err != nil {
		fmt.Sprintln("Error getting access token: ", err)
		t.FailNow()
	}

	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Sprintln("Error reading response body: ", readErr)
		t.FailNow()
	}

	var response Response
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		t.FailNow()
	}

	os.Unsetenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID")
	os.Unsetenv("PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET")
	os.Setenv("PINGFEDERATE_PROVIDER_ACCESS_TOKEN", response.AccessToken)
}

func TestAccATVirtualHostNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() { getAccessToken(t) },
				Config:    testAccATVirtualHostNames("virtualHostNames"),
				Check:     testAccATGetVirtualHostNames(),
			},
		},
	})
}

func testAccATVirtualHostNames(resourceName string) string {
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
func testAccATGetVirtualHostNames() resource.TestCheckFunc {
	test := func(s *terraform.State) error {
		testClient := acctest.TestClient()
		ctx := acctest.TestAccessTokenContext(os.Getenv("PINGFEDERATE_PROVIDER_ACCESS_TOKEN"))
		_, _, respErr := testClient.VirtualHostNamesAPI.GetVirtualHostNamesSettings(ctx).Execute()
		if respErr != nil {
			return respErr
		}

		return nil
	}
	return test
}
