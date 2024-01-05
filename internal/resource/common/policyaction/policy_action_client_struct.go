package policyaction

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1200/configurationapi"
)

func ClientStruct(types.Object) (*client.PolicyActionAggregation, error) {
	return nil, nil
}
