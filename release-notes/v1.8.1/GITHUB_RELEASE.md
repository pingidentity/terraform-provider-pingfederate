### BREAKING CHANGES

[1ad2e9c1](https://github.com/pingidentity/terraform-provider-pingfederate/commit/1ad2e9c1) `resource/pingfederate_redirect_validation`: Updated the `redirect_validation_local_settings.white_list` and `redirect_validation_local_settings.uri_allow_list` fields to be Lists instead of Sets. This change helps prevent errors due to known terraform issues with nested sets. This change will not require any change to existing configuration or state. These fields are now ordered instead of unordered, so future changes in ordering in configuration will result in planned changes with terraform. [#609](https://github.com/pingidentity/terraform-provider-pingfederate/pull/609)

### BUG FIXES

[1ad2e9c1](https://github.com/pingidentity/terraform-provider-pingfederate/commit/1ad2e9c1) `resource/pingfederate_redirect_validation`: Fixed an issue where boolean attributes inside `white_list` and `uri_allow_list` entries would appear to flip between plans after each apply. Related to a known terraform-plugin-framework bug with `Default` values in nested set attributes: [#867](https://github.com/hashicorp/terraform-plugin-framework/issues/867). [#609](https://github.com/pingidentity/terraform-provider-pingfederate/pull/609)

