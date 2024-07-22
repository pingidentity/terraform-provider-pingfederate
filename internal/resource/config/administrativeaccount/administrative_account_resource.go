package administrativeaccount

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
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

// GetSchema defines the schema for the resource.
func (r *administrativeAccountsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages an administrative account.",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Description: "Indicates whether the account is active or not.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"auditor": schema.BoolAttribute{
				Description: "Indicates whether the account belongs to an Auditor. An Auditor has View-only permissions for all administrative functions. An Auditor cannot have any administrative roles.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"department": schema.StringAttribute{
				Description: "The Department name of the account user.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the account.",
				Optional:    true,
			},
			"email_address": schema.StringAttribute{
				Description: "Email address associated with the account.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for the Account. This field is only applicable during account creation.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"encrypted_password": schema.StringAttribute{
				Description: "Read-only attribute. This field holds the value returned from PingFederate and used for updating an existing Administrative Account.",
				Computed:    true,
				Optional:    false,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"phone_number": schema.StringAttribute{
				Description: "Phone number associated with the account.",
				Optional:    true,
			},
			"roles": schema.SetAttribute{
				Description: "Roles available for an administrator. USER_ADMINISTRATOR - Can create, deactivate or delete accounts and reset passwords. Additionally, install replacement license keys. CRYPTO_ADMINISTRATOR - Can manage local keys and certificates. ADMINISTRATOR - Can configure partner connections and most system settings (except the management of native accounts and the handling of local keys and certificates. EXPRESSION_ADMINISTRATOR - Can add and update OGNL expressions.",
				Required:    true,
				ElementType: types.StringType,
			},
			"username": schema.StringAttribute{
				Description: "Username for the Administrative Account.",
				Required:    true,
			},
		},
	}

	id.ToSchema(&schema)
	resp.Schema = schema
}

func addOptionalAdministrativeAccountFields(ctx context.Context, addRequest *client.AdministrativeAccount, plan administrativeAccountModel, isCreate bool) error {
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

	if isCreate {
		addRequest.Password = plan.Password.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.EncryptedPassword) {
		addRequest.EncryptedPassword = plan.EncryptedPassword.ValueStringPointer()
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

func (r *administrativeAccountsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan administrativeAccountModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createAdministrativeAccount := client.NewAdministrativeAccount(plan.Username.ValueString())
	err := addOptionalAdministrativeAccountFields(ctx, createAdministrativeAccount, plan, true)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to the add request for the administrative account", err.Error())
		return
	}

	apiCreateAdministrativeAccount := r.apiClient.AdministrativeAccountsAPI.AddAccount(config.AuthContext(ctx, r.providerConfig))
	apiCreateAdministrativeAccount = apiCreateAdministrativeAccount.Body(*createAdministrativeAccount)
	administrativeAccountResponse, httpResp, err := r.apiClient.AdministrativeAccountsAPI.AddAccountExecute(apiCreateAdministrativeAccount)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the administrative account", err, httpResp)
		return
	}

	// Read the response into the state
	var state administrativeAccountModel

	readAdministrativeAccountResponse(ctx, administrativeAccountResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *administrativeAccountsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state administrativeAccountModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadAdministrativeAccount, httpResp, err := r.apiClient.AdministrativeAccountsAPI.GetAccount(config.AuthContext(ctx, r.providerConfig), state.Username.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Administrative Account", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Administrative Account", err, httpResp)
		}
		return
	}

	// Read the response into the state
	readAdministrativeAccountResponse(ctx, apiReadAdministrativeAccount, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *administrativeAccountsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan administrativeAccountModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state administrativeAccountModel
	req.State.Get(ctx, &state)
	updateAdministrativeAccount := r.apiClient.AdministrativeAccountsAPI.UpdateAccount(config.AuthContext(ctx, r.providerConfig), plan.Username.ValueString())
	createUpdateRequest := client.NewAdministrativeAccount(plan.Username.ValueString())
	err := addOptionalAdministrativeAccountFields(ctx, createUpdateRequest, plan, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to the add request for the administrative account", err.Error())
		return
	}

	updateAdministrativeAccount = updateAdministrativeAccount.Body(*createUpdateRequest)
	updateAdministrativeAccountResponse, httpResp, err := r.apiClient.AdministrativeAccountsAPI.UpdateAccountExecute(updateAdministrativeAccount)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the administrative account", err, httpResp)
		return
	}

	// Read the response
	readAdministrativeAccountResponse(ctx, updateAdministrativeAccountResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *administrativeAccountsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state administrativeAccountModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := r.apiClient.AdministrativeAccountsAPI.DeleteAccount(config.AuthContext(ctx, r.providerConfig), state.Username.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting an administrative account", err, httpResp)
	}
}

func (r *administrativeAccountsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("username"), req, resp)
}
