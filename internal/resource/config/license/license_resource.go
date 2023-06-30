package license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &licenseResource{}
	_ resource.ResourceWithConfigure = &licenseResource{}
)

// LicenseResource is a helper function to simplify the provider implementation.
func LicenseResource() resource.Resource {
	return &licenseResource{}
}

// licenseResource is the resource implementation.
type licenseResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type licenseResourceModel struct {
	Id       types.String `tfsdk:"id"`
	FileData types.String `tfsdk:"file_data"`
}

// GetSchema defines the schema for the resource.
func (r *licenseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	licenseResourceSchema(ctx, req, resp, false)
}

func licenseResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a License.",
		Attributes: map[string]schema.Attribute{
			"file_data": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *licenseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license"
}

func (r *licenseResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readLicenseResponse(ctx context.Context, r *client.LicenseView, state *licenseResourceModel, expectedValues *licenseResourceModel, createPlan types.String) {
	LicenseFileData := createPlan
	state.Id = types.StringValue("id")
	state.FileData = types.StringValue(LicenseFileData.ValueString())
}

func (r *licenseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan licenseResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createLicense := client.NewLicenseFile(plan.FileData.ValueString())
	requestJson, err := createLicense.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateLicense := r.apiClient.LicenseApi.UpdateLicense(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateLicense = apiCreateLicense.Body(*createLicense)
	licenseResponse, httpResp, err := r.apiClient.LicenseApi.UpdateLicenseExecute(apiCreateLicense)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the License", err, httpResp)
		return
	}
	responseJson, err := licenseResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state licenseResourceModel

	readLicenseResponse(ctx, licenseResponse, &state, &state, plan.FileData)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *licenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLicense(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLicense(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state licenseResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadLicense, httpResp, err := apiClient.LicenseApi.GetLicense(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a License", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadLicense.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLicenseResponse(ctx, apiReadLicense, &state, &state, state.FileData)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *licenseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLicense(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLicense(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan licenseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state licenseResourceModel
	req.State.Get(ctx, &state)
	updateLicense := apiClient.LicenseApi.UpdateLicense(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewLicenseFile(plan.FileData.ValueString())
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateLicense = updateLicense.Body(*createUpdateRequest)
	updateLicenseResponse, httpResp, err := apiClient.LicenseApi.UpdateLicenseExecute(updateLicense)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating License", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateLicenseResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readLicenseResponse(ctx, updateLicenseResponse, &state, &state, plan.FileData)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *licenseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *licenseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLicenseLocation(ctx, req, resp)
}

func importLicenseLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
