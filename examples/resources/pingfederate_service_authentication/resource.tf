resource "pingfederate_service_authentication" "serviceAuthentication" {
  attribute_query = {
    id            = "heuristics"
    shared_secret = "mysharedsecret"
  }
  jmx = {
    id            = "heuristics"
    shared_secret = "mysharedsecret"
  }
}