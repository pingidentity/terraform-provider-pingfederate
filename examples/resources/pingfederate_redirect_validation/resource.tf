resource "pingfederate_redirect_validation" "redirectValidationExample" {
  redirect_validation_local_settings = {
    enable_target_resource_validation_for_sso           = true,
    enable_target_resource_validation_for_slo           = true,
    enable_target_resource_validation_for_idp_discovery = true,
    enable_in_error_resource_validation                 = true,
    white_list = [
      {
        target_resource_sso      = true,
        target_resource_slo      = true,
        in_error_resource        = true,
        idp_discovery            = true,
        valid_domain             = "bxretail.org",
        valid_path               = "/callback",
        allow_query_and_fragment = true,
        require_https            = true
      },
      {
        target_resource_sso      = false,
        target_resource_slo      = true,
        in_error_resource        = true,
        idp_discovery            = true,
        valid_domain             = "bxretail.org",
        valid_path               = "/redirect",
        allow_query_and_fragment = false,
        require_https            = true
      },
    ]
    uri_allow_list = [
      {
        allow_query_and_fragment = true
        idp_discovery            = false
        in_error_resource        = true
        target_resource_slo      = true
        target_resource_sso      = true
        valid_uri                = "https://auth.bxretail.org/*/callback"
      },
    ]
  }
  redirect_validation_partner_settings = {
    enable_wreply_validation_slo = true
  }
}
