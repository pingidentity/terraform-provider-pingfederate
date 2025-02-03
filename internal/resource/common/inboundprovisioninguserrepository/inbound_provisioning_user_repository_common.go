// Copyright © 2025 Ping Identity Corporation

package inboundprovisioninguserrepository

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

func IdentityStoreInboundProvisioningUserRepositoryAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"identity_store_provisioner_ref": types.ObjectType{AttrTypes: resourcelink.AttrType()},
	}

}

func LdapInboundProvisioningUserRepositoryAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"data_store_ref":         types.ObjectType{AttrTypes: resourcelink.AttrType()},
		"base_dn":                types.StringType,
		"unique_user_id_filter":  types.StringType,
		"unique_group_id_filter": types.StringType,
	}
}

func ElemAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"identity_store": types.ObjectType{
			AttrTypes: IdentityStoreInboundProvisioningUserRepositoryAttrType(),
		},
		"ldap": types.ObjectType{
			AttrTypes: LdapInboundProvisioningUserRepositoryAttrType(),
		},
	}
}
