package keypairssigningcsr_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

var fileDataInitial, fileDataUpdated, fileDataCa string

const signingCaId = "signingcsrtestca"

func TestAccKeypairsSigningCsrResponseResource(t *testing.T) {
	fileDataInitial = os.Getenv("PF_TF_ACC_TEST_SIGNING_CSR_RESPONSE_1")
	fileDataUpdated = os.Getenv("PF_TF_ACC_TEST_SIGNING_CSR_RESPONSE_2")
	fileDataCa = os.Getenv("PF_TF_ACC_TEST_CA_CERTIFICATE")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if fileDataInitial == "" {
				t.Fatal("PF_TF_ACC_TEST_SIGNING_CSR_RESPONSE_1 must be set for TestAccKeypairsSigningCsrResponseResource")
			}
			if fileDataUpdated == "" {
				t.Fatal("PF_TF_ACC_TEST_SIGNING_CSR_RESPONSE_2 must be set for TestAccKeypairsSigningCsrResponseResource")
			}
			if fileDataCa == "" {
				t.Fatal("PF_TF_ACC_TEST_CA_CERTIFICATE must be set for TestAccKeypairsSigningCsrResponseResource")
			}
			keypairsSigningCsrResponse_ImportCA(t)
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Initial CSR response import on create
				Config: keypairsSigningCsrResponse_HCL(fileDataInitial),
				Check:  keypairsSigningCsrResponse_CheckComputedValuesInitial(),
			},
			{
				// Importing a second CSR response
				Config: keypairsSigningCsrResponse_HCL(fileDataUpdated),
				Check:  keypairsSigningCsrResponse_CheckComputedValuesUpdated(),
			},
		},
	})
}

func keypairsSigningCsrResponse_ImportCA(t *testing.T) {
	testClient := acctest.TestClient()
	trustCaImportReq := testClient.CertificatesCaAPI.ImportTrustedCA(acctest.TestBasicAuthContext())
	trustCaImportReq = trustCaImportReq.Body(client.X509File{
		Id:       utils.Pointer(signingCaId),
		FileData: fileDataCa,
	})
	_, httpResp, err := testClient.CertificatesCaAPI.ImportTrustedCAExecute(trustCaImportReq)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 422) {
		errorMsg := "Failed to import test CA: " + err.Error()
		if httpResp != nil {
			body, internalErr := io.ReadAll(httpResp.Body)
			if internalErr == nil {
				errorMsg += " - Detail: " + string(body)
			}
		}
		t.Error(errorMsg)
		t.FailNow()
	}
}

func keypairsSigningCsrResponse_HCL(fileData string) string {
	return fmt.Sprintf(`
resource "pingfederate_keypairs_signing_csr_response" "example" {
  keypair_id = "419x9yg43rlawqwq9v6az997k"
  file_data  = "%s"
}
`, fileData)
}

func keypairsSigningCsrResponse_CheckComputedValuesInitial() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_signing_csr_response.example", "crypto_provider"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "id", "419x9yg43rlawqwq9v6az997k"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "serial_number", "169806312604756394519182484033336305508"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "subject_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "valid_from", "2024-09-20T19:44:51Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "expires", "2034-04-21T19:44:51Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "version", "3"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "sha1_fingerprint", "E938446F2F9DF707356192A70D105C43D5F0E797"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "sha256_fingerprint", "FF0F885E342BA337F8B44916ECD21D041DD787A29D117D4AB6A1AC121E27CAD7"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "status", "VALID"),
	)
}

func keypairsSigningCsrResponse_CheckComputedValuesUpdated() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckNoResourceAttr("pingfederate_keypairs_signing_csr_response.example", "crypto_provider"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "id", "419x9yg43rlawqwq9v6az997k"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "serial_number", "115908580996287481987637564242695711780"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "subject_dn", "C=US, O=CDR, OU=PING, L=AUSTIN, ST=TEXAS"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "issuer_dn", "CN=Example Authority, O=Example Corporation, C=US"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "valid_from", "2024-09-20T19:53:08Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "expires", "2034-04-21T19:53:08Z"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "version", "3"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "sha1_fingerprint", "255B0418A07C0B189FD810CE5A55B2B2561D937F"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "sha256_fingerprint", "DDEB7B4FB9D8389E3B29FD8F234CCDD39331EAAD4A896FD10F0CE02B6C446DDD"),
		resource.TestCheckResourceAttr("pingfederate_keypairs_signing_csr_response.example", "status", "VALID"),
	)
}
