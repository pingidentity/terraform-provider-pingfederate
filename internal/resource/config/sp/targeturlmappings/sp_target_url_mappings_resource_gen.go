// Code generated by ping-terraform-plugin-framework-generator

package sptargeturlmappings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/configvalidators"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &spTargetUrlMappingsResource{}
	_ resource.ResourceWithConfigure   = &spTargetUrlMappingsResource{}
	_ resource.ResourceWithImportState = &spTargetUrlMappingsResource{}
)

func SpTargetUrlMappingsResource() resource.Resource {
	return &spTargetUrlMappingsResource{}
}

type spTargetUrlMappingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *spTargetUrlMappingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sp_target_url_mappings"
}

func (r *spTargetUrlMappingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type spTargetUrlMappingsResourceModel struct {
	Items types.List `tfsdk:"items"`
}

func (r *spTargetUrlMappingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	itemsRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	itemsAttrTypes := map[string]attr.Type{
		"ref":  types.ObjectType{AttrTypes: itemsRefAttrTypes},
		"type": types.StringType,
		"url":  types.StringType,
	}
	itemsElementType := types.ObjectType{AttrTypes: itemsAttrTypes}
	itemsDefault, diags := types.ListValue(itemsElementType, []attr.Value{})
	resp.Diagnostics.Append(diags...)
	resp.Schema = schema.Schema{
		Description: "Resource to manage the mappings between URLs and adapter or connection instances.",
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ref": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Required:    true,
									Description: "The ID of the resource.",
								},
							},
							Required:    true,
							Description: "The adapter or connection instance mapped for this URL.",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The URL mapping type. Options are `SP_ADAPTER` or `SP_CONNECTION`.",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"SP_ADAPTER",
									"SP_CONNECTION",
								),
							},
						},
						"url": schema.StringAttribute{
							Required:    true,
							Description: "The URL that will be compared against the target URL. Use a wildcard (*) to match multiple URLs to the same adapter or connection instance.",
							Validators: []validator.String{
								configvalidators.ValidUrl(),
							},
						},
					},
				},
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(itemsDefault),
				Description: "The actual list of SP connection URL mappings. The order of the items in this list determines the order in which the mappings are evaluated.",
			},
		},
	}
}

func (model *spTargetUrlMappingsResourceModel) buildClientStruct() *client.SpUrlMappings {
	result := &client.SpUrlMappings{}
	// items
	result.Items = []client.SpUrlMapping{}
	for _, itemsElement := range model.Items.Elements() {
		itemsValue := client.SpUrlMapping{}
		itemsAttrs := itemsElement.(types.Object).Attributes()
		if !itemsAttrs["ref"].IsNull() {
			itemsRefValue := &client.ResourceLink{}
			itemsRefAttrs := itemsAttrs["ref"].(types.Object).Attributes()
			itemsRefValue.Id = itemsRefAttrs["id"].(types.String).ValueString()
			itemsValue.Ref = itemsRefValue
		}
		itemsValue.Type = itemsAttrs["type"].(types.String).ValueStringPointer()
		itemsValue.Url = itemsAttrs["url"].(types.String).ValueStringPointer()
		result.Items = append(result.Items, itemsValue)
	}

	return result
}

func (state *spTargetUrlMappingsResourceModel) readClientResponse(response *client.SpUrlMappings) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// items
	itemsRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	itemsAttrTypes := map[string]attr.Type{
		"ref":  types.ObjectType{AttrTypes: itemsRefAttrTypes},
		"type": types.StringType,
		"url":  types.StringType,
	}
	itemsElementType := types.ObjectType{AttrTypes: itemsAttrTypes}
	var itemsValues []attr.Value
	for _, itemsResponseValue := range response.Items {
		var itemsRefValue types.Object
		if itemsResponseValue.Ref == nil {
			itemsRefValue = types.ObjectNull(itemsRefAttrTypes)
		} else {
			itemsRefValue, diags = types.ObjectValue(itemsRefAttrTypes, map[string]attr.Value{
				"id": types.StringValue(itemsResponseValue.Ref.Id),
			})
			respDiags.Append(diags...)
		}
		itemsValue, diags := types.ObjectValue(itemsAttrTypes, map[string]attr.Value{
			"ref":  itemsRefValue,
			"type": types.StringPointerValue(itemsResponseValue.Type),
			"url":  types.StringPointerValue(itemsResponseValue.Url),
		})
		respDiags.Append(diags...)
		itemsValues = append(itemsValues, itemsValue)
	}
	itemsValue, diags := types.ListValue(itemsElementType, itemsValues)
	respDiags.Append(diags...)

	state.Items = itemsValue
	return respDiags
}

func (r *spTargetUrlMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data spTargetUrlMappingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic, since this is a singleton resource
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.SpTargetUrlMappingsAPI.UpdateSpUrlMappings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SpTargetUrlMappingsAPI.UpdateSpUrlMappingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the spTargetUrlMappings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spTargetUrlMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data spTargetUrlMappingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.SpTargetUrlMappingsAPI.GetSpUrlMappings(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the spTargetUrlMappings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the spTargetUrlMappings", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spTargetUrlMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data spTargetUrlMappingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.SpTargetUrlMappingsAPI.UpdateSpUrlMappings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SpTargetUrlMappingsAPI.UpdateSpUrlMappingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the spTargetUrlMappings", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *spTargetUrlMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// This resource has no identifier attributes, so the value passed in here doesn't matter. Just return an empty state struct.
	var emptyState spTargetUrlMappingsResourceModel
	emptyState.setNullObjectValues()
	resp.Diagnostics.Append(resp.State.Set(ctx, &emptyState)...)
}
