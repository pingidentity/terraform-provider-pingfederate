package defaultref

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func ToSchema() schema.Attribute {
	return resourcelink.CompleteSingleNestedAttribute(
		false,
		false,
		true,
		"Reference to the default, if one is defined.",
	)
}
