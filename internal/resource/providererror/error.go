// Copyright © 2025 Ping Identity Corporation

package providererror

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const (
	// Common provider error summaries
	InvalidProviderConfiguration    = "Invalid provider configuration"
	InvalidAttributeConfiguration   = "Invalid attribute configuration"
	InvalidResourceConfiguration    = "Invalid resource configuration"
	InvalidProductVersionAttribute  = "Invalid product_version for attribute"
	InvalidProductVersionResource   = "Invalid product_version for resource"
	InvalidResourceIdForImport      = "Invalid resource id for import"
	InternalProviderError           = "Internal provider error"
	PingFederateValidationError     = "PingFederate validation error"
	PingFederateAPIError            = "PingFederate API error"
	ConfigurationWarning            = "Plugin configuration warning"
	ConfigurationCannotBeResetError = "Configuration cannot be returned to original state"
	ConflictingValueReturnedError   = "PingFederate returned conflicting value"
)

func WarnConfigurationCannotBeReset(resourceName string, diags *diag.Diagnostics) {
	diags.AddWarning(ConfigurationCannotBeResetError,
		fmt.Sprintf("The %s resource has been destroyed but cannot be returned to its original state. The resource has been removed from Terraform state but the configuration remains applied to the environment", resourceName))
}
