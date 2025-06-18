// Copyright © 2025 Ping Identity Corporation

package oauthclientsettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	emptyStringSetDefault, _ = types.SetValue(types.StringType, nil)

	refListElemType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id": types.StringType,
		},
	}

	emptyRefListDefault, _ = types.ListValue(refListElemType, nil)
)

// Some validation has to be done in ModifyPlan because it depends on the PF version
func (r *oauthClientSettingsResource) validatePf121Config(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *oauthClientSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}

	if internaltypes.IsDefined(plan.DynamicClientRegistration) {
		attrs := plan.DynamicClientRegistration.Attributes()
		// If require_offline_access_scope_to_issue_refresh_tokens is not set to "YES", then offline_access_require_consent_prompt has to have the default value of "SERVER_DEFAULT"
		if !attrs["require_offline_access_scope_to_issue_refresh_tokens"].IsUnknown() && attrs["require_offline_access_scope_to_issue_refresh_tokens"].(types.String).ValueString() != "YES" &&
			!attrs["offline_access_require_consent_prompt"].IsUnknown() && attrs["offline_access_require_consent_prompt"].(types.String).ValueString() != "SERVER_DEFAULT" {
			resp.Diagnostics.AddError("'dynamic_client_registration.offline_access_require_consent_prompt' must be set to 'SERVER_DEFAULT' when 'dynamic_client_registration.require_offline_access_scope_to_issue_refresh_tokens' is not set to 'YES'", "")
		}

		// Validate overriding server default for refresh_token_rolling_interval_type.
		resp.Diagnostics.Append(validateOverride("refresh_token_rolling_interval_type", attrs["refresh_token_rolling_interval_type"].(types.String), map[string]attr.Value{
			"refresh_token_rolling_interval_time_unit": attrs["refresh_token_rolling_interval_time_unit"],
		})...)
	}
}

func (r *oauthClientSettingsResource) setVersionDependentDefaults(ctx context.Context, plan *oauthClientSettingsResourceModel, versionAtLeast1210 bool, resp *resource.ModifyPlanResponse) {
	if plan == nil || !internaltypes.IsDefined(plan.DynamicClientRegistration) {
		return
	}

	attrs := plan.DynamicClientRegistration.Attributes()
	if attrs["require_offline_access_scope_to_issue_refresh_tokens"].IsUnknown() {
		if versionAtLeast1210 {
			attrs["require_offline_access_scope_to_issue_refresh_tokens"] = types.StringValue("SERVER_DEFAULT")
		} else {
			attrs["require_offline_access_scope_to_issue_refresh_tokens"] = types.StringNull()
		}
	}
	if attrs["offline_access_require_consent_prompt"].IsUnknown() {
		if versionAtLeast1210 {
			attrs["offline_access_require_consent_prompt"] = types.StringValue("SERVER_DEFAULT")
		} else {
			attrs["offline_access_require_consent_prompt"] = types.StringNull()
		}
	}
	var diags diag.Diagnostics
	plan.DynamicClientRegistration, diags = types.ObjectValue(plan.DynamicClientRegistration.AttributeTypes(ctx), attrs)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (r *oauthClientSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config *oauthClientSettingsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if config == nil {
		return
	}

	if internaltypes.IsDefined(config.DynamicClientRegistration) {
		attrs := config.DynamicClientRegistration.Attributes()

		// Validate overriding server default for client_secret_retention_period_type
		resp.Diagnostics.Append(validateOverride("client_secret_retention_period_type", attrs["client_secret_retention_period_type"].(types.String), map[string]attr.Value{
			"client_secret_retention_period_override": attrs["client_secret_retention_period_override"],
		})...)

		// Validate overriding server default for device_flow_setting_type
		resp.Diagnostics.Append(validateOverride("device_flow_setting_type", attrs["device_flow_setting_type"].(types.String), map[string]attr.Value{
			"pending_authorization_timeout_override": attrs["pending_authorization_timeout_override"],
			"user_authorization_url_override":        attrs["user_authorization_url_override"],
			"device_polling_interval_override":       attrs["device_polling_interval_override"],
		})...)

		// Validate overriding server default for lockout_max_malicious_actions
		resp.Diagnostics.Append(validateOverride("lockout_max_malicious_actions_type", attrs["lockout_max_malicious_actions_type"].(types.String), map[string]attr.Value{
			"lockout_max_malicious_actions": attrs["lockout_max_malicious_actions"],
		})...)

		// Validate overriding server default for persistent_grant_expiration_type
		resp.Diagnostics.Append(validateOverride("persistent_grant_expiration_type", attrs["persistent_grant_expiration_type"].(types.String), map[string]attr.Value{
			"persistent_grant_expiration_time":      attrs["persistent_grant_expiration_time"],
			"persistent_grant_expiration_time_unit": attrs["persistent_grant_expiration_time_unit"],
		})...)

		// Validate overriding server default for persistent_grant_idle_timeout_type
		resp.Diagnostics.Append(validateOverride("persistent_grant_idle_timeout_type", attrs["persistent_grant_idle_timeout_type"].(types.String), map[string]attr.Value{
			"persistent_grant_idle_timeout":           attrs["persistent_grant_idle_timeout"],
			"persistent_grant_idle_timeout_time_unit": attrs["persistent_grant_idle_timeout_time_unit"],
		})...)

		// Validate overriding server default for refresh_token_rolling_grace_period_type
		resp.Diagnostics.Append(validateOverride("refresh_token_rolling_grace_period_type", attrs["refresh_token_rolling_grace_period_type"].(types.String), map[string]attr.Value{
			"refresh_token_rolling_grace_period": attrs["refresh_token_rolling_grace_period"],
		})...)

		// Validate overriding server default for refresh_token_rolling_interval_type.
		resp.Diagnostics.Append(validateOverride("refresh_token_rolling_interval_type", attrs["refresh_token_rolling_interval_type"].(types.String), map[string]attr.Value{
			"refresh_token_rolling_interval": attrs["refresh_token_rolling_interval"],
		})...)

		// retain_client_secret must be true to configure client_secret_retention_period_type and client_secret_retention_period_override
		if !attrs["retain_client_secret"].IsUnknown() {
			if !attrs["retain_client_secret"].(types.Bool).ValueBool() {
				if internaltypes.IsDefined(attrs["client_secret_retention_period_type"]) && attrs["client_secret_retention_period_type"].(types.String).ValueString() != "SERVER_DEFAULT" {
					resp.Diagnostics.AddError("'dynamic_client_registration.client_secret_retention_period_type' cannot be configured unless 'dynamic_client_registration.retain_client_secret' is set to true", "")
				}
				if internaltypes.IsDefined(attrs["client_secret_retention_period_override"]) {
					resp.Diagnostics.AddError("'dynamic_client_registration.client_secret_retention_period_override' cannot be configured unless 'dynamic_client_registration.retain_client_secret' is set to true", "")
				}
			} else {
				// Validate overriding server default for client_secret_retention_period_type
				resp.Diagnostics.Append(validateOverride("client_secret_retention_period_type", attrs["client_secret_retention_period_type"].(types.String), map[string]attr.Value{
					"client_secret_retention_period_override": attrs["client_secret_retention_period_override"],
				})...)
			}
		}
	}
}

func validateOverride(typeAttrName string, typeAttr types.String, overridingAttrs map[string]attr.Value) diag.Diagnostics {
	var respDiags diag.Diagnostics
	if typeAttr.IsUnknown() {
		return respDiags
	}
	if typeAttr.ValueString() == "OVERRIDE_SERVER_DEFAULT" {
		// Each of the overriding attributes must be set
		for attrName, attr := range overridingAttrs {
			if attr.IsNull() {
				respDiags.AddError(fmt.Sprintf("The 'dynamic_client_registration.%s' attribute must be configured when 'dynamic_client_registration.%s' is set to 'OVERRIDE_SERVER_DEFAULT'", attrName, typeAttrName), "")
			}
		}
	} else {
		// Each of the overriding attributes must not be set
		for attrName, attr := range overridingAttrs {
			if internaltypes.IsDefined(attr) {
				respDiags.AddError(fmt.Sprintf("The 'dynamic_client_registration.%s' attribute cannot be configured unless 'dynamic_client_registration.%s' is set to 'OVERRIDE_SERVER_DEFAULT'", attrName, typeAttrName), "")
			}
		}
	}
	return respDiags
}

func (state *oauthClientSettingsResourceModel) readClientResponseCheckPfIgnoredAttrs(response *client.ClientSettings, checkIgnoredAttrs bool) diag.Diagnostics {
	var respDiags diag.Diagnostics
	if checkIgnoredAttrs {
		respDiags.Append(state.reportPfIgnoredAttrs(response)...)
	}
	respDiags.Append(state.readClientResponse(response)...)
	return respDiags
}

func (state *oauthClientSettingsResourceModel) reportPfIgnoredAttrs(response *client.ClientSettings) diag.Diagnostics {
	if !internaltypes.IsDefined(state.DynamicClientRegistration) || !internaltypes.IsDefined(state.DynamicClientRegistration.Attributes()["oidc_policy"]) ||
		response.DynamicClientRegistration == nil || response.DynamicClientRegistration.OidcPolicy == nil {
		return nil
	}
	var respDiags diag.Diagnostics
	// If the response from PF ignored any of the values passed in, report an error
	oidcPolicyAttributes := state.DynamicClientRegistration.Attributes()["oidc_policy"].(types.Object).Attributes()
	stateTokenContentEncryptionAlgorithm := oidcPolicyAttributes["id_token_content_encryption_algorithm"].(types.String)
	stateTokenEncryptionAlgorithm := oidcPolicyAttributes["id_token_encryption_algorithm"].(types.String)
	if !stateTokenContentEncryptionAlgorithm.IsNull() {
		responseTokenContentEncryptionAlgorithm := ""
		if response.DynamicClientRegistration.OidcPolicy.IdTokenContentEncryptionAlgorithm != nil {
			responseTokenContentEncryptionAlgorithm = *response.DynamicClientRegistration.OidcPolicy.IdTokenContentEncryptionAlgorithm
		}
		addPfIgnoredError(&respDiags, "dynamic_client_registration.oidc_policy.id_token_content_encryption_algorithm", responseTokenContentEncryptionAlgorithm, stateTokenContentEncryptionAlgorithm.ValueString())
	}
	if !stateTokenEncryptionAlgorithm.IsNull() {
		responseTokenEncryptionAlgorithm := ""
		if response.DynamicClientRegistration.OidcPolicy.IdTokenEncryptionAlgorithm != nil {
			responseTokenEncryptionAlgorithm = *response.DynamicClientRegistration.OidcPolicy.IdTokenEncryptionAlgorithm
		}
		addPfIgnoredError(&respDiags, "dynamic_client_registration.oidc_policy.id_token_encryption_algorithm", responseTokenEncryptionAlgorithm, stateTokenEncryptionAlgorithm.ValueString())
	}
	return respDiags
}

func addPfIgnoredError(respDiags *diag.Diagnostics, attrName, responseValue, requestValue string) {
	if requestValue != responseValue {
		respDiags.AddError(fmt.Sprintf("PingFederate failed to save the provided value for %s.", attrName),
			fmt.Sprintf("PingFederate returned \"%s\" after request to set %s to \"%s\"", responseValue, attrName, requestValue))
	}
}

func (r *oauthClientSettingsResource) buildDefaultClientStruct() *client.ClientSettings {
	result := &client.ClientSettings{}
	return result
}

func (r *oauthClientSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this resource will reset the PF config back to its default value.
	clientData := r.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.OauthClientSettingsAPI.UpdateOauthClientSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.OauthClientSettingsAPI.UpdateOauthClientSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the oauthClientSettings", err, httpResp)
	}
}
