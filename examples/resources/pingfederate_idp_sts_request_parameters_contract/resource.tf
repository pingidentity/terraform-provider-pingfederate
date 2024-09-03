resource "pingfederate_idp_sts_request_parameters_contract" "requestParametersContract" {
  contract_id = "mycontract"
  name        = "My Contract"
  parameters = [
    "firstparam",
    "secondparam",
    "thirdparam"
  ]
}