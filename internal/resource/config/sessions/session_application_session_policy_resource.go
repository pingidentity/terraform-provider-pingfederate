package sessions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
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
		Description: "Manages a SessionApplicationSessionPolicy.",
		Attributes: map[string]schema.Attribute{
			// Add necessary attributes here
			"idle_timeout_mins": schema.Int64Attribute{
				Description: "The idle timeout period, in minutes. If set to -1, the idle timeout will be set to the maximum timeout. The default is 60.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown()},
			},
			"max_timeout_mins": schema.Int64Attribute{
				Description: "The maximum timeout period, in minutes. If set to -1, sessions do not expire. The default is 480.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown()},
			},
		},
	}

	config.AddCommonSchema(&schema)
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

func readSessionApplicationSessionPolicyResponse(ctx context.Context, r *client.ApplicationSessionPolicy, state *sessionApplicationSessionPolicyResourceModel, expectedValues *sessionApplicationSessionPolicyResourceModel) {
	//TODO placeholder?
	state.Id = types.StringValue("id")
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
	_, requestErr := createSessionApplicationSessionPolicy.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Session Application Session Policy: %s", requestErr.Error())
	}

	apiCreateSessionApplicationSessionPolicy := r.apiClient.SessionAPI.UpdateApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSessionApplicationSessionPolicy = apiCreateSessionApplicationSessionPolicy.Body(*createSessionApplicationSessionPolicy)
	sessionApplicationSessionPolicyResponse, httpResp, err := r.apiClient.SessionAPI.UpdateApplicationPolicyExecute(apiCreateSessionApplicationSessionPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Session Application Session Policy", err, httpResp)
		return
	}
	_, responseErr := sessionApplicationSessionPolicyResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Session Application Session Policy: %s", responseErr.Error())
	}

	// Read the response into the state
	var state sessionApplicationSessionPolicyResourceModel

	readSessionApplicationSessionPolicyResponse(ctx, sessionApplicationSessionPolicyResponse, &state, &plan)
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
	// Log response JSON
	_, responseErr := apiReadSessionApplicationSessionPolicy.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Session Application Session Policy: %s", responseErr.Error())
	}

	// Read the response into the state
	readSessionApplicationSessionPolicyResponse(ctx, apiReadSessionApplicationSessionPolicy, &state, &state)

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

	// Get the current state to see how any attributes are changing
	var state sessionApplicationSessionPolicyResourceModel
	req.State.Get(ctx, &state)
	updateSessionApplicationSessionPolicy := r.apiClient.SessionAPI.UpdateApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewApplicationSessionPolicy()
	err := addOptionalSessionApplicationSessionPolicyFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Session Application Session Policy", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Session Application Session Policy: %s", requestErr.Error())
	}
	updateSessionApplicationSessionPolicy = updateSessionApplicationSessionPolicy.Body(*createUpdateRequest)
	updateSessionApplicationSessionPolicyResponse, httpResp, err := r.apiClient.SessionAPI.UpdateApplicationPolicyExecute(updateSessionApplicationSessionPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Session Application Session Policy", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateSessionApplicationSessionPolicyResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Session Application Session Policy: %s", responseErr.Error())
	}
	// Read the response
	readSessionApplicationSessionPolicyResponse(ctx, updateSessionApplicationSessionPolicyResponse, &state, &plan)

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
