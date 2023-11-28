resource "pingfederate_oauth_open_id_connect_policy" "oauthOIDCPolicyExample" {
		id:                        oauthOpenIdConnectPoliciesId,
		name:                      "initialName",
		includeOptionalAttributes: false,
}

data "pingfederate_oauth_open_id_connect_policy" "myOauthOIDCPolicyExample" {

}