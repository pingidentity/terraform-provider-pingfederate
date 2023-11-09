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
const fileData2 = "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"

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
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
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
				ResourceName:      "pingfederate_key_pair_signing_import." + resourceName,
				ImportStateId:     keyPairsSigningImportId,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccKeyPairsSigningImport(resourceName string, resourceModel keyPairsSigningImportResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_key_pair_signing_import" "%[1]s" {
  key_pair_signing_import_id = "%[2]s"
  file_data                  = "%[3]s"
  format                     = "%[4]s"
  password                   = "%[5]s"
}

data "pingfederate_key_pair_signing_import" "%[1]s" {
  id = pingfederate_key_pair_signing_import.%[1]s.id
}`, resourceName,
		resourceModel.id,
		resourceModel.fileData,
		resourceModel.format,
		resourceModel.password,
		fileData2,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedKeyPairsSigningImportAttributes(config keyPairsSigningImportResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "KeyPairsSigningImport"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()

		response, _, err := testClient.KeyPairsSigningAPI.GetSigningKeyPair(ctx, keyPairsSigningImportId).Execute()
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
	_, err := testClient.KeyPairsSigningAPI.DeleteSigningKeyPair(ctx, keyPairsSigningImportId).Execute()
	if err == nil {
		return acctest.ExpectedDestroyError("KeyPairsSigningImport", keyPairsSigningImportId)
	}
	return nil
}
