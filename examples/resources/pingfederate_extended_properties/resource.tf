resource "pingfederate_extended_properties" "extendedPropertiesExample" {
  items = [
    {
      name         = "extended_attribute_1"
      description  = "attribute_1_description"
      multi_valued = false
    },
    {
      name         = "extended_attribute_2"
      description  = "attribute_2_description"
      multi_valued = true
    }
  ]
}
