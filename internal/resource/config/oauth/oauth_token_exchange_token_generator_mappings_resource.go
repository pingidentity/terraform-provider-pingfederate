package oauth

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthTokenExchangeTokenGeneratorMappingsResource{}
	_ resource.ResourceWithConfigure   = &oauthTokenExchangeTokenGeneratorMappingsResource{}
	_ resource.ResourceWithImportState = &oauthTokenExchangeTokenGeneratorMappingsResource{}
)

// OauthTokenExchangeTokenGeneratorMappingsResource is a helper function to simplify the provider implementation.
func OauthTokenExchangeTokenGeneratorMappingsResource() resource.Resource {
	return &oauthTokenExchangeTokenGeneratorMappingsResource{}
}

// oauthTokenExchangeTokenGeneratorMappingsResource is the resource implementation.
type oauthTokenExchangeTokenGeneratorMappingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthTokenExchangeTokenGeneratorMappingsResourceModel struct {
	AttributeSources                 types.Set    `tfsdk:"attribute_sources"`
	AttributeContractFulfillment     types.Map    `tfsdk:"attribute_contract_fulfillment"`
	IssuanceCriteria                 types.Object `tfsdk:"issuance_criteria"`
	Id                               types.String `tfsdk:"id"`
	SourceId                         types.String `tfsdk:"source_id"`
	TargetId                         types.String `tfsdk:"target_id"`
	LicenseConnectionGroupAssignment types.String `tfsdk:"license_connection_group_assignment"`
}

// GetSchema defines the schema for the resource.
func (r *oauthTokenExchangeTokenGeneratorMappingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The id of the Token Exchange Processor policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The persistent, unique ID for the Token Exchange Processor policy. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[a-zA-Z0-9_]{1,32}$"),
						"The Token Exchange Processor policy ID must be less than 33 characters, contain no spaces, and be alphanumeric.",
					),
				},
			},
			"attribute_sources": schema.SetNestedAttribute{
				Description: "A list of configured data stores to look up attributes from.",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{},
				},
				// Validators: []validator.String{
				// 	stringvalidator.OneOf([]string{"LDAP", "PING_ONE_LDAP_GATEWAY", "JDBC", "CUSTOM"}...),
				// },
			},
			"attribute_contract_fulfillment": schema.MapNestedAttribute{
				Description: " A list of mappings from attribute names to their fulfillment values.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source": schema.SingleNestedAttribute{
							Description: "Source containing key that is meant to reference a source from which an attribute can be retrieved.",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The source type of this key.",
									Computed:    true,
									Optional:    true,
								},
								"id": schema.StringAttribute{
									Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
									Required:    true,
								},
							},
						},
						"value": schema.StringAttribute{
							Description: "The value for this attribute.",
							Required:    true,
						},
					},
				},
			},
			"issuance_criteria": schema.SingleNestedAttribute{
				Description: "The attribute update policy for authentication sources.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"conditional_criteria": schema.SetNestedAttribute{
						Description: "The source type of this key.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source": schema.SingleNestedAttribute{
									Description: "Source containing key that is meant to reference a source from which an attribute can be retrieved.",
									Required:    true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The source type of this key.",
											Computed:    true,
											Optional:    true,
										},
										"id": schema.StringAttribute{
											Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
											Required:    true,
										},
									},
								},
								"attribute_name": schema.StringAttribute{
									Description: " The name of the attribute to use in this issuance criterion.",
									Computed:    true,
									Optional:    true,
								},
								"condition": schema.StringAttribute{
									Description: "The condition that will be applied to the source attribute's value and the expected value.",
									Computed:    true,
									Optional:    true,
								},
								"value": schema.StringAttribute{
									Description: "The expected value of this issuance criterion.",
									Computed:    true,
									Optional:    true,
								},
								"error_result": schema.StringAttribute{
									Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs..",
									Computed:    true,
									Optional:    true,
								},
							},
						},
					},
					"expression_criteria": schema.SetNestedAttribute{
						Description: "A list of expression issuance criteria where the OGNL expressions must evaluate to true in order for the transaction to continue.",
						Computed:    true,
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"expression": schema.StringAttribute{
									Description: "The OGNL expression to evaluate.",
									Computed:    true,
									Optional:    true,
								},
								"error_result": schema.StringAttribute{
									Description: " The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
									Computed:    true,
									Optional:    true,
								},
							},
						},
					},
				},
			},
			"source_id": schema.StringAttribute{
				Description: "The id of the Token Exchange Processor policy.",
				Required:    true,
			},
			"target_id": schema.StringAttribute{
				Description: "The id of the Token Generator",
				Required:    true,
			},
			"license_connection_group_assignment": schema.StringAttribute{
				Description: "The license connection group",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func addOptionalOauthTokenExchangeTokenGeneratorMappingsFields(ctx context.Context, addRequest *client.ProcessorPolicyToGeneratorMapping, plan oauthTokenExchangeTokenGeneratorMappingsResourceModel) error {

	if internaltypes.IsDefined(plan.AttributeSources) {
		// if attributeSource.type is jbdc, ldap
		// call their respctive model  and add it to the request.
		var slice []correct_type
		//you may need to build the slice using a client method here, otherwise use a primitive type if applicable
		addRequest.AttributeSources = slice
	}

	if internaltypes.IsDefined(plan.AttributeContractFulfillment) {
		addRequest.AttributeContractFulfillment = client.NewAttributeFulfillmentValueWithDefaults()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContractFulfillment, false)), addRequest.AttributeContractFulfillment)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.IssuanceCriteria) {
		addRequest.IssuanceCriteria = client.NewIssuanceCriteria()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.IssuanceCriteria, false)), addRequest.IssuanceCriteria)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.SourceId) {
		addRequest.SourceId = plan.SourceId.ValueString()
	}

	if internaltypes.IsDefined(plan.TargetId) {
		addRequest.TargetId = plan.TargetId.ValueString()
	}

	if internaltypes.IsDefined(plan.LicenseConnectionGroupAssignment) {
		addRequest.LicenseConnectionGroupAssignment = plan.LicenseConnectionGroupAssignment.ValueStringPointer()
	}

	return nil

}

// Metadata returns the resource type name.
func (r *oauthTokenExchangeTokenGeneratorMappingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_token_exchange_token_generator_mappings"
}

func (r *oauthTokenExchangeTokenGeneratorMappingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthTokenExchangeTokenGeneratorMappingsResponse(ctx context.Context, r *client.ProcessorPolicyToGeneratorMapping, state *oauthTokenExchangeTokenGeneratorMappingsResourceModel) {
	// state.AttributeSources = internaltypes.GetCorrectMethodFromInternalTypesForThis(r.AttributeSources)
	// state.AttributeContractFulfillment = You will need to figure out what needs to go into the object(r.AttributeContractFulfillment)
	// state.IssuanceCriteria = (r.IssuanceCriteria)
	// state.Id = internaltypes.StringTypeOrNil(r.Id)
	// state.SourceId = internaltypes.StringTypeOrNil(r.SourceId)
	// state.TargetId = internaltypes.StringTypeOrNil(r.TargetId)
	// state.LicenseConnectionGroupAssignment = internaltypes.StringTypeOrNil(r.LicenseConnectionGroupAssignment)
}

func (r *oauthTokenExchangeTokenGeneratorMappingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthTokenExchangeTokenGeneratorMappingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthTokenExchangeTokenGeneratorMappings := client.NewProcessorPolicyToGeneratorMappingsWithDefaults()
	err := addOptionalOauthTokenExchangeTokenGeneratorMappingsFields(ctx, createOauthTokenExchangeTokenGeneratorMappings, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthTokenExchangeTokenGeneratorMappings", err.Error())
		return
	}
	requestJson, err := createOauthTokenExchangeTokenGeneratorMappings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateOauthTokenExchangeTokenGeneratorMappings := r.apiClient.OauthTokenExchangeProcessorApi.CreateOauthTokenExchangeProcessorPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthTokenExchangeTokenGeneratorMappings = apiCreateOauthTokenExchangeTokenGeneratorMappings.Body(*createOauthTokenExchangeTokenGeneratorMappings)
	oauthTokenExchangeTokenGeneratorMappingsResponse, httpResp, err := r.apiClient.OauthTokenExchangeProcessorApi.CreateOauthTokenExchangeProcessorPolicyExecute(apiCreateOauthTokenExchangeTokenGeneratorMappings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OauthTokenExchangeTokenGeneratorMappings", err, httpResp)
		return
	}
	responseJson, err := oauthTokenExchangeTokenGeneratorMappingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthTokenExchangeTokenGeneratorMappingsResourceModel

	readOauthTokenExchangeTokenGeneratorMappingsResponse(ctx, oauthTokenExchangeTokenGeneratorMappingsResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthTokenExchangeTokenGeneratorMappingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthTokenExchangeTokenGeneratorMappingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthTokenExchangeTokenGeneratorMappings, httpResp, err := r.apiClient.OauthTokenExchangeProcessorApi.GetOauthTokenExchangeProcessorPolicyById(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OauthTokenExchangeTokenGeneratorMappings", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  OauthTokenExchangeTokenGeneratorMappings", err, httpResp)
		}
	}
	// Log response JSON
	responseJson, err := apiReadOauthTokenExchangeTokenGeneratorMappings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthTokenExchangeTokenGeneratorMappingsResponse(ctx, apiReadOauthTokenExchangeTokenGeneratorMappings, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthTokenExchangeTokenGeneratorMappingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan oauthTokenExchangeTokenGeneratorMappingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthTokenExchangeTokenGeneratorMappingsResourceModel
	req.State.Get(ctx, &state)
	updateOauthTokenExchangeTokenGeneratorMappings := r.apiClient.OauthTokenExchangeProcessorApi.UpdateOauthTokenExchangeProcessorPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	createUpdateRequest := client.NewProcessorPolicyToGeneratorMappingWithDefaults()
	err := addOptionalOauthTokenExchangeTokenGeneratorMappingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthTokenExchangeTokenGeneratorMappings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateOauthTokenExchangeTokenGeneratorMappings = updateOauthTokenExchangeTokenGeneratorMappings.Body(*createUpdateRequest)
	updateOauthTokenExchangeTokenGeneratorMappingsResponse, httpResp, err := r.apiClient.OauthTokenExchangeProcessorApi.UpdateOauthTokenExchangeProcessorPolicyExecute(updateOauthTokenExchangeTokenGeneratorMappings)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OauthTokenExchangeTokenGeneratorMappings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateOauthTokenExchangeTokenGeneratorMappingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readOauthTokenExchangeTokenGeneratorMappingsResponse(ctx, updateOauthTokenExchangeTokenGeneratorMappingsResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthTokenExchangeTokenGeneratorMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthTokenExchangeTokenGeneratorMappingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.OauthTokenExchangeProcessorApi.DeleteOauthTokenExchangeProcessorPolicyy(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting OauthTokenExchangeTokenGeneratorMappings", err, httpResp)
	}

}
func (r *oauthTokenExchangeTokenGeneratorMappingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
