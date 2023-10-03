package administrativeaccount

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
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
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
	schema := schema.Schema{
		Description: "Manages a AdministrativeAccount.",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Description: "Indicates whether the account is active or not.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auditor": schema.BoolAttribute{
				Description: "Indicates whether the account belongs to an Auditor. An Auditor has View-only permissions for all administrative functions. An Auditor cannot have any administrative roles.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"department": schema.StringAttribute{
				Description: "The Department name of account user.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the account.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email_address": schema.StringAttribute{
				Description: "Email address associated with the account.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password for the Account. This field is only applicable during a POST operation.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"phone_number": schema.StringAttribute{
				Description: "Phone number associated with the account.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"roles": schema.SetAttribute{
				Description: "Roles available for an administrator. USER_ADMINISTRATOR - Can create, deactivate or delete accounts and reset passwords. Additionally, install replacement license keys. CRYPTO_ADMINISTRATOR - Can manage local keys and certificates. ADMINISTRATOR - Can configure partner connections and most system settings (except the management of native accounts and the handling of local keys and certificates. EXPRESSION_ADMINISTRATOR - Can add and update OGNL expressions.",
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
					setplanmodifier.RequiresReplace(),
				},
				ElementType: types.StringType,
			},
			"username": schema.StringAttribute{
				Description: "Username for the Administrative Account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	config.AddCommonSchema(&schema)
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
	resp.TypeName = req.ProviderTypeName + "_administrative_account"
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
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Administrative Account", err.Error())
		return
	}

	_, requestErr := createAdministrativeAccount.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of an Administrative Account: %s", requestErr.Error())
	}

	apiCreateAdministrativeAccount := r.apiClient.AdministrativeAccountsAPI.AddAccount(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateAdministrativeAccount = apiCreateAdministrativeAccount.Body(*createAdministrativeAccount)
	administrativeAccountResponse, httpResp, err := r.apiClient.AdministrativeAccountsAPI.AddAccountExecute(apiCreateAdministrativeAccount)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Administrative Account", err, httpResp)
		return
	}

	_, responseErr := administrativeAccountResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an Administrative Account: %s", responseErr.Error())
	}

	// Read the response into the state
	var state administrativeAccountResourceModel

	readAdministrativeAccountResponse(ctx, administrativeAccountResponse, &state, &plan, plan.Password)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *administrativeAccountsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state administrativeAccountResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAdministrativeAccount, httpResp, err := r.apiClient.AdministrativeAccountsAPI.GetAccount(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Username.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting an Administrative Account", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting an Administrative Account", err, httpResp)
		}
		return
	}
	// Log response JSON
	_, responseErr := apiReadAdministrativeAccount.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an Administrative Account: %s", responseErr.Error())
	}

	// Read the response into the state
	readAdministrativeAccountResponse(ctx, apiReadAdministrativeAccount, &state, &state, state.Password)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *administrativeAccountsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan administrativeAccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state administrativeAccountResourceModel
	req.State.Get(ctx, &state)
	updateAdministrativeAccount := r.apiClient.AdministrativeAccountsAPI.UpdateAccount(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Username.ValueString())
	createUpdateRequest := client.NewAdministrativeAccount(plan.Username.ValueString())
	err := addOptionalAdministrativeAccountFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Administrative Account", err.Error())
		return
	}

	_, requestErr := createUpdateRequest.MarshalJSON()
	if requestErr != nil {
		diags.AddError("There was an issue retrieving the request of an Administrative Account: %s", requestErr.Error())
	}

	updateAdministrativeAccount = updateAdministrativeAccount.Body(*createUpdateRequest)
	updateAdministrativeAccountResponse, httpResp, err := r.apiClient.AdministrativeAccountsAPI.UpdateAccountExecute(updateAdministrativeAccount)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating Administrative Account", err, httpResp)
		return
	}
	// Log response JSON
	_, responseErr := updateAdministrativeAccountResponse.MarshalJSON()
	if responseErr != nil {
		diags.AddError("There was an issue retrieving the response of an Administrative Account: %s", responseErr.Error())
	}
	// Read the response
	readAdministrativeAccountResponse(ctx, updateAdministrativeAccountResponse, &state, &plan, state.Password)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *administrativeAccountsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state administrativeAccountResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AdministrativeAccountsAPI.DeleteAccount(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Username.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a Administrative Account", err, httpResp)
		return
	}
}

func (r *administrativeAccountsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("username"), req, resp)
}
