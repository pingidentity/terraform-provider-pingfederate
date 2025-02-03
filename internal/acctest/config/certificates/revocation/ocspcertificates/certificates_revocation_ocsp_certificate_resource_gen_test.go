// Copyright © 2025 Ping Identity Corporation

// Code generated by ping-terraform-plugin-framework-generator

package certificatesrevocationocspcertificates_test

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

const certificatesRevocationOcspCertificateCertificateId = "ocspcertid"

var fileDataInitial, fileDataUpdated string

func TestAccCertificatesRevocationOcspCertificate_RemovalDrift(t *testing.T) {
	fileDataInitial = os.Getenv("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1")
	fileDataUpdated = os.Getenv("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_2")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if fileDataInitial == "" {
				t.Fatal("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1 must be set for TestAccCertificatesRevocationOcspCertificate_RemovalDrift")
			}
			if fileDataUpdated == "" {
				t.Fatal("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_2 must be set for TestAccCertificatesRevocationOcspCertificate_RemovalDrift")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: certificatesRevocationOcspCertificate_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: certificatesRevocationOcspCertificate_InitialHCL(),
			},
			{
				// Delete the resource on the service, outside of terraform, verify that a non-empty plan is generated
				PreConfig: func() {
					certificatesRevocationOcspCertificate_Delete(t)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCertificatesRevocationOcspCertificate_MinimalMaximal(t *testing.T) {
	fileDataInitial = os.Getenv("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1")
	fileDataUpdated = os.Getenv("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_2")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.ConfigurationPreCheck(t)
			if fileDataInitial == "" {
				t.Fatal("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_1 must be set for TestAccCertificatesRevocationOcspCertificate_MinimalMaximal")
			}
			if fileDataUpdated == "" {
				t.Fatal("PF_TF_ACC_TEST_CERTIFICATE_CA_FILE_DATA_2 must be set for TestAccCertificatesRevocationOcspCertificate_MinimalMaximal")
			}
		},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		CheckDestroy: certificatesRevocationOcspCertificate_CheckDestroy,
		Steps: []resource.TestStep{
			{
				// Create the resource with a minimal model
				Config: certificatesRevocationOcspCertificate_InitialHCL(),
				Check:  certificatesRevocationOcspCertificate_CheckComputedValuesInitial(),
			},
			{
				// Delete the minimal model
				Config:  certificatesRevocationOcspCertificate_InitialHCL(),
				Destroy: true,
			},
			{
				// Re-create with a complete model
				Config: certificatesRevocationOcspCertificate_UpdatedHCL(),
				Check:  certificatesRevocationOcspCertificate_CheckComputedValuesUpdated(),
			},
			{
				// Back to minimal model
				Config: certificatesRevocationOcspCertificate_InitialHCL(),
				Check:  certificatesRevocationOcspCertificate_CheckComputedValuesInitial(),
			},
			{
				// Back to complete model
				Config: certificatesRevocationOcspCertificate_UpdatedHCL(),
				Check:  certificatesRevocationOcspCertificate_CheckComputedValuesUpdated(),
			},
			{
				// Test importing the resource
				Config:                               certificatesRevocationOcspCertificate_UpdatedHCL(),
				ResourceName:                         "pingfederate_certificates_revocation_ocsp_certificate.example",
				ImportStateId:                        certificatesRevocationOcspCertificateCertificateId,
				ImportStateVerifyIdentifierAttribute: "certificate_id",
				ImportState:                          true,
				ImportStateVerify:                    true,
				// File data is not returned by the API
				ImportStateVerifyIgnore: []string{
					"file_data",
				},
			},
		},
	})
}

// Minimal HCL with only required values set
func certificatesRevocationOcspCertificate_InitialHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_certificates_revocation_ocsp_certificate" "example" {
  certificate_id = "%s"
  file_data      = "%s"
}
`, certificatesRevocationOcspCertificateCertificateId, fileDataInitial)
}

// Maximal HCL with all values set where possible
func certificatesRevocationOcspCertificate_UpdatedHCL() string {
	return fmt.Sprintf(`
resource "pingfederate_certificates_revocation_ocsp_certificate" "example" {
  certificate_id = "%s"
  file_data      = "%s"
}
`, certificatesRevocationOcspCertificateCertificateId, fileDataUpdated)
}

// Validate any computed values when applying initial HCL
func certificatesRevocationOcspCertificate_CheckComputedValuesInitial() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "expires", "2024-05-29T15:59:19Z"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "issuer_dn", "EMAILADDRESS=test@gmail.com, CN=terraformtest, OU=Devops, O=ping Identity, L=san Jose, ST=SJC, C=US"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "key_size", "2048"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "serial_number", "16677565866115840610"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "sha1_fingerprint", "A98434F2CD96AF202E50DDDD8FD6D9354CAC2B80"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "sha256_fingerprint", "0D82B801AAA7CCE20C752EEA02D02296AF258EAA3FF7D565164A0F880EEA910B"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "signature_algorithm", "SHA256withRSA"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "status", "EXPIRED"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "subject_dn", "EMAILADDRESS=test@gmail.com, CN=terraformtest, OU=Devops, O=ping Identity, L=san Jose, ST=SJC, C=US"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "valid_from", "2023-05-30T15:59:19Z"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "version", "1"),
	)
}

// Validate any computed values when applying updated HCL
func certificatesRevocationOcspCertificate_CheckComputedValuesUpdated() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "expires", "2018-03-18T15:40:19Z"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "issuer_dn", "CN=foo, OU=foo, O=foo, L=foo, ST=foo, C=US"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "key_algorithm", "RSA"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "key_size", "1024"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "serial_number", "13299021239615735660"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "sha1_fingerprint", "42DA9CF8F12D70582D97B937E09E667C83A9A0E4"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "sha256_fingerprint", "0547159D0FEE3C5C332518796C441957D8920B7144FE6D2205B6D2E6A6814B9A"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "signature_algorithm", "SHA1withRSA"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "status", "EXPIRED"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "subject_alternative_names.#", "0"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "subject_dn", "CN=foo, OU=foo, O=foo, L=foo, ST=foo, C=US"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "valid_from", "2013-03-19T15:40:19Z"),
		resource.TestCheckResourceAttr("pingfederate_certificates_revocation_ocsp_certificate.example", "version", "3"),
	)
}

// Delete the resource
func certificatesRevocationOcspCertificate_Delete(t *testing.T) {
	testClient := acctest.TestClient()
	_, err := testClient.CertificatesRevocationAPI.DeleteOcspCertificateById(acctest.TestBasicAuthContext(), certificatesRevocationOcspCertificateCertificateId).Execute()
	if err != nil {
		t.Fatalf("Failed to delete config: %v", err)
	}
}

// Test that any objects created by the test are destroyed
func certificatesRevocationOcspCertificate_CheckDestroy(s *terraform.State) error {
	testClient := acctest.TestClient()
	_, err := testClient.CertificatesRevocationAPI.DeleteOcspCertificateById(acctest.TestBasicAuthContext(), certificatesRevocationOcspCertificateCertificateId).Execute()
	if err == nil {
		return fmt.Errorf("certificates_revocation_ocsp_certificate still exists after tests. Expected it to be destroyed")
	}
	return nil
}
