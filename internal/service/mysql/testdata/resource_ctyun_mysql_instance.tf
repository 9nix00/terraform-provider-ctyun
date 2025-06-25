resource "ctyun_mysql_instance" "%[1]s" {
  cycle_type = "%[2]s"
  vpc_id = "%[3]s"
  host_type = "%[4]s"
  subnet_id = "%[5]s"
  security_group_id = "%[6]s"
  name = "%[7]s"
  password = "%[8]s"
  %[9]s // cycle_count
  %[10]s // auto_renew
  prod_id = "%[11]s"
  cpu_type = "%[12]s"
  os_type = "%[13]s"
  %[14]s  //write_port
  instance_series = "%[15]s"
  storage_type = "%[16]s"
  storage_space = %[17]d
  prod_performance_spec = "%[18]s"
  %[19]s  // availability_zone_info
  %[20]s // running_control
  %[21]s // backup_storage_space
}

