resource "pingfederate_extended_properties" "extendedPropertiesExample" {
  items = [
    {
      name         = "extended_attribute"
      description  = "attribute_1_description"
      multi_valued = false
    },
    {
      name         = "extended_attribute_3"
      description  = "attribute_3_description"
      multi_valued = true
    }
  ]
}
