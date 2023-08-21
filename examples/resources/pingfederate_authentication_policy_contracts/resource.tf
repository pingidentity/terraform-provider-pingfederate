resource "pingfederate_authentication_policy_contracts" "authenticationPolicyContractsExample" {
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
}
