resource "ctyun_private_nat_snat" "%[1]s"{
  nat_gateway_id = "%[2]s"
  source_subnet_id= "%[3]s"
  snat_ips = %[4]s
  description = "%[5]s"
}
