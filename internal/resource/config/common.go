package config

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Get schema elements common to all resources
func AddCommonSchema(s *schema.Schema, idRequired bool) {
	// If ID is required (for instantiable config objects) then set it as Required and
	// require replace when changing. Otherwise, mark it as Computed.
	if idRequired {
		s.Attributes["id"] = schema.StringAttribute{
			Description: "The persistent, unique ID for the resource. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		}
	} else {
		s.Attributes["id"] = schema.StringAttribute{
			Description: "Placeholder name of this object required by Terraform.",
			Computed:    true,
			Required:    false,
			Optional:    false,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
	}
}

func SetAllAttributesToOptionalAndComputed(s *schema.Schema, exemptAttributes []string) {
	for key, attribute := range s.Attributes {
		// If more attribute types are used by this provider, this method will need to be updated
		if !internaltypes.StringSliceContains(exemptAttributes, key) {
			stringAttr, ok := attribute.(schema.StringAttribute)
			anyOk := ok
			if ok && (!stringAttr.Computed || !stringAttr.Optional) {
				stringAttr.Required = false
				stringAttr.Optional = true
				stringAttr.Computed = true
				stringAttr.PlanModifiers = append(stringAttr.PlanModifiers, stringplanmodifier.UseStateForUnknown())
				s.Attributes[key] = stringAttr
				continue
			}
			setAttr, ok := attribute.(schema.SetAttribute)
			anyOk = ok || anyOk
			if ok && (!setAttr.Computed || !setAttr.Optional) {
				setAttr.Required = false
				setAttr.Optional = true
				setAttr.Computed = true
				setAttr.PlanModifiers = append(setAttr.PlanModifiers, setplanmodifier.UseStateForUnknown())
				s.Attributes[key] = setAttr
				continue
			}
			boolAttr, ok := attribute.(schema.BoolAttribute)
			anyOk = ok || anyOk
			if ok && (!boolAttr.Computed || !boolAttr.Optional) {
				boolAttr.Required = false
				boolAttr.Optional = true
				boolAttr.Computed = true
				boolAttr.PlanModifiers = append(boolAttr.PlanModifiers, boolplanmodifier.UseStateForUnknown())
				s.Attributes[key] = boolAttr
				continue
			}
			intAttr, ok := attribute.(schema.Int64Attribute)
			anyOk = ok || anyOk
			if ok && (!intAttr.Computed || !intAttr.Optional) {
				intAttr.Required = false
				intAttr.Optional = true
				intAttr.Computed = true
				intAttr.PlanModifiers = append(intAttr.PlanModifiers, int64planmodifier.UseStateForUnknown())
				s.Attributes[key] = intAttr
				continue
			}
			floatAttr, ok := attribute.(schema.Float64Attribute)
			anyOk = ok || anyOk
			if ok && (!floatAttr.Computed || !floatAttr.Optional) {
				floatAttr.Required = false
				floatAttr.Optional = true
				floatAttr.Computed = true
				floatAttr.PlanModifiers = append(floatAttr.PlanModifiers, float64planmodifier.UseStateForUnknown())
				s.Attributes[key] = floatAttr
				continue
			}
			if !anyOk {
				//lintignore:R009
				panic("No valid schema attribute type found when setting attributes to computed: " + key)
			}
		}
	}
}
