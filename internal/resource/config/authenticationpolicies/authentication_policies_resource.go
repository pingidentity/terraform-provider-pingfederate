package authenticationpolicies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/authenticationpolicytreenode"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authenticationPoliciesResource{}
	_ resource.ResourceWithConfigure   = &authenticationPoliciesResource{}
	_ resource.ResourceWithImportState = &authenticationPoliciesResource{}

	authnSelectionTreesAttrTypes = map[string]attr.Type{
		"id":                                 types.StringType,
		"name":                               types.StringType,
		"description":                        types.StringType,
		"authentication_api_application_ref": types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"enabled":                            types.BoolType,
		"root_node":                          types.ObjectType{AttrTypes: authenticationpolicytreenode.GetRootNodeAttrTypes()},
		"handle_failures_locally":            types.BoolType,
	}

	defaultAuthenticationSourcesAttrTypes = map[string]attr.Type{
		"type":       types.StringType,
		"source_ref": types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
	}

	defaultAuthenticationSourcesEmptyList, _ = types.ListValue(types.ObjectType{AttrTypes: defaultAuthenticationSourcesAttrTypes}, []attr.Value{})

	emptyStringList, _ = types.ListValue(types.StringType, []attr.Value{})
)

type authenticationPoliciesModel struct {
	Id                           types.String `tfsdk:"id"`
	AuthnSelectionTrees          types.List   `tfsdk:"authn_selection_trees"`
	DefaultAuthenticationSources types.List   `tfsdk:"default_authentication_sources"`
	FailIfNoSelection            types.Bool   `tfsdk:"fail_if_no_selection"`
	TrackedHttpParameters        types.List   `tfsdk:"tracked_http_parameters"`
}

// authenticationPoliciesResource is a helper function to simplify the provider implementation.
func AuthenticationPoliciesResource() resource.Resource {
	return &authenticationPoliciesResource{}
}

// authenticationPoliciesResource is the resource implementation.
type authenticationPoliciesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *authenticationPoliciesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Authentication Policies",
		Attributes: map[string]schema.Attribute{
			"authn_selection_trees": schema.ListNestedAttribute{
				Optional:    true,
				Description: "The list of authentication policy trees.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The authentication policy tree id. ID is unique.",
						},
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "The authentication policy name. Name is unique.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "A description for the authentication policy.",
						},
						"authentication_api_application_ref": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "Authentication API Application Id to be used in this policy branch. If the value is not specified, no Authentication API Application will be used.",
							Attributes:  resourcelink.ToSchemaNoLocation(),
						},
						"enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(true),
							Description: "Whether or not this authentication policy tree is enabled. Default is true.",
						},
						"root_node": authenticationpolicytreenode.ToSchema("A node inside the authentication policy tree."),
						"handle_failures_locally": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
							Description: "If a policy ends in failure keep the user local.",
						},
					},
				},
			},
			"default_authentication_sources": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(defaultAuthenticationSourcesEmptyList),
				Description: "The default authentication sources.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of this authentication source.",
							Validators: []validator.String{
								stringvalidator.OneOf("IDP_ADAPTER", "IDP_CONNECTION"),
							},
						},
						"source_ref": schema.SingleNestedAttribute{
							Required:    true,
							Attributes:  resourcelink.ToSchemaNoLocation(),
							Description: "A reference to the authentication source.",
						},
					},
				},
			},
			"fail_if_no_selection": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Fail if policy finds no authentication source.",
			},
			"tracked_http_parameters": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(emptyStringList),
				Description: "The HTTP request parameters to track and make available to authentication sources, selectors, and contract mappings throughout the authentication policy.",
				ElementType: types.StringType,
			},
		},
	}
	id.ToSchema(&schema)
	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *authenticationPoliciesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authentication_policies"
}

func (r *authenticationPoliciesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAuthenticationPoliciesResponse(ctx context.Context, r *client.AuthenticationPolicy, state *authenticationPoliciesModel, existingId *string) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	if existingId != nil {
		state.Id = types.StringValue(*existingId)
	} else {
		state.Id = id.GenerateUUIDToState(existingId)
	}

	state.FailIfNoSelection = types.BoolPointerValue(r.FailIfNoSelection)

	state.DefaultAuthenticationSources, respDiags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: defaultAuthenticationSourcesAttrTypes}, r.DefaultAuthenticationSources)
	diags.Append(respDiags...)

	if r.AuthnSelectionTrees != nil {
		authnSelectionTreesToState := []attr.Value{}

		for _, authnSelectionTree := range r.AuthnSelectionTrees {

			authenticationApiApplicationRef, respDiags := resourcelink.ToStateNoLocation(authnSelectionTree.AuthenticationApiApplicationRef)
			diags.Append(respDiags...)

			rootNode, respDiags := authenticationpolicytreenode.ToState(ctx, authnSelectionTree.RootNode)
			diags.Append(respDiags...)

			authnSelectionTreeAttrValues := map[string]attr.Value{
				"id":                                 types.StringPointerValue(authnSelectionTree.Id),
				"name":                               types.StringPointerValue(authnSelectionTree.Name),
				"description":                        types.StringPointerValue(authnSelectionTree.Description),
				"enabled":                            types.BoolPointerValue(authnSelectionTree.Enabled),
				"handle_failures_locally":            types.BoolPointerValue(authnSelectionTree.HandleFailuresLocally),
				"authentication_api_application_ref": authenticationApiApplicationRef,
				"root_node":                          rootNode,
			}

			authenticationSelectionTreeToState, respDiags := types.ObjectValue(authnSelectionTreesAttrTypes, authnSelectionTreeAttrValues)
			diags.Append(respDiags...)

			authnSelectionTreesToState = append(authnSelectionTreesToState, authenticationSelectionTreeToState)
		}

		state.AuthnSelectionTrees, respDiags = types.ListValue(types.ObjectType{AttrTypes: authnSelectionTreesAttrTypes}, authnSelectionTreesToState)
		diags.Append(respDiags...)
	} else {
		state.AuthnSelectionTrees = types.ListNull(types.ObjectType{AttrTypes: authnSelectionTreesAttrTypes})
	}

	state.TrackedHttpParameters = internaltypes.GetStringList(r.TrackedHttpParameters)

	return diags
}

func addOptionalAuthenticationPolicyFields(addRequest *client.AuthenticationPolicy, plan authenticationPoliciesModel) error {
	addRequest.FailIfNoSelection = plan.FailIfNoSelection.ValueBoolPointer()

	if internaltypes.IsDefined(plan.AuthnSelectionTrees) {
		addRequest.AuthnSelectionTrees = []client.AuthenticationPolicyTree{}
		for _, authnSelectionTree := range plan.AuthnSelectionTrees.Elements() {
			authnSelectionTreeObj, ok := authnSelectionTree.(types.Object)
			if !ok {
				return fmt.Errorf("authn_selection_trees must be a list of objects")
			}
			authenticationPolicyTree := client.AuthenticationPolicyTree{}
			authnSelectionTreeObjElements := authnSelectionTreeObj.Attributes()
			if id, ok := authnSelectionTreeObjElements["id"]; ok {
				authenticationPolicyTree.Id = id.(types.String).ValueStringPointer()
			}
			if name, ok := authnSelectionTreeObjElements["name"]; ok {
				authenticationPolicyTree.Name = name.(types.String).ValueStringPointer()
			}
			if description, ok := authnSelectionTreeObjElements["description"]; ok {
				authenticationPolicyTree.Description = description.(types.String).ValueStringPointer()
			}
			if authenticationApiApplicationRef, ok := authnSelectionTreeObjElements["authentication_api_application_ref"]; ok {
				authenticationApiApplicationRefObj, ok := authenticationApiApplicationRef.(types.Object)
				if !ok {
					return fmt.Errorf("authentication_api_application_ref must be an object")
				}
				authenticationApiApplicationRef, err := resourcelink.ClientStruct(authenticationApiApplicationRefObj)
				if err != nil {
					return err
				}
				authenticationPolicyTree.AuthenticationApiApplicationRef = authenticationApiApplicationRef
			}
			if enabled, ok := authnSelectionTreeObjElements["enabled"]; ok {
				authenticationPolicyTree.Enabled = enabled.(types.Bool).ValueBoolPointer()
			}
			if rootNode, ok := authnSelectionTreeObjElements["root_node"]; ok {
				rootNodeObj, err := authenticationpolicytreenode.ClientStruct(rootNode.(types.Object))
				if err != nil {
					return err
				}
				authenticationPolicyTree.RootNode = rootNodeObj
			}
			if handleFailuresLocally, ok := authnSelectionTreeObjElements["handle_failures_locally"]; ok {
				authenticationPolicyTree.HandleFailuresLocally = handleFailuresLocally.(types.Bool).ValueBoolPointer()
			}
			addRequest.AuthnSelectionTrees = append(addRequest.AuthnSelectionTrees, authenticationPolicyTree)
		}
	}

	if internaltypes.IsDefined(plan.DefaultAuthenticationSources) {
		addRequest.DefaultAuthenticationSources = []client.AuthenticationSource{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.DefaultAuthenticationSources, true)), &addRequest.DefaultAuthenticationSources)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.TrackedHttpParameters) {
		addRequest.TrackedHttpParameters = []string{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.TrackedHttpParameters, true)), &addRequest.TrackedHttpParameters)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *authenticationPoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state authenticationPoliciesModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newPolicy := client.NewAuthenticationPolicy()
	err := addOptionalAuthenticationPolicyFields(newPolicy, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Authentication Policies", err.Error())
		return
	}

	apiCreatePolicy := r.apiClient.AuthenticationPoliciesAPI.UpdateDefaultAuthenticationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreatePolicy = apiCreatePolicy.Body(*newPolicy)
	policyResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.UpdateDefaultAuthenticationPolicyExecute(apiCreatePolicy)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Authentication Policies", err, httpResp)
		return
	}

	diags = readAuthenticationPoliciesResponse(ctx, policyResponse, &state, nil)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *authenticationPoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticationPoliciesModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.GetDefaultAuthenticationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the Authentication Policies", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the Authentication Policies", err, httpResp)
		}
		return
	}

	// Read the response into the state
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = readAuthenticationPoliciesResponse(ctx, policyResponse, &state, id)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *authenticationPoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state authenticationPoliciesModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePolicyRequest := r.apiClient.AuthenticationPoliciesAPI.UpdateDefaultAuthenticationPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	updatedPolicies := client.NewAuthenticationPolicy()
	err := addOptionalAuthenticationPolicyFields(updatedPolicies, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for the Authentication Policies", err.Error())
		return
	}

	updatePolicyRequest = updatePolicyRequest.Body(*updatedPolicies)
	updateResponse, httpResp, err := r.apiClient.AuthenticationPoliciesAPI.UpdateDefaultAuthenticationPolicyExecute(updatePolicyRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Authentication Policies", err, httpResp)
		return
	}

	// Read the response
	id, diags := id.GetID(ctx, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	readResponseDiags := readAuthenticationPoliciesResponse(ctx, updateResponse, &state, id)
	resp.Diagnostics.Append(readResponseDiags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// This resource is a put operation, so we do not need to implement the Delete method
func (r *authenticationPoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *authenticationPoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
