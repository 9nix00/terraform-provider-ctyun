resource "ctyun_mysql_instance" "%[1]s" {
  cycle_type = "%[2]s"
  prod_version = "%[3]s"
  vpc_id = "%[4]s"
  host_type = "%[5]s"
  subnet_id = "%[6]s"
  security_group_id = "%[7]s"
  name = "%[8]s"
  password = "%[9]s"
  %[10]s // cycle_count
  %[11]s // auto_renew
  prod_id = "%[12]s"
  cpu_type = "%[13]s"
  os_type = "%[14]s"
  %[15]s  //write_port
  instance_series = "%[16]s"
  storage_type = "%[17]s"
  storage_space = %[18]d
  prod_performance_spec = "%[19]s"
  %[20]s  // availability_zone_info
  %[21]s // running_control
  %[22]s // backup_storage_space
}

