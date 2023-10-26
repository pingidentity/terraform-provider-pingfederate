package scopegroupentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
)

var scopeGroupAttrTypes = map[string]attr.Type{
	"name":        basetypes.StringType{},
	"description": basetypes.StringType{},
	"scopes":      basetypes.SetType{ElemType: types.StringType},
}

func ToState(con context.Context, scopeGroups []client.ScopeGroupEntry) (basetypes.ListValue, diag.Diagnostics) {
	toStateScopeGroups := []client.ScopeGroupEntry{}
	for _, scopeGroup := range scopeGroups {
		scopeGroupEntry := client.ScopeGroupEntry{}
		scopeGroupEntry.Name = scopeGroup.Name
		scopeGroupEntry.Description = scopeGroup.Description
		scopeGroupEntry.Scopes = scopeGroup.Scopes
		toStateScopeGroups = append(toStateScopeGroups, scopeGroupEntry)
	}

	return types.ListValueFrom(con, types.ObjectType{AttrTypes: scopeGroupAttrTypes}, toStateScopeGroups)
}
