package common

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func CreateMapStringSchemaAttribute() map[string]schema.Attribute {
	return make(map[string]schema.Attribute)
}
