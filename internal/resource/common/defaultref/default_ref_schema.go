package defaultref

import (
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func ToSchema(description string) schema.Attribute {
	return resourcelink.CompleteSingleNestedAttribute(
		false,
		false,
		true,
		description,
	)
}
