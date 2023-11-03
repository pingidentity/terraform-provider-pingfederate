resource "pingfederate_virtual_host_names" "myVirtualHostNamesExample" {
  virtual_host_names = ["example1", "example2"]
}

data "pingfederate_virtual_host_names" "myVirtualHostNamesExample" {
  virtual_host_names = pingfederate_virtual_host_names.myVirtualHostNamesExample.virtual_host_names
}
