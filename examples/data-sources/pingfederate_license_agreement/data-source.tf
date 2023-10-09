terraform {
  required_version = ">=1.1"
  required_providers {
    pingfederate = {
      version = "~> 0.0.1"
      source  = "pingidentity/pingfederate"
    }
  }
}

provider "pingfederate" {
  username               = "administrator"
  password               = "2FederateM0re"
  https_host             = "https://localhost:9999"
  insecure_trust_all_tls = true
}

resource "pingfederate_license_agreement" "myLicenseAgreement" {
  accepted = true
}

data "pingfederate_license_agreement" "myLicenseAgreement" {
  accepted = pingfederate_license_agreement.myLicenseAgreement.accepted
}
resource "pingfederate_license_agreement" "licenseAgreementExample" {
  accepted = data.pingfederate_license_agreement.myLicenseAgreement.accepted
}