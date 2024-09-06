resource "pingfederate_oauth_authorization_detail_processor" "processor" {
  processor_id = "detailProcessor"
  configuration = {
    fields = [
      {
        name  = "field1",
        value = "value1"
      }
    ]
  }
  name = "Detail Processor"
  plugin_descriptor_ref = {
    id = "com.example.MyAuthorizationDetailProcessor"
  }
}