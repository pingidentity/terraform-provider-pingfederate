# Pull Request Checklist

The following provides the steps to check/run to prepare for creating a PR to the `main` branch. PRs that follow these checklists will merge faster than PRs that do not.

*Note: This checklist is designed to support both human contributors and automated code review tools.*

## For Automated Code Review

This checklist includes specific verification criteria marked with *Verification* that can be programmatically checked to support both manual and automated review processes

## PR Planning & Structure

- [ ] **PR Scope**. To ensure maintainer reviews are as quick and efficient as possible, please separate support for entities into separate PRs. For example, support for a `pingfederate_oauth_client` resource and data source can go in the same PR, however support for `pingfederate_oauth_client` and `pingfederate_idp_adapter` should be separated. It's acceptable to merge support for different entities into the same PR where structural changes are being made.
  - *Verification*: Check that files modified are logically related (same functional area directory, related functionality)

- [ ] **PR Title**. To assist the maintainers in assessing PRs for priority, please provide a descriptive title of the functionality being supported. For example: `Add support for OAuth Access Token Managers`
  - *Verification*: Title should be descriptive and match the type of changes (Add/Update/Fix/Remove)

- [ ] **PR Description**. Please follow the provided PR description template and check relevant boxes. Include a clear description of:
  - What functionality is being added/changed
  - Why the change is needed (e.g., to fix an issue - include the issue number as reference)
  - Any breaking changes
  - *Verification*: Check that PR description template boxes are completed and description sections are filled

## Code Development

### Architecture & Design

- [ ] **Code implementation**. New code should adhere to the [Provider Design](provider-design.md) guide for implementation patterns.
  - *Verification*: 
    - New resources are in `internal/resource/config/<functional_area>/resource_*.go`
    - New data sources are in `internal/datasource/<functional_area>/data_source_*.go`
    - Resources follow naming convention `pingfederate_<name>` aligned with API endpoints
    - Resources are registered in `internal/provider/provider.go`

- [ ] **PingFederate Go Client Usage**. All PingFederate API interactions must use the PingFederate Go Client rather than direct API calls
  - *Verification*: 
    - No direct HTTP calls to PingFederate APIs (check for `http.Client`, `http.Get`, `http.Post`, etc.)
    - Uses PingFederate Go Client from `github.com/pingidentity/pingfederate-go-client/`
    - Follows established patterns for API error handling and response processing

### Code Quality

- [ ] **Dependencies Check**. Ensure go.mod and go.sum are properly maintained:

```shell
go mod tidy
```
*Verification*: Run command and verify exit code 0

- [ ] **Build**. Verify the provider builds successfully with your changes:

```shell
make generate
```
*Verification*: Run command and verify exit code 0

- [ ] **Code Formatting**. Ensure code is properly formatted:

```shell
make fmt
```
*Verification*: Run command and verify no files are modified (clean git status)

- [ ] **Code Linting**. Run all linting checks to ensure code quality and consistency:

```shell
make lint
```
*Verification*: Command must exit with code 0

This includes:
- Go vet checks (run separately with `make vet`)
- golangci-lint for Go code quality
- Terraform provider-specific linting (tfproviderlint)
- Terraform linting (tflint)
- Terraform format checking
- Import organization checks

## Testing

### Unit Tests

- [ ] **Unit Tests**. Where a code function performs work internally to a module, but has an external scope (i.e., a function with an initial capital letter `func MyFunction`), unit tests should ideally be created. Not all functions require a unit test, if in doubt please ask:

```shell
make test
```
*Verification*: Run command and verify exit code 0

### Acceptance Tests

- [ ] **Acceptance Tests**. Where a new resource or data source is being created, or existing resources or data sources are being updated, acceptance tests will need to be created or modified according to the [acceptance test strategy](/contributing/acceptance-test-strategy.md)
  - *Verification*:
    - New resources have corresponding `*_test.go` files in same directory
    - Test files include basic CRUD test functions and error condition testing
    - Tests follow naming convention `TestAcc<Resource>_*`
    - Test configurations use proper Docker container setup

Example: To run specific tests based on a regex filter (preferred):
```shell
TF_ACC=1 go test -v -timeout 240s -run ^TestAccOAuthClient $(go list ./internal/...)
```

To run the full suite of tests locally (expect 1hr+):
```shell
make testacc
```

- [ ] **Test Environment**. Ensure you have access to a local PingFederate Docker container for acceptance testing. Use the provided Docker compose setup for consistent testing environments.

- [ ] **Terraform Format in Tests**. Ensure embedded Terraform code in tests is properly formatted:

```shell
make terrafmtlint
```
*Verification*: Run command and verify exit code 0

## Documentation

### Code Documentation

- [ ] **Schema Descriptions**. Each resource, data source and each schema field should have a `description` tag to describe its purpose. These descriptions are used in the autogenerated documentation
  - *Verification*: 
    - All schema attributes have non-empty `Description` fields
    - Resource and data source structs have `Description` fields
    - Descriptions are clear and follow consistent formatting

- [ ] **Custom Errors**. If required, implement appropriate custom error or warning messages for better user experience when API errors or validation errors occur. Include instruction on how the reader can address the root of the warning or error. Most API level errors do not need custom error handling.
  - *Verification*: Custom error functions include actionable guidance for users

- [ ] **Retry Logic**. Implement appropriate retry conditions for eventual consistency following the patterns in the [Provider Design](provider-design.md) guide. Note that typical network level and HTTP code retries should be handled by the PingFederate Go Client instead of the Terraform code (for example, if a `502` error code is returned).
  - *Verification*: Uses appropriate retry functions and error handling patterns

### Examples

- [ ] **Terraform HCL Examples**. New/modified resources and data sources should have appropriate Terraform HCL examples created/altered and stored in the `examples` directory. These are used in the documentation autogeneration routine. Further information and examples of how these are created can be found in the [examples README.md](../examples/README.md). Note, to include multiple examples, the relevant resource template will need to be modifed.
  - *Verification*:
    - Examples exist in `examples/resources/<resource_name>/` or `examples/data-sources/<data_source_name>/`
    - Example files have `.tf` extension and valid Terraform syntax
    - Examples demonstrate both minimal and comprehensive usage where applicable

- [ ] **Data carrying resource examples**. Ensure that any resource that carries data in a PingFederate environment has the `lifecycle.prevent_destroy` meta argument added to prompt users to consider destroy prevention to avoid data loss. For example:
  - *Verification*: Resources that manage user data include `lifecycle.prevent_destroy = false` with explanatory comment

### Documentation Generation

- [ ] **Generate Documentation**. The Terraform documentation is autogenerated. After code (including `description` fields) and examples are finished, generate the documentation:

```shell
make generate
```
*Verification*: 
- Run command and verify exit code 0
- Check that new/modified documentation files are created in `docs/` directory
- Verify documentation follows expected format and structure

## Security & Compliance

- [ ] **Security Scan**. Ensure your code passes security scanning (this will be automatically checked in CI, but you can run locally if gosec is installed)
  - *Verification*: No obvious security issues like hardcoded secrets, SQL injection vectors, or unsafe operations

- [ ] **Sensitive Data**. Ensure no sensitive data (API keys, tokens, etc.) are committed to the repository
  - *Verification*: 
    - No API keys, passwords, or tokens in code or test files
    - Sensitive test data uses environment variables
    - No `.env` files or similar containing credentials

- [ ] **Input Validation**. Implement appropriate input validation for all user-provided data
  - *Verification*: Schema includes appropriate validation rules (e.g., `ValidateFunc`, `ValidateDiagFunc`)

## Changelog

- [ ] **Changelog Entry**. Add a new line to the `CHANGELOG.md` file describing your changes. Follow the existing format and place the entry in the appropriate section (e.g., `## [Unreleased]`). Include:
  - Brief description of the change
  - Type of change (new resource, enhancement, bug fix, etc.)
  - Any breaking changes or migration notes
  - *Verification*:
    - New line added to `CHANGELOG.md` in the unreleased section
    - Entry follows the established format and style
    - Change type is clearly indicated

## Final Checks

- [ ] **All Make Targets**. Run the comprehensive development check (excluding time-intensive tests):

```shell
make devchecknotest
```
*Verification*: Run command and verify exit code 0

- [ ] **CI Compatibility**. Verify your changes will pass automated CI checks by ensuring all the above steps pass locally
  - *Verification*: All previous verification steps completed successfully

- [ ] **Breaking Changes**. If your PR introduces breaking changes, ensure they are:
  - Clearly documented in the PR description
  - Included in the changelog entry
  - Follow the project's versioning strategy
  - *Verification*: 
    - Breaking changes are documented in PR description
    - Changelog entry uses `breaking-change` type if applicable
    - Backward compatibility tests exist for deprecated functionality

## Additional Notes

- The maintainers may run additional tests in different PingFederate configurations
- Large PRs may take longer to review - consider breaking them into smaller, focused changes
- If you're unsure about any step, please ask questions in your PR or create an issue for discussion

---

## Documentation-Only Changes

If you are making documentation-only changes (templates, guides, or examples), you can use this simplified checklist:

- [ ] **Template Changes**. If modifying documentation templates in the `templates/` directory, ensure they follow the existing patterns

- [ ] **Guide Updates**. New or updated guides should be clear, well-structured, and include practical examples

- [ ] **Example Updates**. Ensure any Terraform examples are syntactically correct and follow best practices

- [ ] **Generate Documentation**. After making template or example changes, regenerate the documentation:

```shell
make generate
```

Documentation changes are generally merged quicker than code changes as there is less to review.
