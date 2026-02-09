// Copyright © 2026 Ping Identity Corporation

package config

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

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
				stringAttr.PlanModifiers = append(stringAttr.PlanModifiers, stringplanmodifier.UseNonNullStateForUnknown())
				s.Attributes[key] = stringAttr
				continue
			}
			setAttr, ok := attribute.(schema.SetAttribute)
			anyOk = ok || anyOk
			if ok && (!setAttr.Computed || !setAttr.Optional) {
				setAttr.Required = false
				setAttr.Optional = true
				setAttr.Computed = true
				setAttr.PlanModifiers = append(setAttr.PlanModifiers, setplanmodifier.UseNonNullStateForUnknown())
				s.Attributes[key] = setAttr
				continue
			}
			listAttr, ok := attribute.(schema.ListAttribute)
			anyOk = ok || anyOk
			if ok && (!listAttr.Computed || !listAttr.Optional) {
				listAttr.Required = false
				listAttr.Optional = true
				listAttr.Computed = true
				listAttr.PlanModifiers = append(listAttr.PlanModifiers, listplanmodifier.UseNonNullStateForUnknown())
				s.Attributes[key] = listAttr
				continue
			}
			boolAttr, ok := attribute.(schema.BoolAttribute)
			anyOk = ok || anyOk
			if ok && (!boolAttr.Computed || !boolAttr.Optional) {
				boolAttr.Required = false
				boolAttr.Optional = true
				boolAttr.Computed = true
				boolAttr.PlanModifiers = append(boolAttr.PlanModifiers, boolplanmodifier.UseNonNullStateForUnknown())
				s.Attributes[key] = boolAttr
				continue
			}
			intAttr, ok := attribute.(schema.Int64Attribute)
			anyOk = ok || anyOk
			if ok && (!intAttr.Computed || !intAttr.Optional) {
				intAttr.Required = false
				intAttr.Optional = true
				intAttr.Computed = true
				intAttr.PlanModifiers = append(intAttr.PlanModifiers, int64planmodifier.UseNonNullStateForUnknown())
				s.Attributes[key] = intAttr
				continue
			}
			floatAttr, ok := attribute.(schema.Float64Attribute)
			anyOk = ok || anyOk
			if ok && (!floatAttr.Computed || !floatAttr.Optional) {
				floatAttr.Required = false
				floatAttr.Optional = true
				floatAttr.Computed = true
				floatAttr.PlanModifiers = append(floatAttr.PlanModifiers, float64planmodifier.UseNonNullStateForUnknown())
				s.Attributes[key] = floatAttr
				continue
			}
			if !anyOk {
				return
			}
		}
	}
}
