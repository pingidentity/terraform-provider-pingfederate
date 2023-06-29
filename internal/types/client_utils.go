package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
)

func ToStateAuthenticationPolicyContract(r *client.AuthenticationPolicyContract) (basetypes.SetValue, basetypes.SetValue) {
	var attrType = map[string]attr.Type{"name": types.StringType}
	clientCoreAttributes := r.GetCoreAttributes()
	var caSlice = []attr.Value{}
	cAobjSlice := types.ObjectType{AttrTypes: attrType}
	for i := 0; i < len(clientCoreAttributes); i++ {
		cAname := clientCoreAttributes[i].GetName()
		cAnameVal := map[string]attr.Value{"name": types.StringValue(cAname)}
		newCaObj, _ := types.ObjectValue(attrType, cAnameVal)
		caSlice = append(caSlice, newCaObj)
	}
	caSliceOfObj, _ := types.SetValue(cAobjSlice, caSlice)

	clientExtAttributes := r.GetExtendedAttributes()
	var eaSlice = []attr.Value{}
	eAobjSlice := types.ObjectType{AttrTypes: attrType}
	for i := 0; i < len(clientExtAttributes); i++ {
		eAname := clientExtAttributes[i].GetName()
		eAnameVal := map[string]attr.Value{"name": types.StringValue(eAname)}
		newEaObj, _ := types.ObjectValue(attrType, eAnameVal)
		eaSlice = append(eaSlice, newEaObj)
	}
	eaSliceOfObj, _ := types.SetValue(eAobjSlice, eaSlice)

	return caSliceOfObj, eaSliceOfObj
}

func AttrValueToBoolPointer(val attr.Value) *bool {
	boolPointerVal := ConvertToPrimitive(val).(bool)
	return &boolPointerVal
}

func ToStateRedirectValidation(r *client.RedirectValidationSettings, diags diag.Diagnostics) (basetypes.ObjectValue, basetypes.ObjectValue) {
	whiteListAttrTypes := map[string]attr.Type{
		"target_resource_sso":      basetypes.BoolType{},
		"target_resource_slo":      basetypes.BoolType{},
		"in_error_resource":        basetypes.BoolType{},
		"idp_discovery":            basetypes.BoolType{},
		"valid_domain":             basetypes.StringType{},
		"valid_path":               basetypes.StringType{},
		"allow_query_and_fragment": basetypes.BoolType{},
		"require_https":            basetypes.BoolType{},
	}

	whiteListAttrs := r.GetRedirectValidationLocalSettings().WhiteList
	var whiteListSliceAttrVal = []attr.Value{}
	whiteListSliceType := types.ObjectType{AttrTypes: whiteListAttrTypes}
	for i := 0; i < len(whiteListAttrs); i++ {
		whiteListAttrValues := map[string]attr.Value{
			"target_resource_sso":      types.BoolPointerValue(whiteListAttrs[i].TargetResourceSSO),
			"target_resource_slo":      types.BoolPointerValue(whiteListAttrs[i].TargetResourceSLO),
			"in_error_resource":        types.BoolPointerValue(whiteListAttrs[i].InErrorResource),
			"idp_discovery":            types.BoolPointerValue(whiteListAttrs[i].IdpDiscovery),
			"valid_domain":             types.StringValue(whiteListAttrs[i].ValidDomain),
			"valid_path":               types.StringPointerValue(whiteListAttrs[i].ValidPath),
			"allow_query_and_fragment": types.BoolPointerValue(whiteListAttrs[i].AllowQueryAndFragment),
			"require_https":            types.BoolPointerValue(whiteListAttrs[i].RequireHttps),
		}
		whiteListObj, _ := types.ObjectValue(whiteListAttrTypes, whiteListAttrValues)
		whiteListSliceAttrVal = append(whiteListSliceAttrVal, whiteListObj)
	}
	whiteListSlice, _ := types.SetValue(whiteListSliceType, whiteListSliceAttrVal)

	redirectValidationLocalSettingsAttrTypes := map[string]attr.Type{
		"enable_target_resource_validation_for_sso":           basetypes.BoolType{},
		"enable_target_resource_validation_for_slo":           basetypes.BoolType{},
		"enable_target_resource_validation_for_idp_discovery": basetypes.BoolType{},
		"enable_in_error_resource_validation":                 basetypes.BoolType{},
		"white_list":                                          basetypes.SetType{ElemType: basetypes.ObjectType{AttrTypes: whiteListAttrTypes}},
	}

	redirectValidationLocalSettings := r.GetRedirectValidationLocalSettings()
	redirectValidationLocalSettingsAttrVals := map[string]attr.Value{
		"enable_target_resource_validation_for_sso":           types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForSSO()),
		"enable_target_resource_validation_for_slo":           types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForSLO()),
		"enable_target_resource_validation_for_idp_discovery": types.BoolValue(redirectValidationLocalSettings.GetEnableTargetResourceValidationForIdpDiscovery()),
		"enable_in_error_resource_validation":                 types.BoolValue(redirectValidationLocalSettings.GetEnableInErrorResourceValidation()),
		"white_list":                                          whiteListSlice,
	}
	redirectValidationLocalSettingsObjVal := MaptoObjValue(redirectValidationLocalSettingsAttrTypes, redirectValidationLocalSettingsAttrVals, diags)

	redirectValidationPartnerSettingsAttrTypes := map[string]attr.Type{
		"enable_wreply_validation_slo": basetypes.BoolType{},
	}

	redirectValidationPartnerSettingsSlo := r.GetRedirectValidationPartnerSettings().EnableWreplyValidationSLO
	redirectValidationPartnerSettingsAttrVals := map[string]attr.Value{
		"enable_wreply_validation_slo": types.BoolPointerValue(redirectValidationPartnerSettingsSlo),
	}

	redirectValidationPartnerSettingsObjVal := MaptoObjValue(redirectValidationPartnerSettingsAttrTypes, redirectValidationPartnerSettingsAttrVals, diags)

	return redirectValidationLocalSettingsObjVal, redirectValidationPartnerSettingsObjVal
}

func ToRequestResourceLink(con context.Context, planObj basetypes.ObjectValue) *client.ResourceLink {
	objValues := planObj.Attributes()
	objId := objValues["id"]
	objLoc := objValues["location"]
	idStrValue := objId.(basetypes.StringValue)
	locStrValue := objLoc.(basetypes.StringValue)
	newLink := client.NewResourceLinkWithDefaults()
	newLink.SetId(idStrValue.ValueString())
	newLink.SetLocation(locStrValue.ValueString())

	return newLink
}

func ToStateResourceLink(r *client.ResourceLink, diags diag.Diagnostics) basetypes.ObjectValue {
	attrTypes := map[string]attr.Type{
		"id":       basetypes.StringType{},
		"location": basetypes.StringType{},
	}

	getId := r.GetId()
	getLocation := r.GetLocation()
	attrValues := map[string]attr.Value{
		"id":       StringTypeOrNil(&getId, false),
		"location": StringTypeOrNil(&getLocation, false),
	}

	linkObjectValue := MaptoObjValue(attrTypes, attrValues, diags)
	return linkObjectValue
}
