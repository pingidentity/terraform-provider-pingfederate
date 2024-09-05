package sessionapplicationsessionpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sessionApplicationPolicyResource{}
	_ resource.ResourceWithConfigure   = &sessionApplicationPolicyResource{}
	_ resource.ResourceWithImportState = &sessionApplicationPolicyResource{}
)

// SessionApplicationPolicyResource is a helper function to simplify the provider implementation.
func SessionApplicationPolicyResource() resource.Resource {
	return &sessionApplicationPolicyResource{}
}

// sessionApplicationPolicyResource is the resource implementation.
type sessionApplicationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the resource.
func (r *sessionApplicationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the settings for an application session policy.",
		Attributes: map[string]schema.Attribute{
			"idle_timeout_mins": schema.Int64Attribute{
				Description: "The idle timeout period, in minutes. If set to `-1`, the idle timeout will be set to the maximum timeout. The default is `60`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(60),
			},
			"max_timeout_mins": schema.Int64Attribute{
				Description: "The maximum timeout period, in minutes. If set to `-1`, sessions do not expire. The default is `480`.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(480),
			},
		},
	}

	id.ToSchemaDeprecated(&schema, true)
	resp.Schema = schema
}

func addOptionalSessionApplicationPolicyFields(ctx context.Context, addRequest *client.ApplicationSessionPolicy, plan sessionApplicationPolicyModel) error {
	if internaltypes.IsDefined(plan.IdleTimeoutMins) {
		addRequest.IdleTimeoutMins = plan.IdleTimeoutMins.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxTimeoutMins) {
		addRequest.MaxTimeoutMins = plan.MaxTimeoutMins.ValueInt64Pointer()
	}
	return nil

}

func (r *sessionApplicationPolicyModel) buildDefaultClientStruct() *client.ApplicationSessionPolicy {
	return &client.ApplicationSessionPolicy{
		IdleTimeoutMins: utils.Pointer(int64(60)),
		MaxTimeoutMins:  utils.Pointer(int64(480)),
	}
}

// Metadata returns the resource type name.
func (r *sessionApplicationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_application_policy"
}

func (r *sessionApplicationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func (r *sessionApplicationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sessionApplicationPolicyModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSessionApplicationPolicy := client.NewApplicationSessionPolicy()
	err := addOptionalSessionApplicationPolicyFields(ctx, createSessionApplicationPolicy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Session Application Policy", err.Error())
		return
	}

	apiCreateSessionApplicationPolicy := r.apiClient.SessionAPI.UpdateApplicationPolicy(config.AuthContext(ctx, r.providerConfig))
	apiCreateSessionApplicationPolicy = apiCreateSessionApplicationPolicy.Body(*createSessionApplicationPolicy)
	sessionApplicationPolicyResponse, httpResp, err := r.apiClient.SessionAPI.UpdateApplicationPolicyExecute(apiCreateSessionApplicationPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Session Application Policy", err, httpResp)
		return
	}

	// Read the response into the state
	var state sessionApplicationPolicyModel
	readSessionApplicationPolicyResponse(ctx, sessionApplicationPolicyResponse, &state, nil)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *sessionApplicationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sessionApplicationPolicyModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionApplicationPolicy, httpResp, err := r.apiClient.SessionAPI.GetApplicationPolicy(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Session Application Policy", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Session Application Policy", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionApplicationPolicyResponse(ctx, apiReadSessionApplicationPolicy, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sessionApplicationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sessionApplicationPolicyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateSessionApplicationPolicy := r.apiClient.SessionAPI.UpdateApplicationPolicy(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewApplicationSessionPolicy()
	err := addOptionalSessionApplicationPolicyFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Session Application Policy", err.Error())
		return
	}

	updateSessionApplicationPolicy = updateSessionApplicationPolicy.Body(*createUpdateRequest)
	updateSessionApplicationPolicyResponse, httpResp, err := r.apiClient.SessionAPI.UpdateApplicationPolicyExecute(updateSessionApplicationPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Session Application Policy", err, httpResp)
		return
	}

	// Get the current state to see how any attributes are changing
	var state sessionApplicationPolicyModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionApplicationPolicyResponse(ctx, updateSessionApplicationPolicyResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *sessionApplicationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this method will reset the resource to the default configuration in PingFederate.
	var data sessionApplicationPolicyModel
	defaultClientStruct := data.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.SessionAPI.UpdateApplicationPolicy(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*defaultClientStruct)
	_, httpResp, err := r.apiClient.SessionAPI.UpdateApplicationPolicyExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the Session Application Policy", err, httpResp)
	}
}

func (r *sessionApplicationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
