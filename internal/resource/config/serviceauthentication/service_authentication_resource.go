package serviceauthentication

import "github.com/hashicorp/terraform-plugin-framework/types"

func (state *serviceAuthenticationResourceModel) readClientResponseSharedSecret(existingParentValue types.Object) types.String {
	if existingParentValue.IsNull() || existingParentValue.IsUnknown() {
		return types.StringNull()
	}

	// Get the existing sharedSecret value
	attrs := existingParentValue.Attributes()
	sharedSecret, ok := attrs["shared_secret"]
	if !ok {
		return types.StringNull()
	}

	return types.StringValue(sharedSecret.(types.String).ValueString())
}
