package serversettingssystemkeys

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serverSettingsSystemKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &serverSettingsSystemKeysDataSource{}
)

// ServerSettingsSystemKeysDataSource is a helper function to simplify the provider implementation.
func ServerSettingsSystemKeysDataSource() datasource.DataSource {
	return &serverSettingsSystemKeysDataSource{}
}

// serverSettingsSystemKeysDataSource is the resource implementation.
type serverSettingsSystemKeysDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type serverSettingsSystemKeysDataSourceModel struct {
	Id       types.String `tfsdk:"id"`
	Current  types.Object `tfsdk:"current"`
	Previous types.Object `tfsdk:"previous"`
	Pending  types.Object `tfsdk:"pending"`
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
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
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
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
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
					"key_data": schema.StringAttribute{
						Description: "The clear text system key base 64 encoded. The system key must be 32 bytes before base 64 encoding",
						Computed:    true,
						Optional:    false,
					},
				},
			},
		},
	}

	id.ToDataSourceSchema(&schema)
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

func readServerSettingsSystemKeysDataSource(ctx context.Context, r *client.SystemKeys, state *serverSettingsSystemKeysDataSourceModel, existingId *string) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Id = types.StringValue("server_settings_system_keys_id")
	currentAttrs := r.GetCurrent()
	currentAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(currentAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(currentAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(currentAttrs.GetKeyData()),
	}
	currentAttrsObjVal, respDiags := types.ObjectValue(systemKeyAttrTypes, currentAttrVals)
	diags = append(diags, respDiags...)

	previousAttrs := r.GetPrevious()
	previousAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(previousAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(previousAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(previousAttrs.GetKeyData()),
	}
	previousAttrsObjVal, respDiags := types.ObjectValue(systemKeyAttrTypes, previousAttrVals)
	diags = append(diags, respDiags...)

	pendingAttrs := r.GetPending()
	pendingAttrVals := map[string]attr.Value{
		"creation_date":      types.StringValue(pendingAttrs.GetCreationDate().Format(time.RFC3339Nano)),
		"encrypted_key_data": types.StringValue(pendingAttrs.GetEncryptedKeyData()),
		"key_data":           types.StringValue(pendingAttrs.GetKeyData()),
	}
	pendingAttrsObjVal, respDiags := types.ObjectValue(systemKeyAttrTypes, pendingAttrVals)
	diags = append(diags, respDiags...)

	state.Current = currentAttrsObjVal
	state.Pending = pendingAttrsObjVal
	state.Previous = previousAttrsObjVal
	return diags
}

func (r *serverSettingsSystemKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serverSettingsSystemKeysDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadServerSettingsSystemKeys, httpResp, err := r.apiClient.ServerSettingsAPI.GetSystemKeys(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Settings System Keys", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readServerSettingsSystemKeysDataSource(ctx, apiReadServerSettingsSystemKeys, &state, nil)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
