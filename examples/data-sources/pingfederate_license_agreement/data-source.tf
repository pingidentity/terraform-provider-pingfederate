resource "pingfederate_license_agreement" "myLicenseAgreement" {
  accepted = true
}

data "pingfederate_license_agreement" "myLicenseAgreement" {
  depends_on = [
    pingfederate_license_agreement.myLicenseAgreement
  ]
}
resource "pingfederate_license_agreement" "licenseAgreementExample" {
  accepted = data.pingfederate_license_agreement.myLicenseAgreement.accepted
}