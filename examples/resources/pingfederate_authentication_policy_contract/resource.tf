resource "pingfederate_authentication_policy_contract" "authenticationPolicyContractExample" {
  contract_id         = "MyAuthenticationPolicyContract"
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
}
