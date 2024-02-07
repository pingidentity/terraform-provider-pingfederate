package administrativeaccount

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &administrativeAccountDataSource{}
	_ datasource.DataSourceWithConfigure = &administrativeAccountDataSource{}
)

// Create a Administrative Account data source
func AdministrativeAccountDataSource() datasource.DataSource {
	return &administrativeAccountDataSource{}
}

// administrativeAccountDataSource is the datasource implementation.
type administrativeAccountDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *administrativeAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes an administrative account.",
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Description: "Indicates whether the account is active or not.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auditor": schema.BoolAttribute{
				Description: "Indicates whether the account belongs to an Auditor. An Auditor has View-only permissions for all administrative functions. An Auditor cannot have any administrative roles.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"department": schema.StringAttribute{
				Description: "The Department name of the account user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"email_address": schema.StringAttribute{
				Description: "Email address associated with the account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encrypted_password": schema.StringAttribute{
				Description: "For GET requests, this field contains the encrypted account password. For POST and PUT requests, if you wish to re-use the password from an API response to this endpoint, this field should be passed back unchanged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for the Account. This field is only applicable during a POST operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"phone_number": schema.StringAttribute{
				Description: "Phone number associated with the account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"roles": schema.SetAttribute{
				Description: "Roles available for an administrator. USER_ADMINISTRATOR - Can create, deactivate or delete accounts and reset passwords. Additionally, install replacement license keys. CRYPTO_ADMINISTRATOR - Can manage local keys and certificates. ADMINISTRATOR - Can configure partner connections and most system settings (except the management of native accounts and the handling of local keys and certificates. EXPRESSION_ADMINISTRATOR - Can add and update OGNL expressions.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"username": schema.StringAttribute{
				Description: "Username for the Administrative Account.",
				Required:    true,
			},
		},
	}
	id.ToDataSourceSchema(&schema)
	resp.Schema = schema
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

// Read resource information
func (r *administrativeAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state administrativeAccountModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadAdministrativeAccount, httpResp, err := r.apiClient.AdministrativeAccountsAPI.GetAccount(config.AuthContext(ctx, r.providerConfig), state.Username.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the administrative account", err, httpResp)
		return
	}

	// Read the response into the state
	readAdministrativeAccountResponse(ctx, apiReadAdministrativeAccount, &state, nil)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
