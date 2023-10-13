package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &virtualHostNamesResource{}
	_ resource.ResourceWithConfigure   = &virtualHostNamesResource{}
	_ resource.ResourceWithImportState = &virtualHostNamesResource{}
)

// VirtualHostNamesResource is a helper function to simplify the provider implementation.
func VirtualHostNamesResource() resource.Resource {
	return &virtualHostNamesResource{}
}

// virtualHostNamesResource is the resource implementation.
type virtualHostNamesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type virtualHostNamesResourceModel struct {
	Id               types.String `tfsdk:"id"`
	VirtualHostNames types.Set    `tfsdk:"virtual_host_names"`
}

// GetSchema defines the schema for the resource.
func (r *virtualHostNamesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a VirtualHostNames.",
		Attributes: map[string]schema.Attribute{
			"virtual_host_names": schema.SetAttribute{
				Description: "List of virtual host names.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown()},
			},
		},
	}

	id.Schema(&schema)
	resp.Schema = schema
}

func addOptionalVirtualHostNamesFields(ctx context.Context, addRequest *client.VirtualHostNameSettings, plan virtualHostNamesResourceModel) error {
	if internaltypes.IsDefined(plan.VirtualHostNames) {
		var slice []string
		plan.VirtualHostNames.ElementsAs(ctx, &slice, false)
		addRequest.VirtualHostNames = slice
	}
	return nil

}

// Metadata returns the resource type name.
func (r *virtualHostNamesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_host_names"
}

func (r *virtualHostNamesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readVirtualHostNamesResponse(ctx context.Context, r *client.VirtualHostNameSettings, state *virtualHostNamesResourceModel) {
	state.Id = id.GenerateUUIDToState(state.Id)
	state.VirtualHostNames = internaltypes.GetStringSet(r.VirtualHostNames)
}

func (r *virtualHostNamesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualHostNamesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createVirtualHostNames := client.NewVirtualHostNameSettings()
	err := addOptionalVirtualHostNamesFields(ctx, createVirtualHostNames, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Host Names", err.Error())
		return
	}
	_, requestErr := createVirtualHostNames.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Virtual Host Names: %s", requestErr.Error())
	}

	apiCreateVirtualHostNames := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettings(ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateVirtualHostNames = apiCreateVirtualHostNames.Body(*createVirtualHostNames)
	virtualHostNamesResponse, httpResp, err := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettingsExecute(apiCreateVirtualHostNames)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Host Names", err, httpResp)
		return
	}
	_, responseErr := virtualHostNamesResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Virtual Host Names: %s", requestErr.Error())
	}

	// Read the response into the state
	var state virtualHostNamesResourceModel

	readVirtualHostNamesResponse(ctx, virtualHostNamesResponse, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *virtualHostNamesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualHostNamesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadVirtualHostNames, httpResp, err := r.apiClient.VirtualHostNamesAPI.GetVirtualHostNamesSettings(ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting a Virtual Host Names", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Virtual Host Names", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadVirtualHostNames.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Virtual Host Names: %s", responseErr.Error())
	}

	// Read the response into the state
	readVirtualHostNamesResponse(ctx, apiReadVirtualHostNames, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *virtualHostNamesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan virtualHostNamesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state virtualHostNamesResourceModel
	req.State.Get(ctx, &state)
	updateVirtualHostNames := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettings(ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewVirtualHostNameSettings()
	err := addOptionalVirtualHostNamesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Virtual Host Names", err.Error())
		return
	}
	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of Virtual Host Names: %s", requestErr.Error())
	}
	updateVirtualHostNames = updateVirtualHostNames.Body(*createUpdateRequest)
	updateVirtualHostNamesResponse, httpResp, err := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettingsExecute(updateVirtualHostNames)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Virtual Host Names", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateVirtualHostNamesResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of Virtual Host Names: %s", responseErr.Error())
	}
	// Read the response
	readVirtualHostNamesResponse(ctx, updateVirtualHostNamesResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *virtualHostNamesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *virtualHostNamesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
