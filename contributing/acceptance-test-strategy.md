# Acceptance Testing Strategy

This document outlines the comprehensive testing strategy for the terraform-provider-pingfederate, based on analysis of existing test patterns and Terraform provider best practices. All new resources and data sources should follow these testing patterns to ensure consistency, reliability, and maintainability.

## Overview

The provider uses HashiCorp's terraform-plugin-testing framework with acceptance tests that interact with real PingFederate server instances. Tests are organized by API configuration area and follow consistent naming and structure patterns.

## Test Organization

### File Structure
- **Resource tests**: Located in `internal/acctest/config/<config_area>/*_test.go`
- **Data source tests**: Co-located with resource tests when applicable
- **Test helpers**: Located in the same package structure

### Package Structure
- Tests are placed in `<config_area>_test` packages (e.g., `keypairssslserver_test`)
- Use shared test utilities from `internal/acctest`
- Follow consistent configuration patterns for test resources

## Core Testing Patterns

### 1. Standard Test Functions

Every resource should implement these core test functions:

#### **Removal Drift Tests**
Tests resource behavior when the underlying PingFederate resource is deleted outside Terraform:
```go
func TestAcc<Resource>_RemovalDrift(t *testing.T) {
    // Test resource removal detection
    // Use PreConfig to delete resource outside Terraform
    // Verify ExpectNonEmptyPlan is true
}
```

#### **Minimal/Maximal Configuration Tests**
Tests resource lifecycle with minimal and maximal configurations:
```go
func TestAcc<Resource>_<Scenario>MinimalMaximal(t *testing.T) {
    // Test minimal → maximal → transitions
    // Test all schema attributes
    // Import testing where supported
}
```

### 2. Schema Testing Strategy

#### **Minimal Schema Testing**
- Test with only required fields set
- Validate API defaults are properly handled
- Use `TestCheckNoResourceAttr` for optional fields that should be absent

```go
minimalCheck := resource.ComposeTestCheckFunc(
    resource.TestCheckResourceAttr(resourceFullName, "id", expectedId),
    resource.TestCheckResourceAttr(resourceFullName, "name", name),
    resource.TestCheckNoResourceAttr(resourceFullName, "description"),
    resource.TestCheckNoResourceAttr(resourceFullName, "optional_field"),
)
```

#### **Maximal Schema Testing**
- Test with all fields (required and optional) set
- Validate complex nested objects and arrays
- Test all supported enum values and configurations

```go
fullCheck := resource.ComposeTestCheckFunc(
    resource.TestCheckResourceAttr(resourceFullName, "id", expectedId),
    resource.TestCheckResourceAttr(resourceFullName, "name", name),
    resource.TestCheckResourceAttr(resourceFullName, "description", "Test description"),
    resource.TestCheckResourceAttr(resourceFullName, "complex_field.nested_field", "expected_value"),
    resource.TestCheckResourceAttr(resourceFullName, "array_field.#", "2"),
)
```

#### **Transition Testing**
Test that resources can successfully transition between different configurations:
- Minimal → Maximal → Minimal
- Different enum values
- Adding/removing optional blocks
- Updating immutable fields (should force replacement)

#### **Backward Compatibility Testing**
When deprecating schema fields, implement comprehensive backward compatibility tests:
- **Dual Functionality**: Test that both deprecated and new fields work simultaneously
- **Gradual Migration**: Validate that users can migrate from old to new fields incrementally
- **Deprecation Warnings**: Ensure deprecated fields generate appropriate warnings
- **Legacy Support**: Test that existing configurations continue to work unchanged

```go
// Example: Testing deprecated field alongside new field
func TestAcc<Resource>_BackwardCompatibility_DeprecatedField(t *testing.T) {
    // Test old field only (legacy behavior)
    // Test new field only (current behavior)  
    // Test both fields together (migration period)
    // Test transition from old to new field
}
```

### 3. Pre-Check Functions

The `internal/acctest` package provides a `ConfigurationPreCheck` function to validate test requirements before execution. This function ensures tests only run when appropriate dependencies and environment variables are configured.

#### **Essential Pre-Check**

**Basic Configuration Authentication**
```go
acctest.ConfigurationPreCheck(t)
```
- **Purpose**: Validates core PingFederate server connection and credentials
- **Required Environment Variables**: 
  - `PINGFEDERATE_PROVIDER_HTTPS_HOST`
  - `PINGFEDERATE_PROVIDER_INSECURE_TRUST_ALL_TLS`
  - `PINGFEDERATE_PROVIDER_X_BYPASS_EXTERNAL_VALIDATION_HEADER`
  - `PINGFEDERATE_PROVIDER_PRODUCT_VERSION`
  - Authentication variables (username/password, OAuth, or access token)
- **When to Use**: Every acceptance test

#### **Authentication Methods**

The provider supports three authentication methods for testing:

**Username/Password Authentication**
- `PINGFEDERATE_PROVIDER_USERNAME`
- `PINGFEDERATE_PROVIDER_PASSWORD`

**OAuth Authentication**
- `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_ID`
- `PINGFEDERATE_PROVIDER_OAUTH_CLIENT_SECRET`
- `PINGFEDERATE_PROVIDER_OAUTH_TOKEN_URL`
- `PINGFEDERATE_PROVIDER_OAUTH_SCOPES`

**Access Token Authentication**
- `PINGFEDERATE_PROVIDER_ACCESS_TOKEN`

#### **Optional Configuration Variables**

**API Configuration**
- `PINGFEDERATE_PROVIDER_ADMIN_API_PATH` (default: `/pf-admin-api/v1`)
- `PINGFEDERATE_PROVIDER_CA_CERTIFICATE_PEM_FILES`

#### **Pre-Check Usage Patterns**

**Standard Test Pattern**
```go
PreCheck: func() {
    acctest.ConfigurationPreCheck(t)
},
```

**Version-Specific Test Pattern**
```go
PreCheck: func() {
    acctest.ConfigurationPreCheck(t)
    if !acctest.VersionAtLeast(version.PingFederate1220) {
        t.Skipf("Test requires PingFederate 12.2.0 or later")
    }
},
```

**Environment-Dependent Test Pattern**
```go
PreCheck: func() {
    acctest.ConfigurationPreCheck(t)
    if os.Getenv("REQUIRED_TEST_ENVIRONMENT_VAR") == "" {
        t.Skip("Test requires REQUIRED_TEST_ENVIRONMENT_VAR to be set")
    }
},
```

### 4. Resource Configuration Strategy

#### **Standard Resource Configuration**
Most tests use direct PingFederate server configuration without environment isolation:
```go
PreCheck: func() {
    acctest.ConfigurationPreCheck(t)
},
```

Use direct resource configuration in test HCL without environment wrappers.

#### **Test Isolation**
For tests that modify shared/global configuration:
- Use unique resource identifiers to avoid conflicts
- Clean up resources in CheckDestroy functions
- Consider test execution order and dependencies

### 5. Import Testing

Resources that support import must test import functionality:
```go
{
    ResourceName: resourceFullName,
    ImportStateIdFunc: func() resource.ImportStateIdFunc {
        return func(s *terraform.State) (string, error) {
            rs, ok := s.RootModule().Resources[resourceFullName]
            if !ok {
                return "", fmt.Errorf("Resource Not found: %s", resourceFullName)
            }
            return rs.Primary.ID, nil
        }
    }(),
    ImportState:       true,
    ImportStateVerify: true,
    ImportStateVerifyIgnore: []string{
        "sensitive_field",
        "computed_field",
    },
}
```

**Note**: Many PingFederate resources are generated resources that do not support import. The test files will indicate this with comments like "// This resource does not support import".

### 6. Error Testing

#### **API Error Validation**
Test expected API error conditions:
```go
{
    Config:      testAccConfig_InvalidConfiguration(resourceName),
    ExpectError: regexp.MustCompile("Expected error message pattern"),
}
```

#### **Import Error Testing**
Test invalid import scenarios (when import is supported):
```go
{
    ImportState:   true,
    ImportStateId: "invalid-id-format",
    ExpectError:   regexp.MustCompile("Unexpected Import Identifier"),
}
```

#### **Configuration Validation**
Test Terraform configuration validation:
```go
{
    Config:      testAccConfig_ConflictingAttributes(resourceName),
    ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
}
```

## Data Source Testing

### Standard Data Source Tests

#### **By Name Lookup**
```go
func TestAcc<DataSource>DataSource_ByNameFull(t *testing.T) {
    // Test data source lookup by name
    // Use TestCheckResourceAttrPair to compare with resource
}
```

#### **By ID Lookup**
```go
func TestAcc<DataSource>DataSource_ByIDFull(t *testing.T) {
    // Test data source lookup by ID
    // Use TestCheckResourceAttrPair to compare with resource
}
```

#### **Not Found Scenarios**
```go
func TestAcc<DataSource>DataSource_NotFound(t *testing.T) {
    // Test data source behavior with non-existent resources
    // Should return appropriate error
}
```

**Note**: Many data sources in the PingFederate provider are co-located with their corresponding resources in the same test file and use resource dependencies for testing.

## Advanced Testing Patterns

### 1. Complex Resource Types

For resources with multiple configuration modes or types:
- Test each mode's minimal and maximal schema separately
- Test transitions between compatible configurations
- Test error conditions for incompatible configuration changes

### 2. Version-Specific Functionality

For features available only in specific PingFederate versions:
- Use `acctest.VersionAtLeast()` to conditionally run tests
- Test both supported and unsupported versions
- Validate proper error handling for unsupported versions

### 3. Certificate and Key Management

For resources involving certificates and keys:
- Use environment variables for test certificate data
- Test both generation and import scenarios where applicable
- Validate certificate attributes and computed values

### 4. Immutable Field Testing

For resources with immutable fields:
- Test that changes force resource replacement
- Validate error messages for unsupported updates
- Test workarounds for resources with customer data dependencies

### 5. Schema Deprecation and Backward Compatibility

When deprecating schema fields or changing resource behavior:

#### **Deprecation Strategy**
- **Simultaneous Support**: Both deprecated and new functionality must work during the deprecation period
- **Gradual Migration**: Users should be able to migrate incrementally without breaking changes
- **Clear Warnings**: Deprecated fields should generate helpful deprecation warnings
- **Documentation**: Update both code comments and user documentation

#### **Required Deprecation Tests**
```go
func TestAcc<Resource>_Deprecation_<FieldName>(t *testing.T) {
    // Test deprecated field works (legacy configuration)
    // Test new field works (current configuration)
    // Test both fields work together (migration period)
    // Test migration path from deprecated to new field
    // Validate deprecation warnings are generated
}
```

#### **Backward Compatibility Validation**
- Existing user configurations must continue to work unchanged
- No functional regression during deprecation period
- Smooth migration path from old to new schema
- Proper handling of edge cases during transition

### 6. Collection Type Testing

For ambiguous collection fields:
- Test ordering for `list` types
- Test uniqueness for `set` types
- Validate proper handling of duplicates
- Test empty collections

## Test Configuration Patterns

### Configuration Helpers

Use consistent configuration helper functions:
```go
func testAcc<Resource>Config_Minimal(resourceName, name string) string {
    return fmt.Sprintf(`
resource "pingfederate_<resource>" "%[1]s" {
  name = "%[2]s"
  // Only required fields
}
`, resourceName, name)
}

func testAcc<Resource>Config_Full(resourceName, name string) string {
    return fmt.Sprintf(`
resource "pingfederate_<resource>" "%[1]s" {
  name        = "%[2]s"
  description = "Test description"
  // All fields including optional
}
`, resourceName, name)
}
```

### Test Helper Functions

Implement standard helper functions in test packages:
```go
func <Resource>_CheckDestroy(s *terraform.State) error
func <Resource>_Delete(t *testing.T, resourceID string)
func <Resource>_CheckComputedValues<Scenario>() resource.TestCheckFunc
```

## Quality Standards

### Test Coverage Requirements

All resources and data sources must have:
- [ ] **Removal drift detection**
- [ ] **Minimal schema validation** with API defaults
- [ ] **Maximal schema validation** with all optional fields
- [ ] **Schema transition testing** (minimal ↔ maximal)
- [ ] **Import functionality testing** (if supported)
- [ ] **Error condition testing**
- [ ] **Data source lookup testing** (if applicable)
- [ ] **Backward compatibility testing** (when deprecating fields)

#### **Additional Requirements for Schema Changes**
When modifying existing resources:
- [ ] **Deprecation testing** for any removed or changed fields
- [ ] **Migration path validation** from old to new schema
- [ ] **Dual functionality testing** during deprecation periods
- [ ] **Warning validation** for deprecated field usage

### Test Reliability

- Use `t.Parallel()` for parallel execution.  Where parallel execution is not achievable, additional PingFederate instances must be used.
  - The use of `t.Parallel()` is meant to replicate customer's use of Terraform where `plan` and `apply` are multithreaded activities by default.
- Implement proper cleanup in CheckDestroy functions
- Handle dependencies between test resources appropriately
- Use appropriate timeouts for resource operations
- Validate test stability across multiple runs

### Maintenance Considerations

- Keep test configurations DRY with helper functions
- Use descriptive test and variable names
- Document complex test scenarios
- Regular review and update of test patterns
- Ensure tests remain valid as PingFederate APIs evolve

## Examples

### Complete Resource Test Structure
```go
func TestAcc<Resource>_RemovalDrift(t *testing.T) { /* ... */ }
func TestAcc<Resource>_<Scenario>MinimalMaximal(t *testing.T) {
    // Minimal test step
    // Maximal test step  
    // Transition tests
    // Import test (if supported)
    // Error tests
}
func TestAcc<Resource>_<SpecificScenario>(t *testing.T) { /* ... */ }

// For resources with deprecated fields
func TestAcc<Resource>_Deprecation_<FieldName>(t *testing.T) {
    // Legacy configuration test
    // New configuration test
    // Migration path test
    // Dual functionality test
}
```

### Complete Data Source Test Structure
```go
func TestAcc<DataSource>DataSource_ByNameFull(t *testing.T) { /* ... */ }
func TestAcc<DataSource>DataSource_ByIDFull(t *testing.T) { /* ... */ }
func TestAcc<DataSource>DataSource_NotFound(t *testing.T) { /* ... */ }
```

This comprehensive testing strategy ensures that all provider functionality is thoroughly validated, providing confidence in the reliability and correctness of the terraform-provider-pingfederate.