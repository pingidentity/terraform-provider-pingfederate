resource "pingfederate_oauth_authorization_detail_processor" "detailProcessor" {
  processor_id = "myprocessor"
  configuration = {
    fields = [
      {
        name  = "field1"
        value = "value1"
      }
    ]
  }
  name = "My processor"
  plugin_descriptor_ref = {
    id = "com.example.AuthorizationDetailProcessor"
  }
}

resource "pingfederate_oauth_authorization_detail_type" "detailType" {
  type_id = "myDetailType"
  active  = true
  authorization_detail_processor_ref = {
    id = pingfederate_oauth_authorization_detail_processor.detailProcessor.processor_id
  }
  description = "This is my detail type"
  type        = "mytype"
}