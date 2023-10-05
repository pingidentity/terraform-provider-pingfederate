package administrativeaccount

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &administrativeAccountDataSource{}
	_ datasource.DataSourceWithConfigure = &administrativeAccountDataSource{}
)

// Create a Administrative Account data source
func NewAdministrativeAccountDataSource() datasource.DataSource {
	return &administrativeAccountDataSource{}
}

// administrativeAccountDataSource is the datasource implementation.
type administrativeAccountDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type administrativeAccountDataSourceModel struct {
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

// GetSchema defines the schema for the datasource.
func (r *administrativeAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Administrative Account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed attribute tied to the username property of this resource.",
				Optional:    false,
				Computed:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Indicates whether the account is active or not.",
				Optional:    true,
				Computed:    true,
			},
			"auditor": schema.BoolAttribute{
				Description: "Indicates whether the account belongs to an Auditor. An Auditor has View-only permissions for all administrative functions. An Auditor cannot have any administrative roles.",
				Optional:    true,
				Computed:    true,
			},
			"department": schema.StringAttribute{
				Description: "The Department name of account user.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the account.",
				Optional:    true,
			},
			"email_address": schema.StringAttribute{
				Description: "Email address associated with the account.",
				Optional:    true,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for the Account. This field is only applicable during a POST operation.",
				Required:    true,
				Sensitive:   true,
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
	config.AddCommonDataSourceSchema(&schemaDef)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *administrativeAccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administrative_account"
}

// Configure adds the provider configured client to the data source.
func (r *administrativeAccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a AdministrativeAccountResponse object into the model struct
func readAdministrativeAccountResponseDataSource(ctx context.Context, r *client.AdministrativeAccount, state *administrativeAccountDataSourceModel, expectedValues *administrativeAccountDataSourceModel, passwordPlan basetypes.StringValue) {
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

// Read resource information
func (r *administrativeAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state administrativeAccountDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAdministrativeAccount, httpResp, err := r.apiClient.AdministrativeAccountsAPI.GetAccount(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Username.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Administrative Account", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, responseErr := apiReadAdministrativeAccount.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	} else {
		diags.AddError("There was an issue retrieving the response of an Administrative Account: %s", responseErr.Error())
	}

	// Read the response into the state
	readAdministrativeAccountResponseDataSource(ctx, apiReadAdministrativeAccount, &state, &state, state.Password)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
