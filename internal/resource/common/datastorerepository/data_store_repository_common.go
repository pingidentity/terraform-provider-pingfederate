package datastorerepository

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/attributecontractfulfillment"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func AttrType(key string) map[string]attr.Type {
	return map[string]attr.Type{
		key: types.ObjectType{
			AttrTypes: ElemAttrType(),
		},
	}
}

func JdbcDataStoreRepositoryAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"data_store_ref":                   types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"jit_repository_attribute_mapping": attributecontractfulfillment.MapType(),
		"sql_method": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"table": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"schema":           types.StringType,
						"table_name":       types.StringType,
						"unique_id_column": types.StringType,
					},
				},
				"stored_procedure": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"schema":           types.StringType,
						"stored_procedure": types.StringType,
					},
				},
			},
		},
	}

}

func LdapDataStoreRepositoryAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"data_store_ref":                   types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"jit_repository_attribute_mapping": attributecontractfulfillment.MapType(),
		"base_dn":                          types.StringType,
		"unique_user_id_filter":            types.StringType,
	}
}

func ElemAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"jdbc": types.ObjectType{
			AttrTypes: JdbcDataStoreRepositoryAttrType(),
		},
		"ldap": types.ObjectType{
			AttrTypes: LdapDataStoreRepositoryAttrType(),
		},
	}
}
