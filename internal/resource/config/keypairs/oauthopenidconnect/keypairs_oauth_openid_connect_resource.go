package keypairsoauthopenidconnect

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
)

func (r *keypairsOauthOpenidConnectResource) setConditionalDefaults(ctx context.Context, isVersionAtLeast1201 bool, plan *keypairsOauthOpenidConnectResourceModel, resp *resource.ModifyPlanResponse) {
	if isVersionAtLeast1201 {
		// RSA key id lists default to empty lists
		if plan.RsaAlgorithmActiveKeyIds.IsUnknown() {
			plan.RsaAlgorithmActiveKeyIds = emptyRsaKeyListDefault
		}
		if plan.RsaAlgorithmPreviousKeyIds.IsUnknown() {
			plan.RsaAlgorithmPreviousKeyIds = emptyRsaKeyListDefault
		}
	}
	// Nothing else can be set if static_jwks_enabled is set to false
	if plan.StaticJwksEnabled.ValueBool() {
		// If an active cert ref is set, then corresponding publish_x5c_parameter attribute defaults to false
		if plan.P256publishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.P256activeCertRef) {
				plan.P256publishX5cParameter = types.BoolValue(false)
			} else {
				plan.P256publishX5cParameter = types.BoolNull()
			}
		}
		if plan.P256decryptionPublishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.P256decryptionActiveCertRef) {
				plan.P256decryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.P256decryptionPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.P384publishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.P384activeCertRef) {
				plan.P384publishX5cParameter = types.BoolValue(false)
			} else {
				plan.P384publishX5cParameter = types.BoolNull()
			}
		}
		if plan.P384decryptionPublishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.P384decryptionActiveCertRef) {
				plan.P384decryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.P384decryptionPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.P521publishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.P521activeCertRef) {
				plan.P521publishX5cParameter = types.BoolValue(false)
			} else {
				plan.P521publishX5cParameter = types.BoolNull()
			}
		}
		if plan.P521decryptionPublishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.P521decryptionActiveCertRef) {
				plan.P521decryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.P521decryptionPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.RsaPublishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.RsaActiveCertRef) {
				plan.RsaPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.RsaPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.RsaDecryptionPublishX5cParameter.IsUnknown() {
			if internaltypes.IsDefined(plan.RsaDecryptionActiveCertRef) {
				plan.RsaDecryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.RsaDecryptionPublishX5cParameter = types.BoolNull()
			}
		}
	} else {
		// Set the computed and optional fields to null
		if plan.RsaAlgorithmActiveKeyIds.IsUnknown() {
			plan.RsaAlgorithmActiveKeyIds = types.ListNull(types.ObjectType{AttrTypes: rsaKeyIdAttrTypes})
		}
		if plan.RsaAlgorithmPreviousKeyIds.IsUnknown() {
			plan.RsaAlgorithmPreviousKeyIds = types.ListNull(types.ObjectType{AttrTypes: rsaKeyIdAttrTypes})
		}
		if plan.P256publishX5cParameter.IsUnknown() {
			plan.P256publishX5cParameter = types.BoolNull()
		}
		if plan.P256decryptionPublishX5cParameter.IsUnknown() {
			plan.P256decryptionPublishX5cParameter = types.BoolNull()
		}
		if plan.P384publishX5cParameter.IsUnknown() {
			plan.P384publishX5cParameter = types.BoolNull()
		}
		if plan.P384decryptionPublishX5cParameter.IsUnknown() {
			plan.P384decryptionPublishX5cParameter = types.BoolNull()
		}
		if plan.P521publishX5cParameter.IsUnknown() {
			plan.P521publishX5cParameter = types.BoolNull()
		}
		if plan.P521decryptionPublishX5cParameter.IsUnknown() {
			plan.P521decryptionPublishX5cParameter = types.BoolNull()
		}
		if plan.RsaPublishX5cParameter.IsUnknown() {
			plan.RsaPublishX5cParameter = types.BoolNull()
		}
		if plan.RsaDecryptionPublishX5cParameter.IsUnknown() {
			plan.RsaDecryptionPublishX5cParameter = types.BoolNull()
		}
	}

	resp.Plan.Set(ctx, plan)
}

func (r *keypairsOauthOpenidConnectResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config *keypairsOauthOpenidConnectResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if config == nil {
		return
	}

	if config.StaticJwksEnabled.ValueBool() {
		// rsa_active_cert_ref must be set
		if !internaltypes.IsDefined(config.RsaActiveCertRef) {
			resp.Diagnostics.AddError("The rsa_active_cert_ref attribute must be set when static_jwks_enabled is set to true", "")
		}
		validateActiveAndPreviousCertRef("p256", config.P256activeCertRef, config.P256previousCertRef, resp)
		validateActiveAndPreviousCertRef("p256_decryption_", config.P256decryptionActiveCertRef, config.P256decryptionPreviousCertRef, resp)
		validateActiveAndPreviousCertRef("p384", config.P384activeCertRef, config.P384previousCertRef, resp)
		validateActiveAndPreviousCertRef("p384_decryption_", config.P384decryptionActiveCertRef, config.P384decryptionPreviousCertRef, resp)
		validateActiveAndPreviousCertRef("p521", config.P521activeCertRef, config.P521previousCertRef, resp)
		validateActiveAndPreviousCertRef("p521_decryption_", config.P521decryptionActiveCertRef, config.P521decryptionPreviousCertRef, resp)
		validateActiveAndPreviousCertRef("rsa_", config.RsaActiveCertRef, config.RsaPreviousCertRef, resp)
		validateActiveAndPreviousCertRef("rsa_decryption_", config.RsaDecryptionActiveCertRef, config.RsaDecryptionPreviousCertRef, resp)
	} else {
		// Nothing else can be set if static_jwks_enabled is not set to true
		addValidateConfigErrorIfDefined("p256_active_cert_ref", internaltypes.IsDefined(config.P256activeCertRef), resp)
		addValidateConfigErrorIfDefined("p256_active_key_id", internaltypes.IsDefined(config.P256activeKeyId), resp)
		addValidateConfigErrorIfDefined("p256_decryption_active_cert_ref", internaltypes.IsDefined(config.P256decryptionActiveCertRef), resp)
		addValidateConfigErrorIfDefined("p256_decryption_active_key_id", internaltypes.IsDefined(config.P256decryptionActiveKeyId), resp)
		addValidateConfigErrorIfDefined("p256_decryption_previous_cert_ref", internaltypes.IsDefined(config.P256decryptionPreviousCertRef), resp)
		addValidateConfigErrorIfDefined("p256_decryption_previous_key_id", internaltypes.IsDefined(config.P256decryptionPreviousKeyId), resp)
		addValidateConfigErrorIfDefined("p256_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.P256decryptionPublishX5cParameter), resp)
		addValidateConfigErrorIfDefined("p256_previous_cert_ref", internaltypes.IsDefined(config.P256previousCertRef), resp)
		addValidateConfigErrorIfDefined("p256_previous_key_id", internaltypes.IsDefined(config.P256previousKeyId), resp)
		addValidateConfigErrorIfDefined("p256_publish_x5c_parameter", internaltypes.IsDefined(config.P256publishX5cParameter), resp)
		addValidateConfigErrorIfDefined("p384_active_cert_ref", internaltypes.IsDefined(config.P384activeCertRef), resp)
		addValidateConfigErrorIfDefined("p384_active_key_id", internaltypes.IsDefined(config.P384activeKeyId), resp)
		addValidateConfigErrorIfDefined("p384_decryption_active_cert_ref", internaltypes.IsDefined(config.P384decryptionActiveCertRef), resp)
		addValidateConfigErrorIfDefined("p384_decryption_active_key_id", internaltypes.IsDefined(config.P384decryptionActiveKeyId), resp)
		addValidateConfigErrorIfDefined("p384_decryption_previous_cert_ref", internaltypes.IsDefined(config.P384decryptionPreviousCertRef), resp)
		addValidateConfigErrorIfDefined("p384_decryption_previous_key_id", internaltypes.IsDefined(config.P384decryptionPreviousKeyId), resp)
		addValidateConfigErrorIfDefined("p384_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.P384decryptionPublishX5cParameter), resp)
		addValidateConfigErrorIfDefined("p384_previous_cert_ref", internaltypes.IsDefined(config.P384previousCertRef), resp)
		addValidateConfigErrorIfDefined("p384_previous_key_id", internaltypes.IsDefined(config.P384previousKeyId), resp)
		addValidateConfigErrorIfDefined("p384_publish_x5c_parameter", internaltypes.IsDefined(config.P384publishX5cParameter), resp)
		addValidateConfigErrorIfDefined("p521_active_cert_ref", internaltypes.IsDefined(config.P521activeCertRef), resp)
		addValidateConfigErrorIfDefined("p521_active_key_id", internaltypes.IsDefined(config.P521activeKeyId), resp)
		addValidateConfigErrorIfDefined("p521_decryption_active_cert_ref", internaltypes.IsDefined(config.P521decryptionActiveCertRef), resp)
		addValidateConfigErrorIfDefined("p521_decryption_active_key_id", internaltypes.IsDefined(config.P521decryptionActiveKeyId), resp)
		addValidateConfigErrorIfDefined("p521_decryption_previous_cert_ref", internaltypes.IsDefined(config.P521decryptionPreviousCertRef), resp)
		addValidateConfigErrorIfDefined("p521_decryption_previous_key_id", internaltypes.IsDefined(config.P521decryptionPreviousKeyId), resp)
		addValidateConfigErrorIfDefined("p521_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.P521decryptionPublishX5cParameter), resp)
		addValidateConfigErrorIfDefined("p521_previous_cert_ref", internaltypes.IsDefined(config.P521previousCertRef), resp)
		addValidateConfigErrorIfDefined("p521_previous_key_id", internaltypes.IsDefined(config.P521previousKeyId), resp)
		addValidateConfigErrorIfDefined("p521_publish_x5c_parameter", internaltypes.IsDefined(config.P521publishX5cParameter), resp)
		addValidateConfigErrorIfDefined("rsa_active_cert_ref", internaltypes.IsDefined(config.RsaActiveCertRef), resp)
		addValidateConfigErrorIfDefined("rsa_active_key_id", internaltypes.IsDefined(config.RsaActiveKeyId), resp)
		addValidateConfigErrorIfDefined("rsa_algorithm_active_key_ids", internaltypes.IsDefined(config.RsaAlgorithmActiveKeyIds), resp)
		addValidateConfigErrorIfDefined("rsa_algorithm_previous_key_ids", internaltypes.IsDefined(config.RsaAlgorithmPreviousKeyIds), resp)
		addValidateConfigErrorIfDefined("rsa_decryption_active_cert_ref", internaltypes.IsDefined(config.RsaDecryptionActiveCertRef), resp)
		addValidateConfigErrorIfDefined("rsa_decryption_active_key_id", internaltypes.IsDefined(config.RsaDecryptionActiveKeyId), resp)
		addValidateConfigErrorIfDefined("rsa_decryption_previous_cert_ref", internaltypes.IsDefined(config.RsaDecryptionPreviousCertRef), resp)
		addValidateConfigErrorIfDefined("rsa_decryption_previous_key_id", internaltypes.IsDefined(config.RsaDecryptionPreviousKeyId), resp)
		addValidateConfigErrorIfDefined("rsa_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.RsaDecryptionPublishX5cParameter), resp)
		addValidateConfigErrorIfDefined("rsa_previous_cert_ref", internaltypes.IsDefined(config.RsaPreviousCertRef), resp)
		addValidateConfigErrorIfDefined("rsa_previous_key_id", internaltypes.IsDefined(config.RsaPreviousKeyId), resp)
		addValidateConfigErrorIfDefined("rsa_publish_x5c_parameter", internaltypes.IsDefined(config.RsaPublishX5cParameter), resp)
	}
}

func addValidateConfigErrorIfDefined(attrName string, isDefined bool, resp *resource.ValidateConfigResponse) {
	if isDefined {
		resp.Diagnostics.AddError(fmt.Sprintf("The %s attribute cannot be set when static_jwks_enabled is set to false", attrName), "")
	}
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
			resp.Diagnostics.AddError(fmt.Sprintf("The %[1]sactive_cert_ref.id and %[1]sprevious_cert_ref.id attributes must be different.", prefix), fmt.Sprintf("active id: %s, previous id: %s", activeId, previousId))
		}
	} else if internaltypes.IsDefined(previous) {
		// active must be set to set the previous cert ref
		resp.Diagnostics.AddError(fmt.Sprintf("The %[1]sactive_cert_ref attribute must be set when %[1]sprevious_cert_ref is set.", prefix), "")
	}
}

func (m *keypairsOauthOpenidConnectResourceModel) setNullObjectValues() {
	certRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	m.P256activeCertRef = types.ObjectNull(certRefAttrTypes)
	m.P256decryptionActiveCertRef = types.ObjectNull(certRefAttrTypes)
	m.P256decryptionPreviousCertRef = types.ObjectNull(certRefAttrTypes)
	m.P256previousCertRef = types.ObjectNull(certRefAttrTypes)
	m.P384activeCertRef = types.ObjectNull(certRefAttrTypes)
	m.P384decryptionActiveCertRef = types.ObjectNull(certRefAttrTypes)
	m.P384decryptionPreviousCertRef = types.ObjectNull(certRefAttrTypes)
	m.P384previousCertRef = types.ObjectNull(certRefAttrTypes)
	m.P521activeCertRef = types.ObjectNull(certRefAttrTypes)
	m.P521decryptionActiveCertRef = types.ObjectNull(certRefAttrTypes)
	m.P521decryptionPreviousCertRef = types.ObjectNull(certRefAttrTypes)
	m.P521previousCertRef = types.ObjectNull(certRefAttrTypes)
	m.RsaActiveCertRef = types.ObjectNull(certRefAttrTypes)
	// rsa_algorithm_active_key_ids
	rsaAlgorithmKeyIdsAttrTypes := map[string]attr.Type{
		"key_id":       types.StringType,
		"rsa_alg_type": types.StringType,
	}
	m.RsaAlgorithmActiveKeyIds = types.ListNull(types.ObjectType{AttrTypes: rsaAlgorithmKeyIdsAttrTypes})
	m.RsaAlgorithmPreviousKeyIds = types.ListNull(types.ObjectType{AttrTypes: rsaAlgorithmKeyIdsAttrTypes})
	m.RsaDecryptionActiveCertRef = types.ObjectNull(certRefAttrTypes)
	m.RsaDecryptionPreviousCertRef = types.ObjectNull(certRefAttrTypes)
	m.RsaPreviousCertRef = types.ObjectNull(certRefAttrTypes)
}
