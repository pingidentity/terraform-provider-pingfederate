package scopeentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

var exclusiveScopeAttrTypes = map[string]attr.Type{
	"name":        basetypes.StringType{},
	"description": basetypes.StringType{},
	"dynamic":     basetypes.BoolType{},
}

func ToState(con context.Context, scopes []client.ScopeEntry) (basetypes.ListValue, diag.Diagnostics) {
	toStateScopes := []client.ScopeEntry{}
	for _, scope := range scopes {
		scopeEntry := client.ScopeEntry{}
		scopeEntry.Name = scope.Name
		scopeEntry.Description = scope.Description
		scopeEntry.Dynamic = scope.Dynamic
		toStateScopes = append(toStateScopes, scopeEntry)
	}

	return types.ListValueFrom(con, types.ObjectType{AttrTypes: exclusiveScopeAttrTypes}, toStateScopes)
}
