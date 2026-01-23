// Copyright © 2025 Ping Identity Corporation

package virtualhostnames

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &virtualHostNamesResource{}
	_ resource.ResourceWithConfigure   = &virtualHostNamesResource{}
	_ resource.ResourceWithImportState = &virtualHostNamesResource{}

	virtualHostNamesDefault, _ = types.SetValue(types.StringType, nil)
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

// GetSchema defines the schema for the resource.
func (r *virtualHostNamesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages virtual host names settings.",
		Attributes: map[string]schema.Attribute{
			"virtual_host_names": schema.SetAttribute{
				Description: "List of virtual host names.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Default:     setdefault.StaticValue(virtualHostNamesDefault),
			},
		},
	}
	resp.Schema = schema
}

func addOptionalVirtualHostNamesFields(ctx context.Context, addRequest *client.VirtualHostNameSettings, plan virtualHostNamesModel) error {
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

func (m *virtualHostNamesModel) buildDefaultClientStruct() *client.VirtualHostNameSettings {
	return &client.VirtualHostNameSettings{
		VirtualHostNames: []string{},
	}
}

func (r *virtualHostNamesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualHostNamesModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createVirtualHostNames := client.NewVirtualHostNameSettings()
	err := addOptionalVirtualHostNamesFields(ctx, createVirtualHostNames, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Virtual Host Names settings: "+err.Error())
		return
	}

	apiCreateVirtualHostNames := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettings(config.AuthContext(ctx, r.providerConfig))
	apiCreateVirtualHostNames = apiCreateVirtualHostNames.Body(*createVirtualHostNames)
	virtualHostNamesResponse, httpResp, err := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettingsExecute(apiCreateVirtualHostNames)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Virtual Host Names settings", err, httpResp)
		return
	}

	// Read the response into the state
	var state virtualHostNamesModel
	readVirtualHostNamesResponse(ctx, virtualHostNamesResponse, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *virtualHostNamesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualHostNamesModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadVirtualHostNames, httpResp, err := r.apiClient.VirtualHostNamesAPI.GetVirtualHostNamesSettings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Virtual Host Names", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting a Virtual Host Names settings", err, httpResp)
		}
		return
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
	var plan virtualHostNamesModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateVirtualHostNames := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettings(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewVirtualHostNameSettings()
	err := addOptionalVirtualHostNamesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError(providererror.InternalProviderError, "Failed to add optional properties to add request for Virtual Host Names settings: "+err.Error())
		return
	}

	updateVirtualHostNames = updateVirtualHostNames.Body(*createUpdateRequest)
	updateVirtualHostNamesResponse, httpResp, err := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettingsExecute(updateVirtualHostNames)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Virtual Host Names settings", err, httpResp)
		return
	}

	// Read the response
	var state virtualHostNamesModel
	readVirtualHostNamesResponse(ctx, updateVirtualHostNamesResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *virtualHostNamesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this delete will reset the configuration back to the "default" value used by PingFederate.
	var model virtualHostNamesModel
	clientData := model.buildDefaultClientStruct()
	apiUpdateRequest := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	_, httpResp, err := r.apiClient.VirtualHostNamesAPI.UpdateVirtualHostNamesSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the virtual host names settings", err, httpResp)
	}
}

func (r *virtualHostNamesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState virtualHostNamesModel
	emptyState.VirtualHostNames = types.SetNull(types.StringType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
