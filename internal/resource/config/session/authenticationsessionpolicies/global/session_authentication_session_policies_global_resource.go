package sessionauthenticationsessionpoliciesglobal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sessionAuthenticationSessionPoliciesGlobalResource{}
	_ resource.ResourceWithConfigure   = &sessionAuthenticationSessionPoliciesGlobalResource{}
	_ resource.ResourceWithImportState = &sessionAuthenticationSessionPoliciesGlobalResource{}
)

// SessionAuthenticationSessionPoliciesGlobalResource is a helper function to simplify the provider implementation.
func SessionAuthenticationSessionPoliciesGlobalResource() resource.Resource {
	return &sessionAuthenticationSessionPoliciesGlobalResource{}
}

// sessionAuthenticationSessionPoliciesGlobalResource is the resource implementation.
type sessionAuthenticationSessionPoliciesGlobalResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the global settings for authentication session policies.",
		Attributes: map[string]schema.Attribute{
			"enable_sessions": schema.BoolAttribute{
				Description: "Determines whether authentication sessions are enabled globally.",
				Required:    true,
			},
			"persistent_sessions": schema.BoolAttribute{
				Description: "Determines whether authentication sessions are persistent by default. Persistent sessions are linked to a persistent cookie and stored in a data store. This field is ignored if `enable_sessions` is `false`. Default values is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"hash_unique_user_key_attribute": schema.BoolAttribute{
				Description: "Determines whether to hash the value of the unique user key attribute. Default value is `false`.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"idle_timeout_mins": schema.Int64Attribute{
				Description: "The idle timeout period, in minutes. If set to `-1`, the idle timeout will be set to the maximum timeout. The default is `60`.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(60),
			},
			"idle_timeout_display_unit": schema.StringAttribute{
				Description: "The display unit for the idle timeout period in the PingFederate administrative console. When the display unit is `HOURS` or `DAYS`, the timeout value in minutes must correspond to a whole number value for the specified unit. Supported values are `MINUTES`, `HOURS`, and `DAYS`. Default value is `MINUTES`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("MINUTES"),
				Validators: []validator.String{
					stringvalidator.OneOf("MINUTES", "HOURS", "DAYS"),
				},
			},
			"max_timeout_mins": schema.Int64Attribute{
				Description: "The maximum timeout period, in minutes. If set to `-1`, sessions do not expire. The default is `480`.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(480),
			},
			"max_timeout_display_unit": schema.StringAttribute{
				Description: "The display unit for the maximum timeout period in the PingFederate administrative console. When the display unit is `HOURS` or `DAYS`, the timeout value in minutes must correspond to a whole number value for the specified unit. Supported values are `MINUTES`, `HOURS`, and `DAYS`. Default value is `MINUTES`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("MINUTES"),
				Validators: []validator.String{
					stringvalidator.OneOf("MINUTES", "HOURS", "DAYS"),
				},
			},
		},
	}

	id.ToSchemaDeprecated(&schema, true)
	resp.Schema = schema
}

func addOptionalSessionAuthenticationSessionPoliciesGlobalFields(ctx context.Context, addRequest *client.GlobalAuthenticationSessionPolicy, plan sessionAuthenticationSessionPoliciesGlobalModel) error {
	if internaltypes.IsDefined(plan.EnableSessions) {
		addRequest.EnableSessions = plan.EnableSessions.ValueBool()
	}
	if internaltypes.IsDefined(plan.PersistentSessions) {
		addRequest.PersistentSessions = plan.PersistentSessions.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.HashUniqueUserKeyAttribute) {
		addRequest.HashUniqueUserKeyAttribute = plan.HashUniqueUserKeyAttribute.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IdleTimeoutMins) {
		addRequest.IdleTimeoutMins = plan.IdleTimeoutMins.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IdleTimeoutDisplayUnit) {
		addRequest.IdleTimeoutDisplayUnit = plan.IdleTimeoutDisplayUnit.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxTimeoutMins) {
		addRequest.MaxTimeoutMins = plan.MaxTimeoutMins.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxTimeoutDisplayUnit) {
		addRequest.MaxTimeoutDisplayUnit = plan.MaxTimeoutDisplayUnit.ValueStringPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_authentication_session_policies_global"
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (m *sessionAuthenticationSessionPoliciesGlobalModel) buildDefaultClientStruct() *client.GlobalAuthenticationSessionPolicy {
	return &client.GlobalAuthenticationSessionPolicy{
		EnableSessions:             false,
		PersistentSessions:         utils.Pointer(false),
		HashUniqueUserKeyAttribute: utils.Pointer(false),
		IdleTimeoutMins:            utils.Pointer(int64(60)),
		IdleTimeoutDisplayUnit:     utils.Pointer("MINUTES"),
		MaxTimeoutMins:             utils.Pointer(int64(480)),
		MaxTimeoutDisplayUnit:      utils.Pointer("MINUTES"),
	}
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sessionAuthenticationSessionPoliciesGlobalModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSessionAuthenticationSessionPoliciesGlobal := client.NewGlobalAuthenticationSessionPolicy(plan.EnableSessions.ValueBool())
	err := addOptionalSessionAuthenticationSessionPoliciesGlobalFields(ctx, createSessionAuthenticationSessionPoliciesGlobal, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the global authentication session policy", err.Error())
		return
	}

	apiCreateSessionAuthenticationSessionPoliciesGlobal := r.apiClient.SessionAPI.UpdateGlobalPolicy(config.AuthContext(ctx, r.providerConfig))
	apiCreateSessionAuthenticationSessionPoliciesGlobal = apiCreateSessionAuthenticationSessionPoliciesGlobal.Body(*createSessionAuthenticationSessionPoliciesGlobal)
	sessionAuthenticationSessionPoliciesGlobalResponse, httpResp, err := r.apiClient.SessionAPI.UpdateGlobalPolicyExecute(apiCreateSessionAuthenticationSessionPoliciesGlobal)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the the global authentication session policy", err, httpResp)
		return
	}

	// Read the response into the state
	var state sessionAuthenticationSessionPoliciesGlobalModel
	readSessionAuthenticationSessionPoliciesGlobalResponse(ctx, sessionAuthenticationSessionPoliciesGlobalResponse, &state, nil)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sessionAuthenticationSessionPoliciesGlobalModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionAuthenticationSessionPoliciesGlobal, httpResp, err := r.apiClient.SessionAPI.GetGlobalPolicy(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "the global authentication session policy", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a the global authentication session policy", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionAuthenticationSessionPoliciesGlobalResponse(ctx, apiReadSessionAuthenticationSessionPoliciesGlobal, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sessionAuthenticationSessionPoliciesGlobalModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateSessionAuthenticationSessionPoliciesGlobal := r.apiClient.SessionAPI.UpdateGlobalPolicy(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewGlobalAuthenticationSessionPolicy(plan.EnableSessions.ValueBool())
	err := addOptionalSessionAuthenticationSessionPoliciesGlobalFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the global authentication session policy", err.Error())
		return
	}

	updateSessionAuthenticationSessionPoliciesGlobal = updateSessionAuthenticationSessionPoliciesGlobal.Body(*createUpdateRequest)
	updateSessionAuthenticationSessionPoliciesGlobalResponse, httpResp, err := r.apiClient.SessionAPI.UpdateGlobalPolicyExecute(updateSessionAuthenticationSessionPoliciesGlobal)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the global authentication session policy", err, httpResp)
		return
	}

	// Read the response
	var state sessionAuthenticationSessionPoliciesGlobalModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionAuthenticationSessionPoliciesGlobalResponse(ctx, updateSessionAuthenticationSessionPoliciesGlobalResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *sessionAuthenticationSessionPoliciesGlobalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	var model sessionAuthenticationSessionPoliciesGlobalModel
	clientData := model.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.SessionAPI.UpdateGlobalPolicy(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.SessionAPI.UpdateGlobalPolicyExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the global authentication session policy", err, httpResp)
	}
}

func (r *sessionAuthenticationSessionPoliciesGlobalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
