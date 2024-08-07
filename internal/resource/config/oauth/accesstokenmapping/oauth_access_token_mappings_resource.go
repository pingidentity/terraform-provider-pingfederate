package oauthaccesstokenmapping

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributesources"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/issuancecriteria"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAccessTokenMappingsResource{}
	_ resource.ResourceWithConfigure   = &oauthAccessTokenMappingsResource{}
	_ resource.ResourceWithImportState = &oauthAccessTokenMappingsResource{}

	accessTokenMappingContext = map[string]attr.Type{
		"type":        types.StringType,
		"context_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}
)

// OauthAccessTokenMappingsResource is a helper function to simplify the provider implementation.
func OauthAccessTokenMappingsResource() resource.Resource {
	return &oauthAccessTokenMappingsResource{}
}

// oauthAccessTokenMappingsResource is the resource implementation.
type oauthAccessTokenMappingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAccessTokenMappingsResourceModel struct {
	Id                           types.String `tfsdk:"id"`
	Context                      types.Object `tfsdk:"context"`
	AccessTokenManagerRef        types.Object `tfsdk:"access_token_manager_ref"`
	AttributeSources             types.Set    `tfsdk:"attribute_sources"`
	AttributeContractFulfillment types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria             types.Object `tfsdk:"issuance_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *oauthAccessTokenMappingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an OAuth Access Token Mapping",
		Attributes: map[string]schema.Attribute{
			"context": schema.SingleNestedAttribute{
				Description: "The context of the OAuth Access Token Mapping. This property cannot be changed after the mapping is created.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The Access Token Mapping Context type.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("DEFAULT", "PCV", "IDP_CONNECTION", "IDP_ADAPTER", "AUTHENTICATION_POLICY_CONTRACT", "CLIENT_CREDENTIALS", "TOKEN_EXCHANGE_PROCESSOR_POLICY"),
						},
					},
					"context_ref": schema.SingleNestedAttribute{
						Description: "Reference to the associated Access Token Mapping Context instance.",
						Computed:    true,
						Optional:    true,
						Attributes:  resourcelink.ToSchema(),
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"access_token_manager_ref": schema.SingleNestedAttribute{
				Description: "Reference to the access token manager this mapping is associated with. This property cannot be changed after the mapping is created.",
				Required:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: resourcelink.ToSchema(),
			},
			"attribute_sources": attributesources.ToSchema(0, true),
			"attribute_contract_fulfillment": schema.MapNestedAttribute{
				Description: "Defines how an attribute in an attribute contract should be populated.",
				Required:    true,
				Optional:    false,
				Computed:    false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": sourcetypeidkey.ToSchema(false),
						"value": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The value for this attribute.",
							Default:     stringdefault.StaticString(""),
						},
					},
				},
			},
			"issuance_criteria": issuancecriteria.ToSchema(),
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *oauthAccessTokenMappingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_access_token_mapping"
}

func (r *oauthAccessTokenMappingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthAccessTokenMappingsResponse(ctx context.Context, r *client.AccessTokenMapping, state *oauthAccessTokenMappingsResourceModel) diag.Diagnostics {
	var diags, objDiags diag.Diagnostics

	state.Id = types.StringPointerValue(r.Id)

	contextRefObjValue, objDiags := resourcelink.ToState(ctx, &r.Context.ContextRef)
	diags.Append(objDiags...)
	contextAttrValue := map[string]attr.Value{
		"type":        types.StringValue(r.Context.Type),
		"context_ref": contextRefObjValue,
	}
	state.Context, objDiags = types.ObjectValue(accessTokenMappingContext, contextAttrValue)
	diags.Append(objDiags...)
	state.AccessTokenManagerRef, objDiags = resourcelink.ToState(ctx, &r.AccessTokenManagerRef)
	diags.Append(objDiags...)
	state.AttributeSources, objDiags = attributesources.ToState(ctx, r.AttributeSources)
	diags.Append(objDiags...)
	state.AttributeContractFulfillment, objDiags = attributecontractfulfillment.ToState(ctx, &r.AttributeContractFulfillment)
	diags.Append(objDiags...)
	state.IssuanceCriteria, objDiags = issuancecriteria.ToState(ctx, r.IssuanceCriteria)
	diags.Append(objDiags...)

	// make sure all object type building appends diags
	return diags
}

func addOptionalOauthAccessTokenMappingsFields(addRequest *client.AccessTokenMapping, plan oauthAccessTokenMappingsResourceModel) error {
	var err error
	if internaltypes.IsDefined(plan.AttributeSources) {
		addRequest.AttributeSources = []client.AttributeSourceAggregation{}
		addRequest.AttributeSources, err = attributesources.ClientStruct(plan.AttributeSources)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.IssuanceCriteria) {
		addRequest.IssuanceCriteria, err = issuancecriteria.ClientStruct(plan.IssuanceCriteria)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *oauthAccessTokenMappingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model oauthAccessTokenMappingsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if internaltypes.IsDefined(model.Context) {
		modelContextType := model.Context.Attributes()["type"].(types.String).ValueString()
		modelContextContextRef := model.Context.Attributes()["context_ref"].(types.Object)
		if (modelContextType == "DEFAULT" || modelContextType == "CLIENT_CREDENTIALS") && internaltypes.IsDefined(modelContextContextRef) {
			resp.Diagnostics.AddError("Invalid attribute combination",
				"context_ref is not required for the Access Token Mapping Context type: "+modelContextType)
		}
	}
}

func (r *oauthAccessTokenMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAccessTokenMappingsResourceModel
	var err error

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var hasObjectErrMap = make(map[error]bool)
	accessTokenMappingContext := &client.AccessTokenMappingContext{}
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Context, true)), accessTokenMappingContext)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	accessTokenManagerRef, err := resourcelink.ClientStruct(plan.AccessTokenManagerRef)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	attributeContractFulfillment, err := attributecontractfulfillment.ClientStruct(plan.AttributeContractFulfillment)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	for errorVal, hasErr := range hasObjectErrMap {
		if hasErr {
			resp.Diagnostics.AddError("Failed to create item for request object:", errorVal.Error())
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAccessTokenMappings := client.NewAccessTokenMapping(*accessTokenMappingContext, *accessTokenManagerRef, attributeContractFulfillment)

	err = addOptionalOauthAccessTokenMappingsFields(createOauthAccessTokenMappings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Access Token Mapping", err.Error())
		return
	}
	apiCreateOauthAccessTokenMappings := r.apiClient.OauthAccessTokenMappingsAPI.CreateMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAccessTokenMappings = apiCreateOauthAccessTokenMappings.Body(*createOauthAccessTokenMappings)
	oauthAccessTokenMappingsResponse, httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.CreateMappingExecute(apiCreateOauthAccessTokenMappings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OAuth Access Token Mapping", err, httpResp)
		return
	}

	// Read the response into the state
	var state oauthAccessTokenMappingsResourceModel

	diags = readOauthAccessTokenMappingsResponse(ctx, oauthAccessTokenMappingsResponse, &state)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAccessTokenMappingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAccessTokenMappings, httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.GetMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "OAuth Access Token Mapping", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the OAuth Access Token Mapping", err, httpResp)
		}
		return
	}

	// Read the response into the state
	diags = readOauthAccessTokenMappingsResponse(ctx, apiReadOauthAccessTokenMappings, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAccessTokenMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan oauthAccessTokenMappingsResourceModel
	var err error
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var hasObjectErrMap = make(map[error]bool)
	accessTokenMappingContext := &client.AccessTokenMappingContext{}
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Context, true)), accessTokenMappingContext)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	accessTokenManagerRef, err := resourcelink.ClientStruct(plan.AccessTokenManagerRef)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	attributeContractFulfillment, err := attributecontractfulfillment.ClientStruct(plan.AttributeContractFulfillment)
	if err != nil {
		hasObjectErrMap[err] = true
	}

	for errorVal, hasErr := range hasObjectErrMap {
		if hasErr {
			resp.Diagnostics.AddError("Failed to create item for request object:", errorVal.Error())
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}
	updateOauthAccessTokenMappings := client.NewAccessTokenMapping(*accessTokenMappingContext, *accessTokenManagerRef, attributeContractFulfillment)

	err = addOptionalOauthAccessTokenMappingsFields(updateOauthAccessTokenMappings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OAuth Access Token Mapping", err.Error())
		return
	}

	apiUpdateOauthAccessTokenMappings := r.apiClient.OauthAccessTokenMappingsAPI.UpdateMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	apiUpdateOauthAccessTokenMappings = apiUpdateOauthAccessTokenMappings.Body(*updateOauthAccessTokenMappings)
	updateOauthAccessTokenMappingsResponse, httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.UpdateMappingExecute(apiUpdateOauthAccessTokenMappings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OAuth Access Token Mapping", err, httpResp)
		return
	}

	// Read the response
	var state oauthAccessTokenMappingsResourceModel
	diags = readOauthAccessTokenMappingsResponse(ctx, updateOauthAccessTokenMappingsResponse, &state)
	resp.Diagnostics.Append(diags...)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthAccessTokenMappingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthAccessTokenMappingsAPI.DeleteMapping(config.AuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an OAuth Access Token Mapping", err, httpResp)
	}
}

func (r *oauthAccessTokenMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
