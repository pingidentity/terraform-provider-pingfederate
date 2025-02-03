// Copyright © 2025 Ping Identity Corporation

package oauthaccesstokenmanager

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/resourcelink"
)

var (
	selectionSettingsAttrType = map[string]attr.Type{
		"resource_uris": types.SetType{ElemType: types.StringType},
	}

	accessControlSettingsAttrType = map[string]attr.Type{
		"restrict_clients": types.BoolType,
		"allowed_clients":  types.ListType{ElemType: types.ObjectType{AttrTypes: resourcelink.AttrType()}},
	}

	sessionValidationSettingsAttrType = map[string]attr.Type{
		"include_session_id":              types.BoolType,
		"check_valid_authn_session":       types.BoolType,
		"check_session_revocation_status": types.BoolType,
		"update_authn_session_activity":   types.BoolType,
	}

	resourceUrisDefault, _      = types.SetValue(types.StringType, nil)
	selectionSettingsDefault, _ = types.ObjectValue(selectionSettingsAttrType, map[string]attr.Value{
		"resource_uris": resourceUrisDefault,
	})

	allowedClientsDefault, _        = types.ListValue(types.ObjectType{AttrTypes: resourcelink.AttrType()}, nil)
	accessControlSettingsDefault, _ = types.ObjectValue(accessControlSettingsAttrType, map[string]attr.Value{
		"restrict_clients": types.BoolValue(false),
		"allowed_clients":  allowedClientsDefault,
	})

	sessionValidationSettingsDefault, _ = types.ObjectValue(sessionValidationSettingsAttrType, map[string]attr.Value{
		"include_session_id":              types.BoolValue(false),
		"check_valid_authn_session":       types.BoolValue(false),
		"check_session_revocation_status": types.BoolValue(false),
		"update_authn_session_activity":   types.BoolValue(false),
	})
)
