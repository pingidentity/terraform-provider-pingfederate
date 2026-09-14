### BUG FIXES

- `resource/pingfederate_oauth_access_token_manager`: Fixed an issue where `sequence_number` could not be configured. The attribute is now Optional in addition to Computed, allowing the sequence number (embedded in the `pi.atm` claim of issued access tokens) to be pinned to a specific value. If not configured, a server-assigned value is used, and the previously-assigned value is retained if the attribute is later removed from configuration.
- Added validation to prevent inconsistent result errors due to empty `expression_criteria` and trailing newlines for `EXPRESSION` values in `attribute_contract_fulfillment`
- `resource/pingfederate_oauth_token_exchange_processor_policy`: Fixed an issue where creation could silently fail when many policies were created concurrently. The create is now verified and re-issued if the policy was not persisted.

