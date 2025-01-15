// Code generated by ping-terraform-plugin-framework-generator

package secretmanagers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &secretManagerResource{}
	_ resource.ResourceWithConfigure   = &secretManagerResource{}
	_ resource.ResourceWithImportState = &secretManagerResource{}

	customId = "manager_id"
)

func SecretManagerResource() resource.Resource {
	return &secretManagerResource{}
}

type secretManagerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *secretManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_manager"
}

func (r *secretManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type secretManagerResourceModel struct {
	Configuration       types.Object `tfsdk:"configuration"`
	Id                  types.String `tfsdk:"id"`
	ManagerId           types.String `tfsdk:"manager_id"`
	Name                types.String `tfsdk:"name"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
}

func (r *secretManagerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage secret manager plugin instances.",
		Attributes: map[string]schema.Attribute{
			"configuration": pluginconfiguration.ToSchema(),
			"manager_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the plugin instance. This field is immutable and will trigger a replacement plan if changed.<br>Note: Ignored when specifying a connection's adapter override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
				},
				Optional:    true,
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource. This field is immutable and will trigger a replacement plan if changed.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Required:    true,
				Description: "Reference to the plugin descriptor for this instance. This field is immutable and will trigger a replacement plan if changed. Note: Ignored when specifying a connection's adapter override.",
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (r *secretManagerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *secretManagerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	var state *secretManagerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if state == nil {
		return
	}
	var respDiags diag.Diagnostics
	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}

func (model *secretManagerResourceModel) buildClientStruct() (*client.SecretManager, diag.Diagnostics) {
	result := &client.SecretManager{}
	var respDiags diag.Diagnostics
	var err error
	// configuration
	configurationValue, err := pluginconfiguration.ClientStruct(model.Configuration)
	if err != nil {
		respDiags.AddError(providererror.InternalProviderError, "Error building client struct for configuration: "+err.Error())
	} else {
		result.Configuration = *configurationValue
	}

	// manager_id
	result.Id = model.ManagerId.ValueString()
	// name
	result.Name = model.Name.ValueString()
	// parent_ref
	if !model.ParentRef.IsNull() {
		parentRefValue := &client.ResourceLink{}
		parentRefAttrs := model.ParentRef.Attributes()
		parentRefValue.Id = parentRefAttrs["id"].(types.String).ValueString()
		result.ParentRef = parentRefValue
	}

	// plugin_descriptor_ref
	pluginDescriptorRefValue := client.ResourceLink{}
	pluginDescriptorRefAttrs := model.PluginDescriptorRef.Attributes()
	pluginDescriptorRefValue.Id = pluginDescriptorRefAttrs["id"].(types.String).ValueString()
	result.PluginDescriptorRef = pluginDescriptorRefValue

	return result, respDiags
}

func (state *secretManagerResourceModel) readClientResponse(response *client.SecretManager, isImportRead bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringValue(response.Id)
	// configuration
	configurationValue, diags := pluginconfiguration.ToState(state.Configuration, &response.Configuration, isImportRead)
	respDiags.Append(diags...)

	state.Configuration = configurationValue
	// manager_id
	state.ManagerId = types.StringValue(response.Id)
	// name
	state.Name = types.StringValue(response.Name)
	// parent_ref
	parentRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var parentRefValue types.Object
	if response.ParentRef == nil {
		parentRefValue = types.ObjectNull(parentRefAttrTypes)
	} else {
		parentRefValue, diags = types.ObjectValue(parentRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.ParentRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.ParentRef = parentRefValue
	// plugin_descriptor_ref
	pluginDescriptorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	pluginDescriptorRefValue, diags := types.ObjectValue(pluginDescriptorRefAttrTypes, map[string]attr.Value{
		"id": types.StringValue(response.PluginDescriptorRef.Id),
	})
	respDiags.Append(diags...)

	state.PluginDescriptorRef = pluginDescriptorRefValue
	return respDiags
}

func (r *secretManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data secretManagerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.SecretManagersAPI.CreateSecretManager(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SecretManagersAPI.CreateSecretManagerExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the secretManager", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *secretManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var data secretManagerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.SecretManagersAPI.GetSecretManager(config.AuthContext(ctx, r.providerConfig), data.ManagerId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "Secret Manager", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the secretManager", err, httpResp, &customId)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, isImportRead)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *secretManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data secretManagerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.SecretManagersAPI.UpdateSecretManager(config.AuthContext(ctx, r.providerConfig), data.ManagerId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.SecretManagersAPI.UpdateSecretManagerExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the secretManager", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *secretManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data secretManagerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.SecretManagersAPI.DeleteSecretManager(config.AuthContext(ctx, r.providerConfig), data.ManagerId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the secretManager", err, httpResp, &customId)
	}
}

func (r *secretManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to manager_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("manager_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}
