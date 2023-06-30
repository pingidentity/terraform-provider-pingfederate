package administrativeaccounts

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &administrativeAccountsResource{}
	_ resource.ResourceWithConfigure   = &administrativeAccountsResource{}
	_ resource.ResourceWithImportState = &administrativeAccountsResource{}
)

// AdministrativeAccountResource is a helper function to simplify the provider implementation.
func AdministrativeAccountResource() resource.Resource {
	return &administrativeAccountsResource{}
}

// administrativeAccountsResource is the resource implementation.
type administrativeAccountsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type administrativeAccountResourceModel struct {
	Active       types.Bool   `tfsdk:"active"`
	Auditor      types.Bool   `tfsdk:"auditor"`
	Department   types.String `tfsdk:"department"`
	Description  types.String `tfsdk:"description"`
	EmailAddress types.String `tfsdk:"email_address"`
	Id           types.String `tfsdk:"id"`
	Password     types.String `tfsdk:"password"`
	PhoneNumber  types.String `tfsdk:"phone_number"`
	Roles        types.Set    `tfsdk:"roles"`
	Username     types.String `tfsdk:"username"`
}

// GetSchema defines the schema for the resource.
func (r *administrativeAccountsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	administrativeAccountResourceSchema(ctx, req, resp, false)
}

func administrativeAccountResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a AdministrativeAccount.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed attribute tied to the username property of this resource.",
				Optional:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auditor": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"department": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"phone_number": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"roles": schema.SetAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
					setplanmodifier.RequiresReplace(),
				},
				ElementType: types.StringType,
			},
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	// Set attribtues in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"username", "password", "roles"})
	}
	resp.Schema = schema
}
func addOptionalAdministrativeAccountFields(ctx context.Context, addRequest *client.AdministrativeAccount, plan administrativeAccountResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsDefined(plan.Active) {
		addRequest.Active = plan.Active.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Auditor) {
		addRequest.Auditor = plan.Auditor.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Department) {
		addRequest.Department = plan.Department.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EmailAddress) {
		addRequest.EmailAddress = plan.EmailAddress.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.PhoneNumber) {
		addRequest.PhoneNumber = plan.PhoneNumber.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Roles) {
		var slice []string
		plan.Roles.ElementsAs(ctx, &slice, false)
		addRequest.Roles = slice
	}
	return nil
}

// Metadata returns the resource type name.
func (r *administrativeAccountsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administrative_accounts"
}

func (r *administrativeAccountsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readAdministrativeAccountResponse(ctx context.Context, r *client.AdministrativeAccount, state *administrativeAccountResourceModel, expectedValues *administrativeAccountResourceModel, passwordPlan basetypes.StringValue) {
	state.Id = types.StringValue(r.Username)
	state.Username = types.StringValue(r.Username)
	state.Password = types.StringValue(passwordPlan.ValueString())
	state.Active = types.BoolValue(*r.Active)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Auditor = types.BoolValue(*r.Auditor)
	state.PhoneNumber = internaltypes.StringTypeOrNil(r.PhoneNumber, false)
	state.EmailAddress = internaltypes.StringTypeOrNil(r.EmailAddress, false)
	state.Department = internaltypes.StringTypeOrNil(r.Department, false)
	state.Roles = internaltypes.GetStringSet(r.Roles)
}

func (r *administrativeAccountsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan administrativeAccountResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAdministrativeAccount := client.NewAdministrativeAccount(plan.Username.ValueString())
	err := addOptionalAdministrativeAccountFields(ctx, createAdministrativeAccount, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AdministrativeAccount", err.Error())
		return
	}
	requestJson, err := createAdministrativeAccount.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateAdministrativeAccount := r.apiClient.AdministrativeAccountsApi.AddAccount(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAdministrativeAccount = apiCreateAdministrativeAccount.Body(*createAdministrativeAccount)
	administrativeAccountResponse, httpResp, err := r.apiClient.AdministrativeAccountsApi.AddAccountExecute(apiCreateAdministrativeAccount)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the AdministrativeAccount", err, httpResp)
		return
	}
	responseJson, err := administrativeAccountResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state administrativeAccountResourceModel

	readAdministrativeAccountResponse(ctx, administrativeAccountResponse, &state, &plan, plan.Password)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *administrativeAccountsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAdministrativeAccount(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAdministrativeAccount(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state administrativeAccountResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAdministrativeAccount, httpResp, err := apiClient.AdministrativeAccountsApi.GetAccount(config.ProviderBasicAuthContext(ctx, providerConfig), state.Username.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a AdministrativeAccount", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadAdministrativeAccount.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAdministrativeAccountResponse(ctx, apiReadAdministrativeAccount, &state, &state, state.Password)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *administrativeAccountsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAdministrativeAccount(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAdministrativeAccount(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan administrativeAccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state administrativeAccountResourceModel
	req.State.Get(ctx, &state)
	updateAdministrativeAccount := apiClient.AdministrativeAccountsApi.UpdateAccount(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Username.ValueString())
	createUpdateRequest := client.NewAdministrativeAccount(plan.Username.ValueString())
	err := addOptionalAdministrativeAccountFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for AdministrativeAccount", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateAdministrativeAccount = updateAdministrativeAccount.Body(*createUpdateRequest)
	updateAdministrativeAccountResponse, httpResp, err := apiClient.AdministrativeAccountsApi.UpdateAccountExecute(updateAdministrativeAccount)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating AdministrativeAccount", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateAdministrativeAccountResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readAdministrativeAccountResponse(ctx, updateAdministrativeAccountResponse, &state, &plan, state.Password)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *administrativeAccountsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteAdministrativeAccount(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteAdministrativeAccount(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state administrativeAccountResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.AdministrativeAccountsApi.DeleteAccount(config.ProviderBasicAuthContext(ctx, providerConfig), state.Username.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a AdministrativeAccount", err, httpResp)
		return
	}

}

func (r *administrativeAccountsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAdministrativeAccountLocation(ctx, req, resp)
}
func importAdministrativeAccountLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import username and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("username"), req, resp)
}
