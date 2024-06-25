package pluginconfiguration

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
)

func ClientStruct(configurationObj types.Object) (*client.PluginConfiguration, error) {
	configuration := client.NewPluginConfiguration()
	configErr := json.Unmarshal([]byte(internaljson.FromValue(configurationObj, true)), configuration)
	if configErr != nil {
		return nil, configErr
	}
	return configuration, nil
}
