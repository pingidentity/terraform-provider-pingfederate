package attributesources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func CommonAttributeSourceAttrType() map[string]attr.Type {
	commonAttrSourceAttrType := make(map[string]attr.Type)
	commonAttrSourceAttrType["type"] = basetypes.StringType{}
	commonAttrSourceAttrType["data_store_ref"] = basetypes.ObjectType{AttrTypes: resourcelink.ResourceLinkStateAttrType()}
	commonAttrSourceAttrType["id"] = basetypes.StringType{}
	commonAttrSourceAttrType["description"] = basetypes.StringType{}
	commonAttrSourceAttrType["attribute_contract_fulfillment"] = attributecontractfulfillment.AttributeContractFulfillmentAttrType()
	return commonAttrSourceAttrType
}

func CustomAttributeSourceAttrType() map[string]attr.Type {
	customAttrSourceAttrType := CommonAttributeSourceAttrType()
	customAttrSourceAttrType["filter_fields"] = basetypes.ListType{ElemType: basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"value": basetypes.StringType{},
			"name":  basetypes.StringType{},
		},
	}}
	return customAttrSourceAttrType
}

func CustomAttributeSourceAttrVal(con context.Context, customAttrSourceValueFromClient client.CustomAttributeSource) basetypes.ObjectValue {
	customAttrSourceObj, _ := types.ObjectValueFrom(con, CustomAttributeSourceAttrType(), customAttrSourceValueFromClient)
	return customAttrSourceObj
}

func JdbcAttributeSourceAttrType() map[string]attr.Type {
	jdbcAttributeSourceAttrType := CommonAttributeSourceAttrType()
	jdbcAttributeSourceAttrType["schema"] = basetypes.StringType{}
	jdbcAttributeSourceAttrType["table"] = basetypes.StringType{}
	jdbcAttributeSourceAttrType["column_names"] = basetypes.ListType{ElemType: basetypes.StringType{}}
	jdbcAttributeSourceAttrType["filter"] = basetypes.StringType{}

	return jdbcAttributeSourceAttrType
}

func JdbcAttributeSourceAttrVal(con context.Context, jdbcAttrSourceValueFromClient client.JdbcAttributeSource) basetypes.ObjectValue {
	jdbcAttrSourceObj, _ := types.ObjectValueFrom(con, JdbcAttributeSourceAttrType(), jdbcAttrSourceValueFromClient)
	return jdbcAttrSourceObj
}

func LdapAttributeSourceAttrType() map[string]attr.Type {
	ldapAttrSourceAttrType := CommonAttributeSourceAttrType()
	ldapAttrSourceAttrType["base_dn"] = basetypes.StringType{}
	ldapAttrSourceAttrType["search_scope"] = basetypes.StringType{}
	ldapAttrSourceAttrType["search_filter"] = basetypes.StringType{}
	ldapAttrSourceAttrType["search_attributes"] = basetypes.ListType{ElemType: basetypes.StringType{}}
	ldapAttrSourceAttrType["binary_attribute_settings"] = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"binary_encoding": basetypes.StringType{},
		},
	}
	ldapAttrSourceAttrType["member_of_nested_group"] = basetypes.BoolType{}
	return ldapAttrSourceAttrType
}

func LdapAttributeSourceAttrVal(con context.Context, ldapAttrSourceValueFromClient client.LdapAttributeSource) basetypes.ObjectValue {
	ldapAttrSourceObj, _ := types.ObjectValueFrom(con, LdapAttributeSourceAttrType(), ldapAttrSourceValueFromClient)
	return ldapAttrSourceObj
}
