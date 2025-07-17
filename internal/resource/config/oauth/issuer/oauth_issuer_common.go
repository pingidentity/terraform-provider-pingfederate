// Copyright © 2025 Ping Identity Corporation

package oauthissuer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1230/configurationapi"
)

type oauthIssuerModel struct {
	Id          types.String `tfsdk:"id"`
	IssuerId    types.String `tfsdk:"issuer_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Host        types.String `tfsdk:"host"`
	Path        types.String `tfsdk:"path"`
}

// Read a OauthIssuerResponse object into the model struct
func readOauthIssuerResponse(ctx context.Context, r *client.Issuer, state *oauthIssuerModel) {
	state.Id = types.StringPointerValue(r.Id)
	state.IssuerId = types.StringPointerValue(r.Id)
	state.Name = types.StringValue(r.Name)
	if r.Description != nil && *r.Description != "" {
		state.Description = types.StringPointerValue(r.Description)
	} else {
		state.Description = types.StringNull()
	}
	state.Host = types.StringValue(r.Host)
	if r.Path != nil && *r.Path != "" {
		state.Path = types.StringPointerValue(r.Path)
	} else {
		state.Path = types.StringNull()
	}
}
