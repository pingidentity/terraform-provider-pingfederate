// Copyright © 2025 Ping Identity Corporation

package license

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &licenseResource{}
	_ resource.ResourceWithConfigure = &licenseResource{}

	licenseGroupsAttrTypes = map[string]attr.Type{
		"name":             types.StringType,
		"connection_count": types.Int64Type,
		"start_date":       types.StringType,
		"end_date":         types.StringType,
	}

	featuresAttrTypes = map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
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
	FileData types.String `tfsdk:"file_data"`
	// Computed attributes
	Name                types.String `tfsdk:"name"`
	MaxConnections      types.Int64  `tfsdk:"max_connections"`
	UsedConnections     types.Int64  `tfsdk:"used_connections"`
	Tier                types.String `tfsdk:"tier"`
	IssueDate           types.String `tfsdk:"issue_date"`
	ExpirationDate      types.String `tfsdk:"expiration_date"`
	EnforcementType     types.String `tfsdk:"enforcement_type"`
	Version             types.String `tfsdk:"version"`
	Product             types.String `tfsdk:"product"`
	Organization        types.String `tfsdk:"organization"`
	GracePeriod         types.Int64  `tfsdk:"grace_period"`
	NodeLimit           types.Int64  `tfsdk:"node_limit"`
	LicenseGroups       types.List   `tfsdk:"license_groups"`
	OauthEnabled        types.Bool   `tfsdk:"oauth_enabled"`
	WsTrustEnabled      types.Bool   `tfsdk:"ws_trust_enabled"`
	ProvisioningEnabled types.Bool   `tfsdk:"provisioning_enabled"`
	BridgeMode          types.Bool   `tfsdk:"bridge_mode"`
	Features            types.List   `tfsdk:"features"`
}

// GetSchema defines the schema for the resource.
func (r *licenseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a license summary object.",
		Attributes: map[string]schema.Attribute{
			"file_data": schema.StringAttribute{
				Description: "The license file data. This field is immutable and will trigger a replacement plan if changed.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the person the license was issued to.",
				Computed:    true,
			},
			"max_connections": schema.Int64Attribute{
				Description: "Maximum number of connections that can be created under this license (if applicable).",
				Computed:    true,
			},
			"used_connections": schema.Int64Attribute{
				Description: "Number of used connections under this license.",
				Computed:    true,
			},
			"tier": schema.StringAttribute{
				Description: "The tier value from the license file. The possible values are FREE, PERPETUAL or SUBSCRIPTION.",
				Computed:    true,
			},
			"issue_date": schema.StringAttribute{
				Description: "The issue date value from the license file.",
				Computed:    true,
			},
			"expiration_date": schema.StringAttribute{
				Description: "The expiration date value from the license file (if applicable).",
				Computed:    true,
			},
			"enforcement_type": schema.StringAttribute{
				Description: "The enforcement type is a 3-bit binary value, expressed as a decimal digit. The bits from left to right are: 1: Shutdown on expire. 2: Notify on expire. 4: Enforce minor version. if all three enforcements are active, the enforcement type will be 7 (1 + 2 + 4); if only the first two are active, you have an enforcement type of 3 (1 + 2).",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The Ping Identity product version from the license file.",
				Computed:    true,
			},
			"product": schema.StringAttribute{
				Description: "The Ping Identity product value from the license file.",
				Computed:    true,
			},
			"organization": schema.StringAttribute{
				Description: "The organization value from the license file.",
				Computed:    true,
			},
			"grace_period": schema.Int64Attribute{
				Description: "Number of days provided as grace period, past the expiration date (if applicable).",
				Computed:    true,
			},
			"node_limit": schema.Int64Attribute{
				Description: "Maximum number of clustered nodes allowed under this license (if applicable).",
				Computed:    true,
			},
			"license_groups": schema.ListNestedAttribute{
				Description: "License connection groups, if applicable.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Group name from the license file.",
							Computed:    true,
						},
						"connection_count": schema.Int64Attribute{
							Description: "Maximum number of connections permitted under the group.",
							Computed:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "Start date for the group.",
							Computed:    true,
						},
						"end_date": schema.StringAttribute{
							Description: "End date for the group.",
							Computed:    true,
						},
					},
				},
			},
			"oauth_enabled": schema.BoolAttribute{
				Description: "Indicates whether OAuth role is enabled for this license.",
				Computed:    true,
			},
			"ws_trust_enabled": schema.BoolAttribute{
				Description: "Indicates whether WS-Trust role is enabled for this license.",
				Computed:    true,
			},
			"provisioning_enabled": schema.BoolAttribute{
				Description: "Indicates whether Provisioning role is enabled for this license.",
				Computed:    true,
			},
			"bridge_mode": schema.BoolAttribute{
				Description: "Indicates whether this license is a bridge license or not.",
				Computed:    true,
			},
			"features": schema.ListNestedAttribute{
				Description: "Other licence features, if applicable.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the license feature.",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the license feature.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
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

func readLicenseResponse(ctx context.Context, r *client.LicenseView, state *licenseResourceModel, planFileData types.String) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.FileData = types.StringValue(planFileData.ValueString())

	state.Name = types.StringPointerValue(r.Name)
	state.MaxConnections = types.Int64PointerValue(r.MaxConnections)
	state.UsedConnections = types.Int64PointerValue(r.UsedConnections)
	state.Tier = types.StringPointerValue(r.Tier)
	if r.IssueDate != nil {
		state.IssueDate = types.StringValue(r.IssueDate.Format(time.RFC3339))
	} else {
		state.IssueDate = types.StringNull()
	}
	if r.ExpirationDate != nil {
		state.ExpirationDate = types.StringValue(r.ExpirationDate.Format(time.RFC3339))
	} else {
		state.ExpirationDate = types.StringNull()
	}
	state.EnforcementType = types.StringPointerValue(r.EnforcementType)
	state.Version = types.StringPointerValue(r.Version)
	state.Product = types.StringPointerValue(r.Product)
	state.Organization = types.StringPointerValue(r.Organization)
	state.GracePeriod = types.Int64PointerValue(r.GracePeriod)
	state.NodeLimit = types.Int64PointerValue(r.NodeLimit)
	state.OauthEnabled = types.BoolValue(*r.OauthEnabled)
	state.WsTrustEnabled = types.BoolValue(*r.WsTrustEnabled)
	state.ProvisioningEnabled = types.BoolValue(*r.ProvisioningEnabled)
	state.BridgeMode = types.BoolValue(*r.BridgeMode)

	state.LicenseGroups, respDiags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: licenseGroupsAttrTypes}, r.LicenseGroups)
	diags.Append(respDiags...)

	state.Features, respDiags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: featuresAttrTypes}, r.Features)
	diags.Append(respDiags...)
	return diags
}

func (r *licenseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan licenseResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createLicense := client.NewLicenseFile(plan.FileData.ValueString())
	apiCreateLicense := r.apiClient.LicenseAPI.UpdateLicense(config.AuthContext(ctx, r.providerConfig))
	apiCreateLicense = apiCreateLicense.Body(*createLicense)
	licenseResponse, httpResp, err := r.apiClient.LicenseAPI.UpdateLicenseExecute(apiCreateLicense)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the license", err, httpResp)
		return
	}

	// Read the response into the state
	var state licenseResourceModel
	diags = readLicenseResponse(ctx, licenseResponse, &state, plan.FileData)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *licenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state licenseResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadLicense, httpResp, err := r.apiClient.LicenseAPI.GetLicense(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "License", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the license", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readLicenseResponse(ctx, apiReadLicense, &state, state.FileData)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *licenseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan licenseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateLicense := r.apiClient.LicenseAPI.UpdateLicense(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewLicenseFile(plan.FileData.ValueString())
	updateLicense = updateLicense.Body(*createUpdateRequest)
	updateLicenseResponse, httpResp, err := r.apiClient.LicenseAPI.UpdateLicenseExecute(updateLicense)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the license summary", err, httpResp)
		return
	}

	// Read the response
	var state licenseResourceModel
	diags = readLicenseResponse(ctx, updateLicenseResponse, &state, plan.FileData)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *licenseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	providererror.WarnConfigurationCannotBeReset("pingfederate_license", &resp.Diagnostics)
}
