// Copyright © 2025 Ping Identity Corporation

package scopeentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1220/configurationapi"
)

var exclusiveScopeAttrTypes = map[string]attr.Type{
	"name":        types.StringType,
	"description": types.StringType,
	"dynamic":     types.BoolType,
}

func ToState(con context.Context, scopes []client.ScopeEntry) (basetypes.SetValue, diag.Diagnostics) {
	toStateScopes := []client.ScopeEntry{}
	for _, scope := range scopes {
		scopeEntry := client.ScopeEntry{}
		scopeEntry.Name = scope.Name
		scopeEntry.Description = scope.Description
		scopeEntry.Dynamic = scope.Dynamic
		toStateScopes = append(toStateScopes, scopeEntry)
	}

	return types.SetValueFrom(con, types.ObjectType{AttrTypes: exclusiveScopeAttrTypes}, toStateScopes)
}

func AttrTypes() map[string]attr.Type {
	return exclusiveScopeAttrTypes
}
