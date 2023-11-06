package license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &licenseAgreementResource{}
	_ resource.ResourceWithConfigure = &licenseAgreementResource{}
)

// LicenseAgreementResource is a helper function to simplify the provider implementation.
func LicenseAgreementResource() resource.Resource {
	return &licenseAgreementResource{}
}

// licenseAgreementResource is the resource implementation.
type licenseAgreementResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type licenseAgreementResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LicenseAgreementUrl types.String `tfsdk:"license_agreement_url"`
	Accepted            types.Bool   `tfsdk:"accepted"`
}

// GetSchema defines the schema for the resource.
func (r *licenseAgreementResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a LicenseAgreement.",
		Attributes: map[string]schema.Attribute{
			"license_agreement_url": schema.StringAttribute{
				Description: "URL to license agreement",
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accepted": schema.BoolAttribute{
				Description: "Indicates whether license agreement has been accepted. The default value is false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalLicenseAgreementFields(ctx context.Context, addRequest *client.LicenseAgreementInfo, plan licenseAgreementResourceModel) error {

	if internaltypes.IsDefined(plan.LicenseAgreementUrl) {
		addRequest.LicenseAgreementUrl = plan.LicenseAgreementUrl.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Accepted) {
		addRequest.Accepted = plan.Accepted.ValueBoolPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *licenseAgreementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license_agreement"
}

func (r *licenseAgreementResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readLicenseAgreementResponse(ctx context.Context, r *client.LicenseAgreementInfo, state *licenseAgreementResourceModel, existingId *string) {
	state.Id = id.GenerateUUIDToState(existingId)
	state.LicenseAgreementUrl = internaltypes.StringTypeOrNil(r.LicenseAgreementUrl, false)
	state.Accepted = internaltypes.BoolTypeOrNil(r.Accepted)
}

func (r *licenseAgreementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan licenseAgreementResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createLicenseAgreement := client.NewLicenseAgreementInfo()
	err := addOptionalLicenseAgreementFields(ctx, createLicenseAgreement, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for License Agreement", err.Error())
		return
	}

	apiCreateLicenseAgreement := r.apiClient.LicenseAPI.UpdateLicenseAgreement(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateLicenseAgreement = apiCreateLicenseAgreement.Body(*createLicenseAgreement)
	licenseAgreementResponse, httpResp, err := r.apiClient.LicenseAPI.UpdateLicenseAgreementExecute(apiCreateLicenseAgreement)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the License Agreement", err, httpResp)
		return
	}

	// Read the response into the state
	var state licenseAgreementResourceModel
	readLicenseAgreementResponse(ctx, licenseAgreementResponse, &state, nil)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *licenseAgreementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state licenseAgreementResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadLicenseAgreement, httpResp, err := r.apiClient.LicenseAPI.GetLicenseAgreement(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the License Agreement", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the License Agreement", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readLicenseAgreementResponse(ctx, apiReadLicenseAgreement, &state, id)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *licenseAgreementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan licenseAgreementResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateLicenseAgreement := r.apiClient.LicenseAPI.UpdateLicenseAgreement(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewLicenseAgreementInfo()
	err := addOptionalLicenseAgreementFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for License Agreement", err.Error())
		return
	}

	updateLicenseAgreement = updateLicenseAgreement.Body(*createUpdateRequest)
	updateLicenseAgreementResponse, httpResp, err := r.apiClient.LicenseAPI.UpdateLicenseAgreementExecute(updateLicenseAgreement)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating License Agreement", err, httpResp)
		return
	}

	// Get the current state to see how any attributes are changing
	var state licenseAgreementResourceModel
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readLicenseAgreementResponse(ctx, updateLicenseAgreementResponse, &state, id)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *licenseAgreementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *licenseAgreementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
