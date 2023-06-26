package sessionApplicationSessionPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
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
	sessionApplicationSessionPolicyResourceSchema(ctx, req, resp, false)
}

func sessionApplicationSessionPolicyResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
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

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{""})
	}
	config.AddCommonSchema(&schema, false)
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
	resp.TypeName = req.ProviderTypeName + "_session_applicationsessionpolicy"
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for SessionApplicationSessionPolicy", err.Error())
		return
	}
	requestJson, err := createSessionApplicationSessionPolicy.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateSessionApplicationSessionPolicy := r.apiClient.SessionApi.UpdateApplicationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateSessionApplicationSessionPolicy = apiCreateSessionApplicationSessionPolicy.Body(*createSessionApplicationSessionPolicy)
	sessionApplicationSessionPolicyResponse, httpResp, err := r.apiClient.SessionApi.UpdateApplicationPolicyExecute(apiCreateSessionApplicationSessionPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the SessionApplicationSessionPolicy", err, httpResp)
		return
	}
	responseJson, err := sessionApplicationSessionPolicyResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state sessionApplicationSessionPolicyResourceModel

	readSessionApplicationSessionPolicyResponse(ctx, sessionApplicationSessionPolicyResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sessionApplicationSessionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSessionApplicationSessionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSessionApplicationSessionPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state sessionApplicationSessionPolicyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadSessionApplicationSessionPolicy, httpResp, err := apiClient.SessionApi.GetApplicationPolicy(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a SessionApplicationSessionPolicy", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadSessionApplicationSessionPolicy.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSessionApplicationSessionPolicyResponse(ctx, apiReadSessionApplicationSessionPolicy, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sessionApplicationSessionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSessionApplicationSessionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSessionApplicationSessionPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
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
	updateSessionApplicationSessionPolicy := apiClient.SessionApi.UpdateApplicationPolicy(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewApplicationSessionPolicy()
	err := addOptionalSessionApplicationSessionPolicyFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for SessionApplicationSessionPolicy", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateSessionApplicationSessionPolicy = updateSessionApplicationSessionPolicy.Body(*createUpdateRequest)
	updateSessionApplicationSessionPolicyResponse, httpResp, err := apiClient.SessionApi.UpdateApplicationPolicyExecute(updateSessionApplicationSessionPolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating SessionApplicationSessionPolicy", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateSessionApplicationSessionPolicyResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readSessionApplicationSessionPolicyResponse(ctx, updateSessionApplicationSessionPolicyResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *sessionApplicationSessionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *sessionApplicationSessionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
