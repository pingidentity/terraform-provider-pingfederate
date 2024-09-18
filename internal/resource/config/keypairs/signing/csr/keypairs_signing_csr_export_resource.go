package keypairssigningcsr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource              = &keypairsSigningCsrExportResource{}
	_ resource.ResourceWithConfigure = &keypairsSigningCsrExportResource{}

	customId = "keypair_id"
)

func KeypairsSigningCsrExportResource() resource.Resource {
	return &keypairsSigningCsrExportResource{}
}

type keypairsSigningCsrExportResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *keypairsSigningCsrExportResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs_signing_csr_export"
}

func (r *keypairsSigningCsrExportResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keypairsSigningCsrExportResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	KeypairId           types.String `tfsdk:"keypair_id"`
	ExportedCsr         types.String `tfsdk:"exported_csr"`
	ExportTriggerValues types.Map    `tfsdk:"export_trigger_values"`
}

func (r *keypairsSigningCsrExportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource to generate a new certificate signing request (CSR) for a key pair.",
		Attributes: map[string]schema.Attribute{
			"keypair_id": schema.StringAttribute{
				Description: "The ID of the keypair.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"exported_csr": schema.StringAttribute{
				Description: "The exported PEM-encoded certificate signing request.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"export_trigger_values": schema.MapAttribute{
				Description: "A meta-argument map of values that, if any values are changed, will force export of a new CSR. Adding values to and removing values from the map will not trigger an export. This parameter can be used to control time-based exports using Terraform.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

// Export a new CSR via RequiresReplace when the trigger values change
func (r *keypairsSigningCsrExportResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Destruction plan
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state types.Map
	var planValues, stateValues map[string]attr.Value

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("export_trigger_values"), &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planValues = plan.Elements()

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("export_trigger_values"), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateValues = state.Elements()

	for k, v := range planValues {
		if stateValue, ok := stateValues[k]; ok && (v == types.StringUnknown() || !stateValue.Equal(v)) {
			resp.RequiresReplace = path.Paths{path.Root("export_trigger_values")}
			break
		}
	}
}

func (r *keypairsSigningCsrExportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data keypairsSigningCsrExportResourceModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	exportRequest := r.apiClient.KeyPairsSigningAPI.ExportCsr(config.AuthContext(ctx, r.providerConfig), data.KeypairId.ValueString())
	responseData, httpResp, err := exportRequest.Execute()
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while generating the certificate signing request.", err, httpResp, &customId)
		return
	}

	// Set the exported metadata
	data.Id = types.StringValue(data.KeypairId.ValueString())
	data.ExportedCsr = types.StringValue(responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *keypairsSigningCsrExportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// PingFederate provides no read endpoint for this resource, so we'll just maintain whatever is in state
	resp.State.Raw = req.State.Raw
}

func (r *keypairsSigningCsrExportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This will only happen when adding or removing export trigger values.
	// Just copy the existing state and export_trigger_values into state.
	var plan, state keypairsSigningCsrExportResourceModel

	// Read Terraform config data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.ExportTriggerValues = plan.ExportTriggerValues

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *keypairsSigningCsrExportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// There is no way to delete an exported CSR
	providererror.WarnConfigurationCannotBeReset("pingfederate_keypairs_signing_csr_export", &resp.Diagnostics)
}
