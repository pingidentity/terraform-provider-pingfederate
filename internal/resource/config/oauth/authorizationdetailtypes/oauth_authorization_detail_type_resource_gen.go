// Code generated by ping-terraform-plugin-framework-generator

package oauthauthorizationdetailtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	_ resource.Resource                = &oauthAuthorizationDetailTypeResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthorizationDetailTypeResource{}
	_ resource.ResourceWithImportState = &oauthAuthorizationDetailTypeResource{}
)

func OauthAuthorizationDetailTypeResource() resource.Resource {
	return &oauthAuthorizationDetailTypeResource{}
}

type oauthAuthorizationDetailTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *oauthAuthorizationDetailTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_authorization_detail_type"
}

func (r *oauthAuthorizationDetailTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type oauthAuthorizationDetailTypeResourceModel struct {
	Active                          types.Bool   `tfsdk:"active"`
	AuthorizationDetailProcessorRef types.Object `tfsdk:"authorization_detail_processor_ref"`
	Description                     types.String `tfsdk:"description"`
	Id                              types.String `tfsdk:"id"`
	Type                            types.String `tfsdk:"type"`
	TypeId                          types.String `tfsdk:"type_id"`
}

func (r *oauthAuthorizationDetailTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage an OAuth authorization detail type.",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether or not this authorization detail type is active. Defaults to `true`.",
			},
			"authorization_detail_processor_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Required:    true,
				Description: "The authentication detail processor used to process this type.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the authorization detail type.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The authorization detail type.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the authorization detail type. The ID will be system-assigned if not specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (model *oauthAuthorizationDetailTypeResourceModel) buildClientStruct() (*client.AuthorizationDetailType, diag.Diagnostics) {
	result := &client.AuthorizationDetailType{}
	// active
	result.Active = model.Active.ValueBoolPointer()
	// authorization_detail_processor_ref
	authorizationDetailProcessorRefValue := client.ResourceLink{}
	authorizationDetailProcessorRefAttrs := model.AuthorizationDetailProcessorRef.Attributes()
	authorizationDetailProcessorRefValue.Id = authorizationDetailProcessorRefAttrs["id"].(types.String).ValueString()
	result.AuthorizationDetailProcessorRef = authorizationDetailProcessorRefValue

	// description
	result.Description = model.Description.ValueStringPointer()
	// type
	result.Type = model.Type.ValueString()
	// type_id
	result.Id = model.TypeId.ValueStringPointer()
	return result, nil
}

func (state *oauthAuthorizationDetailTypeResourceModel) readClientResponse(response *client.AuthorizationDetailType) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringPointerValue(response.Id)
	// active
	state.Active = types.BoolPointerValue(response.Active)
	// authorization_detail_processor_ref
	authorizationDetailProcessorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	authorizationDetailProcessorRefValue, diags := types.ObjectValue(authorizationDetailProcessorRefAttrTypes, map[string]attr.Value{
		"id": types.StringValue(response.AuthorizationDetailProcessorRef.Id),
	})
	respDiags.Append(diags...)

	state.AuthorizationDetailProcessorRef = authorizationDetailProcessorRefValue
	// description
	state.Description = types.StringPointerValue(response.Description)
	// type
	state.Type = types.StringValue(response.Type)
	// type_id
	state.TypeId = types.StringPointerValue(response.Id)
	return respDiags
}

func (r *oauthAuthorizationDetailTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data oauthAuthorizationDetailTypeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.OauthAuthorizationDetailTypesAPI.AddAuthorizationDetailType(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthAuthorizationDetailTypesAPI.AddAuthorizationDetailTypeExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the oauthAuthorizationDetailType", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthAuthorizationDetailTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data oauthAuthorizationDetailTypeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.OauthAuthorizationDetailTypesAPI.GetAuthorizationDetailTypeById(config.AuthContext(ctx, r.providerConfig), data.TypeId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "oauthAuthorizationDetailType", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the oauthAuthorizationDetailType", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthAuthorizationDetailTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data oauthAuthorizationDetailTypeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.OauthAuthorizationDetailTypesAPI.UpdateAuthorizationDetailType(config.AuthContext(ctx, r.providerConfig), data.TypeId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthAuthorizationDetailTypesAPI.UpdateAuthorizationDetailTypeExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the oauthAuthorizationDetailType", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthAuthorizationDetailTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data oauthAuthorizationDetailTypeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.OauthAuthorizationDetailTypesAPI.DeleteAuthorizationDetailType(config.AuthContext(ctx, r.providerConfig), data.TypeId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the oauthAuthorizationDetailType", err, httpResp)
	}
}

func (r *oauthAuthorizationDetailTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to type_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("type_id"), req, resp)
}
