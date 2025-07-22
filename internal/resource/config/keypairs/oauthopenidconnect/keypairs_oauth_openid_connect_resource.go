// Copyright © 2025 Ping Identity Corporation

package keypairsoauthopenidconnect

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	rsaKeyIdAttrTypes = map[string]attr.Type{
		"key_id":       types.StringType,
		"rsa_alg_type": types.StringType,
	}
	emptyRsaKeyListDefault, _ = types.SetValue(types.ObjectType{AttrTypes: rsaKeyIdAttrTypes}, nil)
)

func (r *keypairsOauthOpenidConnectResource) setConditionalDefaults(ctx context.Context, isVersionAtLeast1201, isVersionAtLeast1230 bool, plan *keypairsOauthOpenidConnectResourceModel, resp *resource.ModifyPlanResponse) {
	if isVersionAtLeast1201 {
		// RSA key id lists default to empty sets
		if plan.RsaAlgorithmActiveKeyIds.IsUnknown() {
			plan.RsaAlgorithmActiveKeyIds = emptyRsaKeyListDefault
		}
		if plan.RsaAlgorithmPreviousKeyIds.IsUnknown() {
			plan.RsaAlgorithmPreviousKeyIds = emptyRsaKeyListDefault
		}
	}
	if plan.DynamicKeyCertificateInformation.IsUnknown() {
		dynamicKeyCertificateInformationAttrTypes := map[string]attr.Type{
			"city":              types.StringType,
			"country":           types.StringType,
			"organization":      types.StringType,
			"organization_unit": types.StringType,
			"state":             types.StringType,
		}
		if isVersionAtLeast1230 {
			var diags diag.Diagnostics
			plan.DynamicKeyCertificateInformation, diags = types.ObjectValue(dynamicKeyCertificateInformationAttrTypes, map[string]attr.Value{
				"city":              types.StringNull(),
				"country":           types.StringNull(),
				"organization":      types.StringNull(),
				"organization_unit": types.StringNull(),
				"state":             types.StringNull(),
			})
			resp.Diagnostics.Append(diags...)
		} else {
			plan.DynamicKeyCertificateInformation = types.ObjectNull(dynamicKeyCertificateInformationAttrTypes)
		}
	}
	if plan.PublishDynamicKeyX5cs.IsUnknown() {
		if isVersionAtLeast1230 {
			plan.PublishDynamicKeyX5cs = types.BoolValue(false)
		} else {
			plan.PublishDynamicKeyX5cs = types.BoolNull()
		}
	}
	// Nothing else can be set if static_jwks_enabled is set to false
	if plan.StaticJwksEnabled.ValueBool() {
		// If an active cert ref is set, then corresponding publish_x5c_parameter attribute defaults to false
		if plan.P256publishX5cParameter.IsUnknown() && !plan.P256activeCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.P256activeCertRef) {
				plan.P256publishX5cParameter = types.BoolValue(false)
			} else {
				plan.P256publishX5cParameter = types.BoolNull()
			}
		}
		if plan.P256decryptionPublishX5cParameter.IsUnknown() && !plan.P256decryptionActiveCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.P256decryptionActiveCertRef) {
				plan.P256decryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.P256decryptionPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.P384publishX5cParameter.IsUnknown() && !plan.P384activeCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.P384activeCertRef) {
				plan.P384publishX5cParameter = types.BoolValue(false)
			} else {
				plan.P384publishX5cParameter = types.BoolNull()
			}
		}
		if plan.P384decryptionPublishX5cParameter.IsUnknown() && !plan.P384decryptionActiveCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.P384decryptionActiveCertRef) {
				plan.P384decryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.P384decryptionPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.P521publishX5cParameter.IsUnknown() && !plan.P521activeCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.P521activeCertRef) {
				plan.P521publishX5cParameter = types.BoolValue(false)
			} else {
				plan.P521publishX5cParameter = types.BoolNull()
			}
		}
		if plan.P521decryptionPublishX5cParameter.IsUnknown() && !plan.P521decryptionActiveCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.P521decryptionActiveCertRef) {
				plan.P521decryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.P521decryptionPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.RsaPublishX5cParameter.IsUnknown() && !plan.RsaActiveCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.RsaActiveCertRef) {
				plan.RsaPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.RsaPublishX5cParameter = types.BoolNull()
			}
		}
		if plan.RsaDecryptionPublishX5cParameter.IsUnknown() && !plan.RsaDecryptionActiveCertRef.IsUnknown() {
			if internaltypes.IsDefined(plan.RsaDecryptionActiveCertRef) {
				plan.RsaDecryptionPublishX5cParameter = types.BoolValue(false)
			} else {
				plan.RsaDecryptionPublishX5cParameter = types.BoolNull()
			}
		}
	} else {
		// Set the computed and optional fields to null
		if plan.RsaAlgorithmActiveKeyIds.IsUnknown() {
			plan.RsaAlgorithmActiveKeyIds = types.SetNull(types.ObjectType{AttrTypes: rsaKeyIdAttrTypes})
		}
		if plan.RsaAlgorithmPreviousKeyIds.IsUnknown() {
			plan.RsaAlgorithmPreviousKeyIds = types.SetNull(types.ObjectType{AttrTypes: rsaKeyIdAttrTypes})
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

	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (config *keypairsOauthOpenidConnectResourceModel) validatePlan() diag.Diagnostics {
	var respDiags diag.Diagnostics

	if config.StaticJwksEnabled.ValueBool() {
		// rsa_active_cert_ref must be set
		if config.RsaActiveCertRef.IsNull() {
			respDiags.AddError("The rsa_active_cert_ref attribute must be set when static_jwks_enabled is set to true", "")
		}
		validateActiveAndPreviousCertRef("p256", config.P256activeCertRef, config.P256previousCertRef, &respDiags)
		validateActiveAndPreviousCertRef("p256_decryption_", config.P256decryptionActiveCertRef, config.P256decryptionPreviousCertRef, &respDiags)
		validateActiveAndPreviousCertRef("p384", config.P384activeCertRef, config.P384previousCertRef, &respDiags)
		validateActiveAndPreviousCertRef("p384_decryption_", config.P384decryptionActiveCertRef, config.P384decryptionPreviousCertRef, &respDiags)
		validateActiveAndPreviousCertRef("p521", config.P521activeCertRef, config.P521previousCertRef, &respDiags)
		validateActiveAndPreviousCertRef("p521_decryption_", config.P521decryptionActiveCertRef, config.P521decryptionPreviousCertRef, &respDiags)
		validateActiveAndPreviousCertRef("rsa_", config.RsaActiveCertRef, config.RsaPreviousCertRef, &respDiags)
		validateActiveAndPreviousCertRef("rsa_decryption_", config.RsaDecryptionActiveCertRef, config.RsaDecryptionPreviousCertRef, &respDiags)
	} else if !config.StaticJwksEnabled.IsUnknown() {
		// Nothing else can be set if static_jwks_enabled is not set to true
		addValidateConfigErrorIfDefined("p256_active_cert_ref", internaltypes.IsDefined(config.P256activeCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p256_active_key_id", internaltypes.IsDefined(config.P256activeKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p256_decryption_active_cert_ref", internaltypes.IsDefined(config.P256decryptionActiveCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p256_decryption_active_key_id", internaltypes.IsDefined(config.P256decryptionActiveKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p256_decryption_previous_cert_ref", internaltypes.IsDefined(config.P256decryptionPreviousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p256_decryption_previous_key_id", internaltypes.IsDefined(config.P256decryptionPreviousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p256_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.P256decryptionPublishX5cParameter), &respDiags)
		addValidateConfigErrorIfDefined("p256_previous_cert_ref", internaltypes.IsDefined(config.P256previousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p256_previous_key_id", internaltypes.IsDefined(config.P256previousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p256_publish_x5c_parameter", internaltypes.IsDefined(config.P256publishX5cParameter), &respDiags)
		addValidateConfigErrorIfDefined("p384_active_cert_ref", internaltypes.IsDefined(config.P384activeCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p384_active_key_id", internaltypes.IsDefined(config.P384activeKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p384_decryption_active_cert_ref", internaltypes.IsDefined(config.P384decryptionActiveCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p384_decryption_active_key_id", internaltypes.IsDefined(config.P384decryptionActiveKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p384_decryption_previous_cert_ref", internaltypes.IsDefined(config.P384decryptionPreviousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p384_decryption_previous_key_id", internaltypes.IsDefined(config.P384decryptionPreviousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p384_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.P384decryptionPublishX5cParameter), &respDiags)
		addValidateConfigErrorIfDefined("p384_previous_cert_ref", internaltypes.IsDefined(config.P384previousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p384_previous_key_id", internaltypes.IsDefined(config.P384previousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p384_publish_x5c_parameter", internaltypes.IsDefined(config.P384publishX5cParameter), &respDiags)
		addValidateConfigErrorIfDefined("p521_active_cert_ref", internaltypes.IsDefined(config.P521activeCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p521_active_key_id", internaltypes.IsDefined(config.P521activeKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p521_decryption_active_cert_ref", internaltypes.IsDefined(config.P521decryptionActiveCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p521_decryption_active_key_id", internaltypes.IsDefined(config.P521decryptionActiveKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p521_decryption_previous_cert_ref", internaltypes.IsDefined(config.P521decryptionPreviousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p521_decryption_previous_key_id", internaltypes.IsDefined(config.P521decryptionPreviousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p521_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.P521decryptionPublishX5cParameter), &respDiags)
		addValidateConfigErrorIfDefined("p521_previous_cert_ref", internaltypes.IsDefined(config.P521previousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("p521_previous_key_id", internaltypes.IsDefined(config.P521previousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("p521_publish_x5c_parameter", internaltypes.IsDefined(config.P521publishX5cParameter), &respDiags)
		addValidateConfigErrorIfDefined("rsa_active_cert_ref", internaltypes.IsDefined(config.RsaActiveCertRef), &respDiags)
		addValidateConfigErrorIfDefined("rsa_active_key_id", internaltypes.IsDefined(config.RsaActiveKeyId), &respDiags)
		addValidateConfigErrorIfDefined("rsa_algorithm_active_key_ids", internaltypes.IsDefined(config.RsaAlgorithmActiveKeyIds) && len(config.RsaAlgorithmActiveKeyIds.Elements()) > 0, &respDiags)
		addValidateConfigErrorIfDefined("rsa_algorithm_previous_key_ids", internaltypes.IsDefined(config.RsaAlgorithmPreviousKeyIds) && len(config.RsaAlgorithmPreviousKeyIds.Elements()) > 0, &respDiags)
		addValidateConfigErrorIfDefined("rsa_decryption_active_cert_ref", internaltypes.IsDefined(config.RsaDecryptionActiveCertRef), &respDiags)
		addValidateConfigErrorIfDefined("rsa_decryption_active_key_id", internaltypes.IsDefined(config.RsaDecryptionActiveKeyId), &respDiags)
		addValidateConfigErrorIfDefined("rsa_decryption_previous_cert_ref", internaltypes.IsDefined(config.RsaDecryptionPreviousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("rsa_decryption_previous_key_id", internaltypes.IsDefined(config.RsaDecryptionPreviousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("rsa_decryption_publish_x5c_parameter", internaltypes.IsDefined(config.RsaDecryptionPublishX5cParameter), &respDiags)
		addValidateConfigErrorIfDefined("rsa_previous_cert_ref", internaltypes.IsDefined(config.RsaPreviousCertRef), &respDiags)
		addValidateConfigErrorIfDefined("rsa_previous_key_id", internaltypes.IsDefined(config.RsaPreviousKeyId), &respDiags)
		addValidateConfigErrorIfDefined("rsa_publish_x5c_parameter", internaltypes.IsDefined(config.RsaPublishX5cParameter), &respDiags)
	}
	return respDiags
}

func addValidateConfigErrorIfDefined(attrName string, isDefined bool, respDiags *diag.Diagnostics) {
	if isDefined {
		respDiags.AddError(fmt.Sprintf("The %s attribute cannot be set when static_jwks_enabled is set to false", attrName), "")
	}
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
			respDiags.AddError(fmt.Sprintf("The %[1]sactive_cert_ref.id and %[1]sprevious_cert_ref.id attributes must be different.", prefix), fmt.Sprintf("active id: %s, previous id: %s", activeId.ValueString(), previousId.ValueString()))
		}
	} else if active.IsNull() && internaltypes.IsDefined(previous) {
		// active must be set to set the previous cert ref
		respDiags.AddError(fmt.Sprintf("The %[1]sactive_cert_ref attribute must be set when %[1]sprevious_cert_ref is set.", prefix), "")
	}
}

func (m *keypairsOauthOpenidConnectResourceModel) setNullObjectValues() {
	// dynamic_key_certificate_information
	dynamicKeyCertificateInformationAttrTypes := map[string]attr.Type{
		"city":              types.StringType,
		"country":           types.StringType,
		"organization":      types.StringType,
		"organization_unit": types.StringType,
		"state":             types.StringType,
	}
	m.DynamicKeyCertificateInformation = types.ObjectNull(dynamicKeyCertificateInformationAttrTypes)
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
	m.RsaAlgorithmActiveKeyIds = types.SetNull(types.ObjectType{AttrTypes: rsaAlgorithmKeyIdsAttrTypes})
	m.RsaAlgorithmPreviousKeyIds = types.SetNull(types.ObjectType{AttrTypes: rsaAlgorithmKeyIdsAttrTypes})
	m.RsaDecryptionActiveCertRef = types.ObjectNull(certRefAttrTypes)
	m.RsaDecryptionPreviousCertRef = types.ObjectNull(certRefAttrTypes)
	m.RsaPreviousCertRef = types.ObjectNull(certRefAttrTypes)
}

func (r *keypairsOauthOpenidConnectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this resource will reset to the PingFederate default value.
	defaultClientStruct := &client.OAuthOidcKeysSettings{
		StaticJwksEnabled: false,
	}
	apiUpdateRequest := r.apiClient.KeyPairsOauthOpenIdConnectAPI.UpdateOAuthOidcKeysSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*defaultClientStruct)
	_, httpResp, err := r.apiClient.KeyPairsOauthOpenIdConnectAPI.UpdateOAuthOidcKeysSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the keypairsOauthOpenidConnect", err, httpResp)
	}
}
