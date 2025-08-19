// Copyright © 2025 Ping Identity Corporation

package extendedproperties

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &extendedPropertiesResource{}
	_ resource.ResourceWithConfigure   = &extendedPropertiesResource{}
	_ resource.ResourceWithImportState = &extendedPropertiesResource{}

	extendedPropertyAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":         types.StringType,
			"description":  types.StringType,
			"multi_valued": types.BoolType,
		},
	}
)

// ExtendedPropertiesResource is a helper function to simplify the provider implementation.
func ExtendedPropertiesResource() resource.Resource {
	return &extendedPropertiesResource{}
}

// extendedPropertiesResource is the resource implementation.
type extendedPropertiesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type extendedPropertiesResourceModel struct {
	Items types.Set `tfsdk:"items"`
}

// GetSchema defines the schema for the resource.
func (r *extendedPropertiesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Extended Properties definitions",
		Attributes: map[string]schema.Attribute{
			"items": schema.SetNestedAttribute{
				Description: "A collection of Extended Properties definitions.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The property name.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"description": schema.StringAttribute{
							Description: "The property description.",
							Optional:    true,
						},
						"multi_valued": schema.BoolAttribute{
							Description: "Indicates whether the property should allow multiple values. Default value is `false`.",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
	resp.Schema = schema
}

func addExtendedPropertiesFields(addRequest *client.ExtendedProperties, plan extendedPropertiesResourceModel) {
	// items
	if !plan.Items.IsNull() && !plan.Items.IsUnknown() {
		addRequest.Items = []client.ExtendedProperty{}
		for _, itemsElement := range plan.Items.Elements() {
			itemsValue := client.ExtendedProperty{}
			itemsAttrs := itemsElement.(types.Object).Attributes()
			itemsValue.Description = itemsAttrs["description"].(types.String).ValueStringPointer()
			itemsValue.MultiValued = itemsAttrs["multi_valued"].(types.Bool).ValueBoolPointer()
			itemsValue.Name = itemsAttrs["name"].(types.String).ValueString()
			addRequest.Items = append(addRequest.Items, itemsValue)
		}
	}
}

// Metadata returns the resource type name.
func (r *extendedPropertiesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extended_properties"
}

func (r *extendedPropertiesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readExtendedPropertiesResponse(ctx context.Context, r *client.ExtendedProperties, state *extendedPropertiesResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.Items, diags = types.SetValueFrom(ctx, extendedPropertyAttrType, r.GetItems())

	// make sure all object type building appends diags
	return diags
}

func (r *extendedPropertiesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan extendedPropertiesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createExtendedProperties := client.NewExtendedProperties()
	addExtendedPropertiesFields(createExtendedProperties, plan)

	apiCreateExtendedProperties := r.apiClient.ExtendedPropertiesAPI.UpdateExtendedProperties(config.AuthContext(ctx, r.providerConfig))
	apiCreateExtendedProperties = apiCreateExtendedProperties.Body(*createExtendedProperties)
	extendedPropertiesResponse, httpResp, err := r.apiClient.ExtendedPropertiesAPI.UpdateExtendedPropertiesExecute(apiCreateExtendedProperties)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the extended properties", err, httpResp)
		return
	}

	// Read the response into the state
	var state extendedPropertiesResourceModel

	diags = readExtendedPropertiesResponse(ctx, extendedPropertiesResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *extendedPropertiesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state extendedPropertiesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadExtendedProperties, httpResp, err := r.apiClient.ExtendedPropertiesAPI.GetExtendedProperties(config.AuthContext(ctx, r.providerConfig)).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Extended Properties", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the extended properties", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readExtendedPropertiesResponse(ctx, apiReadExtendedProperties, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *extendedPropertiesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan extendedPropertiesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateExtendedProperties := r.apiClient.ExtendedPropertiesAPI.UpdateExtendedProperties(config.AuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewExtendedProperties()
	addExtendedPropertiesFields(createUpdateRequest, plan)

	updateExtendedProperties = updateExtendedProperties.Body(*createUpdateRequest)
	updateExtendedPropertiesResponse, httpResp, err := r.apiClient.ExtendedPropertiesAPI.UpdateExtendedPropertiesExecute(updateExtendedProperties)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating extended properties", err, httpResp)
		return
	}

	// Read the response
	var state extendedPropertiesResourceModel
	diags = readExtendedPropertiesResponse(ctx, updateExtendedPropertiesResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *extendedPropertiesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *extendedPropertiesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState extendedPropertiesResourceModel
	emptyState.Items = types.SetNull(extendedPropertyAttrType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
