# Provider Design

This document provides an architectural and design overview of the PingFederate Terraform Provider, outlining its structure, design principles, and development patterns.

## Overview

The PingFederate Terraform Provider is a comprehensive infrastructure-as-code solution for managing PingFederate identity and access management resources. Built using HashiCorp's Terraform Plugin Framework, it provides full lifecycle management of PingFederate configuration through a well-structured, maintainable codebase.

## Architecture

### High-Level Architecture

The provider follows a layered architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     Terraform Core                          │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│              Provider Entry Point (main.go)                 │
│  - Version management                                       │
│  - Plugin Framework v6 server setup                        │
│  - Debug mode support                                       │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│           Provider Core (internal/provider)                 │
│  - Terraform Plugin Framework Provider                      │
│  - Configuration schema and validation                      │
│  - Resource and data source registration                    │
│  - Client configuration and authentication                  │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│             Resource Layer (internal/resource)              │
│  ├── config/        - PingFederate configuration resources  │
│  ├── common/        - Shared resource utilities             │
│  ├── api/           - API interaction utilities             │
│  ├── configvalidators/ - Custom validators                  │
│  ├── planmodifiers/ - Plan modification logic               │
│  └── providererror/ - Error handling utilities              │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│            Data Source Layer (internal/datasource)          │
│  - Data source implementations                              │
│  - Common data source utilities                             │
│  - Query and filtering logic                                │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│          Utilities & Types (internal/utils, types)          │
│  - Custom Terraform types                                   │
│  - Conversion utilities                                     │
│  - Common helper functions                                  │
│  - Version information                                      │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│           PingFederate Go Client                            │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                PingFederate Admin APIs                      │
│  - Configuration API                                        │
│  - Administrative APIs                                      │
│  - Authentication endpoints                                 │
└─────────────────────────────────────────────────────────────┘
```

### Core Components

#### 1. Provider Entry Point (`main.go`)
- **Purpose**: Application entry point and provider server initialization
- **Key Features**:
  - Version management for releases
  - Debug mode support for development
  - Terraform Plugin Framework Protocol v6 server setup
  - Direct provider factory usage (no multiplexing)

#### 2. Provider Core (`internal/provider/`)
- **Purpose**: Main provider implementation using Terraform Plugin Framework
- **Key Features**:
  - Provider configuration schema and validation
  - Authentication and client configuration
  - Resource and data source registration
  - PingFederate-specific provider metadata

#### 3. Resource Layer (`internal/resource/`)
Organized by PingFederate functional areas:

- **`config/`**: PingFederate configuration resources organized by API endpoints
  - Administrative accounts, authentication policies, OAuth settings
  - IDP/SP connections, certificates, and key pairs
  - Server settings, session management, and notifications
- **`common/`**: Shared utilities for resource implementations
- **`api/`**: API interaction and HTTP utilities
- **`configvalidators/`**: Custom validation logic
- **`planmodifiers/`**: Custom plan modification behaviors
- **`providererror/`**: Standardized error handling

#### 4. Data Source Layer (`internal/datasource/`)
- **Purpose**: Data source implementations for querying PingFederate configuration
- **Key Features**:
  - Read-only access to PingFederate resources
  - Common data source patterns and utilities
  - Filtering and query capabilities

#### 5. Utilities & Types (`internal/utils/`, `internal/types/`, `internal/version/`)
- **Custom Types**: PingFederate-specific Terraform types
- **Conversion Utilities**: Type conversion and data transformation helpers
- **Version Management**: Provider version information
- **Common Helpers**: Shared utility functions across the provider

### Design Principles

#### 1. **Separation of Concerns**
- **Functional Organization**: PingFederate resources organized by API functional areas
- **Layer Separation**: Distinct layers for provider logic, resource logic, and API interaction
- **Single Responsibility**: Each component has a clearly defined purpose

#### 2. **Plugin Framework First**
- **Modern Implementation**: Built exclusively on Terraform Plugin Framework
- **Type Safety**: Leverages Framework's type-safe attribute system
- **Enhanced Features**: Better validation, plan modification, and diagnostics

#### 3. **Code Generation Strategy**
- **Generated Resources**: Many resources use generated code from API specifications
- **Manual Override**: Complex resources with custom business logic remain hand-coded
- **Consistency**: Generated code ensures uniform patterns and reduces maintenance overhead

#### 4. **API-Driven Organization**
- **Endpoint Alignment**: Resource structure mirrors PingFederate Admin API organization
- **Configuration Focus**: Resources manage PingFederate configuration objects
- **REST Patterns**: Standard CRUD operations aligned with REST API patterns

#### 5. **Testability**
- **Acceptance Tests**: Comprehensive integration testing against real PingFederate instances
- **Docker Integration**: Local PingFederate containers for consistent testing
- **Test Isolation**: Independent test environments and cleanup procedures

#### 6. **Consistency**
- **Naming Conventions**: Standardized resource and data source naming following PingFederate API patterns
- **File Organization**: Consistent directory structure mirroring API organization
- **Code Patterns**: Reusable patterns for common operations and validations

## Resource Organization Reference

The `internal/resource/config/` package organizes PingFederate resources by functional areas, closely mirroring the PingFederate Admin API structure. This organization provides intuitive navigation and clear separation of concerns.

### Configuration Resources Structure

Resources are organized into functional domains that align with PingFederate's administrative areas:

**Authentication & Identity:**
- `authenticationapi/` - Authentication API configuration
- `authenticationpolicies/` - Authentication policies and fragments
- `authenticationpolicycontract/` - Authentication policy contracts
- `authenticationselector/` - Authentication selectors
- `idp/` - Identity Provider configurations (adapters, connections, tokens)
- `sp/` - Service Provider configurations (adapters, connections, mappings)

**Security & Access Control:**
- `oauth/` - OAuth 2.0 and OpenID Connect configuration
- `certificates/` - Certificate management and validation
- `keypairs/` - Key pair management for signing and SSL
- `captchaproviders/` - CAPTCHA provider configuration

**System Configuration:**
- `serversettings/` - Server-wide settings and configuration
- `cluster/` - Clustering and high availability settings
- `session/` - Session management and policies
- `license/` - License and agreement management

**Integration & Connectivity:**
- `datastore/` - Data store connections
- `pingoneconnection/` - PingOne service integration
- `notificationpublishers/` - Notification and event publishing
- `passwordcredentialvalidator/` - Password validation services

### Resource Implementation Patterns

Resources follow consistent patterns whether manually implemented or generated:

**File Naming:**
- Resources: `{resource_name}_resource.go` or `{resource_name}_resource_gen.go`
- Data Sources: `{resource_name}_data_source.go` or `{resource_name}_data_source_gen.go`
- Tests: `{resource_name}_test.go`

**Generated vs Manual:**
- Generated files include `_gen.go` suffix and header comment indicating generation
- Manual files handle complex business logic and custom validation requirements
- Both follow the same interface patterns and testing standards

### Common Resource Utilities

The `internal/resource/common/` package provides shared utilities used across resource implementations:

- **Attribute Patterns**: Common attribute types and validation patterns
- **Configuration Mapping**: Utilities for mapping between Terraform and API types
- **Import Handling**: Standardized import and state management
- **Plugin Configuration**: Common plugin configuration patterns

### Why This Organization?

**API Alignment**: Resource organization directly mirrors PingFederate's Admin API structure, making it intuitive for users familiar with PingFederate administration.

**Logical Grouping**: Related functionality is grouped together, making it easier to find and maintain related resources.

**Scalability**: The structure accommodates new PingFederate features and API endpoints without major reorganization.

**Developer Experience**: Clear organization helps developers quickly locate relevant code and understand the relationship between resources.

## API Utilities Reference

The `internal/resource/api/` package provides utilities for making API calls to PingFederate with consistent error handling, retry logic, and response parsing. These utilities abstract the complexity of HTTP operations and provide standardized patterns for CRUD operations.

**Why use API utilities:**
- Ensures consistent error handling across all PingFederate API interactions
- Provides automatic retry logic for transient failures
- Standardizes HTTP response handling and error messaging
- Handles authentication and session management transparently
- Includes built-in timeout and logging configuration
- Supports custom error handling for domain-specific business logic

### Core Functions

The package provides HTTP operation utilities that serve as the primary interface for executing API calls against PingFederate's Admin API. These functions handle authentication, request/response processing, error parsing, and diagnostic generation uniformly.

### Error Handling

The API utilities include built-in error handlers for common scenarios such as resource not found conditions, validation errors with enhanced messaging, and authentication/authorization failures. Custom error handlers can be created for domain-specific business logic requirements.

### Authentication Integration

The package provides seamless integration with PingFederate's authentication mechanisms including basic authentication support, OAuth token management, and session handling for long-running operations.

### Configuration Support

The API utilities automatically handle different PingFederate deployment configurations with support for custom endpoints and ports, SSL/TLS configuration including certificate validation, and request timeout and retry configuration.

### Integration Patterns

Developers should use the API utilities for all HTTP interactions with PingFederate, following established CRUD operation patterns. The utilities integrate with the provider's configuration validation and provide consistent diagnostic handling across all resource types.

## Validation Utilities Reference

The `internal/resource/configvalidators/` package provides validation functions and validators for PingFederate-specific requirements and common data formats. These utilities ensure consistent validation across all resources and data sources and provide a centralized location for validation logic.

**Why use validation utilities:**
- Ensures consistent validation behavior across all resources
- Provides PingFederate-specific validation for configuration requirements
- Supports Plugin Framework validation patterns
- Centralizes validation logic for maintainability
- Includes comprehensive validation for PingFederate resource identifiers
- Supports cross-field validation for complex configuration relationships

### Custom Validators

The package includes validators specifically designed for PingFederate configuration including resource identifier validation, configuration relationship validation, and plugin configuration validation. These validators are compatible with the Plugin Framework's validation system.

### Configuration Validation

The package provides validation functions for PingFederate-specific configuration patterns including OAuth and OIDC configuration validation, certificate and key pair validation, authentication policy validation, and connection configuration validation.

### Cross-Field Validation

Specialized utilities handle complex validation scenarios where multiple fields must be validated together including mutually exclusive configuration options, dependent field validation, and conditional requirement validation based on other field values.

### Integration Patterns

Developers should use the validation utilities consistently across all resources, apply appropriate validators in resource schemas using the Plugin Framework patterns, and leverage the centralized validation logic rather than creating custom validation code.

## Utility Functions Reference

The `internal/utils/` package provides helper functions for common data transformations, type conversions, and utility operations used throughout the provider. These utilities standardize common operations and reduce code duplication across implementations.

**Why use utility functions:**
- Standardizes common data transformation operations
- Provides type-safe conversions between different data types
- Includes specialized functions for working with PingFederate API responses
- Offers secure random generation for testing and configuration
- Supports schema composition and reusability patterns
- Reduces code duplication across resource implementations

### Type Conversion Utilities

Conversion functions are available for transforming data between Terraform Plugin Framework types and PingFederate Go client types. These utilities provide type-safe conversions for common scenarios including string/numeric conversions and object/map transformations.

### Configuration Utilities

The package includes functions for handling PingFederate configuration patterns including plugin configuration processing, attribute mapping and transformation, and configuration validation helpers.

### String Utilities

String manipulation functions support common PingFederate configuration requirements including secure random string generation for test data and configuration, and string formatting and validation utilities.

### Schema Utilities

Functions for working with Terraform schemas include utilities for merging schema attribute maps and schema composition patterns that facilitate the creation of reusable schema components.

### Integration Patterns

Developers should leverage these utilities for consistent data transformations, use conversion utilities when working with PingFederate API responses, apply string utilities for secure test data creation, and utilize schema utilities for creating reusable schema components.

## Provider Documentation

The PingFederate Terraform Provider uses an automated documentation generation system that creates comprehensive documentation from multiple sources. Understanding the relationship between templates, examples, and generated documentation is essential for maintaining high-quality provider documentation.

### Documentation Architecture

The provider documentation system consists of three main components that work together to generate the final documentation published at `/docs`:

#### Templates (`/templates`)

Template files define the structure and content for documentation pages. These templates use Go template syntax and are organized by type:

- **Structure**: `/templates/{type}/`
- **Types**: `resources/`, `guides/`
- **Content**: Template files provide structure for resource documentation pages

#### Examples (`/examples`)
Example configurations demonstrate real-world usage patterns and are referenced by templates during documentation generation:

- **Structure**: `/examples/{type}/`
- **Types**: `data-sources/`, `resources/`, `doc-examples/`
- **Content**: Complete Terraform configuration examples showing practical usage
- **Organization**: Examples are organized to match resource structure for easy reference

#### Generated Documentation (`/docs`)

The final documentation is generated by running `make generate`, which processes templates and examples to create comprehensive documentation:

- **Generation Process**: Combines template content with example configurations and schema information
- **Output Structure**: `/docs/{type}/`
- **Content**: Complete documentation pages with descriptions, schemas, and examples

### Documentation Generation Process

#### Make Target

```bash
make generate
```

This command processes all templates and examples to generate the complete documentation set in `/docs`. The generation process:

1. **Runs Code Generation**: Executes `go generate` which includes `tfplugindocs` tool
2. **Incorporates Examples**: Includes relevant example configurations from `/examples`
3. **Generates Schema Documentation**: Extracts schema information from resource implementations
4. **Applies Templates**: Uses templates to structure the final documentation
5. **Formats Documentation**: Runs post-processing scripts for consistent formatting

#### Template Processing

Templates use Go template syntax to include dynamic content:

- **Schema Information**: Automatically extracted from resource implementations
- **Example Configurations**: Referenced from corresponding example files
- **Cross-References**: Links to related resources and data sources
- **Metadata**: Provider version information and generation timestamps

### Resource Documentation Alignment

The documentation structure aligns with the provider's resource organization:

#### Resource Organization → Documentation Mapping

- `internal/resource/config/` → `templates/resources/` → `examples/resources/` → `docs/resources/`
- Resource functional areas directly map to documentation structure
- Examples are organized to match resource hierarchy for intuitive navigation

This alignment ensures that documentation organization matches the provider's internal structure, making it easier for developers to locate relevant documentation and examples while maintaining consistency between code organization and user-facing documentation.
