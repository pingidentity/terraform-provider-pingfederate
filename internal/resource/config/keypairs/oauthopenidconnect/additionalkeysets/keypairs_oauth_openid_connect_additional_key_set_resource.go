package keypairsoauthopenidconnectadditionalkeysets

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	rsaKeyIdAttrTypes = map[string]attr.Type{
		"key_id":       types.StringType,
		"rsa_alg_type": types.StringType,
	}
	emptyRsaKeySetDefault, _ = types.SetValue(types.ObjectType{AttrTypes: rsaKeyIdAttrTypes}, nil)

	refAttrTypes = map[string]attr.Type{
		"id": types.StringType,
	}
	signingKeysKeyIdsAttrTypes = map[string]attr.Type{
		"key_id":       types.StringType,
		"rsa_alg_type": types.StringType,
	}
	signingKeysKeyIdsElementType = types.ObjectType{AttrTypes: signingKeysKeyIdsAttrTypes}
	signingKeysAttrTypes         = map[string]attr.Type{
		"p256_active_cert_ref":           types.ObjectType{AttrTypes: refAttrTypes},
		"p256_active_key_id":             types.StringType,
		"p256_previous_cert_ref":         types.ObjectType{AttrTypes: refAttrTypes},
		"p256_previous_key_id":           types.StringType,
		"p256_publish_x5c_parameter":     types.BoolType,
		"p384_active_cert_ref":           types.ObjectType{AttrTypes: refAttrTypes},
		"p384_active_key_id":             types.StringType,
		"p384_previous_cert_ref":         types.ObjectType{AttrTypes: refAttrTypes},
		"p384_previous_key_id":           types.StringType,
		"p384_publish_x5c_parameter":     types.BoolType,
		"p521_active_cert_ref":           types.ObjectType{AttrTypes: refAttrTypes},
		"p521_active_key_id":             types.StringType,
		"p521_previous_cert_ref":         types.ObjectType{AttrTypes: refAttrTypes},
		"p521_previous_key_id":           types.StringType,
		"p521_publish_x5c_parameter":     types.BoolType,
		"rsa_active_cert_ref":            types.ObjectType{AttrTypes: refAttrTypes},
		"rsa_active_key_id":              types.StringType,
		"rsa_algorithm_active_key_ids":   types.SetType{ElemType: signingKeysKeyIdsElementType},
		"rsa_algorithm_previous_key_ids": types.SetType{ElemType: signingKeysKeyIdsElementType},
		"rsa_previous_cert_ref":          types.ObjectType{AttrTypes: refAttrTypes},
		"rsa_previous_key_id":            types.StringType,
		"rsa_publish_x5c_parameter":      types.BoolType,
	}
)

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) setConditionalDefaults(ctx context.Context, isVersionAtLeast1201 bool, plan *keypairsOauthOpenidConnectAdditionalKeySetResourceModel, resp *resource.ModifyPlanResponse) {
	signingKeysAttrs := plan.SigningKeys.Attributes()
	if isVersionAtLeast1201 {
		// RSA key id sets default to empty sets
		if signingKeysAttrs["rsa_algorithm_active_key_ids"].IsUnknown() {
			signingKeysAttrs["rsa_algorithm_active_key_ids"] = emptyRsaKeySetDefault
		}
		if signingKeysAttrs["rsa_algorithm_previous_key_ids"].IsUnknown() {
			signingKeysAttrs["rsa_algorithm_previous_key_ids"] = emptyRsaKeySetDefault
		}
	}
	// If an active cert ref is set, then corresponding publish_x5c_parameter attribute defaults to false
	if signingKeysAttrs["p256_publish_x5c_parameter"].IsUnknown() {
		if internaltypes.IsDefined(signingKeysAttrs["p256_active_cert_ref"]) {
			signingKeysAttrs["p256_publish_x5c_parameter"] = types.BoolValue(false)
		} else {
			signingKeysAttrs["p256_publish_x5c_parameter"] = types.BoolNull()
		}
	}
	if signingKeysAttrs["p384_publish_x5c_parameter"].IsUnknown() {
		if internaltypes.IsDefined(signingKeysAttrs["p384_active_cert_ref"]) {
			signingKeysAttrs["p384_publish_x5c_parameter"] = types.BoolValue(false)
		} else {
			signingKeysAttrs["p384_publish_x5c_parameter"] = types.BoolNull()
		}
	}
	if signingKeysAttrs["p521_publish_x5c_parameter"].IsUnknown() {
		if internaltypes.IsDefined(signingKeysAttrs["p521_active_cert_ref"]) {
			signingKeysAttrs["p521_publish_x5c_parameter"] = types.BoolValue(false)
		} else {
			signingKeysAttrs["p521_publish_x5c_parameter"] = types.BoolNull()
		}
	}
	if signingKeysAttrs["rsa_publish_x5c_parameter"].IsUnknown() {
		if internaltypes.IsDefined(signingKeysAttrs["rsa_active_cert_ref"]) {
			signingKeysAttrs["rsa_publish_x5c_parameter"] = types.BoolValue(false)
		} else {
			signingKeysAttrs["rsa_publish_x5c_parameter"] = types.BoolNull()
		}
	}
	var diags diag.Diagnostics
	plan.SigningKeys, diags = types.ObjectValue(signingKeysAttrTypes, signingKeysAttrs)
	resp.Diagnostics.Append(diags...)

	resp.Plan.Set(ctx, plan)
}

func (m *keypairsOauthOpenidConnectAdditionalKeySetResourceModel) validateActivePreviousCertRefs() diag.Diagnostics {
	var respDiags diag.Diagnostics

	if internaltypes.IsDefined(m.SigningKeys) {
		signingKeysAttrs := m.SigningKeys.Attributes()
		validateActiveAndPreviousCertRef("p256", signingKeysAttrs["p256_active_cert_ref"].(types.Object),
			signingKeysAttrs["p256_previous_cert_ref"].(types.Object), &respDiags)
		validateActiveAndPreviousCertRef("p384", signingKeysAttrs["p384_active_cert_ref"].(types.Object),
			signingKeysAttrs["p384_previous_cert_ref"].(types.Object), &respDiags)
		validateActiveAndPreviousCertRef("p521", signingKeysAttrs["p521_active_cert_ref"].(types.Object),
			signingKeysAttrs["p521_previous_cert_ref"].(types.Object), &respDiags)
		validateActiveAndPreviousCertRef("rsa_", signingKeysAttrs["rsa_active_cert_ref"].(types.Object),
			signingKeysAttrs["rsa_previous_cert_ref"].(types.Object), &respDiags)
	}
	return respDiags
}

func validateActiveAndPreviousCertRef(prefix string, active, previous types.Object, respDiags *diag.Diagnostics) {
	if active.IsUnknown() || previous.IsUnknown() {
		return
	}
	if internaltypes.IsDefined(active) && internaltypes.IsDefined(previous) {
		// The active cert ref, if set, must be different than the previous cert ref for each type
		activeId := active.Attributes()["id"].(types.String)
		previousId := previous.Attributes()["id"].(types.String)
		if !activeId.IsUnknown() && activeId.Equal(previousId) {
			respDiags.AddAttributeError(
				path.Root("signing_keys"),
				providererror.InvalidAttributeConfiguration,
				fmt.Sprintf("The signing_keys.%[1]sactive_cert_ref.id and signing_keys.%[1]sprevious_cert_ref.id attributes must be different. "+
					"active id: %[2]s, previous id: %[3]s", prefix, activeId.ValueString(), previousId.ValueString()),
			)
		}
	} else if !internaltypes.IsDefined(active) && internaltypes.IsDefined(previous) {
		// active must be set to set the previous cert ref
		respDiags.AddAttributeError(
			path.Root("signing_keys"),
			providererror.InvalidAttributeConfiguration,
			fmt.Sprintf("The signing_keys.%[1]sactive_cert_ref attribute must be set when signing_keys.%[1]sprevious_cert_ref is set.", prefix),
		)
	}
}
