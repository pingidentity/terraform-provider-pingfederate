package administrativeaccount

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1125/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

type administrativeAccountModel struct {
	Active            types.Bool   `tfsdk:"active"`
	Auditor           types.Bool   `tfsdk:"auditor"`
	Department        types.String `tfsdk:"department"`
	Description       types.String `tfsdk:"description"`
	EmailAddress      types.String `tfsdk:"email_address"`
	Id                types.String `tfsdk:"id"`
	EncryptedPassword types.String `tfsdk:"encrypted_password"`
	Password          types.String `tfsdk:"password"`
	PhoneNumber       types.String `tfsdk:"phone_number"`
	Roles             types.Set    `tfsdk:"roles"`
	Username          types.String `tfsdk:"username"`
}

// Read a AdministrativeAccountResponse object into the model struct
func readAdministrativeAccountResponse(ctx context.Context, r *client.AdministrativeAccount, state *administrativeAccountModel, plan *administrativeAccountModel) {
	state.Id = types.StringValue(r.Username)
	state.Username = types.StringValue(r.Username)
	// state.Password and state.EncryptedPassword
	if plan != nil {
		state.Password = types.StringValue(plan.Password.ValueString())
		if internaltypes.IsDefined(plan.EncryptedPassword) {
			state.EncryptedPassword = types.StringValue(plan.EncryptedPassword.ValueString())
		} else {
			state.EncryptedPassword = types.StringPointerValue(r.EncryptedPassword)
		}
	} else {
		state.Password = types.StringValue("")
		state.EncryptedPassword = types.StringPointerValue(r.EncryptedPassword)
	}
	state.Active = types.BoolPointerValue(r.Active)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Auditor = types.BoolPointerValue(r.Auditor)
	state.PhoneNumber = internaltypes.StringTypeOrNil(r.PhoneNumber, false)
	state.EmailAddress = internaltypes.StringTypeOrNil(r.EmailAddress, false)
	state.Department = internaltypes.StringTypeOrNil(r.Department, false)
	state.Roles = internaltypes.GetStringSet(r.Roles)
}
