// Copyright © 2026 Ping Identity Corporation

package sptargeturlmappings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
)

func (m *spTargetUrlMappingsResourceModel) setNullObjectValues() {
	itemsRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	itemsAttrTypes := map[string]attr.Type{
		"ref":  types.ObjectType{AttrTypes: itemsRefAttrTypes},
		"type": types.StringType,
		"url":  types.StringType,
	}
	itemsElementType := types.ObjectType{AttrTypes: itemsAttrTypes}
	m.Items = types.ListNull(itemsElementType)
}

func (r *spTargetUrlMappingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this method will reset the resource to the default configuration in PingFederate.
	var data spTargetUrlMappingsResourceModel
	emptyClientStruct := data.buildClientStruct()
	apiUpdateRequest := r.apiClient.SpTargetUrlMappingsAPI.UpdateSpUrlMappings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(*emptyClientStruct)
	_, httpResp, err := r.apiClient.SpTargetUrlMappingsAPI.UpdateSpUrlMappingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the spTargetUrlMappings", err, httpResp)
	}
}
