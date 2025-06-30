resource "ctyun_postgresql_instance" "%[1]s" {
  cycle_type              = "%[2]s"
  host_type              = "%[3]s"
  prod_id                =  "%[4]s"
  storage_type           = "%[5]s"
  storage_space          =  %[6]d
  name                   = "%[7]s"
  password               = "%[8]s"
  case_sensitive         = %[9]t
  instance_series        = "%[10]s"
  prod_performance_spec  = "%[11]s"
  vpc_id                 = "%[12]s"
  subnet_id              = "%[13]s"
  security_group_id      = "%[14]s"
  availability_zone_info =  %[15]s
  %[16]s  // backup_storage_space
  %[17]s // running_control = start, stop, restart
  os_type = "%[18]s"
  cpu_type = "%[19]s"
  %[20]s // cycle_count
}

