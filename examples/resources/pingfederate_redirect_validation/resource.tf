resource "pingfederate_redirect_validation" "redirectValidationExample" {
	redirect_validation_local_settings = {
    enable_target_resource_validation_for_sso = true,
    enable_target_resource_validation_for_slo = true,
    enable_target_resource_validation_for_idp_discovery = true,
    enable_in_error_resource_validation = true,
    white_list = [
      {
        target_resource_sso = true,
        target_resource_slo = true,
        in_error_resource = true,
        idp_discovery = true,
        valid_domain = "example.com",
        valid_path = "/path",
        allow_query_and_fragment = true,
        require_https = true
      },
      {
        target_resource_sso = false,
        target_resource_slo = true,
        in_error_resource = true,
        idp_discovery = true,
        valid_domain = "example2.com",
        valid_path = "/path2",
        allow_query_and_fragment = false,
        require_https = true
      }
    ]
  }
  redirect_validation_partner_settings = {
    enable_wreply_validation_slo = true
  }
}
