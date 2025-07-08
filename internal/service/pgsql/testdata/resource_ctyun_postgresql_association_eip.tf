resource "ctyun_postgresql_association_eip" "%[1]s" {
  eip_id = "%[2]s"
  eip    = "%[3]s"
  inst_id = "%[4]s"
}
