package providererror

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const (
	// Common provider error summaries
	InvalidProviderConfiguration  = "Invalid provider configuration"
	InvalidAttributeConfiguration = "Invalid attribute configuration"
	InternalProviderError         = "Internal provider error"
)

func WarnConfigurationCannotBeReset(resourceName string, diags *diag.Diagnostics) {
	diags.AddWarning("Configuration cannot be returned to original state",
		fmt.Sprintf("The %s resource has been destroyed but cannot be returned to its original state. The resource has been removed from Terraform state but the configuration remains applied to the environment", resourceName))
}
