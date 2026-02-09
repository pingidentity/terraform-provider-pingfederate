// Copyright © 2026 Ping Identity Corporation

package attributecontractfulfillment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
)

func ToSchema(required, computed, fullyComputed bool) schema.MapNestedAttribute {
	return ToSchemaWithSuffix(required, computed, fullyComputed, "")
}

func ToSchemaWithSuffix(required, computed, fullyComputed bool, descriptionSuffix string) schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Description: "Defines how an attribute in an attribute contract should be populated." + descriptionSuffix,
		Required:    required,
		Optional:    !required,
		Computed:    computed,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"source": sourcetypeidkey.ToSchema(fullyComputed),
				"value": schema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Description: "The value for this attribute.",
				},
			},
		},
		Validators: []validator.Map{
			configvalidators.ValidAttributeContractFulfillment(),
		},
	}
}
