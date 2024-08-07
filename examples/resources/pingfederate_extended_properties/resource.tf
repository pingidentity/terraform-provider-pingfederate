resource "pingfederate_extended_properties" "example" {
  items = [
    {
      name        = "Attribute 1"
      description = "My single valued extended attribute"
    },
    {
      name         = "Attribute 2"
      description  = "My multi-valued extended attribute"
      multi_valued = true
    },
  ]
}