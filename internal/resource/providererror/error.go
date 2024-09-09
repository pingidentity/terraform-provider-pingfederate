package providererror

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const (
	InvalidProviderConfiguration = "Invalid provider configuration"
)

func WarnConfigurationCannotBeReset(resourceName string, diags *diag.Diagnostics) {
	diags.AddWarning("Configuration cannot be returned to original state",
		fmt.Sprintf("The %s resource has been destroyed but cannot be returned to its original state. The resource has been removed from Terraform state but the configuration remains applied to the environment", resourceName))
}
