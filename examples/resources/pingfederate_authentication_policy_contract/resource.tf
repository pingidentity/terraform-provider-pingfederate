resource "pingfederate_authentication_policy_contract" "myAuthenticationPolicyContract" {
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
}
