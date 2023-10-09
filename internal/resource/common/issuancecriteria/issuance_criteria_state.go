package issuancecriteria

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/sourcetypeidkey"
)

func AttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"conditional_criteria": basetypes.ListType{
			ElemType: basetypes.ObjectType{
				AttrTypes: map[string]attr.Type{
					"source": basetypes.ObjectType{
						AttrTypes: sourcetypeidkey.AttrType(),
					},
					"attribute_name": basetypes.StringType{},
					"condition":      basetypes.StringType{},
					"value":          basetypes.StringType{},
					"error_result":   basetypes.StringType{},
				},
			},
		},
		"expression_criteria": basetypes.ListType{
			ElemType: basetypes.ObjectType{
				AttrTypes: map[string]attr.Type{
					"expression":   basetypes.StringType{},
					"error_result": basetypes.StringType{},
				},
			},
		},
	}
}

func ToState(con context.Context, issuanceCriteriaFromClient *client.IssuanceCriteria) basetypes.ObjectValue {
	issuanceCriteriaObj, _ := types.ObjectValueFrom(con, AttrType(), issuanceCriteriaFromClient)
	return issuanceCriteriaObj
}
