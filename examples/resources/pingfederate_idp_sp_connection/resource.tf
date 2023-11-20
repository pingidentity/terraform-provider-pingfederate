resource "pingfederate_idp_sp_connection" "idpSpConnectionExample" {
    connection_id = "myConnection"
    name = "myConnection"
    entity_id = "myEntity"
    type = "SP"
}
