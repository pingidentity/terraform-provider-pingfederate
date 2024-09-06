---
page_title: "pingfederate_oauth_token_exchange_generator_group Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Resource to create and manage an OAuth 2.0 token exchange generator group.
---

# pingfederate_oauth_token_exchange_generator_group (Resource)

Resource to create and manage an OAuth 2.0 token exchange generator group.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `generator_mappings` (Attributes Set) A list of Token Generator mapping into an OAuth 2.0 Token Exchange requested token type. (see [below for nested schema](#nestedatt--generator_mappings))
- `group_id` (String) The Token Exchange Generator group ID. ID is unique.
- `name` (String) The Token Exchange Generator group name. Name is unique.

### Optional

- `resource_uris` (Set of String) The list of resource URI's which map to this Token Exchange Generator group.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--generator_mappings"></a>
### Nested Schema for `generator_mappings`

Required:

- `requested_token_type` (String) The Requested token type
- `token_generator` (Attributes) The Token Generator used to generate the requested token. (see [below for nested schema](#nestedatt--generator_mappings--token_generator))

Optional:

- `default_mapping` (Boolean) Whether this is the default Token Generator Mapping. Defaults to `false`.

<a id="nestedatt--generator_mappings--token_generator"></a>
### Nested Schema for `generator_mappings.token_generator`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "groupId" should be the id of the OAuth 2.0 Token Exchange Generator Group to be imported

```shell
terraform import pingfederate_oauth_token_exchange_generator_group.generatorGroup groupId
```