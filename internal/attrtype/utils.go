package internalmapstringattrtype

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func AddPasswordAttrType(mapStringAttrType map[string]attr.Type) map[string]attr.Type {
	mapStringAttrType["password"] = basetypes.StringType{}
	return mapStringAttrType
}

func AddEncryptedPasswordAttrType(mapStringAttrType map[string]attr.Type) map[string]attr.Type {
	mapStringAttrType["encrypted_password"] = basetypes.StringType{}
	return mapStringAttrType
}
