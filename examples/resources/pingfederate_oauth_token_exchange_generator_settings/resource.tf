resource "pingfederate_oauth_token_exchange_generator_group" "generatorGroup" {
  group_id = "myGroup"
  generator_mappings = [
    {
      default_mapping      = true
      requested_token_type = "urn:ietf:params:oauth:token-type:saml2"
      token_generator = {
        id = pingfederate_sp_token_generator.tokenGenerator.generator_id
      }
    }
  ]
  name = "My generator group"
}

resource "pingfederate_oauth_token_exchange_generator_settings" "tokenExchangeGeneratorSettings" {
  default_generator_group_ref = {
    id = pingfederate_oauth_token_exchange_generator_group.generatorGroup.group_id
  }
}
