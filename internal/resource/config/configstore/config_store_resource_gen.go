// Code generated by ping-terraform-plugin-framework-generator

package configstore

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &configStoreResource{}
	_ resource.ResourceWithConfigure   = &configStoreResource{}
	_ resource.ResourceWithImportState = &configStoreResource{}
)

func ConfigStoreResource() resource.Resource {
	return &configStoreResource{}
}

type configStoreResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *configStoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config_store"
}

func (r *configStoreResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type configStoreResourceModel struct {
	Bundle      types.String `tfsdk:"bundle"`
	Id          types.String `tfsdk:"id"`
	ListValue   types.List   `tfsdk:"list_value"`
	MapValue    types.Map    `tfsdk:"map_value"`
	SettingId   types.String `tfsdk:"setting_id"`
	StringValue types.String `tfsdk:"string_value"`
}

func (r *configStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage bundle settings.",
		Attributes: map[string]schema.Attribute{
			"bundle": schema.StringAttribute{
				Required:    true,
				Description: "This field represents a configuration file that contains a bundle of settings.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"list_value": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The list of values for the configuration setting. This is used when the setting has a list of string values. Exactly one of `list_value`, `map_value`, or `string_value` must be set. Changing the type of the setting will require deletion and recreation of the setting.",
				Validators: []validator.List{
					listvalidator.ExactlyOneOf(
						path.MatchRelative().AtParent().AtName("map_value"),
						path.MatchRelative().AtParent().AtName("string_value"),
					),
				},
			},
			"map_value": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The map of key/value pairs for the configuration setting. This is used when the setting has a map of string keys and values. Exactly one of `list_value`, `map_value`, or `string_value` must be set. Changing the type of the setting will require deletion and recreation of the setting.",
			},
			"setting_id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the configuration setting.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"string_value": schema.StringAttribute{
				Optional:    true,
				Description: "The value of the configuration setting. This is used when the setting has a single string value. Exactly one of `list_value`, `map_value`, or `string_value` must be set. Changing the type of the setting will require deletion and recreation of the setting.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (model *configStoreResourceModel) buildClientStruct() (*client.ConfigStoreSetting, diag.Diagnostics) {
	result := &client.ConfigStoreSetting{}
	// list_value
	if !model.ListValue.IsNull() {
		result.ListValue = []string{}
		for _, listValueElement := range model.ListValue.Elements() {
			result.ListValue = append(result.ListValue, listValueElement.(types.String).ValueString())
		}
	}

	// map_value
	if !model.MapValue.IsNull() {
		result.MapValue = &map[string]string{}
		for key, mapValueElement := range model.MapValue.Elements() {
			(*result.MapValue)[key] = mapValueElement.(types.String).ValueString()
		}
	}

	// setting_id
	result.Id = model.SettingId.ValueString()
	// string_value
	result.StringValue = model.StringValue.ValueStringPointer()
	// type
	model.setType(result)
	return result, nil
}

func (state *configStoreResourceModel) readClientResponse(response *client.ConfigStoreSetting) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// list_value
	if len(response.ListValue) > 0 || !state.ListValue.IsNull() {
		state.ListValue, diags = types.ListValueFrom(context.Background(), types.StringType, response.ListValue)
		respDiags.Append(diags...)
	} else {
		state.ListValue = types.ListNull(types.StringType)
	}
	// map_value
	if response.MapValue == nil {
		state.MapValue = types.MapNull(types.StringType)
	} else {
		state.MapValue, diags = types.MapValueFrom(context.Background(), types.StringType, (*response.MapValue))
		respDiags.Append(diags...)
	}
	// setting_id
	state.SettingId = types.StringValue(response.Id)
	// id
	state.Id = types.StringValue(response.Id)
	// string_value
	state.StringValue = types.StringPointerValue(response.StringValue)
	return respDiags
}

func (r *configStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data configStoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.ConfigStoreAPI.UpdateConfigStoreSetting(config.AuthContext(ctx, r.providerConfig), data.Bundle.ValueString(), data.SettingId.ValueString())
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ConfigStoreAPI.UpdateConfigStoreSettingExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the configStore", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *configStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data configStoreResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.ConfigStoreAPI.GetConfigStoreSetting(config.AuthContext(ctx, r.providerConfig), data.Bundle.ValueString(), data.SettingId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "configStore", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the configStore", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *configStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data configStoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.ConfigStoreAPI.UpdateConfigStoreSetting(config.AuthContext(ctx, r.providerConfig), data.Bundle.ValueString(), data.SettingId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.ConfigStoreAPI.UpdateConfigStoreSettingExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the configStore", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *configStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data configStoreResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.ConfigStoreAPI.DeleteConfigStoreSetting(config.AuthContext(ctx, r.providerConfig), data.Bundle.ValueString(), data.SettingId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the configStore", err, httpResp)
	}
}

func (r *configStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [bundle]/[setting_id]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bundle"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("setting_id"), split[1])...)
}
