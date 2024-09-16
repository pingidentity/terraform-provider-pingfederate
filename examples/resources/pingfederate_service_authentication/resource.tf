resource "pingfederate_service_authentication" "serviceAuthentication" {
  attribute_query = {
    id            = "heuristics"
    shared_secret = var.attribute_query_service_shared_secret
  }
  jmx = {
    id            = "heuristics"
    shared_secret = var.jmx_service_shared_secret
  }
}