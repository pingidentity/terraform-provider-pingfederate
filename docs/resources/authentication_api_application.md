---
page_title: "pingfederate_authentication_api_application Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages an Authentication Api Application
---

# pingfederate_authentication_api_application (Resource)

Manages an Authentication Api Application

## Example Usage

```terraform
resource "pingfederate_authentication_api_application" "authenticationApiApplicationExample" {
  name        = "My Example Application"
  description = "My example application that has the authentication API widget embedded, or implements the authentication API directly."

  url = "https://bxretail.org"
  additional_allowed_origins = [
    "https://bxretail.org",
    "https://bxretail.org/*",
    "https://bxretail.org/cb/*",
    "https://auth.bxretail.org/*",
  ]

  client_for_redirectless_mode_ref = {
    id = pingfederate_oauth_client.oauthClientExample.id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The Authentication API Application Name. Name must be unique.
- `url` (String) The Authentication API Application redirect URL.

### Optional

- `additional_allowed_origins` (Set of String) The domain in the redirect URL is always whitelisted. This field contains a list of additional allowed origin URL's for cross-origin resource sharing.
- `application_id` (String) The persistent, unique ID for the Authentication API application. It can be any combination of `[a-zA-Z0-9._-]`. This property is system-assigned if not specified. This field is immutable and will trigger a replacement plan if changed.
- `client_for_redirectless_mode_ref` (Attributes) The client this application must use if it invokes the authentication API in redirectless mode. No client may be specified if `restrict_access_to_redirectless_mode` is `false` under `pingfederate_authentication_api_settings`. (see [below for nested schema](#nestedatt--client_for_redirectless_mode_ref))
- `description` (String) The Authentication API Application description.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--client_for_redirectless_mode_ref"></a>
### Nested Schema for `client_for_redirectless_mode_ref`

Required:

- `id` (String) The ID of the resource.

## Import

Import is supported using the following syntax:

~> "applicationId" should be the id of the Authentication Api Application to be imported

```shell
terraform import pingfederate_authentication_api_application.authenticationApiApplication applicationId
```