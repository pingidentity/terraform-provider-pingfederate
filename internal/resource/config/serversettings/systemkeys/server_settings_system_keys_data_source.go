package serversettingssystemkeys

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsSystemKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsSystemKeysDataSource{}

	systemKeyDataSourceAttrTypes = map[string]attr.Type{
		"creation_date":      types.StringType,
		"encrypted_key_data": types.StringType,
	}
)

type serverSettingsSystemKeysModel struct {
	Current  types.Object `tfsdk:"current"`
	Previous types.Object `tfsdk:"previous"`
	Pending  types.Object `tfsdk:"pending"`
}

// ServerSettingsSystemKeysDataSource is a helper function to simplify the provider implementation.
func ServerSettingsSystemKeysDataSource() datasource.DataSource {
	return &serverSettingsSystemKeysDataSource{}
}

// serverSettingsSystemKeysDataSource is the resource implementation.
type serverSettingsSystemKeysDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// GetSchema defines the schema for the datasource.
func (r *serverSettingsSystemKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Describes the server system keys.",
		Attributes: map[string]schema.Attribute{
			"current": schema.SingleNestedAttribute{
				Description: "Current SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"previous": schema.SingleNestedAttribute{
				Description: "Previous SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
			"pending": schema.SingleNestedAttribute{
				Description: "Pending SystemKeys Secrets that are used in cryptographic operations to generate and consume internal tokens.",
				Computed:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"creation_date": schema.StringAttribute{
						Description: "Creation time of the key.",
						Computed:    true,
						Optional:    false,
					},
					"encrypted_key_data": schema.StringAttribute{
						Description: "The system key encrypted.",
						Computed:    true,
						Optional:    false,
					},
				},
			},
		},
	}

	resp.Schema = schema
}

// Metadata returns the resource type name.
func (r *serverSettingsSystemKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_settings_system_keys"
}

func (r *serverSettingsSystemKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readServerSettingsSystemKeysDataSourceResponse(ctx context.Context, r *client.SystemKeys, state *serverSettingsSystemKeysModel) diag.Diagnostics {
	var diags diag.Diagnostics
	currentAttrs := r.GetCurrent()
	currentAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(currentAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(currentAttrs.GetEncryptedKeyData()),
	}
	currentAttrsObjVal, respDiags := types.ObjectValue(systemKeyDataSourceAttrTypes, currentAttrVals)
	diags = append(diags, respDiags...)

	previousAttrs := r.GetPrevious()
	previousAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(previousAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(previousAttrs.GetEncryptedKeyData()),
	}
	previousAttrsObjVal, respDiags := types.ObjectValue(systemKeyDataSourceAttrTypes, previousAttrVals)
	diags = append(diags, respDiags...)

	pendingAttrs := r.GetPending()
	pendingAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(pendingAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(pendingAttrs.GetEncryptedKeyData()),
	}
	pendingAttrsObjVal, respDiags := types.ObjectValue(systemKeyDataSourceAttrTypes, pendingAttrVals)
	diags = append(diags, respDiags...)

	state.Current = currentAttrsObjVal
	state.Pending = pendingAttrsObjVal
	state.Previous = previousAttrsObjVal
	return diags
}

func (r *serverSettingsSystemKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsSystemKeysModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadServerSettingsSystemKeys, httpResp, err := r.apiClient.ServerSettingsAPI.GetSystemKeys(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings System Keys", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readServerSettingsSystemKeysDataSourceResponse(ctx, apiReadServerSettingsSystemKeys, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
