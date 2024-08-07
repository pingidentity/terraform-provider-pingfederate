resource "pingfederate_authentication_policy_contract" "example" {
  name = "User"
  extended_attributes = [
    { name = "email" },
    { name = "given_name" },
    { name = "family_name" }
  ]
}