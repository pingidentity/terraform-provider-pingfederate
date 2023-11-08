package license

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/datasource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &licenseDataSource{}
	_ datasource.DataSourceWithConfigure = &licenseDataSource{}
)

// Create a Administrative Account data source
func NewLicenseDataSource() datasource.DataSource {
	return &licenseDataSource{}
}

// licenseDataSource is the datasource implementation.
type licenseDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type licenseDataSourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	MaxConnections      types.Int64  `tfsdk:"max_connections"`
	UsedConnections     types.Int64  `tfsdk:"used_connections"`
	Tier                types.String `tfsdk:"tier"`
	IssueDate           types.String `tfsdk:"issue_date"`
	ExpirationDate      types.String `tfsdk:"expiration_date"`
	EnforcementType     types.String `tfsdk:"enforcement_type"`
	Version             types.String `tfsdk:"version"`
	Product             types.String `tfsdk:"product"`
	Organization        types.String `tfsdk:"organization"`
	GracePeriod         types.Int64  `tfsdk:"grace_period"`
	NodeLimit           types.Int64  `tfsdk:"node_limit"`
	LicenseGroups       types.List   `tfsdk:"license_groups"`
	OauthEnabled        types.Bool   `tfsdk:"oauth_enabled"`
	WsTrustEnabled      types.Bool   `tfsdk:"ws_trust_enabled"`
	ProvisioningEnabled types.Bool   `tfsdk:"provisioning_enabled"`
	BridgeMode          types.Bool   `tfsdk:"bridge_mode"`
	Features            types.List   `tfsdk:"features"`
}

// GetSchema defines the schema for the datasource.
func (r *licenseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a License.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the person the license was issued to.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_connections": schema.Int64Attribute{
				Description: "Maximum number of connections that can be created under this license (if applicable).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"used_connections": schema.Int64Attribute{
				Description: "Number of used connections under this license.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"tier": schema.StringAttribute{
				Description: "The tier value from the license file. The possible values are FREE, PERPETUAL or SUBSCRIPTION.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"issue_date": schema.StringAttribute{
				Description: "The issue date value from the license file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"expiration_date": schema.StringAttribute{
				Description: "The expiration date value from the license file (if applicable).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enforcement_type": schema.StringAttribute{
				Description: "The enforcement type is a 3-bit binary value, expressed as a decimal digit. The bits from left to right are: 1: Shutdown on expire. 2: Notify on expire. 4: Enforce minor version. if all three enforcements are active, the enforcement type will be 7 (1 + 2 + 4); if only the first two are active, you have an enforcement type of 3 (1 + 2).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The Ping Identity product version from the license file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"product": schema.StringAttribute{
				Description: "The Ping Identity product value from the license file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"organization": schema.StringAttribute{
				Description: "The organization value from the license file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"grace_period": schema.Int64Attribute{
				Description: "Number of days provided as grace period, past the expiration date (if applicable).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"node_limit": schema.Int64Attribute{
				Description: "Maximum number of clustered nodes allowed under this license (if applicable).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"license_groups": schema.ListNestedAttribute{
				Description: "License connection groups, if applicable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Group name from the license file.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"connection_count": schema.Int64Attribute{
							Description: "Maximum number of connections permitted under the group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"start_date": schema.StringAttribute{
							Description: "Start date for the group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"end_date": schema.StringAttribute{
							Description: "End date for the group.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
			"oauth_enabled": schema.BoolAttribute{
				Description: "Indicates whether OAuth role is enabled for this license.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ws_trust_enabled": schema.BoolAttribute{
				Description: "Indicates whether WS-Trust role is enabled for this license.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"provisioning_enabled": schema.BoolAttribute{
				Description: "Indicates whether Provisioning role is enabled for this license.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bridge_mode": schema.BoolAttribute{
				Description: "Indicates whether this license is a bridge license or not.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"features": schema.ListNestedAttribute{
				Description: "Other licence features, if applicable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the license feature.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the license feature.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
		},
	}

	id.ToDataSourceSchema(&schemaDef, false, "Unique identifier of a license.")
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *licenseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license"
}

// Configure adds the provider configured client to the data source.
func (r *licenseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readLicenseResponseDataSource(ctx context.Context, r *client.LicenseView, state *licenseDataSourceModel) diag.Diagnostics {
	var diags, respDiags diag.Diagnostics
	state.Id = types.StringValue("id")
	state.Name = internaltypes.StringTypeOrNil(r.Name, false)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.UsedConnections = internaltypes.Int64TypeOrNil(r.UsedConnections)
	state.Tier = internaltypes.StringTypeOrNil(r.Tier, false)
	state.IssueDate = types.StringValue(r.IssueDate.Format(time.RFC3339))
	state.ExpirationDate = types.StringValue(r.ExpirationDate.Format(time.RFC3339))
	state.EnforcementType = internaltypes.StringTypeOrNil(r.EnforcementType, false)
	state.Version = internaltypes.StringTypeOrNil(r.Version, false)
	state.Product = internaltypes.StringTypeOrNil(r.Product, false)
	state.Organization = internaltypes.StringTypeOrNil(r.Organization, false)
	state.GracePeriod = internaltypes.Int64TypeOrNil(r.GracePeriod)
	state.NodeLimit = internaltypes.Int64TypeOrNil(r.NodeLimit)
	state.OauthEnabled = types.BoolValue(*r.OauthEnabled)
	state.WsTrustEnabled = types.BoolValue(*r.WsTrustEnabled)
	state.ProvisioningEnabled = types.BoolValue(*r.ProvisioningEnabled)
	state.BridgeMode = types.BoolValue(*r.BridgeMode)

	licenseGroups := r.LicenseGroups
	licenseGroupsAttrTypes := map[string]attr.Type{
		"name":             basetypes.StringType{},
		"connection_count": basetypes.Int64Type{},
		"start_date":       basetypes.StringType{},
		"end_date":         basetypes.StringType{},
	}
	state.LicenseGroups, respDiags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: licenseGroupsAttrTypes}, licenseGroups)
	diags.Append(respDiags...)

	features := r.Features
	featuresAttrTypes := map[string]attr.Type{
		"name":  basetypes.StringType{},
		"value": basetypes.StringType{},
	}
	state.Features, respDiags = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: featuresAttrTypes}, features)
	diags.Append(respDiags...)

	return diags
}

// Read resource information
func (r *licenseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state licenseDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadLicense, httpResp, err := r.apiClient.LicenseAPI.GetLicense(config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the License", err, httpResp)
		return
	}

	// Read the response into the state
	diags = readLicenseResponseDataSource(ctx, apiReadLicense, &state)
	resp.Diagnostics.Append(diags...)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
