// Code generated by ping-terraform-plugin-framework-generator

package clusterstatus

import (
	"context"
	"io"
	"strings"
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

var (
	_ datasource.DataSource              = &clusterStatusDataSource{}
	_ datasource.DataSourceWithConfigure = &clusterStatusDataSource{}
)

func ClusterStatusDataSource() datasource.DataSource {
	return &clusterStatusDataSource{}
}

type clusterStatusDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *clusterStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_status"
}

func (r *clusterStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type clusterStatusDataSourceModel struct {
	CurrentNodeIndex     types.Int64  `tfsdk:"current_node_index"`
	LastConfigUpdateTime types.String `tfsdk:"last_config_update_time"`
	LastReplicationTime  types.String `tfsdk:"last_replication_time"`
	MixedMode            types.Bool   `tfsdk:"mixed_mode"`
	Nodes                types.List   `tfsdk:"nodes"`
	ReplicationRequired  types.Bool   `tfsdk:"replication_required"`
}

func (r *clusterStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource to retrive information on the current status of the cluster.",
		Attributes: map[string]schema.Attribute{
			"current_node_index": schema.Int64Attribute{
				Computed:    true,
				Description: "Index of the current node in the cluster.",
			},
			"last_config_update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time when the configuration of this node was last updated.",
			},
			"last_replication_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time when configuration changes were last replicated.",
			},
			"mixed_mode": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether there is more than one version of PingFederate in the cluster.",
			},
			"nodes": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Computed:    true,
							Description: "The IP address and port this node is running on.",
						},
						"admin_console_info": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"config_sync_status": schema.StringAttribute{
									Computed:    true,
									Description: "The status of the last configuration synchronization.",
								},
								"config_sync_timestamp": schema.StringAttribute{
									Computed:    true,
									Description: "The timestamp of the last configuration synchronization.",
								},
								"console_role": schema.StringAttribute{
									Computed:    true,
									Description: "For console nodes, indicates whether the node is active or passive.",
								},
								"console_role_last_update_date": schema.StringAttribute{
									Computed:    true,
									Description: "The timestamp of when the administrative console's role was last updated.",
								},
							},
							Computed:    true,
							Description: "The administrative console information when the active/passive administrative console feature is enabled.",
						},
						"configuration_timestamp": schema.StringAttribute{
							Computed:    true,
							Description: "The time stamp of the configuration data retrieved by this node.",
						},
						"index": schema.Int64Attribute{
							Computed:    true,
							Description: "Index of the node within the cluster, or `-1` if an index is not assigned.",
						},
						"mode": schema.StringAttribute{
							Computed:    true,
							Description: "The deployment mode of this node, from a clustering standpoint. `CLUSTERED_DUAL` is not supported.",
						},
						"node_group": schema.StringAttribute{
							Computed:    true,
							Description: "The node group for this node. This field is only populated if adaptive clustering is enabled.",
						},
						"node_tags": schema.StringAttribute{
							Computed:    true,
							Description: "The node tags for this node. This field is only populated for engine nodes.",
						},
						"replication_status": schema.StringAttribute{
							Computed:    true,
							Description: "The replication status of the node.",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "The PingFederate version this node is running on.",
						},
					},
				},
				Computed:    true,
				Description: "List of nodes in the cluster.",
			},
			"replication_required": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether a replication is required to propagate config updates.",
			},
		},
	}
}

func (state *clusterStatusDataSourceModel) readClientResponse(response *client.ClusterStatus) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// current_node_index
	state.CurrentNodeIndex = types.Int64PointerValue(response.CurrentNodeIndex)
	// last_config_update_time
	if response.LastConfigUpdateTime == nil {
		state.LastConfigUpdateTime = types.StringNull()
	} else {
		state.LastConfigUpdateTime = types.StringValue(response.LastConfigUpdateTime.Format(time.RFC3339))
	}
	// last_replication_time
	if response.LastReplicationTime == nil {
		state.LastReplicationTime = types.StringNull()
	} else {
		state.LastReplicationTime = types.StringValue(response.LastReplicationTime.Format(time.RFC3339))
	}
	// mixed_mode
	state.MixedMode = types.BoolPointerValue(response.MixedMode)
	// nodes
	nodesAdminConsoleInfoAttrTypes := map[string]attr.Type{
		"config_sync_status":            types.StringType,
		"config_sync_timestamp":         types.StringType,
		"console_role":                  types.StringType,
		"console_role_last_update_date": types.StringType,
	}
	nodesAttrTypes := map[string]attr.Type{
		"address":                 types.StringType,
		"admin_console_info":      types.ObjectType{AttrTypes: nodesAdminConsoleInfoAttrTypes},
		"configuration_timestamp": types.StringType,
		"index":                   types.Int64Type,
		"mode":                    types.StringType,
		"node_group":              types.StringType,
		"node_tags":               types.StringType,
		"replication_status":      types.StringType,
		"version":                 types.StringType,
	}
	nodesElementType := types.ObjectType{AttrTypes: nodesAttrTypes}
	var nodesValues []attr.Value
	for _, nodesResponseValue := range response.Nodes {
		var nodesAdminConsoleInfoValue types.Object
		if nodesResponseValue.AdminConsoleInfo == nil {
			nodesAdminConsoleInfoValue = types.ObjectNull(nodesAdminConsoleInfoAttrTypes)
		} else {
			var configSyncTimestampValue types.String
			if nodesResponseValue.AdminConsoleInfo.ConfigSyncTimestamp == nil {
				configSyncTimestampValue = types.StringNull()
			} else {
				configSyncTimestampValue = types.StringValue(nodesResponseValue.AdminConsoleInfo.ConfigSyncTimestamp.Format(time.RFC3339))
			}
			var consoleRoleLastUpdateDateValue types.String
			if nodesResponseValue.AdminConsoleInfo.ConsoleRoleLastUpdateDate == nil {
				consoleRoleLastUpdateDateValue = types.StringNull()
			} else {
				consoleRoleLastUpdateDateValue = types.StringValue(nodesResponseValue.AdminConsoleInfo.ConsoleRoleLastUpdateDate.Format(time.RFC3339))
			}
			nodesAdminConsoleInfoValue, diags = types.ObjectValue(nodesAdminConsoleInfoAttrTypes, map[string]attr.Value{
				"config_sync_status":            types.StringPointerValue(nodesResponseValue.AdminConsoleInfo.ConfigSyncStatus),
				"config_sync_timestamp":         configSyncTimestampValue,
				"console_role":                  types.StringPointerValue(nodesResponseValue.AdminConsoleInfo.ConsoleRole),
				"console_role_last_update_date": consoleRoleLastUpdateDateValue,
			})
			respDiags.Append(diags...)
		}

		var configurationTimestampValue types.String
		if nodesResponseValue.ConfigurationTimestamp == nil {
			configurationTimestampValue = types.StringNull()
		} else {
			configurationTimestampValue = types.StringValue(nodesResponseValue.ConfigurationTimestamp.Format(time.RFC3339))
		}
		nodesValue, diags := types.ObjectValue(nodesAttrTypes, map[string]attr.Value{
			"address":                 types.StringPointerValue(nodesResponseValue.Address),
			"admin_console_info":      nodesAdminConsoleInfoValue,
			"configuration_timestamp": configurationTimestampValue,
			"index":                   types.Int64PointerValue(nodesResponseValue.Index),
			"mode":                    types.StringPointerValue(nodesResponseValue.Mode),
			"node_group":              types.StringPointerValue(nodesResponseValue.NodeGroup),
			"node_tags":               types.StringPointerValue(nodesResponseValue.NodeTags),
			"replication_status":      types.StringPointerValue(nodesResponseValue.ReplicationStatus),
			"version":                 types.StringPointerValue(nodesResponseValue.Version),
		})
		respDiags.Append(diags...)
		nodesValues = append(nodesValues, nodesValue)
	}
	nodesValue, diags := types.ListValue(nodesElementType, nodesValues)
	respDiags.Append(diags...)

	state.Nodes = nodesValue
	// replication_required
	state.ReplicationRequired = types.BoolPointerValue(response.ReplicationRequired)
	return respDiags
}

func (state *clusterStatusDataSourceModel) emptyModel() {
	// current_node_index
	state.CurrentNodeIndex = types.Int64Null()
	// last_config_update_time
	state.LastConfigUpdateTime = types.StringNull()
	// last_replication_time
	state.LastReplicationTime = types.StringNull()
	// mixed_mode
	state.MixedMode = types.BoolNull()
	// nodes
	nodesAdminConsoleInfoAttrTypes := map[string]attr.Type{
		"config_sync_status":            types.StringType,
		"config_sync_timestamp":         types.StringType,
		"console_role":                  types.StringType,
		"console_role_last_update_date": types.StringType,
	}
	nodesAttrTypes := map[string]attr.Type{
		"address":                 types.StringType,
		"admin_console_info":      types.ObjectType{AttrTypes: nodesAdminConsoleInfoAttrTypes},
		"configuration_timestamp": types.StringType,
		"index":                   types.Int64Type,
		"mode":                    types.StringType,
		"node_group":              types.StringType,
		"node_tags":               types.StringType,
		"replication_status":      types.StringType,
		"version":                 types.StringType,
	}
	nodesElementType := types.ObjectType{AttrTypes: nodesAttrTypes}
	state.Nodes = types.ListNull(nodesElementType)
	// replication_required
	state.ReplicationRequired = types.BoolNull()
}

func (r *clusterStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read API call logic
	responseData, httpResp, err := r.apiClient.ClusterAPI.GetClusterStatus(config.AuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		// If the error indicates that this PF server is not running in clustered mode, just return an empty model
		isNonClusteredMode := false
		if httpResp != nil {
			body, internalError := io.ReadAll(httpResp.Body)
			if internalError == nil {
				bodyContents := string(body)
				if strings.Contains(bodyContents, "not deployed in clustered mode") {
					var data clusterStatusDataSourceModel
					data.emptyModel()
					// Save updated data into Terraform state
					resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
					isNonClusteredMode = true
				}
			}
		}
		// Otherwise, report the error
		if !isNonClusteredMode {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the cluster status", err, httpResp)
		}
		return
	}

	// Read response into the model
	var data clusterStatusDataSourceModel
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
