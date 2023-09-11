package oauth

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	_ resource.Resource                = &oauthAccessTokenManagersResource{}
	_ resource.ResourceWithConfigure   = &oauthAccessTokenManagersResource{}
	_ resource.ResourceWithImportState = &oauthAccessTokenManagersResource{}
)

// OauthAccessTokenManagersResource is a helper function to simplify the provider implementation.
func OauthAccessTokenManagersResource() resource.Resource {
	return &oauthAccessTokenManagersResource{}
}

// oauthAccessTokenManagersResource is the resource implementation.
type oauthAccessTokenManagersResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAccessTokenManagersResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	PluginDescriptorRef       types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef                 types.Object `tfsdk:"parent_ref"`
	Configuration             types.Object `tfsdk:"configuration"`
	AttributeContract         types.Object `tfsdk:"attribute_contract"`
	SelectionSettings         types.Object `tfsdk:"selection_settings"`
	AccessControlSettings     types.Object `tfsdk:"access_control_setting"`
	SessionValidationSettings types.Object `tfsdk:"session_validation_settings"`
	SequenceNumber            types.Int64  `tfsdk:"sequence_number"`
}

// GetSchema defines the schema for the resource.
func (r *oauthAccessTokenManagersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages Oauth Access Token Managers",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[a-zA-Z0-9_]{1,32}$"),
						"The plugin ID must be less than 33 characters, contain no spaces, and be alphanumeric.",
					),
				},
			},
		},
	}
	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"FIX_ME"})
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

func addOptionalOauthAccessTokenManagersFields(ctx context.Context, addRequest *client.AccessTokenManagers, plan oauthAccessTokenManagersResourceModel) error {

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = client.NewAttributeContract()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, false)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Id) {
		addRequest.Id = plan.Id.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.PluginDescriptorRef) {
		addRequest.PluginDescriptorRef = client.NewPluginDescriptorRef()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), addRequest.PluginDescriptorRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.ParentRef) {
		addRequest.ParentRef = client.NewParentRef()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ParentRef, false)), addRequest.ParentRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.Configuration) {
		addRequest.Configuration = client.NewConfiguration()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, false)), addRequest.Configuration)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.SelectionSettings) {
		addRequest.SelectionSettings = client.NewSelectionSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.SelectionSettings, false)), addRequest.SelectionSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AccessControlSettings) {
		addRequest.AccessControlSettings = client.NewAccessControlSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AccessControlSettings, false)), addRequest.AccessControlSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.SessionValidationSettings) {
		addRequest.SessionValidationSettings = client.NewSessionValidationSettings()
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.SessionValidationSettings, false)), addRequest.SessionValidationSettings)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.SequenceNumber) {
		addRequest.SequenceNumber = plan.SequenceNumber.ValueInt64Pointer()
	}

	return nil

}

// Metadata returns the resource type name.
func (r *oauthAccessTokenManagersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_access_token_managers"
}

func (r *oauthAccessTokenManagersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthAccessTokenManagersResponse(ctx context.Context, r *client.AccessTokenManagers, state *oauthAccessTokenManagersResourceModel) {
	// state.AttributeContract = (r.AttributeContract)
	// state.Id = internaltypes.StringTypeOrNil(r.Id)
	// state.Name = internaltypes.StringTypeOrNil(r.Name)
	// state.PluginDescriptorRef = (r.PluginDescriptorRef)
	// state.ParentRef = (r.ParentRef)
	// state.Configuration = (r.Configuration)
	// state.SelectionSettings = (r.SelectionSettings)
	// state.AccessControlSettings = (r.AccessControlSettings)
	// state.SessionValidationSettings = (r.SessionValidationSettings)
	// state.SequenceNumber = types.Int64Value(r.SequenceNumber)
}

func (r *oauthAccessTokenManagersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAccessTokenManagersResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAccessTokenManagers := client.NewAccessTokenManager()
	err := addOptionalOauthAccessTokenManagersFields(ctx, createOauthAccessTokenManagers, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAccessTokenManagers", err.Error())
		return
	}
	requestJson, err := createOauthAccessTokenManagers.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateOauthAccessTokenManagers := r.apiClient.OauthApi.AddAccessTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAccessTokenManagers = apiCreateOauthAccessTokenManagers.Body(*createOauthAccessTokenManagers)
	oauthAccessTokenManagersResponse, httpResp, err := r.apiClient.OauthApi.AddAccessTokenManagerExecute(apiCreateOauthAccessTokenManagers)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OauthAccessTokenManagers", err, httpResp)
		return
	}
	responseJson, err := oauthAccessTokenManagersResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthAccessTokenManagersResourceModel

	readOauthAccessTokenManagersResponse(ctx, oauthAccessTokenManagersResponse, &state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *oauthAccessTokenManagersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauthAccessTokenManagersResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAccessTokenManagers, httpResp, err := r.apiClient.OauthApi.GetAccessTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.VALUE.ValueString()).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the OauthAccessTokenManagers", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the  OauthAccessTokenManagers", err, httpResp)
		}
	}
	// Log response JSON
	responseJson, err := apiReadOauthAccessTokenManagers.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthAccessTokenManagersResponse(ctx, apiReadOauthAccessTokenManagers, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAccessTokenManagersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan oauthAccessTokenManagersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthAccessTokenManagersResourceModel
	req.State.Get(ctx, &state)
	updateOauthAccessTokenManagers := r.apiClient.OauthApi.UpdateAccessTokenManager(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.VALUE.ValueString())
	createUpdateRequest := client.NewAccessTokenManager()
	err := addOptionalOauthAccessTokenManagersFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAccessTokenManagers", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateOauthAccessTokenManagers = updateOauthAccessTokenManagers.Body(*createUpdateRequest)
	updateOauthAccessTokenManagersResponse, httpResp, err := r.apiClient.OauthApi.UpdateAccessTokenManagerExecute(updateOauthAccessTokenManagers)
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OauthAccessTokenManagers", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateOauthAccessTokenManagersResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readOauthAccessTokenManagersResponse(ctx, updateOauthAccessTokenManagersResponse, &state)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// This config object is edit-only, so Terraform can't delete it.
func (r *oauthAccessTokenManagersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *oauthAccessTokenManagersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
