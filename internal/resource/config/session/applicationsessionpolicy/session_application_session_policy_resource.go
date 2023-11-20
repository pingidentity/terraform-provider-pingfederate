package sessionapplicationsessionpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sessionApplicationSessionPolicyResource{}
	_ resource.ResourceWithConfigure   = &sessionApplicationSessionPolicyResource{}
	_ resource.ResourceWithImportState = &sessionApplicationSessionPolicyResource{}
)

// SessionApplicationSessionPolicyResource is a helper function to simplify the provider implementation.
func SessionApplicationSessionPolicyResource() resource.Resource {
	return &sessionApplicationSessionPolicyResource{}
}

// sessionApplicationSessionPolicyResource is the resource implementation.
type sessionApplicationSessionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type sessionApplicationSessionPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	IdleTimeoutMins types.Int64  `tfsdk:"idle_timeout_mins"`
	MaxTimeoutMins  types.Int64  `tfsdk:"max_timeout_mins"`
}

// GetSchema defines the schema for the resource.
func (r *sessionApplicationSessionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the settings for an application session policy.",
		Attributes: map[string]schema.Attribute{
			"idle_timeout_mins": schema.Int64Attribute{
				Description: "The idle timeout period, in minutes. If set to -1, the idle timeout will be set to the maximum timeout. The default is 60.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(60),
			},
			"max_timeout_mins": schema.Int64Attribute{
				Description: "The maximum timeout period, in minutes. If set to -1, sessions do not expire. The default is 480.",
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(480),
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalSessionApplicationSessionPolicyFields(ctx context.Context, addRequest *client.ApplicationSessionPolicy, plan sessionApplicationSessionPolicyResourceModel) error {
	if internaltypes.IsDefined(plan.IdleTimeoutMins) {
		addRequest.IdleTimeoutMins = plan.IdleTimeoutMins.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxTimeoutMins) {
		addRequest.MaxTimeoutMins = plan.MaxTimeoutMins.ValueInt64Pointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *sessionApplicationSessionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session_application_session_policy"
}

func (r *sessionApplicationSessionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readSessionApplicationSessionPolicyResponse(ctx context.Context, r *client.ApplicationSessionPolicy, state *sessionApplicationSessionPolicyResourceModel, existingId *string) {
	state.Id = id.GenerateUUIDToState(existingId)
	state.IdleTimeoutMins = types.Int64Value(r.GetIdleTimeoutMins())
	state.MaxTimeoutMins = types.Int64Value(r.GetMaxTimeoutMins())
}

func (r *sessionApplicationSessionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sessionApplicationSessionPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSessionApplicationSessionPolicy := client.NewApplicationSessionPolicy()
	err := addOptionalSessionApplicationSessionPolicyFields(ctx, createSessionApplicationSessionPolicy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Session Application Session Policy", err.Error())
		return
	}

	apiCreateSessionApplicationSessionPolicy := r.apiClient.SessionAPI.UpdateApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSessionApplicationSessionPolicy = apiCreateSessionApplicationSessionPolicy.Body(*createSessionApplicationSessionPolicy)
	sessionApplicationSessionPolicyResponse, httpResp, err := r.apiClient.SessionAPI.UpdateApplicationPolicyExecute(apiCreateSessionApplicationSessionPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Session Application Session Policy", err, httpResp)
		return
	}

	// Read the response into the state
	var state sessionApplicationSessionPolicyResourceModel
	readSessionApplicationSessionPolicyResponse(ctx, sessionApplicationSessionPolicyResponse, &state, nil)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *sessionApplicationSessionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sessionApplicationSessionPolicyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionApplicationSessionPolicy, httpResp, err := r.apiClient.SessionAPI.GetApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Session Application Session Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Session Application Session Policy", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionApplicationSessionPolicyResponse(ctx, apiReadSessionApplicationSessionPolicy, &state, id)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sessionApplicationSessionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sessionApplicationSessionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateSessionApplicationSessionPolicy := r.apiClient.SessionAPI.UpdateApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewApplicationSessionPolicy()
	err := addOptionalSessionApplicationSessionPolicyFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Session Application Session Policy", err.Error())
		return
	}

	updateSessionApplicationSessionPolicy = updateSessionApplicationSessionPolicy.Body(*createUpdateRequest)
	updateSessionApplicationSessionPolicyResponse, httpResp, err := r.apiClient.SessionAPI.UpdateApplicationPolicyExecute(updateSessionApplicationSessionPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Session Application Session Policy", err, httpResp)
		return
	}

	// Get the current state to see how any attributes are changing
	var state sessionApplicationSessionPolicyResourceModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readSessionApplicationSessionPolicyResponse(ctx, updateSessionApplicationSessionPolicyResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *sessionApplicationSessionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *sessionApplicationSessionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
