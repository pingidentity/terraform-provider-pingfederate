package inboundprovisioninguserrepository

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func AttrType(key string) map[string]attr.Type {
	return map[string]attr.Type{
		key: types.ObjectType{
			AttrTypes: ElemAttrType(),
		},
	}
}

func IdentityStoreInboundProvisioningUserRepositoryAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                           types.StringType,
		"identity_store_provisioner_ref": types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
	}

}

func LdapInboundProvisioningUserRepositoryAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                   types.StringType,
		"data_store_ref":         types.ObjectType{AttrTypes: resourcelink.AttrTypeNoLocation()},
		"base_dn":                types.StringType,
		"unique_user_id_filter":  types.StringType,
		"unique_group_id_filter": types.StringType,
	}
}

func ElemAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"identity_store_inbound_provisioning_user_repository": types.ObjectType{
			AttrTypes: IdentityStoreInboundProvisioningUserRepositoryAttrType(),
		},
		"ldap_inbound_provisioning_user_repository": types.ObjectType{
			AttrTypes: LdapInboundProvisioningUserRepositoryAttrType(),
		},
	}
}
