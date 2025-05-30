resource "ctyun_nat_dnat" "%[1]s"{
    nat_gateway_id = "%[2]s"
    external_id = "%[3]s"
    external_port = "%[4]d"
    internal_ip = "%[6]s"
    virtual_machine_type = "%[5]d"
    internal_port = "%[7]d"
    protocol = "%[8]s"
}
