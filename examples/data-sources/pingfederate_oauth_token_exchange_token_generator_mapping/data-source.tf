data "pingfederate_oauth_token_exchange_token_generator_mapping" "oauthTokenExchangeTokenGeneratorMappingsExample" {
  mapping_id = "${pingfederate_oauth_token_exchange_token_generator_mapping.oauthTokenExchangeTokenGeneratorMappingsExample.source_id}|${pingfederate_oauth_token_exchange_token_generator_mapping.oauthTokenExchangeTokenGeneratorMappingsExample.target_id}"
}
