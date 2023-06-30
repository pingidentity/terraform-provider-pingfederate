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

const keyPairsSigningImportId = "2"
const fileData = "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"
const password = "2FederateM0re"
const format = "PKCS12"

// Attributes to test with. Add optional properties to test here if desired.
type keyPairsSigningImportResourceModel struct {
	id       string
	fileData string
	format   string
	password string
}

func TestAccKeyPairsSigningImport(t *testing.T) {
	resourceName := "myKeyPairsSigningImport"
	initialResourceModel := keyPairsSigningImportResourceModel{
		id:       keyPairsSigningImportId,
		fileData: fileData,
		format:   format,
		password: password,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		CheckDestroy: testAccCheckKeyPairsSigningImportDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyPairsSigningImport(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedKeyPairsSigningImportAttributes(initialResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccKeyPairsSigningImport(resourceName, initialResourceModel),
				ResourceName:      "pingfederate_key_pairs_signing_import." + resourceName,
				ImportStateId:     keyPairsSigningImportId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccKeyPairsSigningImport(resourceName string, resourceModel keyPairsSigningImportResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_key_pairs_signing_import" "%[1]s" {
  id        = "%[2]s"
  file_data = "%[3]s"
  format    = "%[4]s"
  password  = "%[5]s"
}`, resourceName,
		resourceModel.id,
		resourceModel.fileData,
		resourceModel.format,
		resourceModel.password,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedKeyPairsSigningImportAttributes(config keyPairsSigningImportResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "KeyPairsSigningImport"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()

		response, _, err := testClient.KeyPairsSigningApi.GetSigningKeyPair(ctx, keyPairsSigningImportId).Execute()
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchString(resourceType, &config.id, "id", config.id, *response.Id)
		if err != nil {
			return err
		}

		return nil
	}
}

// Test that any objects created by the test are destroyed
func testAccCheckKeyPairsSigningImportDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	ctx := acctest.TestBasicAuthContext()
	_, err := testClient.KeyPairsSigningApi.DeleteSigningKeyPair(ctx, keyPairsSigningImportId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("KeyPairsSigningImport", keyPairsSigningImportId)
	}
	return nil
}
