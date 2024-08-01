// Code generated by ping-terraform-plugin-framework-generator

package keypairssigningkeyblabla_test

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

const keypairsSigningKeyGenerateKeyId = "keypairssigninggenkeyid"
const keypairsSigningKeyImportKeyId = "keypairssigningimpkeyid"

func TestAccKeypairsSigningKey_RemovalDrift(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: keypairsSigningKey_GenerateCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: keypairsSigningKey_GenerateMinimalHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					keypairsSigningKey_Delete(t, keypairsSigningKeyGenerateKeyId)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKeypairsSigningKey_GenerateMinimalMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: keypairsSigningKey_GenerateCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: keypairsSigningKey_GenerateMinimalHCL(),
				Check:  keypairsSigningKey_CheckComputedValuesGenerateMinimal(),
			},
			{
				// Update with a complete model - this should cause a full replacement
				Config: keypairsSigningKey_GenerateCompleteHCL(),
				Check:  keypairsSigningKey_CheckComputedValuesGenerateComplete(),
			},
			{
				// Test importing the resource
				Config:                               keypairsSigningKey_GenerateCompleteHCL(),
				ResourceName:                         "pingfederate_keypairs_signing_key.example",
				ImportStateId:                        keypairsSigningKeyGenerateKeyId,
				ImportStateVerifyIdentifierAttribute: "key_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// Some attributes are only used on generation/import, and aren't returned individually from the API
				ImportStateVerifyIgnore: []string{
					"city",
					"common_name",
					"country",
					"organization",
					"organization_unit",
					"state",
					"valid_days",
				},
			},
		},
	})
}

var fileDataInitial, fileDataUpdated string

func TestAccKeypairsSigningKey_FileDataMinimalMaximal(t *testing.T) {
	fileDataInitial = os.Getenv("PF_TF_ACC_TEST_SIGNING_KEY_KEYSTORE_1")
	fileDataUpdated = os.Getenv("PF_TF_ACC_TEST_SIGNING_KEY_KEYSTORE_2")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if fileDataInitial == "" {
				t.Fatal("PF_TF_ACC_TEST_SIGNING_KEY_KEYSTORE_1 must be set for TestAccKeypairsSigningKey_FileDataMinimalMaximal")
			}
			if fileDataUpdated == "" {
				t.Fatal("PF_TF_ACC_TEST_SIGNING_KEY_KEYSTORE_2 must be set for TestAccKeypairsSigningKey_FileDataMinimalMaximal")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: keypairsSigningKey_ImportCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: keypairsSigningKey_ImportMinimalHCL(),
				Check:  keypairsSigningKey_CheckComputedValuesImportMinimal(),
			},
			{
				// Update with a complete model - this should cause a full replacement
				Config: keypairsSigningKey_ImportCompleteHCL(),
				Check:  keypairsSigningKey_CheckComputedValuesImportComplete(),
			},
			{
				// Test importing the resource
				Config:                               keypairsSigningKey_ImportCompleteHCL(),
				ResourceName:                         "pingfederate_keypairs_signing_key.example",
				ImportStateId:                        keypairsSigningKeyImportKeyId,
				ImportStateVerifyIdentifierAttribute: "key_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// Some attributes are only used on generation/import, and aren't returned individually from the API
				ImportStateVerifyIgnore: []string{
					"file_data",
					"password",
					"format",
				},
			},
		},
	})
}

// Minimal HCL with only required values set
func keypairsSigningKey_GenerateMinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_signing_key" "example" {
  key_id        = "%s"
  common_name   = "Example"
  country       = "US"
  key_algorithm = "RSA"
  organization  = "Ping Identity"
  valid_days    = 365
}

data "pingfederate_keypairs_signing_key" "example" {
  depends_on = [pingfederate_keypairs_signing_key.example]
  key_id = pingfederate_keypairs_signing_key.example.key_id
}
`, keypairsSigningKeyGenerateKeyId)
}

// Maximal HCL with all values set where possible
func keypairsSigningKey_GenerateCompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_signing_key" "example" {
  key_id                    = "%s"
  city                      = "Austin"
  common_name               = "Example"
  country                   = "US"
  key_algorithm             = "RSA"
  key_size                  = 2048
  organization              = "Ping Identity"
  organization_unit         = "Engineering"
  signature_algorithm       = "SHA256withRSA"
  state                     = "Texas"
  subject_alternative_names = ["example.com"]
  valid_days                = 365
}

data "pingfederate_keypairs_signing_key" "example" {
  depends_on = [pingfederate_keypairs_signing_key.example]
  key_id = pingfederate_keypairs_signing_key.example.key_id
}
`, keypairsSigningKeyGenerateKeyId)
}

// Minimal HCL with only required values set
func keypairsSigningKey_ImportMinimalHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_signing_key" "example" {
  key_id    = "%s"
  file_data = "%s"
  password  = "2FederateM0re"
}

data "pingfederate_keypairs_signing_key" "example" {
  depends_on = [pingfederate_keypairs_signing_key.example]
  key_id = pingfederate_keypairs_signing_key.example.key_id
}
`, keypairsSigningKeyImportKeyId, fileDataInitial)
}

// Maximal HCL with all values set where possible
func keypairsSigningKey_ImportCompleteHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_signing_key" "example" {
  key_id    = "%s"
  file_data = "%s"
  password  = "2FederateM0re"
  format    = "PKCS12"
}

data "pingfederate_keypairs_signing_key" "example" {
  depends_on = [pingfederate_keypairs_signing_key.example]
  key_id = pingfederate_keypairs_signing_key.example.key_id
}
`, keypairsSigningKeyImportKeyId, fileDataUpdated)
}

// Validate any computed values when applying minimal generated key HCL
func keypairsSigningKey_CheckComputedValuesGenerateMinimal() resource.TestCheckFunc {
	testChecks := []resource.TestCheckFunc{}
	for _, prefix := range []string{"", "data."} {
		testChecks = append(testChecks,
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "expires"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "issuer_dn", "CN=Example, O=Ping Identity, C=US"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "key_size", "2048"),
			resource.TestCheckNoResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "rotation_settings"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "serial_number"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "sha1_fingerprint"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "sha256_fingerprint"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "signature_algorithm", "SHA256withRSA"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "status", "VALID"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_alternative_names.#", "0"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_dn", "CN=Example, O=Ping Identity, C=US"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "valid_from"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "version", "3"),
		)
	}

	return resource.ComposeTestCheckFunc(
		testChecks...,
	)
}

// Validate any computed values when applying complete generated key HCL
func keypairsSigningKey_CheckComputedValuesGenerateComplete() resource.TestCheckFunc {
	testChecks := []resource.TestCheckFunc{}
	for _, prefix := range []string{"", "data."} {
		testChecks = append(testChecks,
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "expires"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "issuer_dn", "CN=Example, OU=Engineering, O=Ping Identity, L=Austin, ST=Texas, C=US"),
			resource.TestCheckNoResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "rotation_settings"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "serial_number"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "sha1_fingerprint"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "sha256_fingerprint"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "signature_algorithm", "SHA256withRSA"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "status", "VALID"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_alternative_names.0", "example.com"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_dn", "CN=Example, OU=Engineering, O=Ping Identity, L=Austin, ST=Texas, C=US"),
			resource.TestCheckResourceAttrSet(prefix+"pingfederate_keypairs_signing_key.example", "valid_from"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "version", "3"),
		)
	}

	return resource.ComposeTestCheckFunc(
		testChecks...,
	)
}

// Validate any computed values when applying minimal HCL
func keypairsSigningKey_CheckComputedValuesImportMinimal() resource.TestCheckFunc {
	testChecks := []resource.TestCheckFunc{}
	for _, prefix := range []string{"", "data."} {
		testChecks = append(testChecks,
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "expires", "2044-07-24T15:46:27Z"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "key_size", "2048"),
			resource.TestCheckNoResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "rotation_settings"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "serial_number", "28463092959443571178990831419139562736"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "sha1_fingerprint", "1C83D0C571A1AE934C3C2A4BF7BDC541974497E5"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "sha256_fingerprint", "B9A2940E5E5E06AC2852DD0A32B7192876C3B194577155CE58E1AD5234375EB7"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "signature_algorithm", "SHA256withRSA"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "status", "VALID"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_alternative_names.#", "0"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_dn", "CN=Example Authority, O=Example Corporation, C=US"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "valid_from", "2024-07-29T15:46:27Z"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "version", "3"),
		)
	}

	return resource.ComposeTestCheckFunc(
		testChecks...,
	)
}

// Validate any computed values when applying complete HCL
func keypairsSigningKey_CheckComputedValuesImportComplete() resource.TestCheckFunc {
	testChecks := []resource.TestCheckFunc{}
	for _, prefix := range []string{"", "data."} {
		testChecks = append(testChecks,
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "expires", "2025-08-01T15:16:44Z"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "issuer_dn", "CN=Another Authority, O=Example Corporation, C=US"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "key_size", "2048"),
			resource.TestCheckNoResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "rotation_settings"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "serial_number", "34314007937343527069893005115224475439"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "sha1_fingerprint", "60CB3F8861673E1E814D87D84C8FADDDC37AE270"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "sha256_fingerprint", "8AA7D3C77D5053A9C8781D4F3E123712667E6B9A3E103DB74D035D2751695938"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "signature_algorithm", "SHA256withRSA"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "status", "VALID"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_alternative_names.#", "0"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "subject_dn", "CN=Another Authority, O=Example Corporation, C=US"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "valid_from", "2024-08-01T15:16:44Z"),
			resource.TestCheckResourceAttr(prefix+"pingfederate_keypairs_signing_key.example", "version", "3"),
		)
	}

	return resource.ComposeTestCheckFunc(
		testChecks...,
	)
}

// Delete the resource
func keypairsSigningKey_Delete(t *testing.T, keyId string) {
	testClient := acctest.TestClient()
	_, err := testClient.KeyPairsSigningAPI.DeleteSigningKeyPair(acctest.TestBasicAuthContext(), keyId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func keypairsSigningKey_GenerateCheckDestroy(s *terraform.State) error {
	return keypairsSigningKey_CheckDestroy(s, keypairsSigningKeyGenerateKeyId)
}

func keypairsSigningKey_ImportCheckDestroy(s *terraform.State) error {
	return keypairsSigningKey_CheckDestroy(s, keypairsSigningKeyImportKeyId)
}

func keypairsSigningKey_CheckDestroy(s *terraform.State, keyId string) error {
	testClient := acctest.TestClient()
	_, err := testClient.KeyPairsSigningAPI.DeleteSigningKeyPair(acctest.TestBasicAuthContext(), keyId).Execute()
	if err == nil {
		return fmt.Errorf("keypairs_signing_key still exists after tests. Expected it to be destroyed")
	}
	return nil
}
