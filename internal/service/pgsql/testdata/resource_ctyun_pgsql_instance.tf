resource "ctyun_postgresql_instance" "%[1]s" {
  cycle_type              = "%[2]s"
  host_type              = "%[3]s"
  prod_version           = "%[4]s"
  prod_id                =  "%[5]s"
  storage_type           = "%[6]s"
  storage_space          =  %[7]d
  name                   = "%[8]s"
  password               = "%[9]s"
  case_sensitive         = %[10]t
  instance_series        = "%[11]s"
  prod_performance_spec  = "%[12]s"
  vpc_id                 = "%[13]s"
  subnet_id              = "%[14]s"
  security_group_id      = "%[15]s"
  availability_zone_info =  %[16]s
  %[17]s  // backup_storage_space
  %[18]s // running_control = start, stop, restart
  os_type = "%[19]s"
  cpu_type = "%[20]s"
  %[21]s // cycle_count
}

