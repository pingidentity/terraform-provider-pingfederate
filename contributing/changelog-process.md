# Changelog Process

This project maintains a manually curated [CHANGELOG.md](../CHANGELOG.md) file in the root directory. This document outlines the format, patterns, and rules for updating the changelog.

## Changelog Format and Structure

The changelog follows a consistent format based on the existing patterns in the project:

### Version Header Format
```markdown
# v<major>.<minor>.<patch> <Month> <Day>, <Year>
```

Examples:
- `# v1.6.2 September 19, 2025`
- `# v1.0.0 October 1, 2024`

### Section Order and Structure

Each version entry should follow this section order (only include sections that have content):

1. **Breaking changes** (if any)
2. **Enhancements** (if any)
3. **Resources** (for new resources)
4. **Data Sources** (for new data sources)  
5. **Bug fixes** (if any)
6. **Deprecated** (if any)
7. **Notes** (for dependency updates, etc.)

### Section Headers and Content

#### Breaking Changes
```markdown
### Breaking changes
* Description of breaking change with link to PR. ([#PR_NUMBER](PR_URL))
```

#### Enhancements
```markdown
### Enhancements
* Description of enhancement with link to PR. ([#PR_NUMBER](PR_URL))
```

#### New Resources
```markdown
### Resources
* **New Resource:** `resource_name` ([#PR_NUMBER](PR_URL))
```

#### New Data Sources
```markdown
### Data Sources
* **New Data Source:** `data_source_name` ([#PR_NUMBER](PR_URL))
```

#### Bug Fixes
```markdown
### Bug fixes
* Description of bug fix with link to PR. ([#PR_NUMBER](PR_URL))
```

#### Deprecations
```markdown
### Deprecated
* Description of deprecated feature with replacement guidance and link to PR. ([#PR_NUMBER](PR_URL))
```

#### Notes
```markdown
### Notes
* Description of dependency updates, Go version bumps, etc. ([#PR_NUMBER](PR_URL))
```

## Content Guidelines

### Description Patterns

**Bug Fixes:**
- Start with "Fixed" or describe the resolution
- Be specific about what was fixed
- Include affected resource names when relevant

Examples:
- `Fixed an inconsistent result error that would occur when configuring OAuth Clients...`
- `Fixed the required 'file_data' field not being written to state on import...`

**Enhancements:**
- Start with "Added" for new features
- Be descriptive about the improvement
- Include version information when adding support for new PingFederate versions

Examples:
- `Added support for PingFederate '12.3.0' and implemented new attributes...`
- `Added new 'configuration.sensitive_fields' attributes...`

**Breaking Changes:**
- Clearly state what changed
- Include migration guidance when possible
- Reference upgrade guides when available

Examples:
- `Removed support for PingFederate '11.2.x', in accordance with Ping's end of life policy.`
- `Marked 'sp_idp_connection' resource attributes as required...`

**New Resources/Data Sources:**
- Use exact resource/data source name in backticks
- No additional description needed beyond the name

**Deprecations:**
- State what is deprecated
- Provide replacement guidance
- Mention when it will be removed

Example:
- `The 'pingfederate_key_pair_ssl_server_import' resource has been renamed. Use 'pingfederate_keypairs_ssl_server_key' instead.`

### Link Format

All PR references should use this format:
```markdown
([#PR_NUMBER](https://github.com/pingidentity/terraform-provider-pingfederate/pull/PR_NUMBER))
```

Note: Some older entries have malformed URLs with extra `([https]` - new entries should use the correct format above.

## Process for Adding Changelog Entries

### For each PR/release:

1. **Determine the appropriate version number** following semantic versioning
2. **Add a new version header** at the top of the file
3. **Organize changes into appropriate sections** following the order above
4. **Write clear, descriptive entries** following the content patterns
5. **Include proper PR links** for all entries
6. **Review for consistency** with existing entries

### Version Numbering Guidelines

- **Major version (x.0.0)**: Breaking changes, major milestones
- **Minor version (1.x.0)**: New features, enhancements, new resources
- **Patch version (1.1.x)**: Bug fixes, minor improvements

### Quality Standards

- Use consistent formatting and language
- Ensure all changes are documented
- Group related changes logically
- Maintain chronological order (newest first)
- Include all relevant PR links
- Use proper grammar and spelling

## Example Entry

```markdown
# v1.7.0 October 15, 2025
### Enhancements
* Added support for PingFederate `12.4.0` and implemented new attributes for the new version. ([#550](https://github.com/pingidentity/terraform-provider-pingfederate/pull/550))

### Resources
* **New Resource:** `pingfederate_example_resource` ([#545](https://github.com/pingidentity/terraform-provider-pingfederate/pull/545))

### Bug fixes
* Fixed an issue where the `example_attribute` in `pingfederate_example_resource` could cause repeated plans after a successful `terraform apply`. ([#548](https://github.com/pingidentity/terraform-provider-pingfederate/pull/548))

### Deprecated
* The `pingfederate_old_resource` resource has been deprecated. Use `pingfederate_new_resource` instead. `pingfederate_old_resource` will be removed in a future release. ([#549](https://github.com/pingidentity/terraform-provider-pingfederate/pull/549))
```