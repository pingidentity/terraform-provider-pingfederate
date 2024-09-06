resource "pingfederate_keypairs_signing_key" "signingKey" {
  key_id    = "signingkey"
  file_data = filebase64("./assets/signingkey.p12")
  password  = var.signing_key_password
  format    = "PKCS12"
}

resource "pingfederate_sp_token_generator" "tokenGenerator" {
  generator_id = "myGenerator"
  attribute_contract = {
    core_attributes = [
      {
        name = "SAML_SUBJECT"
      }
    ]
  }
  configuration = {
    fields = [
      {
        name  = "Minutes Before"
        value = "60"
      },
      {
        name  = "Minutes After"
        value = "60"
      },
      {
        name  = "Issuer"
        value = "issuer"
      },
      {
        name  = "Signing Certificate"
        value = pingfederate_keypairs_signing_key.signingKey.key_id
      },
      {
        name  = "Signing Algorithm"
        value = "SHA1"
      },
      {
        name  = "Include Certificate in KeyInfo"
        value = "false"
      },
      {
        name  = "Include Raw Key in KeyValue"
        value = "false"
      },
      {
        name  = "Audience"
        value = "audience"
      },
      {
        name  = "Confirmation Method"
        value = "urn:oasis:names:tc:SAML:2.0:cm:sender-vouches"
      }
    ]
    tables = []
  }
  name = "My token generator"
  plugin_descriptor_ref = {
    id = "org.sourceid.wstrust.generator.saml.Saml20TokenGenerator"
  }
}

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