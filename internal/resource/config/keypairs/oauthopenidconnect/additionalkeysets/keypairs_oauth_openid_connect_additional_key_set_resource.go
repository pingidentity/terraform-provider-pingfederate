package keypairsoauthopenidconnectadditionalkeysets

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	rsaKeyIdAttrTypes = map[string]attr.Type{
		"key_id":       types.StringType,
		"rsa_alg_type": types.StringType,
	}
	emptyRsaKeyListDefault, _ = types.ListValue(types.ObjectType{AttrTypes: rsaKeyIdAttrTypes}, nil)

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
		"rsa_algorithm_active_key_ids":   types.ListType{ElemType: signingKeysKeyIdsElementType},
		"rsa_algorithm_previous_key_ids": types.ListType{ElemType: signingKeysKeyIdsElementType},
		"rsa_previous_cert_ref":          types.ObjectType{AttrTypes: refAttrTypes},
		"rsa_previous_key_id":            types.StringType,
		"rsa_publish_x5c_parameter":      types.BoolType,
	}
)

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) setConditionalDefaults(ctx context.Context, isVersionAtLeast1201 bool, plan *keypairsOauthOpenidConnectAdditionalKeySetResourceModel, resp *resource.ModifyPlanResponse) {
	signingKeysAttrs := plan.SigningKeys.Attributes()
	if isVersionAtLeast1201 {
		// RSA key id lists default to empty lists
		if signingKeysAttrs["rsa_algorithm_active_key_ids"].IsUnknown() {
			signingKeysAttrs["rsa_algorithm_active_key_ids"] = emptyRsaKeyListDefault
		}
		if signingKeysAttrs["rsa_algorithm_previous_key_ids"].IsUnknown() {
			signingKeysAttrs["rsa_algorithm_previous_key_ids"] = emptyRsaKeyListDefault
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

func (r *keypairsOauthOpenidConnectAdditionalKeySetResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config *keypairsOauthOpenidConnectAdditionalKeySetResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if config == nil {
		return
	}
	signingKeysAttrs := config.SigningKeys.Attributes()

	validateActiveAndPreviousCertRef("p256", signingKeysAttrs["p256_active_cert_ref"].(types.Object),
		signingKeysAttrs["p256_previous_cert_ref"].(types.Object), resp)
	validateActiveAndPreviousCertRef("p384", signingKeysAttrs["p384_active_cert_ref"].(types.Object),
		signingKeysAttrs["p384_previous_cert_ref"].(types.Object), resp)
	validateActiveAndPreviousCertRef("p521", signingKeysAttrs["p521_active_cert_ref"].(types.Object),
		signingKeysAttrs["p521_previous_cert_ref"].(types.Object), resp)
	validateActiveAndPreviousCertRef("rsa_", signingKeysAttrs["rsa_active_cert_ref"].(types.Object),
		signingKeysAttrs["rsa_previous_cert_ref"].(types.Object), resp)

}

func validateActiveAndPreviousCertRef(prefix string, active, previous types.Object, resp *resource.ValidateConfigResponse) {
	if internaltypes.IsDefined(active) {
		// The active cert ref, if set, must be different than the previous cert ref for each type
		activeId := active.Attributes()["id"].(types.String).ValueString()
		previousId := ""
		if internaltypes.IsDefined(previous) {
			previousId = previous.Attributes()["id"].(types.String).ValueString()
		}
		if activeId == previousId {
			resp.Diagnostics.AddError(fmt.Sprintf("The signing_keys.%[1]sactive_cert_ref.id and signing_keys.%[1]sprevious_cert_ref.id attributes must be different.", prefix), fmt.Sprintf("active id: %s, previous id: %s", activeId, previousId))
		}
	} else if internaltypes.IsDefined(previous) {
		// active must be set to set the previous cert ref
		resp.Diagnostics.AddError(fmt.Sprintf("The signing_keys.%[1]sactive_cert_ref attribute must be set when signing_keys.%[1]sprevious_cert_ref is set.", prefix), "")
	}
}
