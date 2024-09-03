package clusterstatus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/version"
)

func TestAccClusterStatus(t *testing.T) {
	// Check if the server is running in clustered mode or not
	inCluster := false
	testClient := acctest.TestClient()
	_, _, err := testClient.ClusterAPI.GetClusterStatus(acctest.TestBasicAuthContext()).Execute()
	if err == nil {
		// The API returned a status, so this server must be in clustered mode
		inCluster = true
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.NewTestProvider()),
		},
		Steps: []resource.TestStep{
			{
				// Run the export and validate the results
				Config: clusterStatus_MinimalHCL(),
				Check:  clusterStatus_CheckComputedValues(inCluster),
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
func clusterStatus_CheckComputedValues(inCluster bool) resource.TestCheckFunc {
	if inCluster {
		checks := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "mixed_mode", "false"),
			resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "replication_required"),
			resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "last_config_update_time"),
			resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "nodes.#", "1"),
			resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "nodes.0.address"),
			resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "nodes.0.index"),
			resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.mode", "CLUSTERED_CONSOLE"),
			resource.TestCheckResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.node_group", ""),
			resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "nodes.0.replication_status"),
			resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "nodes.0.version"),
			resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.admin_console_info"),
			resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "nodes.0.node_tags"),
		}
		if acctest.VersionAtLeast(version.PingFederate1210) {
			checks = append(checks, resource.TestCheckResourceAttrSet("data.pingfederate_cluster_status.example", "current_node_index"))
		}
		return resource.ComposeTestCheckFunc(checks...)
	}
	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "mixed_mode"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "replication_required"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "last_config_update_time"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "last_replication_time"),
		resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "nodes"),
	}
	if acctest.VersionAtLeast(version.PingFederate1210) {
		checks = append(checks, resource.TestCheckNoResourceAttr("data.pingfederate_cluster_status.example", "current_node_index"))
	}
	return resource.ComposeTestCheckFunc(checks...)
}
