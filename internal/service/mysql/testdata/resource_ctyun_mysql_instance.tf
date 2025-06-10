resource "ctyun_mysql_instance" "%[1]s" {
  bill_mode = "%[2]s"
  prod_version = "%[3]s"
  vpc_id = "%[4]s"
  host_type = "%[5]s"
  subnet_id = "%[6]s"
  security_group_id = "%[7]s"
  name = "%[8]s"
  password = "%[9]s"
  period = %[10]d
  purchase_count = %[11]d
  auto_renew_status = %[12]d
  prod_id = %[13]d
  cpu_type = "%[14]s"
  os_type = "%[15]s"
  %[16]s  //write_port
  node_type = "%[17]s"
  inst_spec = "%[18]s"
  storage_type = "%[19]s"
  storage_space = %[20]d
  prod_performance_spec = "%[21]s"
  disks = %[22]d
  %[23]s  // availability_zone_info
  start = %[24]t
  restart = %[25]t
  stop = %[26]t
}
