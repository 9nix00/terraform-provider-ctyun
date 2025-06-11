resource "ctyun_postgresql_instance" "%[1]s" {
  bill_mode              = "%[2]s"
  host_type              = "%[3]s"
  prod_version           = "%[4]s"
  prod_id                =  %[5]d
  storage_type           = "%[6]s"
  storage_space          =  %[7]d
  name                   = "%[8]s"
  password               = "%[9]s"
  case_sensitive         = "%[10]s"
  node_type              = "%[11]s"
  inst_spec              = "%[12]s"
  prod_performance_spec  = "%[13]s"
  vpc_id                 = "%[14]s"
  subnet_id              = "%[15]s"
  security_group_id      = "%[16]s"
  availability_zone_info =  %[17]s
}
