package clusterstatus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

func TestAccClusterStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: clusterStatus_MinimalHCL(),
				Check:  clusterStatus_CheckComputedValues(),
			},
		},
	})
}

// Only the ca_id attribute can be set on this resource
func clusterStatus_MinimalHCL() string {
	return `
data "pingfederate_cluster_status" "example" {
}
`
}

// Validate any computed values when applying HCL
func clusterStatus_CheckComputedValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "mixed_mode", "false"),
		resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "replication_required", "true"),
		resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "current_node_index"),
		resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "last_config_update_time"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "last_replication_time"),
		resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "nodes.#", "1"),
		resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "nodes.0.address"),
		resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "nodes.0.index"),
		resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.mode", "CLUSTERED_CONSOLE"),
		resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.node_group", ""),
		resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.replication_status", "OUT_OF_DATE"),
		resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "nodes.0.version"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.admin_console_info"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.configuration_timestamp"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.node_tags"),
	)
}
